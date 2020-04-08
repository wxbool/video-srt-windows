package parse

import (
	"bufio"
	"errors"
	"os"
	"strings"
	"videosrt/app/tool"
)

//允许解析的文件类型
const (
	PARSE_FILE_SRT = 1
)

const (
	PARSE_SRT_NUMBER = 1
	PARSE_SRT_TIME_RANGE = 2
	PARSE_SRT_SUBTITLE_TEXT = 3
	PARSE_SRT_TRIM = 4
)

//字幕行结构
type SrtRows struct {
	Id int //字幕自然行id
	Number string //字幕序号
	TimeStart string //字幕开始时间戳
	TimeStartSecond float64 //字幕开始（秒）
	TimeEnd string //字幕结束时间戳
	TimeEndSecond float64 //字幕结束（秒）
	Text []string //字幕文本
}

//字幕文件属性
type Srt struct {
	File string
	FileType int
	AppointBilingualRows int //指定使用N行字幕文本[双语字幕参数]
}

//字幕解析器
type SubtitleParse struct {
	Srt *Srt
	Rows []*SrtRows
}

//获取解析字幕实例
func NewSubtitleParse(srt *Srt) (*SubtitleParse) {
	parse := new(SubtitleParse)
	if srt.FileType == 0 {
		srt.FileType = PARSE_FILE_SRT //默认文件类型
	}

	parse.Srt = srt
	return parse
}

//解析文本
func (parse *SubtitleParse) Parse() error {
	srtfile := parse.Srt.File
	//校验文件
	if tool.VaildFile(srtfile) == false {
		return errors.New("字幕文件不存在")
	}

	if file, err := os.Open(srtfile); err!=nil {
		return errors.New("打开字幕文件失败")
	} else {

		lineRows := 0
		lineStart := false
		lineRowSrt := new(SrtRows)
		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			lineText := scanner.Text()

			//逐行读取
			if strResult, rowsType , err := ParseSrtRows(lineText , lineStart); err != nil {
				return err
			} else {

				//分支处理
				switch rowsType {
				case PARSE_SRT_NUMBER:
					lineRowSrt = new(SrtRows)
					lineRowSrt.Id = lineRows
					lineRowSrt.Number = strResult

					lineStart = true
					break
				case PARSE_SRT_TIME_RANGE:
					srtTimes := strings.Split(strResult , "-")
					lineRowSrt.TimeStart = srtTimes[0]
					lineRowSrt.TimeEnd = srtTimes[1]
					break
				case PARSE_SRT_SUBTITLE_TEXT:
					lineRowSrt.Text = append(lineRowSrt.Text , strResult)
					break
				case PARSE_SRT_TRIM:
					//结束上一段字幕输入
					if lineStart == true {
						if lineRowSrt.TimeStart != "" && lineRowSrt.TimeEnd != "" && len(lineRowSrt.Text) > 0 {
							if fsecond, err := SrtTimeFormatToSecond(lineRowSrt.TimeStart);err != nil {
								return err
							} else {
								lineRowSrt.TimeStartSecond = fsecond
							}
							if fsecond, err := SrtTimeFormatToSecond(lineRowSrt.TimeEnd);err != nil {
								return err
							} else {
								lineRowSrt.TimeEndSecond = fsecond
							}
							lineRows++;
							//追加行
							parse.Rows = append(parse.Rows , lineRowSrt)
						}
						lineStart = false
					}
					break
				}
			}

		}

		//未结束的最后一行
		if lineStart == true {
			if lineRowSrt.TimeStart != "" && lineRowSrt.TimeEnd != "" && len(lineRowSrt.Text) > 0 {
				if fsecond, err := SrtTimeFormatToSecond(lineRowSrt.TimeStart);err != nil {
					return err
				} else {
					lineRowSrt.TimeStartSecond = fsecond
				}
				if fsecond, err := SrtTimeFormatToSecond(lineRowSrt.TimeEnd);err != nil {
					return err
				} else {
					lineRowSrt.TimeEndSecond = fsecond
				}
				lineRows++;
				//追加行
				parse.Rows = append(parse.Rows , lineRowSrt)
			}
			lineStart = false
		}

	}
	return nil
}
