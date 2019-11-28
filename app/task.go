package app

type VideoMultitask struct {
	MaxConcurrencyNumber int //最大运行并发数
	Total int //任务总数
	QueueFile []string //任务队列
	CurrentIndex int //已处理的下标
	FinishNumber int //已完成的任务数量

	VideoSrt *VideoSrt
}

func NewVideoMultitask(concurrencyNumber int) (*VideoMultitask) {
	app := new(VideoMultitask)
	if concurrencyNumber == 0 {
		//默认并发数
		app.MaxConcurrencyNumber = 2
	} else {
		app.MaxConcurrencyNumber = concurrencyNumber
	}
	return app
}

//设置 任务队列
func (task *VideoMultitask) SetQueueFile(queue []string)  {
	task.QueueFile = queue
	task.Total = len(queue)
}

//设置 任务队列
func (task *VideoMultitask) SetVideoSrt(v *VideoSrt)  {
	task.VideoSrt = v
}

//设置 并发数量
func (task *VideoMultitask) SetMaxConcurrencyNumber(n int)  {
	task.MaxConcurrencyNumber = n
}

//并发运行
func (task *VideoMultitask) Run() {
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
				task.VideoSrt.Run(path)
			}()
		} else {
			task.CurrentIndex++
			path := task.QueueFile[task.CurrentIndex]
			go func() {
				task.VideoSrt.Run(path)
			}()
		}
		number++
	}
}

func (task *VideoMultitask) RunOver() bool {
	//fmt.Println("RunOver：" , task.CurrentIndex)

	if task.CurrentIndex >= (task.Total - 1) {
		//任务队列处理完成
		return true
	}
	//执行
	task.CurrentIndex++
	path := task.QueueFile[task.CurrentIndex]
	go func() {
		task.VideoSrt.Run(path)
	}()
	return false
}


//标记已完成
func (task *VideoMultitask) FinishTask() bool {
	task.FinishNumber++
	//fmt.Println("FinishTask：" , task.FinishNumber)

	if (task.FinishNumber == task.Total) {
		return true
	}
	return false
}