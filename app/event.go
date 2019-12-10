package app

import (
	"errors"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
	"os"
	"path"
	"path/filepath"
	"videosrt/app/tool"
)

type MyMainWindow struct {
	*walk.MainWindow
}

//创建一个提示消息
func (mw *MyMainWindow) NewInformationTips(title string , message string) {
	walk.MsgBox(mw, title , message , walk.MsgBoxIconInformation)
}

//创建一个错误消息
func (mw *MyMainWindow) NewErrormationTips(title string , message string) {
	walk.MsgBox(mw, title , message , walk.MsgBoxIconWarning)
}


// 运行 应用设置 Dialog
func(mw *MyMainWindow) RunAppSetingDialog(owner walk.Form , confirmCall func(*AppSetings))  {
	var setings *AppSetings
	var dlg *walk.Dialog
	var db *walk.DataBinder
	var acceptPB, cancelPB *walk.PushButton

	setings = GetCacheAppSetingsData() //查询缓存数据
	//fmt.Println( setings )
	if setings.MaxConcurrency == 0 {
		setings.MaxConcurrency = 2 //默认并发数
	}
	if setings.OutputType == 0 {
		setings.OutputType = 1 //默认输出文件类型
	}

	Dialog{
		AssignTo:      &dlg,
		Title:         "软件设置",
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		DataBinder: DataBinder{
			AssignTo:       &db,
			Name:           "setings",
			DataSource:     setings,
			ErrorPresenter: ToolTipErrorPresenter{},
		},
		MinSize: Size{450, 200},
		Layout:  VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					//输出文件类型
					Label{
						Text: "输出文件类型:",
					},
					ComboBox{
						Value: Bind("OutputType", SelRequired{}),
						BindingMember: "Id",
						DisplayMember: "Name",
						Model: GetOutputOptionsSelects(),
					},


					Label{
						Text: "任务处理并发数：",
					},
					NumberEdit{
						Value:    Bind("MaxConcurrency", Range{1, 20}),
						Decimals: 0,
					},

					Label{
						Text: "字幕文件输出目录：",
					},
					LineEdit{
						Text: Bind("SrtFileDir"),
					},

					Label{
						ColumnSpan: 2,
						Text: "说明：\r\n“字幕文件输出目录” 若留空，则默认与媒体文件输出到同一目录下",
						TextColor:walk.RGB(190 , 190 , 190),
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						AssignTo: &acceptPB,
						Text:     "保存",
						OnClicked: func() {
							if err := db.Submit(); err != nil {
								log.Fatal(err)
								return
							}
							//目录校验
							if setings.SrtFileDir != "" {
								tmpDir := tool.WinDir(setings.SrtFileDir)
								if !tool.DirExists(tmpDir) {
									mw.NewErrormationTips("错误" , "目录无效/不存在：" + setings.SrtFileDir)
									return;
								}
								setings.SrtFileDir = tmpDir
							}

							//设置缓存
							SetCacheAppSetingsData(setings)

							//设置回调
							confirmCall(setings)

							dlg.Accept()
						},
					},
					PushButton{
						AssignTo:  &cancelPB,
						Text:      "取消",
						OnClicked: func() { dlg.Cancel() },
					},
				},
			},
		},
	}.Run( owner )
}

// 运行 新建语音引擎 Dialog
func(mw *MyMainWindow) RunSpeechEngineSetingDialog(owner walk.Form , confirmCall func())  {
	var engine *AliyunEngineCache
	var dlg *walk.Dialog
	var db *walk.DataBinder
	var acceptPB, cancelPB *walk.PushButton

	engine = new(AliyunEngineCache)

	Dialog{
		AssignTo:      &dlg,
		Title:         "新建语音引擎",
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		DataBinder: DataBinder{
			AssignTo:       &db,
			Name:           "engine",
			DataSource:     engine,
			ErrorPresenter: ToolTipErrorPresenter{},
		},
		MinSize: Size{500, 300},
		Layout:  VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{
						Text: "名称/别名：",
					},
					LineEdit{
						Text: Bind("Alias"),
					},

					Label{
						Text: "AppKey：",
					},
					LineEdit{
						Text: Bind("AppKey"),
					},

					Label{
						Text: "AccessKeyId：",
					},
					LineEdit{
						Text: Bind("AccessKeyId"),
					},

					Label{
						Text: "AccessKeySecret：",
					},
					LineEdit{
						Text: Bind("AccessKeySecret"),
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						AssignTo: &acceptPB,
						Text:     "确定新增",
						OnClicked: func() {
							if err := db.Submit(); err != nil {
								log.Print(err)
								return
							}
							//参数验证
							if (engine.Alias == "") {
								mw.NewInformationTips("提示" , "请填写 名称/别名")
								return
							}
							if (engine.AppKey == "") {
								mw.NewInformationTips("提示" , "请填写 AppKey")
								return
							}
							if (engine.AccessKeyId == "") {
								mw.NewInformationTips("提示" , "请填写 AccessKeyId")
								return
							}
							if (engine.AccessKeySecret == "") {
								mw.NewInformationTips("提示" , "请填写 AccessKeySecret")
								return
							}

							//获取缓存数据
							localData := GetCacheAliyunEngineListData()
							//生成id
							lens := len(localData.Engine)
							if lens == 0 {
								engine.Id = 1
							} else {
								engine.Id = localData.Engine[ lens - 1 ].Id + 1
							}
							//追加数据
							localData.Engine = append(localData.Engine , engine)
							//缓存数据
							SetCacheAliyunEngineListData(localData)

							//调用回调
							confirmCall()

							dlg.Accept()
						},
					},
					PushButton{
						AssignTo:  &cancelPB,
						Text:      "取消",
						OnClicked: func() { dlg.Cancel() },
					},
				},
			},
		},
	}.Run( owner )
}

//运行 Oss对象存储配置 Dialog
func (mw *MyMainWindow) RunObjectStorageSetingDialog(owner walk.Form) {
	var oss *AliyunOssCache
	var dlg *walk.Dialog
	var db *walk.DataBinder
	var acceptPB, cancelPB *walk.PushButton

	oss = GetCacheAliyunOssData() //查询缓存数据

	Dialog{
		AssignTo:      &dlg,
		Title:         "Oss对象存储设置",
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		DataBinder: DataBinder{
			AssignTo:       &db,
			Name:           "oss",
			DataSource:     oss,
			ErrorPresenter: ToolTipErrorPresenter{},
		},
		MinSize: Size{500, 300},
		Layout:  VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{
						Text: "Endpoint：",
					},
					LineEdit{
						Text: Bind("Endpoint"),
					},

					Label{
						Text: "AccessKeyId：",
					},
					LineEdit{
						Text: Bind("AccessKeyId"),
					},

					Label{
						Text: "AccessKeySecret：",
					},
					LineEdit{
						Text: Bind("AccessKeySecret"),
					},

					Label{
						Text: "BucketName：",
					},
					LineEdit{
						Text: Bind("BucketName"),
					},

					Label{
						Text: "BucketDomain：",
					},
					LineEdit{
						Text: Bind("BucketDomain"),
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						AssignTo: &acceptPB,
						Text:     "保存",
						OnClicked: func() {
							if err := db.Submit(); err != nil {
								log.Fatal(err)
								return
							}
							//参数验证
							if (oss.Endpoint == "") {
								mw.NewInformationTips("提示" , "请填写 Endpoint")
								return
							}
							if (oss.AccessKeyId == "") {
								mw.NewInformationTips("提示" , "请填写 AccessKeyId")
								return
							}
							if (oss.AccessKeySecret == "") {
								mw.NewInformationTips("提示" , "请填写 AccessKeySecret")
								return
							}
							if (oss.BucketName == "") {
								mw.NewInformationTips("提示" , "请填写 BucketName")
								return
							}
							if (oss.BucketDomain == "") {
								mw.NewInformationTips("提示" , "请填写 BucketDomain")
								return
							}

							//设置缓存
							SetCacheAliyunOssData(oss)

							dlg.Accept()
						},
					},
					PushButton{
						AssignTo:  &cancelPB,
						Text:      "取消",
						OnClicked: func() { dlg.Cancel() },
					},
				},
			},
		},
	}.Run( owner )
}

//打开 Github
func (mw *MyMainWindow) OpenAboutGithub() {
	tool.OpenUrl("https://github.com/wxbool/video-srt")
}

//打开 Gitee
func (mw *MyMainWindow) OpenAboutGitee() {
	tool.OpenUrl("https://gitee.com/641453620/video-srt")
}


//校验待处理文件
func VaildateHandleFiles(files [] string) ([]string , error) {
	result := []string{}
	allowExts := []string{
		".mp4",".mpeg",".mkv",".wmv",".avi",".m4v",".mov",".flv",".rmvb",".3gp",".f4v",
		".mp3",".wav",".aac",".wma",
	}
	for _,f := range files {
		f = tool.WinDir(f)
		if thisFile, err := os.Stat(f); err != nil {
			if os.IsNotExist(err) {
				return result , errors.New("文件不存在：" + f)
			}
			return result , errors.New("文件校验不通过：" + f)
		} else {
			if thisFile.IsDir() {
				return result , errors.New("不允许操作文件夹：" + f)
			}
			//校验视频格式后缀
			ext := path.Ext(f)
			if !tool.InSliceString(ext , allowExts) {
				return result , errors.New("文件后缀不允许：" + f)
			}
			//允许加入
			result = append(result , f)
		}
	}
	//数量限制
	if len(result) > 300 {
		return result , errors.New("文件数量不允许超过 300 个")
	}
	return result , nil
}


//获取应用根目录
func GetAppRootDir() string {
	if rootDir , err := filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
		return ""
	} else {
		return tool.WinDir(rootDir)
	}
}