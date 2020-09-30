package main

import (
	"logagent/log_transfer/conf"
	"gopkg.in/ini.v1"
	"fmt"
	"logagent/log_transfer/es"
	"logagent/log_transfer/kafka"
)

func main()  {
	//0.加载配置文件
	var cfg = new(conf.LogTransferCfg)
	err := ini.MapTo(cfg, "./conf/cfg.ini")
	if err != nil {
		fmt.Printf("init config failed err:%v \n", err)
		return
	}
	fmt.Printf("cfg:%v\n", cfg)

	//1.初始化ES
	//1.1 初始化es的client
	//1.2 对外提供一个往es写入数据的函数
	err = es.Init(cfg.ESCfg.Address, cfg.ESCfg.ChanSize,cfg.ESCfg.Nums)
	if err != nil {
		fmt.Printf("init es failed err:%v \n", err)
		return
	}

	//2 初始化kafka
	//2.1 连接kafka， 创建分区的消费者
	//2.2 每个分区的消费者分别取出数据， 通过SendToChan将数据发往ES
	err = kafka.Init([]string{cfg.KafkaCfg.Address}, cfg.KafkaCfg.Topic)
	if err != nil {
		fmt.Printf("init kafka consumer failed err:%v \n", err)
		return
	}

	select {
	}
}
