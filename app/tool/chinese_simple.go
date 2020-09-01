package tool

import (
	"regexp"
	"strings"
	"unicode/utf8"
)


func CheckChineseNumber(s string) bool {
	regx := regexp.MustCompile("(.*)([一|二|两|三|四|五|六|七|八|九|十|百|千|万|亿]+)(.*)")
	return regx.MatchString(s)
}

func ChineseNumberToLowercaseLength(s string) int {
	st := GetStringUtf8Length(s)
	regx := regexp.MustCompile("([一|二|两|三|四|五|六|七|八|九|十|百|千|万|亿]+)(.*)")
	s = regx.ReplaceAllString(s , "$1")
	if s == "" || !IsChineseNumber(s) {
		return st
	}
	rst := GetStringUtf8Length(s)
	cha_t := 0
	if st > rst {
		cha_t = st - rst
	}
	s = strings.TrimSpace(s)
	zhTexts := strings.Split(s , "")
	zhTextsLens := len(zhTexts)

	numberPosi := true
	maxBaseNumber := 1

	if s == "十" {
		maxBaseNumber = 2
	} else {
		for i:=0; i<zhTextsLens; i++ {
			if numberPosi == false {
				switch zhTexts[i] {
				case "十":
					maxBaseNumber = 2
					break
				case "百":
					maxBaseNumber = 3
					break
				case "千":
					maxBaseNumber = 4
					break
				case "万":
					maxBaseNumber = 5
					break
				}
				break
			}
			numberPosi = !numberPosi
		}
	}
	return maxBaseNumber+cha_t
}


func IsChineseNumber(text string) bool {
	text = strings.TrimSpace(text)
	if text == "" {
		return false
	}

	zhTexts := strings.Split(text , "")
	zhTextsLens := len(zhTexts)

	numberPosi := true
	for i:=0; i<zhTextsLens; i++ {
		if !ValiChineseNumberChar(zhTexts[i] , !numberPosi) {
			return false
		}

		numberPosi = !numberPosi
	}
	return true
}


func ValiChineseNumberChar(s string , unit bool) bool {
	zh_number := []string{"一","二","两","三","四","五","六","七","八","九","十","百","千","万"}
	zh_unit := []string{"十","百","千","万","亿"}

	if unit {
		for _,v := range zh_unit {
			if v == s {
				return true
			}
		}
		return false
	}

	for _,v := range zh_number {
		if v == s {
			return true
		}
	}
	return false
}


//获取utf8字符长度
func GetStringUtf8Length(s string) int {
	return utf8.RuneCountInString(s)
}


//检测文本是否仅符号
func CheckOnlySymbolText(s string) bool {
	if GetStringUtf8Length(s) > 6 {
		return false
	}
	regx := regexp.MustCompile(`^(\\|\{|\}|\[|\]|（|）|\(|\)|\*|/|~|<|>|_|\-|\+|=|&|%|\$|@|#|—|」|「|！|，|。|。|‍|、|？|；|：|‘|’|”|“|"|'|,|\.|\?|;|:|!|\s)+$`)
	if regx.Match([]byte(s)) {
		return true
	} else {
		return false
	}
}