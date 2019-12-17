package translate

import (
	"errors"
	"github.com/buger/jsonparser"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"videosrt/app/tool"
)

//百度翻译
type BaiduTranslate struct {
	AppId string //appid
	AppSecret string //appsecret
}

//百度翻译结果集
type BaiduTranslateResult struct {
	From string //翻译源语言
	To string //译文语言
	TransResultSrc string //翻译结果（原文）
	TransResultDst string //翻译结果（译文）
	ErrorCode int64 //错误码（仅当出现错误时存在）
	ErrorMsg string //错误消息（仅当出现错误时存在）
}

//常量
const (
	TRANS_API string = "https://fanyi-api.baidu.com/api/trans/vip/translate"
)


//百度api文档
//http://api.fanyi.baidu.com/api/trans/product/apidoc
//支持语言列表 http://api.fanyi.baidu.com/api/trans/product/apidoc#languageList
func (trans *BaiduTranslate) Translate (strings string , from string , to string) (*BaiduTranslateResult , error) {
	params := &url.Values{}

	params.Add("q" , strings)
	params.Add("appid" , trans.AppId)
	params.Add("salt" , strconv.FormatInt(tool.GetIntRandomNumber(10000 , 99999) , 10))
	params.Add("from" , from)
	params.Add("to" , to)
	params.Add("sign" , trans.BuildSign(strings , params.Get("salt")))

	return trans.CallRequest(params)
}

//生成加密sign
func (trans *BaiduTranslate) BuildSign (strings string , salt string) string {
	str := trans.AppId + strings + salt + trans.AppSecret
	return tool.Md5String(str)
}

//发起请求
func (trans *BaiduTranslate) CallRequest (params *url.Values ) (*BaiduTranslateResult , error) {
	url := TRANS_API + "?" +  params.Encode()

	request, e := http.NewRequest(http.MethodGet, url , nil)
	if e != nil {
		return nil,e
	}
	//do request
	response, e := http.DefaultClient.Do(request)
	if e != nil {
		return nil,e
	}
	//content
	content, e := ioutil.ReadAll(response.Body)
	if e != nil {
		return nil,e
	}

	//解析数据
	errorCode , _ := jsonparser.GetString(content , "error_code")
	errorMsg , _ := jsonparser.GetString(content , "error_msg")
	from , _ := jsonparser.GetString(content , "from")
	to , _ := jsonparser.GetString(content , "to")

	errorCodeInt , _ := strconv.Atoi(errorCode)

	result := &BaiduTranslateResult{
		ErrorCode:int64(errorCodeInt),
		ErrorMsg:errorMsg,
		From:from,
		To:to,
	}

	_, _ = jsonparser.ArrayEach(content, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		result.TransResultSrc, _ = jsonparser.GetString(value, "src")
		result.TransResultDst, _ = jsonparser.GetString(value, "dst")
	}, "trans_result")

	//翻译错误校验
	if result.ErrorCode != 0 {
		return nil , errors.New("翻译失败 , 错误提示（" + result.ErrorMsg + "）")
	}

	return result,nil
}