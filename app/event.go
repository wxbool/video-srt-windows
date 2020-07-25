package app

import (
	"errors"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
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

							//去空格
							engine.AppKey = strings.TrimSpace(engine.AppKey)
							engine.AccessKeyId = strings.TrimSpace(engine.AccessKeyId)
							engine.AccessKeySecret = strings.TrimSpace(engine.AccessKeySecret)

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

							//去空格
							engine.BaiduEngine.AppId = strings.TrimSpace(engine.BaiduEngine.AppId)
							engine.BaiduEngine.AppSecret = strings.TrimSpace(engine.BaiduEngine.AppSecret)

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

							//去空格
							engine.TengxunyunEngine.SecretId = strings.TrimSpace(engine.TengxunyunEngine.SecretId)
							engine.TengxunyunEngine.SecretKey = strings.TrimSpace(engine.TengxunyunEngine.SecretKey)

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


//运行 语气词过滤设置 Dialog
func (mw *MyMainWindow) RunGlobalFilterSetingDialog (owner walk.Form , historyWords string , confirmCall func(words string)) {
	var dlg *walk.Dialog
	var db *walk.DataBinder
	var acceptPB, cancelPB *walk.PushButton

	var tmpData = new(AppFilterSetings)
	tmpData.GlobalFilter.Words = historyWords

	Dialog{
		AssignTo:      &dlg,
		Title:         "全局语气词过滤设置",
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		DataBinder: DataBinder{
			AssignTo:       &db,
			Name:           "filter",
			DataSource:     tmpData,
		},
		MinSize: Size{500, 300},
		Layout:  VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{
						ColumnSpan: 1,
						Text:       "过滤语气词:",
					},
					TextEdit{
						ColumnSpan: 1,
						MinSize:    Size{150, 80},
						Text:       Bind("GlobalFilter.Words"),
						VScroll:    true,
					},
					Label{
						ColumnSpan: 2,
						Text: "说明：\r\n“过滤语气词” 支持设置多个，请保持每个词语都单独一行",
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
							confirmCall(tmpData.GlobalFilter.Words)
							//参数验证
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



type DefinedRuleTableRows struct {
	Id int
	Target string //目标规则
	Replace string //替换规则
	Way int //规则类型
}
type DefinedRuleTableModel struct {
	walk.SortedReflectTableModelBase
	maxIndex int
	items []*DefinedRuleTableRows
}
func NewDefinedRuleTableModel () *DefinedRuleTableModel {
	t := new(DefinedRuleTableModel)
	return t
}
func (m *DefinedRuleTableModel) Items() interface{} {
	return m.items
}
func (m *DefinedRuleTableModel) AddRow (row *DefinedRuleTableRows) {
	m.maxIndex++
	row.Id = m.maxIndex;
	m.items = append(m.items , row)
}
func (m *DefinedRuleTableModel) BatchDelRow (indexs []int) {
	id := make([]int , 0)
	for row_i , row_v := range m.items {
		for _ , op_v := range indexs {
			if row_i == op_v {
				id = append(id , row_v.Id)
			}
		}
	}
	for _ , vid := range id {
		m.DelRow(vid)
	}
}
func (m *DefinedRuleTableModel) GetRowIndex (index int) *DefinedRuleTableRows {
	tmp := new(DefinedRuleTableRows)
	for row_i , row_v := range m.items {
		if row_i == index {
			tmp = row_v
			break
		}
	}
	return tmp
}
func (m *DefinedRuleTableModel) DelRow (id int) {
	t := len(m.items)
	for row_i , row_v := range m.items {
		if row_v.Id == id {
			if row_i == 0 {
				if t <= 1 {
					m.items = make([]*DefinedRuleTableRows , 0)
				} else {
					m.items = m.items[row_i+1:]
				}
			} else if row_i+1 >= t {
				m.items = m.items[:row_i]
			} else {
				m.items = append(m.items[:row_i] , m.items[row_i+1:]...)
			}
			break
		}
	}
}
func (tv *DefinedRuleTableModel) SetAndInitFilterRules (rules []*AppDefinedFilterRule)  {
	//初始化
	tv.maxIndex = 0
	tv.items = make([]*DefinedRuleTableRows , 0)

	for _ , v := range rules {
		tv.maxIndex++
		tv.items = append(tv.items , &DefinedRuleTableRows{
			Id:tv.maxIndex,
			Target:v.Target,
			Replace:v.Replace,
			Way:v.Way,
		})
	}
}
func (tv *DefinedRuleTableModel) GetFilterRuleResult () []*AppDefinedFilterRule {
	result := make([]*AppDefinedFilterRule , 0)
	for _ , v := range tv.items {
		result = append(result , &AppDefinedFilterRule{
			Target:v.Target,
			Replace:v.Replace,
			Way:v.Way,
		})
	}
	return result
}


//运行 自定义过滤设置 Dialog
func (mw *MyMainWindow) RunDefinedFilterSetingDialog (owner walk.Form , historyRule []*AppDefinedFilterRule , confirmCall func(rule []*AppDefinedFilterRule)) {
	var dlg *walk.Dialog
	var acceptPB, cancelPB *walk.PushButton
	var tv *walk.TableView

	tableModel := NewDefinedRuleTableModel()
	tableModel.SetAndInitFilterRules(historyRule)

	var currentIndexs []int = make([]int , 0) //选择的项

	Dialog{
		AssignTo:      &dlg,
		Title:         "自定义过滤设置",
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		MinSize: Size{600, 500},
		Layout:  VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					PushButton{
						Text: "新增规则",
						OnClicked: func() {
							copyRow := new(DefinedRuleTableRows)
							if len(currentIndexs) == 1 {
								copyRow = tableModel.GetRowIndex(currentIndexs[0])
							}

							mw.RunNewDefinedFilterRuleDialog(mw , copyRow , func(rule *DefinedRuleTableRows) {
								tableModel.AddRow(rule)
								tv.SetModel(tableModel)
							})
						},
					},
					PushButton{
						Text: "删除规则",
						OnClicked: func() {
							if len(currentIndexs) < 1 {
								mw.NewErrormationTips("错误" , "请选择操作的对象")
								return
							}
							tableModel.BatchDelRow(currentIndexs)
							tv.SetModel(tableModel)
						},
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					TableView{
						Name:"tableView",
						AssignTo:         &tv,
						AlternatingRowBG: true,
						NotSortableByHeaderClick: true,
						MultiSelection:true,
						Columns: []TableViewColumn{
							{Title: "编号", DataMember: "Id" , Width:90},
							{Title: "类型", DataMember: "Way", Width:90 , FormatFunc: func(value interface{}) string {
								switch v := value.(type) {
								case int:
									if v == FILTER_TYPE_STRING {
										return "文本过滤"
									}
									if v == FILTER_TYPE_REGX {
										return "正则过滤"
									}
									return ""
								default:
									return ""
								}
							}},
							{Title: "目标规则", DataMember: "Target", Width:165},
							{Title: "替换规则", DataMember: "Replace", Width:185},
						},
						Model: tableModel,
						OnSelectedIndexesChanged: func() {
							indexs := tv.SelectedIndexes()
							if (len(indexs) > 0) {
								currentIndexs = indexs
							} else {
								currentIndexs = []int{}
							}
						},
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
							confirmCall(tableModel.GetFilterRuleResult())

							//参数验证
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
	}.Run(owner)
}

//新建自定义过滤规则 Dialog
func (mw *MyMainWindow) RunNewDefinedFilterRuleDialog (owner walk.Form , copyRows *DefinedRuleTableRows , confirmCall func(rule *DefinedRuleTableRows)) {
	var dlg *walk.Dialog
	var db *walk.DataBinder
	var acceptPB, cancelPB *walk.PushButton

	var tmpData = new(DefinedRuleTableRows)
	if copyRows.Id != 0 {
		tmpData.Target = copyRows.Target
		tmpData.Replace = copyRows.Replace
		tmpData.Way = copyRows.Way
	} else {
		tmpData.Way = 1 //默认
	}

	Dialog{
		AssignTo:      &dlg,
		Title:         "新增自定义过滤规则",
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		DataBinder: DataBinder{
			AssignTo:       &db,
			Name:           "defined",
			DataSource:     tmpData,
		},
		MinSize: Size{500, 300},
		Layout:  VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{
						ColumnSpan: 1,
						Text:       "目标规则:",
					},
					LineEdit{
						ColumnSpan: 1,
						MinSize:    Size{Width:150},
						Text:       Bind("Target"),
					},
					Label{
						ColumnSpan: 1,
						Text:       "替换规则:",
					},
					LineEdit{
						ColumnSpan: 1,
						MinSize:    Size{Width:150},
						Text:       Bind("Replace"),
					},
					Label{
						Text: "过滤类型:",
					},
					ComboBox{
						Value: Bind("Way", SelRequired{}),
						BindingMember: "Id",
						DisplayMember: "Name",
						Model: GetFilterTypeOptionsSelects(),
					},

					Label{
						ColumnSpan: 2,
						Text: "说明：\r\n1.“目标规则” 填写查找的文本/正则， “替换规则” 填写替换的文本/正则\r\n2.过滤类型为正则时，“替换规则” 允许使用 $1...$9 进行反向引用",
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
							if strings.TrimSpace(tmpData.Target) == "" {
								mw.NewErrormationTips("错误" , "必须填写目标规则噢")
								return
							}
							if tmpData.Way == FILTER_TYPE_REGX {
								//正则规则
								//校验规则
								_, e := regexp.Compile(tmpData.Target)
								if e != nil {
									mw.NewErrormationTips("错误" , "目标正则规则格式校验不通过，请检查是否正确")
									return
								}
							}

							confirmCall(tmpData)
							//参数验证
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