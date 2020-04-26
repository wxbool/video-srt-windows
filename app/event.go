package app

import (
	"errors"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"videosrt/app/aliyun"
	"videosrt/app/tool"
	"videosrt/app/translate"
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
	var seting *AppSetings
	var dlg *walk.Dialog
	var db *walk.DataBinder
	var acceptPB, cancelPB *walk.PushButton

	seting = Setings.GetCacheAppSetingsData() //查询缓存数据
	//fmt.Println( setings )
	if seting.MaxConcurrency == 0 {
		seting.MaxConcurrency = 2 //默认并发数
	}

	Dialog{
		AssignTo:      &dlg,
		Title:         "软件设置",
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		DataBinder: DataBinder{
			AssignTo:       &db,
			Name:           "setings",
			DataSource:     seting,
			ErrorPresenter: ToolTipErrorPresenter{},
		},
		MinSize: Size{450, 220},
		Layout:  VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{
						Text: "任务处理并发数：",
					},
					NumberEdit{
						Value:    Bind("MaxConcurrency", Range{1, 20}),
						Decimals: 0,
					},

					Label{
						Text: "文件输出目录：",
					},
					LineEdit{
						Text: Bind("SrtFileDir"),
					},
					Label{
						ColumnSpan: 2,
						Text: "说明：\r\n“文件输出目录” 若留空，则默认与媒体文件输出到同一目录下",
						TextColor:walk.RGB(190 , 190 , 190),
						MinSize:Size{Height:40},
					},


					Label{
						Text: "关闭OSS临时文件清理:",
					},
					CheckBox{
						Checked: Bind("CloseAutoDeleteOssTempFile"),
					},
					Label{
						Text: "关闭软件新版本提醒:",
					},
					CheckBox{
						Checked: Bind("CloseNewVersionMessage"),
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
							if seting.SrtFileDir != "" {
								tmpDir := tool.WinDir(seting.SrtFileDir)
								if !tool.DirExists(tmpDir) {
									mw.NewErrormationTips("错误" , "目录无效/不存在：" + seting.SrtFileDir)
									return;
								}
								seting.SrtFileDir = tmpDir
							}

							//设置缓存
							Setings.SetCacheAppSetingsData(seting)

							//设置回调
							confirmCall(seting)

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

	engine.Region = aliyun.ALIYUN_CLOUND_REGION_CHA //默认值

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

					Label{
						Text: "服务区域：",
					},
					ComboBox{
						Value: Bind("Region", SelRequired{}),
						BindingMember: "Id",
						DisplayMember: "Name",
						Model: GetAliyunEngineRegionOptionSelects(),
					},

					Label{
						ColumnSpan: 2,
						Text: "说明：\r\n“语音识别” 目前使用的是阿里云语音服务商，请填写相关的引擎配置。",
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
						Text:     "确定新增",
						OnClicked: func() {
							if err := db.Submit(); err != nil {
								log.Print(err)
								return
							}
							//参数验证
							if (engine.Alias == "") {
								mw.NewErrormationTips("提示" , "请填写 名称/别名")
								return
							}
							if (engine.AppKey == "") {
								mw.NewErrormationTips("提示" , "请填写 AppKey")
								return
							}
							if (engine.AccessKeyId == "") {
								mw.NewErrormationTips("提示" , "请填写 AccessKeyId")
								return
							}
							if (engine.AccessKeySecret == "") {
								mw.NewErrormationTips("提示" , "请填写 AccessKeySecret")
								return
							}

							//获取缓存数据
							localData := Engine.GetCacheAliyunEngineListData()
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
							Engine.SetCacheAliyunEngineListData(localData)

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




// 运行 新建[百度]翻译引擎 Dialog
func(mw *MyMainWindow) RunBaiduTranslateEngineSetingDialog(owner walk.Form , confirmCall func())  {
	var engine *TranslateEngineStruct
	var dlg *walk.Dialog
	var db *walk.DataBinder
	var acceptPB, cancelPB *walk.PushButton

	engine = new(TranslateEngineStruct)

	//默认值
	engine.Supplier = TRANSLATE_SUPPLIER_BAIDU
	engine.BaiduEngine.AuthenType = translate.ACCOUNT_COMMON_AUTHEN

	Dialog{
		AssignTo:      &dlg,
		Title:         "新建翻译引擎（百度翻译）",
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
						Text: "AppId：",
					},
					LineEdit{
						Text: Bind("BaiduEngine.AppId"),
					},

					Label{
						Text: "AppSecret：",
					},
					LineEdit{
						Text: Bind("BaiduEngine.AppSecret"),
					},

					Label{
						Text: "账号认证类型：",
					},
					ComboBox{
						Value: Bind("BaiduEngine.AuthenType", SelRequired{}),
						BindingMember: "Id",
						DisplayMember: "Name",
						Model: GetBaiduTranslateAuthenTypeOptionsSelects(),
					},

					TextLabel{
						ColumnSpan: 2,
						Row: 5,
						Text: "\r\n说明：\r\n请填写在 “百度翻译开放平台” 的申请密钥，注意 “标准版” 和 “高级版” 的选择",
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
						Text:     "确定新增",
						OnClicked: func() {
							if err := db.Submit(); err != nil {
								log.Print(err)
								return
							}
							//参数验证
							if (engine.Alias == "") {
								mw.NewErrormationTips("提示" , "请填写 名称/别名")
								return
							}
							if (engine.BaiduEngine.AppId == "") {
								mw.NewErrormationTips("提示" , "请填写 AppId")
								return
							}
							if (engine.BaiduEngine.AppSecret == "") {
								mw.NewErrormationTips("提示" , "请填写 AppSecret")
								return
							}

							//获取缓存数据
							localData := Translate.GetCacheTranslateEngineListData()
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
							Translate.SetCacheTranslateEngineListData(localData)

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

// 运行 新建[腾讯云]翻译引擎 Dialog
func(mw *MyMainWindow) RunTengxunyunTranslateEngineSetingDialog(owner walk.Form , confirmCall func())  {
	var engine *TranslateEngineStruct
	var dlg *walk.Dialog
	var db *walk.DataBinder
	var acceptPB, cancelPB *walk.PushButton

	engine = new(TranslateEngineStruct)

	//默认值
	engine.Supplier = TRANSLATE_SUPPLIER_TENGXUNYUN

	Dialog{
		AssignTo:      &dlg,
		Title:         "新建翻译引擎（腾讯云）",
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
						Text: "SecretId：",
					},
					LineEdit{
						Text: Bind("TengxunyunEngine.SecretId"),
					},

					Label{
						Text: "SecretKey：",
					},
					LineEdit{
						Text: Bind("TengxunyunEngine.SecretKey"),
					},

					TextLabel{
						ColumnSpan: 2,
						Row: 5,
						Text: "\r\n说明：\r\n请填写 “腾讯云” 创建的子用户密钥",
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
						Text:     "确定新增",
						OnClicked: func() {
							if err := db.Submit(); err != nil {
								log.Print(err)
								return
							}
							//参数验证
							if (engine.Alias == "") {
								mw.NewErrormationTips("提示" , "请填写 名称/别名")
								return
							}
							if (engine.TengxunyunEngine.SecretId == "") {
								mw.NewErrormationTips("提示" , "请填写 SecretId")
								return
							}
							if (engine.TengxunyunEngine.SecretKey == "") {
								mw.NewErrormationTips("提示" , "请填写 SecretKey")
								return
							}

							//获取缓存数据
							localData := Translate.GetCacheTranslateEngineListData()
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
							Translate.SetCacheTranslateEngineListData(localData)

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

	oss = Oss.GetCacheAliyunOssData() //查询缓存数据

	Dialog{
		AssignTo:      &dlg,
		Title:         "OSS对象存储设置",
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

					Label{
						ColumnSpan: 2,
						Text: "说明：\r\n“OSS对象存储”目前使用的是阿里云服务，请填写相关的服务配置。",
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

							//去空格
							oss.Endpoint = strings.TrimSpace(oss.Endpoint)
							oss.AccessKeyId = strings.TrimSpace(oss.AccessKeyId)
							oss.AccessKeySecret = strings.TrimSpace(oss.AccessKeySecret)
							oss.BucketName = strings.TrimSpace(oss.BucketName)
							oss.BucketDomain = strings.TrimSpace(oss.BucketDomain)

							//设置缓存
							Oss.SetCacheAliyunOssData(oss)

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
	tool.OpenUrl("https://github.com/wxbool/video-srt-windows")
}

//打开 Gitee
func (mw *MyMainWindow) OpenAboutGitee() {
	tool.OpenUrl("https://gitee.com/641453620/video-srt-windows")
}


//支持的文件后缀校验
func VaildateHandleFiles(files [] string , mediaExt bool , srtExt bool) ([]string , error) {
	result := []string{}
	allowExts := []string{}
	mediaExts := []string{
		".mp4",".mpeg",".mkv",".wmv",".avi",".m4v",".mov",".flv",".rmvb",".3gp",".f4v",
		".mp3",".wav",".aac",".wma",".flac",".m4a",
	}
	srtExts := []string{".srt"}

	if mediaExt {
		allowExts = append(allowExts , mediaExts...)
	}
	if srtExt {
		allowExts = append(allowExts , srtExts...)
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
			if !tool.InSliceString(strings.ToLower(ext) , allowExts) {
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