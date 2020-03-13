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

	IntelligentBlock bool //智能分段处理
	TempDir string //临时文件目录
	AppDir string //应用根目录
	SrtDir string //字幕文件输出目录
	OutputType int //输出文件类型
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

	app.IntelligentBlock = true
	app.TempDir = "temp/audio"
	app.AppDir = appDir
	app.OutputType = OUTPUT_SRT
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
}


//加载翻译设置
func (app *VideoSrt) InitTranslateConfig (translateSettings *VideoSrtTranslateStruct) {
	app.TranslateCfg = translateSettings
}


func (app *VideoSrt) SetSrtDir(dir string)  {
	app.SrtDir = dir
}

func (app *VideoSrt) SetOutputType(output int)  {
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

	//智能分段校验
	if app.OutputType == OUTPUT_STRING {
		app.IntelligentBlock = false //非输出字幕文件 关闭智能分段
	} else {
		app.IntelligentBlock = true
	}

	if video == "" {
		panic("enter a video file waiting to be processed .")
	}

	//校验媒体文件
	if tool.VaildVideo(video) != true {
		panic("the input video file does not exist .")
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
	filelink := UploadAudioToClound(app.AliyunOss , tmpAudio)
	//获取完整链接
	filelink = app.AliyunOss.GetObjectFileUrl(filelink)

	app.Log("上传文件成功 , 识别中 ..." , video)

	//阿里云录音文件识别
	AudioResult := AliyunAudioRecognition(app.AliyunClound, filelink , app.IntelligentBlock)

	app.Log("文件识别成功 , 字幕处理中 ..." , video)

	//输出字幕文件
	AliyunAudioResultMakeSubtitleFile(app , video , AudioResult)

	app.Log("处理完成" , video)

	//删除临时文件
	if tmpAudio != "" {
		if remove := os.Remove(tmpAudio); remove != nil {
			app.Log("错误：删除临时文件失败，请手动删除" , video)
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
			//休眠900毫秒
			time.Sleep(time.Millisecond * 900)
		} else {
			//休眠100毫秒
			time.Sleep(time.Millisecond * 150)
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


//阿里云录音文件识别
func AliyunAudioRecognition(engine aliyun.AliyunClound , filelink string , intelligent_block bool) (AudioResult map[int64][] *aliyun.AliyunAudioRecognitionResult) {
	//创建识别请求
	taskid, client, e := engine.NewAudioFile(filelink)
	if e != nil {
		panic(e)
	}

	AudioResult = make(map[int64][] *aliyun.AliyunAudioRecognitionResult)

	//遍历获取识别结果
	resultError := engine.GetAudioFileResult(taskid , client , func(result []byte) {
		//mylog.WriteLog( string( result ) )

		//结果处理
		statusText, _ := jsonparser.GetString(result, "StatusText") //结果状态
		if statusText == aliyun.STATUS_SUCCESS {

			//智能分段
			if intelligent_block {
				aliyun.AliyunAudioResultWordHandle(result , func(vresult *aliyun.AliyunAudioRecognitionResult) {
					channelId := vresult.ChannelId

					_ , isPresent  := AudioResult[channelId]
					if isPresent {
						//追加
						AudioResult[channelId] = append(AudioResult[channelId] , vresult)
					} else {
						//初始
						AudioResult[channelId] = []*aliyun.AliyunAudioRecognitionResult{}
						AudioResult[channelId] = append(AudioResult[channelId] , vresult)
					}
				})
				return
			}

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


//阿里云录音识别结果集生成字幕文件
func AliyunAudioResultMakeSubtitleFile(app *VideoSrt , video string , AudioResult map[int64][] *aliyun.AliyunAudioRecognitionResult)  {
	var subfileDir string
	if app.SrtDir == "" {
		subfileDir = path.Dir(video)
	} else {
		subfileDir = app.SrtDir
	}

	subfile := tool.GetFileBaseName(video)

	//输出日志
	if app.TranslateCfg.TranslateSwitch {
		app.Log("字幕翻译处理中 ..." , video)

		//百度翻译标准版
		if app.TranslateCfg.Supplier == TRANSLATE_SUPPLIER_BAIDU && app.TranslateCfg.BaiduTranslate.AuthenType == translate.ACCOUNT_COMMON_AUTHEN {
			app.Log("你使用的是 “百度翻译标准版” 账号，翻译速度较慢，请耐心等待 ..." , video)
		}
	}

	oneSoundChannel := false //输出单条音轨
	//根据音轨，输出文件
	for channel,result := range AudioResult {
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
		if app.OutputType == OUTPUT_SRT {
			thisfile += ".srt"
		} else if app.OutputType == OUTPUT_STRING {
			thisfile += ".txt"
		} else if app.OutputType == OUTPUT_LRC {
			thisfile += ".lrc"
		}

		file, e := os.Create(thisfile)
		if e != nil {
			panic(e)
		}

		//文件编码分支
		if app.OutputEncode == OUTPUT_ENCODE_UTF8_BOM {
			if _, e := file.Write([]byte{0xEF, 0xBB, 0xBF});e != nil {
				panic(e)
			}
		}

		//歌词头
		if app.OutputType == OUTPUT_LRC {
			_,_ = file.WriteString("[ar:]\r\n[ti:]\r\n[al:]\r\n[by:]\r\n")
		}

		index := 0
		for _ , data := range result {
			var linestr string
			var datastr string

			//开启翻译
			if app.TranslateCfg.TranslateSwitch {
				//双语字幕
				if app.TranslateCfg.BilingualSubtitleSwitch && app.OutputType == OUTPUT_SRT {
					datastr = data.Text
				} else {
					transResult,e := app.RunTranslate(data.Text)
					if e != nil {
						panic(e) //终止翻译
					}
					datastr = transResult.TransResultDst //译文
				}
			} else {
				datastr = data.Text
			}

			datastr = strings.TrimSpace(datastr)
			//拼接文本
			if app.OutputType == OUTPUT_SRT {
				linestr = MakeSubtitleText(app , index , data.BeginTime , data.EndTime , datastr)
			} else if app.OutputType == OUTPUT_STRING {
				linestr = MakeText(index , data.BeginTime , data.EndTime , datastr)
			} else if app.OutputType == OUTPUT_LRC {
				linestr = MakeMusicLrcText(index , data.BeginTime , data.EndTime , datastr)
			}

			if _, e = file.WriteString(linestr);e != nil {
				panic(e)
			}
			index++
		}

		//close
		_ = file.Close()
	}
}


//拼接字幕字符串
func MakeSubtitleText(app *VideoSrt, index int , startTime int64 , endTime int64 , text string) string {
	var content bytes.Buffer
	content.WriteString(strconv.Itoa(index))
	content.WriteString("\r\n")
	content.WriteString(tool.SubtitleTimeMillisecond(startTime , true))
	content.WriteString(" --> ")
	content.WriteString(tool.SubtitleTimeMillisecond(endTime , true))
	content.WriteString("\r\n")

	//输出双语字幕
	if app.TranslateCfg.TranslateSwitch && app.TranslateCfg.BilingualSubtitleSwitch {
		transResult,e := app.RunTranslate(text)
		if e != nil {
			panic(e) //终止翻译
		}

		content.WriteString(transResult.TransResultSrc)
		content.WriteString("\r\n")
		content.WriteString(transResult.TransResultDst)
	} else {
		content.WriteString(text)
	}

	content.WriteString("\r\n")
	content.WriteString("\r\n")
	return content.String()
}


//拼接文本格式
func MakeText(index int , startTime int64 , endTime int64 , text string) string {
	var content bytes.Buffer
	content.WriteString(text)
	content.WriteString("\r\n")
	content.WriteString("\r\n")
	return content.String()
}


//拼接歌词文本
func MakeMusicLrcText(index int , startTime int64 , endTime int64 , text string) string {
	var content bytes.Buffer
	content.WriteString("[")
	content.WriteString(tool.MusicLrcTextMillisecond(startTime))
	content.WriteString("]")
	content.WriteString(text)
	content.WriteString("\r\n")
	return content.String()
}