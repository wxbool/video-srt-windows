package app

import "strings"

type TaskLog struct {
	content []string
}

func NewTasklog() *TaskLog {
	task := new(TaskLog)
	return task
}

func (t *TaskLog) AppendLogText(s string)  {
	t.content = append(t.content , s)
}

func (t *TaskLog) ClearLogText()  {
	t.content = []string{}
}

func (t *TaskLog) GetString() string {
	return strings.Join(t.content , "\r\n")
}
