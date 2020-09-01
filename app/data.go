package app

import (
	"videosrt/app/aliyun"
	"videosrt/app/datacache"
	"videosrt/app/translate"
)

//输出文件类型
const(
	OUTPUT_SRT = 1 //字幕SRT文件
	OUTPUT_STRING = 2 //普通文本
	OUTPUT_LRC = 3 //LRC文本
)

//输出文件编码
const(
	OUTPUT_ENCODE_UTF8 = 1 //文件编码 utf-8
	OUTPUT_ENCODE_UTF8_BOM = 2 //文件编码 utf-8 带 BOM
)

//翻译语言
const (
	LANGUAGE_ZH = 1 //中文
	LANGUAGE_EN = 2 //英文
	LANGUAGE_JP = 3 //日语
	LANGUAGE_KOR = 4 //韩语
	LANGUAGE_FRA = 5 //法语 fra
	LANGUAGE_DE = 6 //德语 de
	LANGUAGE_SPA = 7 //西班牙语 spa
	LANGUAGE_RU = 8 //俄语 ru
	LANGUAGE_IT = 9 //意大利语 it
	LANGUAGE_TH = 10 //泰语 th
)

//缓存结构
type OssAppStruct struct {
	Data *datacache.AppCache
}
type TranslateEngineAppStruct struct {
	Data *datacache.AppCache
}
type SpeechEngineAppStruct struct {
	Data *datacache.AppCache
}
type AppSetingsAppStruct struct {
	Data *datacache.AppCache
}
type AppFilterAppStruct struct {
	Data *datacache.AppCache
}

var RootDir string
var Oss *OssAppStruct
var Translate *TranslateEngineAppStruct
var Engine *SpeechEngineAppStruct
var Setings *AppSetingsAppStruct
var Filter *AppFilterAppStruct


func init()  {
	RootDir = GetAppRootDir()
	if RootDir == "" {
		panic("应用根目录获取失败")
	}

	Oss = new(OssAppStruct)
	Translate = new(TranslateEngineAppStruct)
	Engine = new(SpeechEngineAppStruct)
	Setings = new(AppSetingsAppStruct)
	Filter = new(AppFilterAppStruct)

	Oss.Data =  datacache.NewAppCahce(RootDir , "oss")
	Translate.Data =  datacache.NewAppCahce(RootDir , "translate_engine")
	Engine.Data =  datacache.NewAppCahce(RootDir , "engine")
	Setings.Data =  datacache.NewAppCahce(RootDir , "setings")
	Filter.Data =  datacache.NewAppCahce(RootDir , "filter")
}


//表单结构
type OperateFrom struct {
	EngineId int //当前语音引擎
	TranslateEngineId int //当前翻译引擎

	TranslateSwitch bool //字幕翻译开关
	BilingualSubtitleSwitch bool //是否输出双语字幕
	InputLanguage int //输入字幕语言
	OutputLanguage int //输出字幕语言
	OutputMainSubtitleInputLanguage bool //双语主字幕（输入语言）

	OutputSrt bool
	OutputLrc bool
	OutputTxt bool

	OutputType *AppSetingsOutput //输出文件类型
	OutputEncode int //输出文件编码
	SoundTrack int //输出音轨
}

//输出类型选项
type OutputSelects struct {
	Id   int
	Name string
}

//输出音轨类型选项
type SoundTrackSelects struct {
	Id   int
	Name string
}

//字幕翻译语言选项列表
type LanguageSelects struct {
	Id   int
	Name string
}

type AppSetingsOutput struct {
	SRT bool
	LRC bool
	TXT bool
}

//应用配置结构
type AppSetings struct {
	CurrentEngineId int //目前语音引擎Id
	CurrentTranslateEngineId int //目前翻译引擎Id
	MaxConcurrency int //任务最大处理并发数
	OutputType *AppSetingsOutput //输出文件类型
	OutputEncode int //输出文件编码
	SrtFileDir string //Srt文件输出目录
	SoundTrack int //输出音轨

	TranslateSwitch bool //字幕翻译开关
	BilingualSubtitleSwitch bool //是否输出双语字幕
	InputLanguage int //输入字幕语言
	OutputLanguage int //输出字幕语言
	OutputMainSubtitleInputLanguage bool //双语主字幕（输入语言）

	CloseIntelligentBlockSwitch bool //关闭智能分段
	CloseNewVersionMessage bool //关闭软件新版本提醒（默认开启）[false开启 true关闭]
	CloseAutoDeleteOssTempFile bool //关闭自动删除临时音频文件（默认开启）[false开启 true关闭]
}


const (
	FILTER_TYPE_STRING = 1 //文本过滤
	FILTER_TYPE_REGX = 2 //正则过滤
)
//自定义过滤器规则
type AppDefinedFilterRule struct {
	Target string //目标规则
	Replace string //替换规则
	Way int //规则类型
}
//应用字幕过滤器结构
type AppFilterSetings struct {
	//通用过滤器
	GlobalFilter struct{
		Switch bool
		Words string //过滤词组
	}
	//自定义过滤器
	DefinedFilter struct{
		Switch bool
		Rule [] *AppDefinedFilterRule
	}
}


//任务文件列表 - 结构
type TaskHandleFile struct {
	Files [] string
}

//根据配置初始化表单
func (from *OperateFrom) Init(setings *AppSetings)  {
	from.OutputType = new(AppSetingsOutput)
	if setings.CurrentEngineId != 0 {
		from.EngineId = setings.CurrentEngineId
	}
	if setings.CurrentTranslateEngineId != 0 {
		from.TranslateEngineId = setings.CurrentTranslateEngineId
	}

	if !setings.OutputType.LRC && !setings.OutputType.SRT && !setings.OutputType.TXT {
		from.OutputType.SRT = true
		from.OutputSrt = true
	} else {
		from.OutputType = setings.OutputType
		if setings.OutputType.SRT {
			from.OutputSrt = true
		}
		if setings.OutputType.TXT {
			from.OutputTxt = true
		}
		if setings.OutputType.LRC {
			from.OutputLrc = true
		}
	}

	if setings.OutputEncode == 0 {
		from.OutputEncode = OUTPUT_ENCODE_UTF8 //默认编码
	} else {
		from.OutputEncode = setings.OutputEncode
	}

	from.OutputMainSubtitleInputLanguage = setings.OutputMainSubtitleInputLanguage

	if setings.SoundTrack == 0 {
		from.SoundTrack = 1 //默认输出音轨一
	} else {
		from.SoundTrack = setings.SoundTrack
	}

	//默认翻译设置
	if setings.InputLanguage == 0 {
		from.InputLanguage = LANGUAGE_ZH
	} else {
		from.InputLanguage = setings.InputLanguage
	}
	if setings.OutputLanguage == 0 {
		from.OutputLanguage = LANGUAGE_ZH
	} else {
		from.OutputLanguage = setings.OutputLanguage
	}

	from.TranslateSwitch = setings.TranslateSwitch
	from.BilingualSubtitleSwitch = setings.BilingualSubtitleSwitch
}

//获取 输出文件选项列表
func GetOutputOptionsSelects() []*OutputSelects {
	return []*OutputSelects{
		&OutputSelects{Id:OUTPUT_SRT , Name:"字幕文件"},
		&OutputSelects{Id:OUTPUT_STRING , Name:"普通文本"},
	}
}

//获取 输出文件编码选项列表
func GetOutputEncodeOptionsSelects() []*OutputSelects {
	return []*OutputSelects{
		&OutputSelects{Id:OUTPUT_ENCODE_UTF8 , Name:"UTF-8"},
		&OutputSelects{Id:OUTPUT_ENCODE_UTF8_BOM , Name:"UTF-8-BOM"},
	}
}

//获取 输出音轨选项列表
func GetSoundTrackSelects() []*SoundTrackSelects {
	return []*SoundTrackSelects{
		&SoundTrackSelects{Id:3 , Name:"全部"},
		&SoundTrackSelects{Id:1 , Name:"音轨一"},
		&SoundTrackSelects{Id:2 , Name:"音轨二"},
	}
}


//获取 允许翻译[输入字幕语言]列表
func GetTranslateInputLanguageOptionsSelects() []*LanguageSelects {
	return []*LanguageSelects{
		&LanguageSelects{Id:LANGUAGE_ZH , Name:"中文"},
		&LanguageSelects{Id:LANGUAGE_EN , Name:"英文"},
		&LanguageSelects{Id:LANGUAGE_JP , Name:"日语"},
		&LanguageSelects{Id:LANGUAGE_KOR , Name:"韩语"},
		&LanguageSelects{Id:LANGUAGE_FRA , Name:"法语"},
		&LanguageSelects{Id:LANGUAGE_DE , Name:"德语"},
		&LanguageSelects{Id:LANGUAGE_SPA , Name:"西班牙语"},
		&LanguageSelects{Id:LANGUAGE_RU , Name:"俄语"},
		&LanguageSelects{Id:LANGUAGE_IT , Name:"意大利语"},
		&LanguageSelects{Id:LANGUAGE_TH , Name:"泰语"},
	}
}

//获取 允许翻译[输出字幕语言]列表
func GetTranslateOutputLanguageOptionsSelects() []*LanguageSelects {
	return []*LanguageSelects{
		&LanguageSelects{Id:LANGUAGE_ZH , Name:"中文"},
		&LanguageSelects{Id:LANGUAGE_EN , Name:"英文"},
		&LanguageSelects{Id:LANGUAGE_JP , Name:"日语"},
		&LanguageSelects{Id:LANGUAGE_KOR , Name:"韩语"},
		&LanguageSelects{Id:LANGUAGE_FRA , Name:"法语"},
		&LanguageSelects{Id:LANGUAGE_DE , Name:"德语"},
		&LanguageSelects{Id:LANGUAGE_SPA , Name:"西班牙语"},
		&LanguageSelects{Id:LANGUAGE_RU , Name:"俄语"},
		&LanguageSelects{Id:LANGUAGE_IT , Name:"意大利语"},
		&LanguageSelects{Id:LANGUAGE_TH , Name:"泰语"},
	}
}


//获取 应用配置
func (setings *AppSetingsAppStruct) GetCacheAppSetingsData() *AppSetings {
	data := new(AppSetings)
	data.OutputType = new(AppSetingsOutput)
	vdata := setings.Data.Get(data)
	if v, ok := vdata.(*AppSetings); ok {
		return v
	}
	return data
}

//设置 应用配置
func (setings *AppSetingsAppStruct) SetCacheAppSetingsData(data *AppSetings)  {
	setings.Data.Set(data)
}






//获取 应用过滤器配置
func (setings *AppFilterAppStruct) GetCacheAppFilterData() *AppFilterSetings {
	data := new(AppFilterSetings)
	vdata := setings.Data.Get(data)
	if v, ok := vdata.(*AppFilterSetings); ok {
		return v
	}
	return data
}
//设置 应用过滤器配置
func (setings *AppFilterAppStruct) SetCacheAppFilterData(data *AppFilterSetings)  {
	setings.Data.Set(data)
}

//过滤类型选项结构
type FilterTypeSelects struct {
	Id   int
	Name string
}
//获取 过滤类型选项列表
func GetFilterTypeOptionsSelects() []*FilterTypeSelects {
	return []*FilterTypeSelects{
		&FilterTypeSelects{Id:FILTER_TYPE_STRING , Name:"文本替换"},
		&FilterTypeSelects{Id:FILTER_TYPE_REGX , Name:"正则替换"},
	}
}




//阿里云OSS - 缓存结构
type AliyunOssCache struct {
	aliyun.AliyunOss
}

//设置 阿里云OSS 缓存
func (oss *OssAppStruct) SetCacheAliyunOssData(data *AliyunOssCache) {
	oss.Data.Set(data)
}

//获取 阿里云OSS 缓存数据
func (oss *OssAppStruct) GetCacheAliyunOssData() *AliyunOssCache {
	data := new(AliyunOssCache)
	vdata := oss.Data.Get(data)
	if v, ok := vdata.(*AliyunOssCache); ok {
		return v
	}
	return data
}










//阿里云语音识别引擎 - 缓存结构
type AliyunEngineCache struct {
	aliyun.AliyunClound
	Id int //Id
	Alias string //别名
}

//阿里云语音识别引擎 - 列表缓存结构
type AliyunEngineListCache struct {
	Engine [] *AliyunEngineCache
}

//语音引擎选项
type EngineSelects struct {
	Id   int
	Name string
}

//阿里云语音引擎区域选项
type AliyunEngineRegionSelects struct {
	Id   int
	Name string
}

//获取 阿里云语音引擎区域选项列表
func GetAliyunEngineRegionOptionSelects() []*BaiduAuthTypeSelects {
	return []*BaiduAuthTypeSelects{
		&BaiduAuthTypeSelects{Id:aliyun.ALIYUN_CLOUND_REGION_CHA , Name:"中国"},
		&BaiduAuthTypeSelects{Id:aliyun.ALIYUN_CLOUND_REGION_INT , Name:"海外"},
	}
}

//获取 引擎选项列表
func (speechEng *SpeechEngineAppStruct) GetEngineOptionsSelects() []*EngineSelects {
	engines := make([]*EngineSelects , 0)
	//获取数据
	data := speechEng.GetCacheAliyunEngineListData()

	for _,v := range data.Engine {
		engines = append(engines , &EngineSelects{
			Id:v.Id,
			Name:v.Alias,
		})
	}
	return engines
}

//根据选择引擎id 获取 引擎数据
func (speechEng *SpeechEngineAppStruct) GetEngineById(id int) (*AliyunEngineCache , bool) {
	//获取数据
	data := speechEng.GetCacheAliyunEngineListData()
	for _,v := range data.Engine {
		if id == v.Id {
			return v , true
		}
	}
	return nil , false
}

//获取 当前引擎id 下标
func (speechEng *SpeechEngineAppStruct) GetCurrentIndex(data []*EngineSelects , id int) int {
	for index,v := range data {
		if v.Id == id {
			return index
		}
	}
	return -1
}

//获取 阿里云语音识别引擎 缓存数据
func (speechEng *SpeechEngineAppStruct) GetCacheAliyunEngineListData() *AliyunEngineListCache {
	data := new(AliyunEngineListCache)
	vdata := speechEng.Data.Get(data)
	if v, ok := vdata.(*AliyunEngineListCache); ok {
		return v
	}
	return data
}

//设置 阿里云语音识别引擎 缓存
func (speechEng *SpeechEngineAppStruct) SetCacheAliyunEngineListData(data *AliyunEngineListCache)  {
	speechEng.Data.Set(data)
}


//根据id 删除 阿里云语音识别引擎 缓存数据
func (speechEng *SpeechEngineAppStruct) RemoveCacheAliyunEngineData(id int) (bool) {
	var ok = false
	var newEngine = make([]*AliyunEngineCache , 0)
	origin := speechEng.GetCacheAliyunEngineListData()

	total := len(origin.Engine)
	for i,engine := range origin.Engine	{
		if engine.Id == id {
			if i == (total - 1) {
				newEngine = origin.Engine[:i]
			} else {
				newEngine = append(origin.Engine[:i] , origin.Engine[i+1:]...)
			}
			ok = true
			break
		}
	}
	if ok {
		origin.Engine = newEngine
		//更新缓存数据
		speechEng.SetCacheAliyunEngineListData(origin)
	}
	return ok
}











//声明引擎供应商
const (
	TRANSLATE_SUPPLIER_BAIDU = 1 //百度翻译
	TRANSLATE_SUPPLIER_TENGXUNYUN = 2 //腾讯云翻译
)

//翻译引擎 - 数据结构
type TranslateEngineStruct struct {
	TengxunyunEngine translate.TengxunyunTranslate
	BaiduEngine translate.BaiduTranslate

	Supplier int //引擎供应商
	Id int //Id
	Alias string //别名
}

//翻译引擎 - 列表缓存结构
type TranslateEngineListCacheStruct struct {
	Engine [] *TranslateEngineStruct
}

//翻译引擎选项 - 数据结构
type TranslateEngineSelects struct {
	Id   int
	Name string
}

//获取 翻译引擎 选项列表
func (translateEng *TranslateEngineAppStruct) GetTranslateEngineOptionsSelects() []*TranslateEngineSelects {
	engines := make([]*TranslateEngineSelects , 0)
	//获取数据
	data := translateEng.GetCacheTranslateEngineListData()

	for _,v := range data.Engine {
		engines = append(engines , &TranslateEngineSelects{
			Id:v.Id,
			Name:v.Alias,
		})
	}
	return engines
}

//根据 选择翻译引擎id 获取数据
func (translateEng *TranslateEngineAppStruct) GetTranslateEngineById(id int) (*TranslateEngineStruct , bool) {
	//获取数据
	data := translateEng.GetCacheTranslateEngineListData()
	for _,v := range data.Engine {
		if id == v.Id {
			return v , true
		}
	}
	return nil , false
}

//获取 当前翻译引擎id 下标
func (translateEng *TranslateEngineAppStruct) GetCurrentTranslateEngineIndex(data []*TranslateEngineSelects , id int) int {
	for index,v := range data {
		if v.Id == id {
			return index
		}
	}
	return -1
}

//设置 翻译引擎 数据缓存
func (translateEng *TranslateEngineAppStruct) SetCacheTranslateEngineListData(data *TranslateEngineListCacheStruct)  {
	translateEng.Data.Set(data)
}

//根据id 删除 翻译引擎 数据缓存
func (translateEng *TranslateEngineAppStruct) RemoveCacheTranslateEngineData(id int) (bool) {
	var ok = false
	var newEngine = make([]*TranslateEngineStruct , 0)
	origin := translateEng.GetCacheTranslateEngineListData()

	total := len(origin.Engine)
	for i,engine := range origin.Engine	{
		if engine.Id == id {
			if i == (total - 1) {
				newEngine = origin.Engine[:i]
			} else {
				newEngine = append(origin.Engine[:i] , origin.Engine[i+1:]...)
			}
			ok = true
			break
		}
	}
	if ok {
		origin.Engine = newEngine
		//更新缓存数据
		translateEng.SetCacheTranslateEngineListData(origin)
	}
	return ok
}

//获取 翻译引擎 数据缓存
func (translateEng *TranslateEngineAppStruct) GetCacheTranslateEngineListData() *TranslateEngineListCacheStruct {
	data := new(TranslateEngineListCacheStruct)
	vdata := translateEng.Data.Get(data)
	if v, ok := vdata.(*TranslateEngineListCacheStruct); ok {
		return v
	}
	return data
}


//获取不同翻译引擎的语音字符标识
func GetLanguageChar(Language int , Supplier int) string {
	if Supplier == TRANSLATE_SUPPLIER_BAIDU {
		switch Language {
		case LANGUAGE_ZH:
			return "zh"
		case LANGUAGE_EN:
			return "en"
		case LANGUAGE_JP:
			return "jp"
		case LANGUAGE_KOR:
			return "kor"
		case LANGUAGE_FRA:
			return "fra"
		case LANGUAGE_DE:
			return "de"
		case LANGUAGE_SPA:
			return "spa"
		case LANGUAGE_RU:
			return "ru"
		case LANGUAGE_IT:
			return "it"
		case LANGUAGE_TH:
			return "th"
		}
	}
	if Supplier == TRANSLATE_SUPPLIER_TENGXUNYUN {
		switch Language {
		case LANGUAGE_ZH:
			return "zh"
		case LANGUAGE_EN:
			return "en"
		case LANGUAGE_JP:
			return "jp"
		case LANGUAGE_KOR:
			return "kr"
		case LANGUAGE_FRA:
			return "fr"
		case LANGUAGE_DE:
			return "de"
		case LANGUAGE_SPA:
			return "es"
		case LANGUAGE_RU:
			return "ru"
		case LANGUAGE_IT:
			return "it"
		case LANGUAGE_TH:
			return "th"
		}
	}
	return ""
}


//百度翻译账号认证类型选项
type BaiduAuthTypeSelects struct {
	Id   int
	Name string
}

//获取 百度翻译账号认证类型
func GetBaiduTranslateAuthenTypeOptionsSelects() []*BaiduAuthTypeSelects {
	return []*BaiduAuthTypeSelects{
		&BaiduAuthTypeSelects{Id:translate.ACCOUNT_COMMON_AUTHEN , Name:"标准版"},
		&BaiduAuthTypeSelects{Id:translate.ACCOUNT_SENIOR_AUTHEN , Name:"高级版"},
	}
}