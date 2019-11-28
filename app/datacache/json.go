package datacache

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func SavetoJson(data interface{} , path string) error {
	string , err := json.Marshal(data)
	if err != nil{
		return err
	}
	err = ioutil.WriteFile(path , string , os.ModePerm)
	if err != nil{
		return err
	}
	return nil
}

func GettoJson(path string , structs interface{}) (error , interface{}) {
	data , err := ioutil.ReadFile(path)
	if err != nil{
		return err , structs
	}
	err = json.Unmarshal(data , structs)
	if err != nil{
		return err , structs
	}
	return nil , structs
}
