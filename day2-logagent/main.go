package main

import (
	"logagent/day2-logagent/kafka"
	"fmt"
	"logagent/day2-logagent/taillog"
	"time"
	"logagent/day2-logagent/config"
	"gopkg.in/ini.v1"
	"logagent/day2-logagent/etcd"
	"sync"
	"logagent/day2-logagent/utils"
)

var (
	cfg = new(config.AppConf)
)

func main(){
	// 0. 加载配置文件
	err := ini.MapTo(cfg, "./config/config.ini")
	if err != nil {
		fmt.Printf("load ini failed, err:%v\n", err)
		return
	}
	// 1. 初始化kafka连接
	err = kafka.Init([]string{cfg.KafkaConf.Address}, cfg.KafkaConf.ChanMaxSize)
	if err != nil {
		fmt.Printf("init Kafka failed,err:%v\n", err)
		return
	}
	fmt.Println("init kafka success.")
	// 2. 初始化ETCD
	// 5 * time.Second
	err = etcd.Init(cfg.EtcdConf.Address, time.Duration(cfg.EtcdConf.Timeout)*time.Second)
	if err != nil {
		fmt.Printf("init etcd failed,err:%v\n", err)
		return
	}
	fmt.Println("init etcd success.")
	//为了实现每个logagent都拉取自己独有的配置， 所以要以自己的IP地址作为分区
	ipstr, err := utils.GetOutboundIP()
	if err != nil {
		panic(err)
	}
	etcdConfKey := fmt.Sprintf(cfg.EtcdConf.Key, ipstr)
	// 2.1 从etcd中获取日志收集项的配置信息
	logEntryConf, err := etcd.GetConf(etcdConfKey)
	if err != nil {
		fmt.Printf("etcd.GetConf failed,err:%v\n", err)
		return
	}
	fmt.Printf("get conf from etcd success, %v\n", logEntryConf)
	// 2.2 派一个哨兵去监视日志收集项的变化（有变化及时通知我的logAgent实现热加载配置）

	for index, value := range logEntryConf{
		fmt.Printf("index:%v value:%v\n", index, value)
	}
	// 3. 收集日志发往Kafka
	taillog.Init(logEntryConf)
	// 因为NewConfChan访问了tskMgr的newConfChan, 这个channel是在taillog.Init(logEntryConf) 执行的初始化
	newConfChan := taillog.NewConfChan() // 从taillog包中获取对外暴露的通道
	var wg sync.WaitGroup
	wg.Add(1)
	go etcd.WatchConf(etcdConfKey, newConfChan) // 哨兵发现最新的配置信息会通知上面的那个通道
	wg.Wait()
	// 3. 具体的业务
	//run()

}


/*//logagent入口程序
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

}*/