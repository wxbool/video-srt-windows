package parse

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)


//解析SRT行文本
func ParseSrtRows(rowText string , lineStart bool) (string , int , error) {
	rowText = strings.TrimSpace(rowText)

	if _,e := strconv.Atoi(rowText);e == nil && lineStart == false { //识别为字幕序号
		return rowText , PARSE_SRT_NUMBER , nil
	}
	r := regexp.MustCompile(`(\d{1,2}:\d{1,2}:\d{1,2}[,|\.]\d{2,4})\s*-->\s*(\d{1,2}:\d{1,2}:\d{1,2}[,|\.]\d{2,4})`)
	if r.Match([]byte(rowText)) == true {
		timeResult := r.FindSubmatch([]byte(rowText))
		if len(timeResult) != 3 {
			return "" , 0 , errors.New("解析字幕时间错误")
		}
		times := strings.Join([]string{string(timeResult[1]) , string(timeResult[2])} , "-")
		return times , PARSE_SRT_TIME_RANGE , nil
	}
	if rowText != "" {
		return rowText , PARSE_SRT_SUBTITLE_TEXT , nil
	} else {
		return rowText , PARSE_SRT_TRIM , nil
	}
}


//校验srt字幕时间格式
func VaildateSrtTimeFormat(time string) bool {
	r := regexp.MustCompile(`^(\d{1,2}:\d{1,2}:\d{1,2}[,|\.]\d{2,4})$`)
	return r.Match([]byte(time))
}


//srt字幕时间文本 -> 秒
func SrtTimeFormatToSecond(time string) (float64 , error) {
	if VaildateSrtTimeFormat(time) == false {
		return 0 , errors.New("字幕时间格式不正确")
	}
	time = strings.Replace(time , "," , "." , 1)
	timeSplit := strings.Split(time , ":")
	if len(timeSplit) != 3 {
		return 0 , errors.New("字幕时间格式不正确")
	}
	var second float64 = 0
	for k,v := range timeSplit {
		if vn, err := strconv.ParseFloat(v , 64); err != nil {
			return 0 , errors.New("字幕时间解析失败")
		} else {
			if k == 0 {
				second += vn * 3600
			} else if k == 1 {
				second += vn * 60
			} else if k == 2 {
				second += vn
			}
		}
	}
	return second , nil
}


func SubString(str string , begin int ,length int) (substr string) {
	// 将字符串的转换成[]rune
	rs := []rune(str)
	lth := len(rs)

	// 简单的越界判断
	if begin < 0 {
		begin = 0
	}
	if begin >= lth {
		begin = lth
	}
	end := begin + length
	if end > lth {
		end = lth
	}
	// 返回子串
	return string(rs[begin:end])
}