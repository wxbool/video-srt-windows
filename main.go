package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
	"runtime"
	"strings"
	"time"
	. "videosrt/app"
	"videosrt/app/ffmpeg"
	"videosrt/app/tool"
)

//应用版本号
const APP_VERSION = "0.1.3"

var AppRootDir string
var mw *MyMainWindow


func init()  {
	//设置可同时执行的最大CPU数
	runtime.GOMAXPROCS(runtime.NumCPU())
	//MY Window
	mw = new(MyMainWindow)

	AppRootDir = GetAppRootDir()
	if AppRootDir == "" {
		panic("应用根目录获取失败")
	}

	//校验ffmpeg环境
	if e := ffmpeg.VailFfmpegLibrary(); e != nil {
		//尝试自动引入 ffmpeg 环境
		ffmpeg.VailTempFfmpegLibrary(AppRootDir)
	}
}


func main() {
	var taskFiles = new(TaskHandleFile)

	var logText *walk.TextEdit
	var operateDb *walk.DataBinder
	var operateFrom = new(OperateFrom)

	var startBtn *walk.PushButton
	var engineOtionsBox *walk.ComboBox
	var dropFilesEdit *walk.TextEdit

	var appSetings = GetCacheAppSetingsData()
	if appSetings.CurrentEngineId != 0 {
		operateFrom.EngineId = appSetings.CurrentEngineId
	}

	var videosrt = NewApp(AppRootDir)
	var tasklog = NewTasklog()

	//注册日志事件
	videosrt.SetLogHandler(func(s string, video string) {
		baseName := tool.GetFileBaseName(video)
		//fmt.Println("日志：" , s , baseName)
		//fmt.Println("\r\n")
		//fmt.Println("\r\n")
		//fmt.Println("\r\n")
		strs := strings.Join([]string{"【" , baseName , "】" , s} , "")
		//追加日志
		tasklog.AppendLogText(strs)
	})
	//字幕输出目录
	videosrt.SetSrtDir(appSetings.SrtFileDir)

	//注册多任务
	var multitask = NewVideoMultitask(appSetings.MaxConcurrency)

	if err := (MainWindow{
		AssignTo: &mw.MainWindow,
		Icon:"./data/img/index.png",
		Title:    "VideoSrt - 识别视频语音生成字幕文件的小工具" + " - " + APP_VERSION,
		ToolBar: ToolBar{
			ButtonStyle: ToolBarButtonImageBeforeText,
			Items: []MenuItem{
				Menu{
					Text:  "新建",
					Image: "./data/img/new.png",
					Items: []MenuItem{
						Action{
							Image:  "./data/img/voice.png",
							Text:   "语音引擎",
							OnTriggered: func() {
								mw.RunSpeechEngineSetingDialog(mw , func() {
									thisData := GetEngineOtionsSelects()
									if appSetings.CurrentEngineId == 0 {
										appSetings.CurrentEngineId = thisData[0].Id

										//更新缓存
										SetCacheAppSetingsData(appSetings)
									}

									//重新加载选项
									_ = engineOtionsBox.SetModel(thisData)
									//重置index
									engIndex := GetCurrentIndex(thisData , appSetings.CurrentEngineId)
									if engIndex != -1 {
										_ = engineOtionsBox.SetCurrentIndex(engIndex)
									}
									operateFrom.EngineId = appSetings.CurrentEngineId
								})
							},
						},
					},
				},
				Menu{
					Text:  "设置",
					Image: "./data/img/setings.png",
					Items: []MenuItem{
						Action{
							Text:    "Oss对象存储设置",
							Image:   "./data/img/oss.png",
							OnTriggered: func() {
								mw.RunObjectStorageSetingDialog(mw)
							},
						},
						Action{
							Text:    "软件设置",
							Image:   "./data/img/app-setings.png",
							OnTriggered: func() {
								mw.RunAppSetingDialog(mw , func(setings *AppSetings) {
									//更新配置
									appSetings.MaxConcurrency = setings.MaxConcurrency
									appSetings.SrtFileDir = setings.SrtFileDir
									appSetings.OutputType = setings.OutputType

									videosrt.SetOutputType( setings.OutputType )
									videosrt.SetSrtDir( setings.SrtFileDir )
									multitask.SetMaxConcurrencyNumber( setings.MaxConcurrency )
								})
							},
						},
					},
				},
				Menu{
					Text:  "关于/支持",
					Image: "./data/img/about.png",
					Items: []MenuItem{
						Action{
							Text:        "github",
							Image:      "./data/img/github.png",
							OnTriggered: mw.OpenAboutGithub,
						},
						Action{
							Text:        "gitee",
							Image:      "./data/img/gitee.png",
							OnTriggered: mw.OpenAboutGitee,
						},
						Action{
							Text:        "最新版本",
							Image:      "./data/img/version.png",
							OnTriggered: func() {
								tool.OpenUrl("https://gitee.com/641453620/video-srt-windows/releases")
							},
						},
					},
				},
			},
		},
		Size: Size{800, 450},
		MinSize: Size{300, 280},
		Layout:  VBox{},
		Children: []Widget{
			HSplitter{
				Children: []Widget{
					Composite{
						DataBinder: DataBinder{
							AssignTo:    &operateDb,
							DataSource:   operateFrom,
						},
						Layout: Grid{Columns: 6},
						Children: []Widget{
							Label{
								Text: "语音引擎:",
							},
							ComboBox{
								AssignTo:&engineOtionsBox,
								Value: Bind("EngineId", SelRequired{}),
								BindingMember: "Id",
								DisplayMember: "Name",
								Model:  GetEngineOtionsSelects(),
								OnCurrentIndexChanged: func() {
									operateDb.Submit()

									//fmt.Println( "OnCurrentIndexChanged：" , operateFrom.EngineId ,appSetings.CurrentEngineId )
									if operateFrom.EngineId == 0 {
										return
									}

									appSetings.CurrentEngineId = operateFrom.EngineId
									//更新缓存
									SetCacheAppSetingsData(appSetings)
								},
							},


							PushButton{
								Text: "删除引擎",
								MaxSize:Size{80 , 55},
								OnClicked: func() {
									var thisEngineOptions = make([]*EngineSelects , 0)
									thisEngineOptions = GetEngineOtionsSelects()
									//删除校验
									if appSetings.CurrentEngineId == 0 {
										mw.NewErrormationTips("错误" , "请先选择要操作的语音引擎")
										return
									}
									if len(thisEngineOptions) <= 1 {
										mw.NewErrormationTips("错误" , "不允许删除最后一个语音引擎")
										return
									}

									//删除引擎
									if ok := RemoveCacheAliyunEngineData(appSetings.CurrentEngineId);ok == false {
										//删除失败
										mw.NewErrormationTips("错误" , "语音引擎删除失败")
										return
									}

									thisEngineOptions = GetEngineOtionsSelects()

									//fmt.Println( "thisEngineOptions" , thisEngineOptions[0] )
									//重新加载列表
									_ = engineOtionsBox.SetModel(thisEngineOptions)

									appSetings.CurrentEngineId = thisEngineOptions[0].Id
									operateFrom.EngineId = appSetings.CurrentEngineId
									//更新缓存
									SetCacheAppSetingsData(appSetings)
									//更新下标
									_ = engineOtionsBox.SetCurrentIndex(0)
								},
							},
						},
					},
				},
			},
			HSplitter{
				Children: []Widget{
					TextEdit{
						AssignTo: &dropFilesEdit,
						ReadOnly: true,
						Text:     "将需要生成字幕的媒体文件，拖入放到这里\r\n\r\n支持的视频格式：.mp4 , .mpeg , .wmv , .avi , .m4v , .mov , .flv , .rmvb , .3gp , .f4v\r\n支持的音频格式：.mp3 , .wav , .aac , .wma",
						TextColor:walk.RGB(136 , 136 , 136),
					},
					TextEdit{
						AssignTo: &logText,
						ReadOnly: true,
						Text:"",
						//HScroll:true,
						VScroll:true,
					},
				},
			},
			PushButton{
				AssignTo: &startBtn,
				Text: "生成文件",
				MinSize:Size{500 , 50},
				OnClicked: func() {

					//设置随机种子
					tool.SetRandomSeed()

					tlens := len(taskFiles.Files)
					if tlens == 0 {
						mw.NewErrormationTips("错误" , "请先拖入要处理的视频文件")
						return
					}
					ossData := GetCacheAliyunOssData()
					if ossData.Endpoint == "" {
						mw.NewErrormationTips("错误" , "请先设置Oss对象配置")
						return
					}
					engineIndex := GetCacheAppSetingsData()
					if engineIndex.CurrentEngineId == 0 {
						mw.NewErrormationTips("错误" , "请先新建/选择语音引擎")
						return
					}
					currentEngine , ok := GetEngineById(engineIndex.CurrentEngineId)
					if !ok {
						mw.NewErrormationTips("错误" , "你选择的语音引擎不存在")
						return
					}
					//加载配置
					videosrt.InitConfig(ossData , currentEngine)
					videosrt.SetSrtDir(appSetings.SrtFileDir)
					if appSetings.OutputType != 0 {
						videosrt.SetOutputType(appSetings.OutputType)
					}

					multitask.SetVideoSrt(videosrt)
					//设置队列
					multitask.SetQueueFile(taskFiles.Files)

					var finish = false

					startBtn.SetEnabled(false)
					startBtn.SetText("任务运行中，请勿关闭软件窗口...")
					//清除log
					tasklog.ClearLogText()
					tasklog.AppendLogText("任务开始... \r\n")

					//运行
					multitask.Run()

					//回调链式执行
					videosrt.SetFailHandler(func(video string) {
						//运行下一任务
						multitask.RunOver()

						//任务完成
						if ok := multitask.FinishTask(); ok && finish == false {
							//延迟结束
							go func() {
								time.Sleep(time.Second)
								finish = true
								startBtn.SetEnabled(true)
								startBtn.SetText("生成文件")

								logText.AppendText("\r\n\r\n任务完成！")

								//清空临时目录
								videosrt.ClearTempDir()
							}()
						}
					})
					videosrt.SetSuccessHandler(func(video string) {
						//运行下一任务
						multitask.RunOver()

						//任务完成
						if ok := multitask.FinishTask(); ok && finish == false {
							//延迟结束
							go func() {
								time.Sleep(time.Second)
								finish = true
								startBtn.SetEnabled(true)
								startBtn.SetText("生成字幕文件")

								logText.AppendText("\r\n\r\n任务完成！")

								//清空临时目录
								videosrt.ClearTempDir()
							}()
						}
					})

					//日志输出
					go func() {
						for finish == false {
							logText.SetText("")
							logText.AppendText(tasklog.GetString())
							time.Sleep(time.Millisecond * 150)
						}
					}()
				},
			},
		},
		OnDropFiles: func(files []string) {
			//检测文件列表
			result , err := VaildateHandleFiles(files)
			if err != nil {
				mw.NewErrormationTips("错误" , err.Error())
				return
			}

			taskFiles.Files = result
			dropFilesEdit.SetText(strings.Join(result, "\r\n"))
		},
	}.Create()); err != nil {
		log.Fatal(err)

		time.Sleep(1 * time.Second)
	}

	//校验依赖库
	if e := ffmpeg.VailFfmpegLibrary(); e != nil {
		mw.NewErrormationTips("错误" , "请先下载并安装 ffmpeg 软件，才可以正常使用软件哦")
		tool.OpenUrl("https://gitee.com/641453620/video-srt")
		return
	}

	mw.Run()
}
