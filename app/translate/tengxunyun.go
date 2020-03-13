package translate

import (
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"
	v20180321 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tmt/v20180321"
	"time"
)

//腾讯云机器翻译
type TengxunyunTranslate struct {
	SecretId string //secretId
	SecretKey string //secretKey
}


//腾讯云翻译结果集
type TengxunyunTranslateResult struct {
	From string //翻译源语言
	To string //译文语言
	TransResultSrc string //翻译结果（原文）
	TransResultDst string //翻译结果（译文）
}


//调起腾讯云机器文本翻译
func (trans *TengxunyunTranslate) TranslateTengxunyun (strings string , from string , to string) (*TengxunyunTranslateResult , error) {
	//fmt.Println("TranslateTengxunyun : " , strings , from , to)
	credential := common.NewCredential(
		trans.SecretId ,
		trans.SecretKey ,
	)

	cpf := profile.NewClientProfile()
	client, _ := v20180321.NewClient(credential, regions.Guangzhou, cpf)
	request := v20180321.NewTextTranslateRequest()

	var SourceText *string = new(string)
	var Source *string = new(string)
	var Target *string = new(string)
	var ProjectId *int64 = new(int64)

	//*SourceText = "这个需求做不了"
	//*Source = "zh"
	//*Target = "en"
	//*ProjectId = 0

	*SourceText = strings
	*Source = from
	*Target = to
	*ProjectId = 0

	request.SourceText = SourceText
	request.Source = Source
	request.Target = Target
	request.ProjectId = ProjectId

	result := new(TengxunyunTranslateResult)

	// 通过client对象调用想要访问的接口，需要传入请求对象
	response, err := client.TextTranslate(request)
	// 处理异常
	if errorObj, ok := err.(*errors.TencentCloudSDKError); ok {
		return result , errorObj
	}
	// 非SDK异常，直接失败。实际代码中可以加入其他的处理。
	if err != nil {
		return result , err
	}

	result.From = *response.Response.Source
	result.To = *response.Response.Target
	result.TransResultSrc = strings
	result.TransResultDst = *response.Response.TargetText

	return result,nil
}


//获取并发请求停顿的时间
func (trans *TengxunyunTranslate) TranslateSleepTime (maxConcurrency int) (time.Duration) {
	if maxConcurrency == 1 {
		return time.Millisecond * 250
	} else if maxConcurrency == 2 {
		return time.Millisecond * 600
	} else {
		return time.Millisecond * time.Duration(maxConcurrency) * 900
	}
}