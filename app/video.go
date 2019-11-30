package app

import (
	"bytes"
	"github.com/buger/jsonparser"
	"os"
	"path"
	"strconv"
	"videosrt/app/aliyun"
	"videosrt/app/ffmpeg"
	"videosrt/app/tool"
)

//主应用
type VideoSrt struct {
	Ffmpeg ffmpeg.Ffmpeg
	AliyunOss aliyun.AliyunOss //oss
	AliyunClound aliyun.AliyunClound //语音识别引擎

	IntelligentBlock bool //智能分段处理
	TempDir string //临时文件目录
	AppDir string //应用根目录
	SrtDir string //字幕文件输出目录

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
	return app
}


//应用加载配置
func (app *VideoSrt) InitConfig(oss *AliyunOssCache , engine *AliyunEngineCache) {
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

func (app *VideoSrt) SetSrtDir(dir string)  {
	app.SrtDir = dir
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
func  (app *VideoSrt) ClearTempDir()  {
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

	if video == "" {
		panic("enter a video file waiting to be processed .")
	}

	//校验视频
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

	//分离视频音频
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
	AliyunAudioResultMakeSubtitleFile(app.SrtDir , video , AudioResult)

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
	engine.GetAudioFileResult(taskid , client , func(result []byte) {
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

	return
}


//阿里云录音识别结果集生成字幕文件
func AliyunAudioResultMakeSubtitleFile(outputDir string , video string , AudioResult map[int64][] *aliyun.AliyunAudioRecognitionResult)  {
	var subfileDir string
	if outputDir == "" {
		subfileDir = path.Dir(video)
	} else {
		subfileDir = outputDir
	}

	subfile := tool.GetFileBaseName(video)

	//输出字幕文件
	for channel,result := range AudioResult {
		thisfile := subfileDir + "/" + subfile + "_channel_" +  strconv.FormatInt(channel , 10) + ".srt"
		//println(thisfile)

		file, e := os.Create(thisfile)
		if e != nil {
			panic(e)
		}

		defer file.Close() //defer
		index := 0
		for _ , data := range result {
			linestr := MakeSubtitleText(index , data.BeginTime , data.EndTime , data.Text)
			file.WriteString(linestr)
			index++
		}
	}
}


//拼接字幕字符串
func MakeSubtitleText(index int , startTime int64 , endTime int64 , text string) string {
	var content bytes.Buffer
	content.WriteString(strconv.Itoa(index))
	content.WriteString("\n")
	content.WriteString(tool.SubtitleTimeMillisecond(startTime))
	content.WriteString(" --> ")
	content.WriteString(tool.SubtitleTimeMillisecond(endTime))
	content.WriteString("\n")
	content.WriteString(text)
	content.WriteString("\n")
	content.WriteString("\n")
	return content.String()
}