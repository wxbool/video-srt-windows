package app

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"time"
	"videosrt/app/parse"
	"videosrt/app/tool"
	"videosrt/app/translate"
)

//字幕行结构
type SrtRows struct {
	Id int //字幕自然行id
	Number string //字幕序号
	TimeStart string //字幕开始时间戳
	TimeStartSecond float64 //字幕开始（秒）
	TimeStartMilliSecond int64 //字幕开始（毫秒）
	TimeEnd string //字幕结束时间戳
	TimeEndSecond float64 //字幕结束（秒）
	TimeEndMilliSecond int64 //字幕结束（毫秒）
	Text string //字幕文本
	TranslateText string //翻译字幕文本
}

//应用翻译配置
type SrtTranslateStruct struct {
	TranslateSwitch bool //字幕翻译开关
	Supplier int //引擎供应商
	BilingualSubtitleSwitch bool //是否输出双语字幕
	InputLanguage int //输入字幕语言
	OutputLanguage int //输出字幕语言
	OutputMainSubtitleInputLanguage bool //双语主字幕（输入语言）

	BaiduTranslate translate.BaiduTranslate //百度翻译
	TengxunyunTranslate translate.TengxunyunTranslate //腾讯云翻译
}

//统一翻译结果
type SrtTranslateResult struct {
	From string //翻译源语言
	To string //译文语言
	TransResultSrc string //翻译结果（原文）
	TransResultDst string //翻译结果（译文）
}

//主应用
type SrtTranslateApp struct {
	AppDir string //应用根目录
	SrtDir string //文件输出目录
	OutputType *AppSetingsOutput //输出文件类型
	OutputEncode int //输出文件编码
	MaxConcurrency int //最大处理并发数
	TranslateCfg *SrtTranslateStruct //翻译配置

	LogHandler func(s string , file string) //日志回调
	SuccessHandler func(file string) //成功回调
	FailHandler func(file string) //失败回调
}

//获取应用
func NewSrtTranslateApp(appDir string) *SrtTranslateApp {
	app := new(SrtTranslateApp)
	app.AppDir = appDir
	app.OutputEncode = OUTPUT_ENCODE_UTF8 //默认输出文件编码
	app.MaxConcurrency = 2

	//实例应用翻译配置
	app.TranslateCfg = new(SrtTranslateStruct)
	return app
}

//加载翻译设置
func (app *SrtTranslateApp) InitTranslateConfig (translateSettings *SrtTranslateStruct) {
	app.TranslateCfg = translateSettings
}

func (app *SrtTranslateApp) SetSrtDir(dir string)  {
	app.SrtDir = dir
}

func (app *SrtTranslateApp) SetOutputType(output *AppSetingsOutput)  {
	app.OutputType = output
}

func (app *SrtTranslateApp) SetOutputEncode(encode int)  {
	app.OutputEncode = encode
}

func (app *SrtTranslateApp) SetMaxConcurrency(number int)  {
	if number == 0 {
		number = 2
	}
	app.MaxConcurrency = number
}

func (app *SrtTranslateApp) SetSuccessHandler(callback func(file string))  {
	app.SuccessHandler = callback
}

func (app *SrtTranslateApp) SetFailHandler(callback func(file string))  {
	app.FailHandler = callback
}

func (app *SrtTranslateApp) SetLogHandler(callback func(s string , file string))  {
	app.LogHandler = callback
}

//输出日志
func (app *SrtTranslateApp) Log(s string , file string)  {
	app.LogHandler(s , file)
}


//应用运行
func (app *SrtTranslateApp) Run(srtfile string) {
	//致命错误捕获
	defer func() {
		//拦截panic
		if err := recover(); err != nil {
			//失败回调
			app.FailHandler(srtfile)

			e , ok := err.(error)
			if ok {
				app.Log("错误：" + e.Error() , srtfile)
			} else {
				panic(err)
			}
		}
	}()

	//校验文件
	if tool.VaildFile(srtfile) != true {
		panic("字幕文件不存在")
	}

	//分析、解析字幕
	srt := &parse.Srt{File:srtfile}
	subtitleParse := parse.NewSubtitleParse(srt)
	if e := subtitleParse.Parse(); e!=nil {
		panic("字幕文件解析错误：" + e.Error())
	}
	//获取字幕结构
	srtRows := app.SrtHandleGetData(subtitleParse)

	app.Log("字幕文件解析成功" , srtfile)

	//字幕翻译
	app.SrtTranslate(srtfile , srtRows)

	//输出文件
	if app.OutputType.SRT {
		app.SrtOutputFile(srtfile , srtRows , OUTPUT_SRT)
	}
	if app.OutputType.LRC {
		app.SrtOutputFile(srtfile , srtRows , OUTPUT_LRC)
	}
	if app.OutputType.TXT {
		app.SrtOutputFile(srtfile , srtRows , OUTPUT_STRING)
	}

	app.Log("处理完成" , srtfile)

	//成功回调
	app.SuccessHandler(srtfile)
}


//字幕结构处理
func (app *SrtTranslateApp) SrtHandleGetData(srt *parse.SubtitleParse) ([] *SrtRows)  {
	data := make([]*SrtRows , 0)
	for _,rows := range srt.Rows {
		temp := new(SrtRows)
		temp.Id = rows.Id
		temp.Number = rows.Number
		temp.TimeStart = rows.TimeStart
		temp.TimeStartSecond = rows.TimeStartSecond
		temp.TimeStartMilliSecond = int64(rows.TimeStartSecond * 1000)
		temp.TimeEnd = rows.TimeEnd
		temp.TimeEndSecond = rows.TimeEndSecond
		temp.TimeEndMilliSecond = int64(rows.TimeEndSecond * 1000)

		line := srt.Srt.AppointBilingualRows
		if line < 0 || line > 1 {
			line = 0
		}
		if len(rows.Text) == 0 {
			continue
		}
		temp.Text = rows.Text[line]
		data = append(data , temp)
	}
	return data
}


//字幕翻译处理
func (app *SrtTranslateApp) SrtTranslate(file string , srtRows []*SrtRows)  {
	//输出日志
	if app.TranslateCfg.TranslateSwitch {
		app.Log("字幕翻译处理中 ..." , file)

		//百度翻译标准版
		if app.TranslateCfg.Supplier == TRANSLATE_SUPPLIER_BAIDU && app.TranslateCfg.BaiduTranslate.AuthenType == translate.ACCOUNT_COMMON_AUTHEN {
			app.Log("你使用的是 “百度翻译标准版” 账号，翻译速度较慢，请耐心等待 ..." , file)
		}
	} else {
		return
	}

	total := len(srtRows)
	var lastrv float64 = 0
	//翻译处理
	index := 0
	for _,data := range srtRows {
		transResult,e := app.RunTranslate(data.Text)
		if e != nil {
			panic(e) //终止翻译
		}
		data.TranslateText = strings.TrimSpace(transResult.TransResultDst) //译文
		index++

		rv := (float64(index)/float64(total))*100
		if (rv - lastrv) > 20 {
			//输出比例
			app.Log("字幕翻译已处理：" + fmt.Sprintf("%.2f" , rv) + "%" , file)
			lastrv = rv
		}
	}
}


//文件输出
func (app *SrtTranslateApp) SrtOutputFile(file string , srtRows []*SrtRows , outputType int)  {
	var subfileDir string
	if app.SrtDir == "" {
		subfileDir = path.Dir(file)
	} else {
		subfileDir = app.SrtDir
	}
	subfile := tool.GetFileBaseName(file)
	//输出文件名
	thisfile := subfileDir + "/" + subfile + "_translaste"
	//文件后缀
	if outputType == OUTPUT_SRT {
		thisfile += ".srt"
	} else if outputType == OUTPUT_STRING {
		thisfile += ".txt"
	} else if outputType == OUTPUT_LRC {
		thisfile += ".lrc"
	}
	//创建文件
	fd, e := os.Create(thisfile)
	if e != nil {
		panic("创建文件失败：" + e.Error())
	}
	defer fd.Close()
	//文件编码分支
	if app.OutputEncode == OUTPUT_ENCODE_UTF8_BOM {
		if _, e := fd.Write([]byte{0xEF, 0xBB, 0xBF});e != nil {
			panic(e)
		}
	}
	//歌词头
	if outputType == OUTPUT_LRC {
		_,_ = fd.WriteString("[ar:]\r\n[ti:]\r\n[al:]\r\n[by:]\r\n")
	}

	//主字幕
	bilingualAsc := app.TranslateCfg.OutputMainSubtitleInputLanguage
	index := 0

	//普通文本容器
	var txtOutputContent bytes.Buffer
	var txtTransalteOutputContent bytes.Buffer

	for _ , data := range srtRows {
		var linestr string

		//字幕、歌词文件处理
		if outputType == OUTPUT_SRT || outputType == OUTPUT_LRC {
			//拼接
			if outputType == OUTPUT_SRT {
				if app.TranslateCfg.TranslateSwitch {
					if app.TranslateCfg.BilingualSubtitleSwitch {
						linestr = MakeSubtitleText(index , data.TimeStartMilliSecond , data.TimeEndMilliSecond , data.Text , data.TranslateText , true , bilingualAsc)
					} else {
						linestr = MakeSubtitleText(index , data.TimeStartMilliSecond , data.TimeEndMilliSecond , data.TranslateText , "" , false , true)
					}
				} else {
					linestr = MakeSubtitleText(index , data.TimeStartMilliSecond , data.TimeEndMilliSecond , data.Text , "" , false , true)
				}
			} else if outputType == OUTPUT_LRC {
				if app.TranslateCfg.TranslateSwitch {
					linestr = MakeMusicLrcText(index , data.TimeStartMilliSecond , data.TimeEndMilliSecond , data.TranslateText)
				} else {
					linestr = MakeMusicLrcText(index , data.TimeStartMilliSecond , data.TimeEndMilliSecond , data.Text)
				}
			}

			//写入行
			if _, e = fd.WriteString(linestr);e != nil {
				panic("写入文件失败：" + e.Error())
			}
		} else if outputType == OUTPUT_STRING {
			//普通文本处理
			txtOutputContent.WriteString(data.Text)
			txtOutputContent.WriteString("\r\n")

			if app.TranslateCfg.TranslateSwitch {
				txtTransalteOutputContent.WriteString(data.TranslateText)
				txtTransalteOutputContent.WriteString("\r\n")
			}
		}

		index++
	}

	//写入文本文件
	if outputType == OUTPUT_STRING {
		txtOutputContent.WriteString("\r\n\r\n\r\n\r\n\r\n")
		if _, e = fd.WriteString(txtOutputContent.String());e != nil {
			panic("写入文件失败：" + e.Error())
		}
		if _, e = fd.WriteString(txtTransalteOutputContent.String());e != nil {
			panic("写入文件失败：" + e.Error())
		}
	}
}


//统一运行翻译结果
func (app *SrtTranslateApp) RunTranslate(s string) (*SrtTranslateResult , error) {
	translateResult := new(SrtTranslateResult)

	if app.TranslateCfg.Supplier == TRANSLATE_SUPPLIER_BAIDU {
		if app.TranslateCfg.BaiduTranslate.AuthenType == translate.ACCOUNT_COMMON_AUTHEN { //百度翻译标准版
			//休眠1010毫秒
			time.Sleep(time.Millisecond * 1010)
		} else {
			//休眠200毫秒
			time.Sleep(time.Millisecond * 200)
		}

		from := GetLanguageChar(app.TranslateCfg.InputLanguage , TRANSLATE_SUPPLIER_BAIDU)
		to := GetLanguageChar(app.TranslateCfg.OutputLanguage , TRANSLATE_SUPPLIER_BAIDU)

		baiduResult,e := app.TranslateCfg.BaiduTranslate.TranslateBaidu(s , from , to)
		if e != nil {
			return translateResult,errors.New("翻译失败！错误信息：" + e.Error())
		}
		translateResult.TransResultDst = baiduResult.TransResultDst
		translateResult.TransResultSrc = baiduResult.TransResultSrc
		translateResult.From = baiduResult.From
		translateResult.To = baiduResult.To

		return translateResult,nil
	} else if app.TranslateCfg.Supplier == TRANSLATE_SUPPLIER_TENGXUNYUN {
		//休眠
		t :=  app.TranslateCfg.TengxunyunTranslate.TranslateSleepTime(app.MaxConcurrency)
		time.Sleep(t)

		from := GetLanguageChar(app.TranslateCfg.InputLanguage , TRANSLATE_SUPPLIER_TENGXUNYUN)
		to := GetLanguageChar(app.TranslateCfg.OutputLanguage , TRANSLATE_SUPPLIER_TENGXUNYUN)

		txResult,e := app.TranslateCfg.TengxunyunTranslate.TranslateTengxunyun(s , from , to)
		if e != nil {
			return translateResult,errors.New("翻译失败！错误信息：" + e.Error())
		}

		translateResult.TransResultDst = txResult.TransResultDst
		translateResult.TransResultSrc = txResult.TransResultSrc
		translateResult.From = txResult.From
		translateResult.To = txResult.To

		return translateResult,nil
	} else {
		return translateResult,errors.New("翻译失败！缺少翻译引擎！")
	}
}