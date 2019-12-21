package aliyun

import (
	"github.com/buger/jsonparser"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)


type AliyunAudioRecognitionResultBlock struct {
	AliyunAudioRecognitionResult
	Blocks []int
}

//阿里云录音录音文件识别 - 智能分段处理
func AliyunAudioResultWordHandle(result [] byte , callback func (vresult *AliyunAudioRecognitionResult)) {
	var audioResult = make(map[int64][] *AliyunAudioRecognitionResultBlock)
	var wordResult = make(map[int64][]*AliyunAudioWord)
	var err error

	//获取录音识别数据集
	_, err = jsonparser.ArrayEach(result, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		text , _ := jsonparser.GetString(value, "Text")
		channelId , _ := jsonparser.GetInt(value, "ChannelId")
		beginTime , _ := jsonparser.GetInt(value, "BeginTime")
		endTime , _ := jsonparser.GetInt(value, "EndTime")
		silenceDuration , _ := jsonparser.GetInt(value, "SilenceDuration")
		speechRate , _ := jsonparser.GetInt(value, "SpeechRate")
		emotionValue , _ := jsonparser.GetInt(value, "EmotionValue")

		vresult := &AliyunAudioRecognitionResultBlock {}
		vresult.Text = text
		vresult.ChannelId = channelId
		vresult.BeginTime = beginTime
		vresult.EndTime = endTime
		vresult.SilenceDuration = silenceDuration
		vresult.SpeechRate = speechRate
		vresult.EmotionValue = emotionValue

		_ , isPresent  := audioResult[channelId]
		if isPresent {
			//追加
			audioResult[channelId] = append(audioResult[channelId] , vresult)
		} else {
			//初始
			audioResult[channelId] = []*AliyunAudioRecognitionResultBlock{}
			audioResult[channelId] = append(audioResult[channelId] , vresult)
		}
	} , "Result", "Sentences")
	if err != nil {
		panic(err)
	}

	//获取词语数据集
	_, err = jsonparser.ArrayEach(result , func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		word , _ := jsonparser.GetString(value, "Word")
		channelId , _ := jsonparser.GetInt(value, "ChannelId")
		beginTime , _ := jsonparser.GetInt(value, "BeginTime")
		endTime , _ := jsonparser.GetInt(value, "EndTime")
		vresult := &AliyunAudioWord{
			Word:word,
			ChannelId:channelId,
			BeginTime:beginTime,
			EndTime:endTime,
		}
		_, isPresent  := wordResult[channelId]
		if isPresent {
			//追加
			wordResult[channelId] = append(wordResult[channelId] , vresult)
		} else {
			//初始
			wordResult[channelId] = []*AliyunAudioWord{}
			wordResult[channelId] = append(wordResult[channelId] , vresult)
		}
	} , "Result" , "Words")
	if err != nil {
		panic(err)
	}


	var symbol = []string{"？","。","，","！","；","?",".",",","!"}
	//数据集处理
	for _ , value := range audioResult {
		for _ , data := range value {
			data.Blocks = GetTextBlock(data.Text)
			data.Text = ReplaceStrs(data.Text , symbol , "")
		}
	}

	//遍历输出
	for _,value := range wordResult {

		var block string = ""
		var blockRune int = 0
		var lastBlock int = 0

		var beginTime int64 = 0
		var blockBool = false

		var ischinese = IsChineseWords(value) //校验中文

		for i , word := range value {
			if blockBool || i == 0 {
				beginTime = word.BeginTime
				blockBool = false
			}

			if ischinese {
				block += word.Word
			} else {
				block += CompleSpace(word.Word) //补全空格
			}
			blockRune = utf8.RuneCountInString(block)

			for channel , p := range audioResult {
				if word.ChannelId != channel {
					continue
				}
				for windex , w := range p {

					if (word.BeginTime >= w.BeginTime && word.EndTime <= w.EndTime) || ((word.BeginTime < w.EndTime && word.EndTime > w.EndTime) && (FindSliceIntCount(w.Blocks , -1) != len(w.Blocks))) {
						flag := false
						early := false

						for t , B := range w.Blocks{
							if (blockRune >= B) && B != -1 {
								flag = true

								//fmt.Println(  block )
								//fmt.Println(  w.Text )
								//fmt.Println(  w.Blocks )
								//fmt.Println(B , lastBlock , (B - lastBlock) , word.Word)
								//fmt.Println(w.Text)

								var thisText = ""
								//容错机制
								if t == (len(w.Blocks) - 1) {
									thisText = SubString(w.Text , lastBlock , 10000)
								} else {
									//下个词提前结束
									if i < len(value)-1 && value[i+1].BeginTime >= w.EndTime{
										thisText = SubString(w.Text , lastBlock , 10000)
										early = true
									} else {
										thisText = SubString(w.Text , lastBlock , (B - lastBlock))
									}
								}

								lastBlock = B
								if early == true {
									//全部设置为-1
									for vt,vb := range w.Blocks{
										if vb != -1 {
											w.Blocks[vt] = -1;
										}
									}
								} else {
									w.Blocks[t] = -1
								}

								vresult := &AliyunAudioRecognitionResult{
									Text:thisText,
									ChannelId:channel,
									BeginTime:beginTime,
									EndTime:word.EndTime,
									SilenceDuration:w.SilenceDuration,
									SpeechRate:w.SpeechRate,
									EmotionValue:w.EmotionValue,
								}
								callback(vresult) //回调传参

								blockBool = true
								break
							}
						}

						//fmt.Println("word.Word:" , word.Word)
						//fmt.Println(block)

						if FindSliceIntCount(w.Blocks , -1) == len(w.Blocks) {
							//全部截取完成
							block = ""
							lastBlock = 0
						}

						//容错机制
						if FindSliceIntCount(w.Blocks , -1) == (len(w.Blocks)-1) && flag == false {
							var thisText = SubString(w.Text , lastBlock , 10000)

							w.Blocks[len(w.Blocks) - 1] = -1
							//vresult
							vresult := &AliyunAudioRecognitionResult{
								Text:thisText,
								ChannelId:channel,
								BeginTime:beginTime,
								EndTime:w.EndTime,
								SilenceDuration:w.SilenceDuration,
								SpeechRate:w.SpeechRate,
								EmotionValue:w.EmotionValue,
							}

							//fmt.Println(  thisText )
							//fmt.Println(  block )
							//fmt.Println(  word.Word , beginTime, w.EndTime , flag  , word.EndTime  )

							callback(vresult) //回调传参

							//覆盖下一段落的时间戳
							if windex < (len(p)-1) {
								beginTime = p[windex+1].BeginTime
							} else {
								beginTime = w.EndTime
							}

							//清除参数
							block = ""
							lastBlock = 0
						}
					}
				}
			}
		}
	}
}



func FindSliceIntCount(slice []int , target int) int {
	c := 0
	for _ , v := range slice {
		if target == v {
			c++
		}
	}
	return c
}


//批量替换多个关键词文本
func ReplaceStrs(strs string , olds []string , s string) string {
	for _ , word := range olds {
		strs = strings.Replace(strs , word , s , -1)
	}
	return strs
}

func StringIndex(strs string , word rune) int {
	strsRune := []rune(strs)
	for i,v := range strsRune {
		if v == word {
			return i
		}
	}
	return -1
}

//补全右边空格
func CompleSpace(s string) string {
	s = strings.TrimLeft(s , " ");
	s = strings.TrimRight(s , " ");
	return s + " ";
}

func IsChineseWords(words []*AliyunAudioWord) bool {
	for _,v := range words {
		if (IsChineseChar(v.Word)){
			return true
		}
	}
	return false
}

func IsChineseChar(str string) bool {
	for _, r := range str {
		if unicode.Is(unicode.Scripts["Han"], r) || (regexp.MustCompile("[\u3002\uff1b\uff0c\uff1a\u201c\u201d\uff08\uff09\u3001\uff1f\u300a\u300b]").MatchString(string(r))) {
			return true
		}
	}
	return false
}

func IndexRunes(strs string , olds []rune) int  {
	min := -1
	for i , word := range olds {
		index := StringIndex(strs , word)
		//println( "ts : " , index)
		if i == 0 {
			min = index
		} else {
			if min == -1 {
				min = index
			} else {
				if index < min && index != -1 {
					min = index
				}
			}
		}
	}
	return min
}

func GetTextBlock(strs string) ([]int) {
	var symbol_zhcn = []rune{'？','。','，','！','；','?','.',',','!'}
	//var symbol_en = []rune{'?','.',',','!'}
	strsRune := []rune(strs)

	blocks := []int{}
	for {
		index := IndexRunes(strs , symbol_zhcn)
		if index == -1 {
			break
		}
		strs = string(strsRune[0:index]) + string(strsRune[(index + 1):])
		strsRune = []rune(strs)
		blocks = append(blocks , index)
	}
	return blocks
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
