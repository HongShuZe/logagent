package main

import (
	"logagent/logagent/kafka"
	"fmt"
	"logagent/logagent/taillog"
	"time"
	"logagent/logagent/config"
	"gopkg.in/ini.v1"
)

var (
	cfg = new(config.AppConf)
)

func main()  {
	//加载配置文件
	err := ini.MapTo(cfg, "./config/config.ini")
	if err != nil {
		fmt.Printf("load ini failed, err:%v\n", err)
		return
	}

	//初始化kafka
	err = kafka.Init([]string{cfg.KafkaConf.Address})
	if err != nil {
		fmt.Printf("init kafka failed, err:%v\n", err)
		return
	}
	fmt.Println("kafka init success")

	//打开日志文件准备收集日志
	err = taillog.Init(cfg.TaillogConf.FileName)
	if err != nil {
		fmt.Printf("init taillog failed, err:%v\n", err)
		return
	}
	fmt.Println("taillog init success")

	run()
}

//logagent入口程序
func run() {
	//读取日志
	for {
		select {
			case line := <- taillog.ReadChan():
				//发送到kafka
				kafka.SendToKafka(cfg.KafkaConf.Topic, line.Text)
		default:
			time.Sleep(time.Second)
		}
	}

}