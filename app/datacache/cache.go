package datacache

import (
	"os"
	"videosrt/app/tool"
)

type AppCache struct {
	RootDir string
	Dir string
	File string
}

func NewAppCahce(rootdir string , file string) *AppCache {
	c := new(AppCache)
	c.Dir = "data/json"
	c.RootDir = rootdir
	c.File = c.RootDir + "/" + c.Dir + "/" + file + ".json"

	c.initDir()
	return c
}

func (app *AppCache) Set(data interface{})  {
	file := app.File
	if err := SavetoJson(data , file); err != nil {
		panic(err)
	}
}


func (app *AppCache) Get(structs interface{}) interface{} {
	file := app.File
	err, data := GettoJson(file, structs)
	if err != nil {
		if os.IsNotExist(err) {
			return structs
		}
		return data
	}
	return data
}


func (app *AppCache) initDir()  {
	fileDir := app.RootDir + "/" + app.Dir
	if !tool.DirExists(fileDir) {
		//创建目录
		if err := tool.CreateDir(fileDir , false); err != nil {
			panic(err)
		}
	}
}