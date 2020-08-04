package app

import (
	"github.com/lxn/walk"
)

type TaskLog struct {
	TextEdit *walk.TextEdit
}

func NewTasklog(textEdit *walk.TextEdit) *TaskLog {
	task := new(TaskLog)
	task.TextEdit = textEdit
	return task
}

func (t *TaskLog) SetTextEdit(textEdit *walk.TextEdit)  {
	t.TextEdit = textEdit
}

func (t *TaskLog) AppendLogText(s string)  {
	t.TextEdit.AppendText(s + "\r\n")
}

func (t *TaskLog) ClearLogText()  {
	t.TextEdit.SetText("")
}