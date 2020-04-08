package app

import (
	"bytes"
	"strconv"
	"videosrt/app/tool"
)

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