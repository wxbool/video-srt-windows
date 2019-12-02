package app

import (
	"videosrt/app/aliyun"
	"videosrt/app/datacache"
)

var RootDir string

var oss,engine,setings *datacache.AppCache

//输出文件类型
const(
	OUTPUT_SRT = 1 //字幕文件
	OUTPUT_STRING = 2 //普通文本
)

func init()  {
	RootDir = GetAppRootDir()
	if RootDir == "" {
		panic("应用根目录获取失败")
	}

	oss = datacache.NewAppCahce(RootDir , "oss")
	engine = datacache.NewAppCahce(RootDir , "engine")
	setings = datacache.NewAppCahce(RootDir , "setings")
}




//设置表单
type OperateFrom struct {
	EngineId int
}

//引擎选项
type EngineSelects struct {
	Id   int
	Name string
}

//输出类型选项
type OutputSelects struct {
	Id   int
	Name string
}

//阿里云OSS - 缓存结构
type AliyunOssCache struct {
	aliyun.AliyunOss
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

//应用配置 - 缓存结构
type AppSetings struct {
	CurrentEngineId int //目前使用引擎Id
	MaxConcurrency int //任务最大处理并发数
	OutputType int //输出文件类型
	SrtFileDir string //Srt文件输出目录
}

//任务文件列表 - 结构
type TaskHandleFile struct {
	Files [] string
}


//获取 阿里云OSS 缓存数据
func GetCacheAliyunOssData() *AliyunOssCache {
	data := new(AliyunOssCache)
	vdata := oss.Get(data)
	if v, ok := vdata.(*AliyunOssCache); ok {
		return v
	}
	return data
}

//设置 阿里云OSS 缓存
func SetCacheAliyunOssData(data *AliyunOssCache) {
	oss.Set(data)
}


//获取 阿里云语音识别引擎 缓存数据
func GetCacheAliyunEngineListData() *AliyunEngineListCache {
	data := new(AliyunEngineListCache)
	vdata := engine.Get(data)
	if v, ok := vdata.(*AliyunEngineListCache); ok {
		return v
	}
	return data
}

//设置 阿里云语音识别引擎 缓存
func SetCacheAliyunEngineListData(data *AliyunEngineListCache)  {
	engine.Set(data)
}

//根据id 删除 阿里云语音识别引擎 缓存数据
func RemoveCacheAliyunEngineData(id int) (bool) {
	var ok = false
	var newEngine = make([]*AliyunEngineCache , 0)
	origin := GetCacheAliyunEngineListData()

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
		SetCacheAliyunEngineListData(origin)
	}
	return ok
}


//获取 应用配置
func GetCacheAppSetingsData() *AppSetings {
	data := new(AppSetings)
	vdata := setings.Get(data)
	if v, ok := vdata.(*AppSetings); ok {
		return v
	}
	return data
}

//设置 应用配置
func SetCacheAppSetingsData(data *AppSetings)  {
	setings.Set(data)
}


//获取 引擎选项列表
func GetEngineOtionsSelects() []*EngineSelects {
	engines := make([]*EngineSelects , 0)
	//获取数据
	data := GetCacheAliyunEngineListData()

	for _,v := range data.Engine {
		engines = append(engines , &EngineSelects{
			Id:v.Id,
			Name:v.Alias,
		})
	}
	return engines
}


//根据选择引擎id 获取 引擎数据
func GetEngineById(id int) (*AliyunEngineCache , bool) {
	//获取数据
	data := GetCacheAliyunEngineListData()
	for _,v := range data.Engine {
		if id == v.Id {
			return v , true
		}
	}
	return nil , false
}


//获取 当前引擎id 下标
func GetCurrentIndex(data []*EngineSelects , id int) int {
	for index,v := range data {
		if v.Id == id {
			return index
		}
	}
	return -1
}


//获取 输出文件选项列表
func GetOutputOptionsSelects() []*OutputSelects {
	return []*OutputSelects{
		&OutputSelects{Id:OUTPUT_SRT , Name:"字幕文件"},
		&OutputSelects{Id:OUTPUT_STRING , Name:"普通文本"},
	}
}