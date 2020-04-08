package app

type TranslateMultitask struct {
	MaxConcurrencyNumber int //最大运行并发数
	Total int //任务总数
	QueueFile []string //任务队列
	CurrentIndex int //已处理的下标
	FinishNumber int //已完成的任务数量

	SrtTranslateApp *SrtTranslateApp
}

func NewTranslateMultitask(concurrencyNumber int) (*TranslateMultitask) {
	app := new(TranslateMultitask)
	if concurrencyNumber == 0 {
		//默认并发数
		app.MaxConcurrencyNumber = 2
	} else {
		app.MaxConcurrencyNumber = concurrencyNumber
	}
	return app
}

//设置 任务队列
func (task *TranslateMultitask) SetQueueFile(queue []string)  {
	task.QueueFile = queue
	task.Total = len(queue)
}

//设置 任务队列
func (task *TranslateMultitask) SetSrtTranslateApp(app *SrtTranslateApp)  {
	task.SrtTranslateApp = app
}

//设置 并发数量
func (task *TranslateMultitask) SetMaxConcurrencyNumber(n int)  {
	task.MaxConcurrencyNumber = n
}

//并发运行
func (task *TranslateMultitask) Run() {
	//初始参数
	task.CurrentIndex = -1
	task.FinishNumber = 0

	number := 1
	//并发调用
	for number <= task.MaxConcurrencyNumber && task.CurrentIndex < (task.Total - 1){
		if task.CurrentIndex == -1 {
			task.CurrentIndex = 0;
			path := task.QueueFile[task.CurrentIndex]
			go func() {
				task.SrtTranslateApp.Run(path)
			}()
		} else {
			task.CurrentIndex++
			path := task.QueueFile[task.CurrentIndex]
			go func() {
				task.SrtTranslateApp.Run(path)
			}()
		}
		number++
	}
}

func (task *TranslateMultitask) RunOver() bool {
	if task.CurrentIndex >= (task.Total - 1) {
		//任务队列处理完成
		return true
	}
	//执行
	task.CurrentIndex++
	path := task.QueueFile[task.CurrentIndex]
	go func() {
		task.SrtTranslateApp.Run(path)
	}()
	return false
}


//标记已完成
func (task *TranslateMultitask) FinishTask() bool {
	task.FinishNumber++

	if (task.FinishNumber == task.Total) {
		return true
	}
	return false
}