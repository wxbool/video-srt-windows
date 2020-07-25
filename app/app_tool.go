package app

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"
	"videosrt/app/tool"
)


//语气词过滤
func ModalWordsFilter(s string , w string) string {
	tmpText := strings.ReplaceAll(s , w , "")
	if strings.TrimSpace(tmpText) == "" || tool.CheckOnlySymbolText(strings.TrimSpace(tmpText)) {
		return ""
	} else {
		//尝试过滤重复语气词
		compile, e := regexp.Compile(w + "{2,}")
		if e != nil {
			return s
		}
		return compile.ReplaceAllString(s , "")
	}
}

//自定义规则过滤
func DefinedWordRuleFilter(s string , rule *AppDefinedFilterRule) string {
	if rule.Way == FILTER_TYPE_STRING {
		//文本过滤
		s = strings.ReplaceAll(s , rule.Target , rule.Replace)
	} else if rule.Way == FILTER_TYPE_REGX {
		//正则过滤
		compile, e := regexp.Compile(rule.Target)
		if e != nil {
			return s
		}
		s = compile.ReplaceAllString(s , rule.Replace)
	}
	if strings.TrimSpace(s) == "" || tool.CheckOnlySymbolText(strings.TrimSpace(s)) {
		return ""
	}
	return s
}


//拼接字幕字符串
func MakeSubtitleText(index int , startTime int64 , endTime int64 , text string , translateText string , bilingualSubtitleSwitch bool , bilingualAsc bool) string {
	var content bytes.Buffer
	content.WriteString(strconv.Itoa(index))
	content.WriteString("\r\n")
	content.WriteString(tool.SubtitleTimeMillisecond(startTime , true))
	content.WriteString(" --> ")
	content.WriteString(tool.SubtitleTimeMillisecond(endTime , true))
	content.WriteString("\r\n")

	//输出双语字幕
	if bilingualSubtitleSwitch {
		if bilingualAsc {
			content.WriteString(text)
			content.WriteString("\r\n")
			content.WriteString(translateText)
		} else {
			content.WriteString(translateText)
			content.WriteString("\r\n")
			content.WriteString(text)
		}
	} else {
		content.WriteString(text)
	}

	content.WriteString("\r\n")
	content.WriteString("\r\n")
	return content.String()
}


//拼接文本格式
func MakeText(index int , startTime int64 , endTime int64 , text string) string {
	var content bytes.Buffer
	content.WriteString(text)
	content.WriteString("\r\n")
	content.WriteString("\r\n")
	return content.String()
}


//拼接歌词文本
func MakeMusicLrcText(index int , startTime int64 , endTime int64 , text string) string {
	var content bytes.Buffer
	content.WriteString("[")
	content.WriteString(tool.MusicLrcTextMillisecond(startTime))
	content.WriteString("]")
	content.WriteString(text)
	content.WriteString("\r\n")
	return content.String()
}