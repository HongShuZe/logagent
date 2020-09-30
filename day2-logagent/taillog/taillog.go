package taillog

import (
	"fmt"
	"github.com/hpcloud/tail"
	"context"
	"logagent/day2-logagent/kafka"
)

//专门从日志文件收集日志模块

//TailTask： 一个日志收集的任务
type TailTask struct {
	path string
	topic string
	instance *tail.Tail
	//为了能实现退出t.run()
	ctx context.Context
	cancelFunc context.CancelFunc
}

/*var (
	tailObj *tail.Tail
	LogChan chan string
)*/

func NewTailTask(path, topic string) (tailObj *TailTask) {
	ctx, cancel := context.WithCancel(context.Background())
	tailObj = &TailTask{
		path:path,
		topic:topic,
		ctx:ctx,
		cancelFunc:cancel,
	}
	tailObj.init() //根据路径去打开对应的日志
	return
}

func (t *TailTask)init() {
	config := tail.Config{
		ReOpen:		true, 				//重新打开
		Follow:		true,				//是否跟随
		Location:	&tail.SeekInfo{Offset:0, Whence:2},	//从文件那个地方开始读
		MustExist:	false,				//文件不存在报错
		Poll:		true,
	}
	var err error
	t.instance, err = tail.TailFile(t.path, config)
	if err != nil {
		fmt.Println("tail file failed, err:", err)
	}
	//当goroutine执行的函数退出时，goroutine就结束了
	go t.run()
}

func (t *TailTask) run()  {
	for {
		select {
		case <-t.ctx.Done():
			fmt.Printf("tail task:%s_%s 结束了", t.path, t.topic)
			return
		case line := <-t.instance.Lines: //从tailObj的通道中一行一行的读取日志数据
			//3.2 发往kafka
			//先把一个`日志数据发送到一个通道中
			fmt.Printf("get log data from %s success, log:%v\n",t.path, line.Text )
			kafka.SendToKafka(t.topic, line.Text)
			//kafka那个包中有单独的goroutine去取日志数据发送到kafka
		}
	}
}