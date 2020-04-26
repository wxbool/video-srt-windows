package app

import (
	"bytes"
	"errors"
	"github.com/buger/jsonparser"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
	"videosrt/app/aliyun"
	"videosrt/app/ffmpeg"
	"videosrt/app/tool"
	"videosrt/app/translate"
)

//应用翻译配置
type VideoSrtTranslateStruct struct {
	TranslateSwitch bool //字幕翻译开关
	Supplier int //引擎供应商
	BilingualSubtitleSwitch bool //是否输出双语字幕
	InputLanguage int //输入字幕语言
	OutputLanguage int //输出字幕语言
	OutputMainSubtitleInputLanguage bool //双语主字幕（输入语言）

	BaiduTranslate translate.BaiduTranslate //百度翻译
	TengxunyunTranslate translate.TengxunyunTranslate //腾讯云翻译
}

//统一翻译结果
type VideoSrtTranslateResult struct {
	From string //翻译源语言
	To string //译文语言
	TransResultSrc string //翻译结果（原文）
	TransResultDst string //翻译结果（译文）
}

//主应用
type VideoSrt struct {
	Ffmpeg ffmpeg.Ffmpeg
	AliyunOss aliyun.AliyunOss //oss
	AliyunClound aliyun.AliyunClound //语音识别引擎

	CloseAutoDeleteOssTempFile bool //关闭自动删除临时音频文件（默认开启）[false开启 true关闭]
	//IntelligentBlock bool //智能分段处理 [废弃2020.4.7]
	TempDir string //临时文件目录
	AppDir string //应用根目录
	SrtDir string //字幕文件输出目录
	OutputType *AppSetingsOutput //输出文件类型
	OutputEncode int //输出文件编码
	SoundTrack int //输出音轨（0输出全部音轨）

	MaxConcurrency int //最大处理并发数

	TranslateCfg *VideoSrtTranslateStruct //翻译配置

	LogHandler func(s string , video string) //日志回调
	SuccessHandler func(video string) //成功回调
	FailHandler func(video string) //失败回调
}


//获取应用
func NewApp(appDir string) *VideoSrt {
	app := new(VideoSrt)

	app.TempDir = "temp/audio"
	app.AppDir = appDir
	app.OutputEncode = OUTPUT_ENCODE_UTF8 //默认输出文件编码
	app.MaxConcurrency = 2

	//实例应用翻译配置
	app.TranslateCfg = new(VideoSrtTranslateStruct)
	return app
}


//应用加载配置
func (app *VideoSrt) InitAppConfig(oss *AliyunOssCache , engine *AliyunEngineCache) {
	//oss
	app.AliyunOss.Endpoint = oss.Endpoint
	app.AliyunOss.AccessKeyId = oss.AccessKeyId
	app.AliyunOss.AccessKeySecret = oss.AccessKeySecret
	app.AliyunOss.BucketName = oss.BucketName
	app.AliyunOss.BucketDomain = oss.BucketDomain

	//engine
	app.AliyunClound.AppKey = engine.AppKey
	app.AliyunClound.AccessKeyId = engine.AccessKeyId
	app.AliyunClound.AccessKeySecret = engine.AccessKeySecret
	app.AliyunClound.Region = engine.Region
}


//加载翻译设置
func (app *VideoSrt) InitTranslateConfig (translateSettings *VideoSrtTranslateStruct) {
	app.TranslateCfg = translateSettings
}


func (app *VideoSrt) SetCloseAutoDeleteOssTempFile(state bool)  {
	app.CloseAutoDeleteOssTempFile = state
}

func (app *VideoSrt) SetSrtDir(dir string)  {
	app.SrtDir = dir
}

func (app *VideoSrt) SetOutputType(output *AppSetingsOutput)  {
	app.OutputType = output
}

func (app *VideoSrt) SetOutputEncode(encode int)  {
	app.OutputEncode = encode
}

func (app *VideoSrt) SetSoundTrack(track int)  {
	app.SoundTrack = track
}

func (app *VideoSrt) SetMaxConcurrency(number int)  {
	if number == 0 {
		number = 2
	}
	app.MaxConcurrency = number
}

func (app *VideoSrt) SetSuccessHandler(callback func(video string))  {
	app.SuccessHandler = callback
}

func (app *VideoSrt) SetFailHandler(callback func(video string))  {
	app.FailHandler = callback
}

func (app *VideoSrt) SetLogHandler(callback func(s string , video string))  {
	app.LogHandler = callback
}

//输出日志
func (app *VideoSrt) Log(s string , v string)  {
	app.LogHandler(s , v)
}


//清空临时目录
func (app *VideoSrt) ClearTempDir()  {
	//临时目录
	tmpAudioDir := app.AppDir + "/" + app.TempDir
	if tool.DirExists(tmpAudioDir) {
		//清空
		if remove := os.RemoveAll(tmpAudioDir); remove != nil {
			app.Log("清空临时文件夹失败，请手动操作" , "警告")
		}
	}
}


//应用运行
func (app *VideoSrt) Run(video string) {
	var tmpAudio string = ""

	//致命错误捕获
	defer func() {
		//拦截panic
		if err := recover(); err != nil {
			//失败回调
			app.FailHandler(video)
			//fmt.Println( err )

			e , ok := err.(error)
			if ok {
				app.Log("错误：" + e.Error() , video)
			} else {
				panic(err)
			}
		}
	}()

	//校验媒体文件
	if tool.VaildFile(video) != true {
		panic("媒体文件不存在")
	}

	tmpAudioDir := app.AppDir + "/" + app.TempDir
	if !tool.DirExists(tmpAudioDir) {
		//创建目录
		if err := tool.CreateDir(tmpAudioDir , false); err != nil {
			panic(err)
		}
	}
	tmpAudioFile := tool.GetRandomCodeString(15) + ".mp3"
	tmpAudio = tmpAudioDir + "/" + tmpAudioFile //临时音频文件

	app.Log("提取音频文件 ..." , video)

	//分离/转换媒体音频
	ExtractVideoAudio(video , tmpAudio)

	app.Log("上传音频文件 ..." , video)

	//上传音频至OSS
	tempfile := UploadAudioToClound(app.AliyunOss , tmpAudio)
	//获取完整链接
	filelink := app.AliyunOss.GetObjectFileUrl(tempfile)

	app.Log("上传文件成功 , 识别中 ..." , video)

	defer func() {
		//清理oss音频文件
		if app.CloseAutoDeleteOssTempFile == false {
			if e := DeleteOssCloundTempAudio(app.AliyunOss , tempfile); e!=nil {
				app.Log("OSS临时音频清理失败，建议手动删除" , video)
			} else {
				app.Log("OSS临时音频清理成功" , video)
			}
		}
	}()

	//阿里云录音文件识别
	AudioResult , IntelligentBlockResult := AliyunAudioRecognition(app.AliyunClound, filelink)

	app.Log("文件识别成功 , 字幕处理中 ..." , video)

	//翻译字幕块
	AliyunAudioResultTranslate(app , video , AudioResult , IntelligentBlockResult)

	//输出文件
	if app.OutputType.SRT {
		AliyunAudioResultMakeSubtitleFile(app , video , OUTPUT_SRT , AudioResult , IntelligentBlockResult)
	}
	if app.OutputType.LRC {
		AliyunAudioResultMakeSubtitleFile(app , video , OUTPUT_LRC , AudioResult , IntelligentBlockResult)
	}
	if app.OutputType.TXT {
		AliyunAudioResultMakeSubtitleFile(app , video , OUTPUT_STRING , AudioResult , IntelligentBlockResult)
	}

	app.Log("处理完成" , video)

	//删除临时文件
	if tmpAudio != "" {
		if remove := os.Remove(tmpAudio); remove != nil {
			app.Log("删除临时文件失败，请手动删除" , video)
		}
	}

	//成功回调
	app.SuccessHandler(video)
}



//统一运行翻译结果
func (app *VideoSrt) RunTranslate(s string) (*VideoSrtTranslateResult , error) {
	translateResult := new(VideoSrtTranslateResult)

	if app.TranslateCfg.Supplier == TRANSLATE_SUPPLIER_BAIDU {
		if app.TranslateCfg.BaiduTranslate.AuthenType == translate.ACCOUNT_COMMON_AUTHEN { //百度翻译标准版
			//休眠1010毫秒
			time.Sleep(time.Millisecond * 1010)
		} else {
			//休眠200毫秒
			time.Sleep(time.Millisecond * 200)
		}

		from := GetLanguageChar(app.TranslateCfg.InputLanguage , TRANSLATE_SUPPLIER_BAIDU)
		to := GetLanguageChar(app.TranslateCfg.OutputLanguage , TRANSLATE_SUPPLIER_BAIDU)

		baiduResult,e := app.TranslateCfg.BaiduTranslate.TranslateBaidu(s , from , to)
		if e != nil {
			return translateResult,errors.New("翻译失败！错误信息：" + e.Error())
		}
		translateResult.TransResultDst = baiduResult.TransResultDst
		translateResult.TransResultSrc = baiduResult.TransResultSrc
		translateResult.From = baiduResult.From
		translateResult.To = baiduResult.To

		return translateResult,nil
	} else if app.TranslateCfg.Supplier == TRANSLATE_SUPPLIER_TENGXUNYUN {
		//休眠
		t :=  app.TranslateCfg.TengxunyunTranslate.TranslateSleepTime(app.MaxConcurrency)
		time.Sleep(t)

		from := GetLanguageChar(app.TranslateCfg.InputLanguage , TRANSLATE_SUPPLIER_TENGXUNYUN)
		to := GetLanguageChar(app.TranslateCfg.OutputLanguage , TRANSLATE_SUPPLIER_TENGXUNYUN)

		txResult,e := app.TranslateCfg.TengxunyunTranslate.TranslateTengxunyun(s , from , to)
		if e != nil {
			return translateResult,errors.New("翻译失败！错误信息：" + e.Error())
		}

		translateResult.TransResultDst = txResult.TransResultDst
		translateResult.TransResultSrc = txResult.TransResultSrc
		translateResult.From = txResult.From
		translateResult.To = txResult.To

		return translateResult,nil
	} else {
		return translateResult,errors.New("翻译失败！缺少翻译引擎！")
	}
}


//提取视频音频文件
func ExtractVideoAudio(video string , tmpAudio string) {
	if err := ffmpeg.ExtractAudio(video , tmpAudio); err != nil {
		panic(err)
	}
}


//上传音频至oss
func UploadAudioToClound(target aliyun.AliyunOss , audioFile string) string {
	name := ""
	//提取文件名称
	if fileInfo, e := os.Stat(audioFile);e != nil {
		panic(e)
	} else {
		name = fileInfo.Name()
	}

	//上传
	if file , e := target.UploadFile(audioFile , name); e != nil {
		panic(e)
	} else {
		return file
	}
}


//清理临时音频文件
func DeleteOssCloundTempAudio (target aliyun.AliyunOss , objectName string) error {
	//删除
	if e := target.DeleteFile(objectName); e != nil {
		return e
	}
	return nil
}


//阿里云录音文件识别
func AliyunAudioRecognition(engine aliyun.AliyunClound , filelink string) (AudioResult map[int64][] *aliyun.AliyunAudioRecognitionResult , IntelligentBlockResult map[int64][] *aliyun.AliyunAudioRecognitionResult) {
	//创建识别请求
	taskid, client, e := engine.NewAudioFile(filelink)
	if e != nil {
		panic(e)
	}

	AudioResult = make(map[int64][] *aliyun.AliyunAudioRecognitionResult)
	IntelligentBlockResult = make(map[int64][] *aliyun.AliyunAudioRecognitionResult)

	//遍历获取识别结果
	resultError := engine.GetAudioFileResult(taskid , client , func(result []byte) {
		//mylog.WriteLog( string( result ) )

		//结果处理
		statusText, _ := jsonparser.GetString(result, "StatusText") //结果状态

		if statusText == aliyun.STATUS_SUCCESS {
			//获取智能分段结果
			aliyun.AliyunAudioResultWordHandle(result , func(vresult *aliyun.AliyunAudioRecognitionResult) {
				channelId := vresult.ChannelId

				_ , isPresent  := IntelligentBlockResult[channelId]
				if isPresent {
					//追加
					IntelligentBlockResult[channelId] = append(IntelligentBlockResult[channelId] , vresult)
				} else {
					//初始
					IntelligentBlockResult[channelId] = []*aliyun.AliyunAudioRecognitionResult{}
					IntelligentBlockResult[channelId] = append(IntelligentBlockResult[channelId] , vresult)
				}
			})

			//获取原始结果
			_, err := jsonparser.ArrayEach(result, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
				text , _ := jsonparser.GetString(value, "Text")
				channelId , _ := jsonparser.GetInt(value, "ChannelId")
				beginTime , _ := jsonparser.GetInt(value, "BeginTime")
				endTime , _ := jsonparser.GetInt(value, "EndTime")
				silenceDuration , _ := jsonparser.GetInt(value, "SilenceDuration")
				speechRate , _ := jsonparser.GetInt(value, "SpeechRate")
				emotionValue , _ := jsonparser.GetInt(value, "EmotionValue")

				vresult := &aliyun.AliyunAudioRecognitionResult {
					Text:text,
					ChannelId:channelId,
					BeginTime:beginTime,
					EndTime:endTime,
					SilenceDuration:silenceDuration,
					SpeechRate:speechRate,
					EmotionValue:emotionValue,
				}

				_ , isPresent  := AudioResult[channelId]
				if isPresent {
					//追加
					AudioResult[channelId] = append(AudioResult[channelId] , vresult)
				} else {
					//初始
					AudioResult[channelId] = []*aliyun.AliyunAudioRecognitionResult{}
					AudioResult[channelId] = append(AudioResult[channelId] , vresult)
				}
			} , "Result", "Sentences")

			if err != nil {
				panic(err)
			}
		}
	})

	if (resultError != nil) {
		panic(resultError)
	}

	return
}


//阿里云识别字幕块翻译处理
func AliyunAudioResultTranslate(app *VideoSrt , video string , AudioResult map[int64][] *aliyun.AliyunAudioRecognitionResult , IntelligentBlockResult map[int64][] *aliyun.AliyunAudioRecognitionResult)  {
	//输出日志
	if app.TranslateCfg.TranslateSwitch {
		app.Log("字幕翻译处理中 ..." , video)

		//百度翻译标准版
		if app.TranslateCfg.Supplier == TRANSLATE_SUPPLIER_BAIDU && app.TranslateCfg.BaiduTranslate.AuthenType == translate.ACCOUNT_COMMON_AUTHEN {
			app.Log("你使用的是 “百度翻译标准版” 账号，翻译速度较慢，请耐心等待 ..." , video)
		}
	} else {
		return
	}

	if app.OutputType.SRT || app.OutputType.LRC {
		for _,result := range IntelligentBlockResult {
			for _ , data := range result {
				transResult,e := app.RunTranslate(data.Text)
				if e != nil {
					panic(e) //终止翻译
				}
				data.TranslateText = strings.TrimSpace(transResult.TransResultDst) //译文
			}
		}
	}

	if app.OutputType.TXT {
		for _,result := range AudioResult {
			for _ , data := range result {
				transResult,e := app.RunTranslate(data.Text)
				if e != nil {
					panic(e) //终止翻译
				}
				data.TranslateText = strings.TrimSpace(transResult.TransResultDst) //译文
			}
		}
	}
}

//阿里云录音识别结果集生成字幕文件
func AliyunAudioResultMakeSubtitleFile(app *VideoSrt , video string , outputType int , AudioResult map[int64][] *aliyun.AliyunAudioRecognitionResult , IntelligentBlockResult map[int64][] *aliyun.AliyunAudioRecognitionResult)  {
	var subfileDir string
	if app.SrtDir == "" {
		subfileDir = path.Dir(video)
	} else {
		subfileDir = app.SrtDir
	}

	subfile := tool.GetFileBaseName(video)

	//注入合适的数据块
	thisAudioResult := make(map[int64][] *aliyun.AliyunAudioRecognitionResult)
	if outputType == OUTPUT_SRT || outputType == OUTPUT_LRC {
		thisAudioResult = IntelligentBlockResult
	} else if outputType == OUTPUT_STRING {
		thisAudioResult = AudioResult
	}

	oneSoundChannel := false //是否输出单条音轨
	//根据音轨，输出文件
	for channel,result := range thisAudioResult {
		soundChannel := channel + 1
		if app.SoundTrack != 3 && app.SoundTrack != 0 {
			oneSoundChannel = true
			if app.SoundTrack != int(soundChannel) {
				//跳过此音轨
				continue
			}
		}

		var thisfile string
		//文件名称
		if oneSoundChannel {
			thisfile = subfileDir + "/" + subfile
		} else {
			thisfile = subfileDir + "/" + subfile + "_channel_" +  strconv.FormatInt(soundChannel , 10)
		}
		//输出文件类型
		if outputType == OUTPUT_SRT {
			thisfile += ".srt"
		} else if outputType == OUTPUT_STRING {
			thisfile += ".txt"
		} else if outputType == OUTPUT_LRC {
			thisfile += ".lrc"
		}

		//创建文件
		file, e := os.Create(thisfile)
		if e != nil {
			panic(e)
		}
		defer file.Close()

		//文件编码分支
		if app.OutputEncode == OUTPUT_ENCODE_UTF8_BOM {
			if _, e := file.Write([]byte{0xEF, 0xBB, 0xBF});e != nil {
				panic(e)
			}
		}
		//歌词头
		if outputType == OUTPUT_LRC {
			_,_ = file.WriteString("[ar:]\r\n[ti:]\r\n[al:]\r\n[by:]\r\n")
		}

		//主字幕
		bilingualAsc := app.TranslateCfg.OutputMainSubtitleInputLanguage
		index := 0

		//普通文本容器
		var txtOutputContent bytes.Buffer
		var txtTransalteOutputContent bytes.Buffer

		for _ , data := range result {
			var linestr string

			//字幕、歌词文件处理
			if outputType == OUTPUT_SRT || outputType == OUTPUT_LRC {
				//拼接
				if outputType == OUTPUT_SRT {
					if app.TranslateCfg.TranslateSwitch {
						if app.TranslateCfg.BilingualSubtitleSwitch {
							linestr = MakeSubtitleText(index , data.BeginTime , data.EndTime , data.Text , data.TranslateText , true , bilingualAsc)
						} else {
							linestr = MakeSubtitleText(index , data.BeginTime , data.EndTime , data.TranslateText , "" , false , true)
						}
					} else {
						linestr = MakeSubtitleText(index , data.BeginTime , data.EndTime , data.Text , "" , false , true)
					}
				} else if outputType == OUTPUT_LRC {
					if app.TranslateCfg.TranslateSwitch {
						linestr = MakeMusicLrcText(index , data.BeginTime , data.EndTime , data.TranslateText)
					} else {
						linestr = MakeMusicLrcText(index , data.BeginTime , data.EndTime , data.Text)
					}
				}

				//写入行
				if _, e = file.WriteString(linestr);e != nil {
					panic(e)
				}
			} else if outputType == OUTPUT_STRING {
				//普通文本处理
				txtOutputContent.WriteString(data.Text)
				txtOutputContent.WriteString("\r\n")

				if app.TranslateCfg.TranslateSwitch {
					txtTransalteOutputContent.WriteString(data.TranslateText)
					txtTransalteOutputContent.WriteString("\r\n")
				}
			}

			index++
		}

		//写入文本文件
		if outputType == OUTPUT_STRING {
			txtOutputContent.WriteString("\r\n\r\n\r\n\r\n\r\n")
			if _, e = file.WriteString(txtOutputContent.String());e != nil {
				panic(e)
			}
			if _, e = file.WriteString(txtTransalteOutputContent.String());e != nil {
				panic(e)
			}
		}
	}
}