package etcd

import (
	"go.etcd.io/etcd/clientv3"
	"time"
	"fmt"
	"context"
	"encoding/json"
)

var (
	cli *clientv3.Client
)

//需要收集的日志配置信息
type LogEntry struct {
	Path string `json:"path"` //日志存放路径
	Topic string `json:"topic"` //日志要发往kafka中那个topic
}

//初始化ETCD的函数
func Init(addr string, timeout time.Duration) (err error) {
	cli, err = clientv3.New(clientv3.Config{
		Endpoints: []string{addr},
		DialTimeout: timeout,
	})
	if err != nil {
		fmt.Printf("connect to etcd failed %v", err)
		return
	}
	return
}

//从ETCD中根据key获取配置项
func GetConf(key string) (logEntryConf []*LogEntry, err error) {
	//get
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, key)
	cancel()
	if err != nil {
		fmt.Printf("get from etcd failed err :%v\n", err)
		return
	}

	for _, ev := range resp.Kvs{
		err = json.Unmarshal(ev.Value, &logEntryConf)
		if err != nil {
			fmt.Printf("unmarshal etcd value failed, err %v", err)
			return
		}
	}
	return
}

//etcd watch
func WatchConf(key string, newConfCH chan<-[]*LogEntry)  {
	ch := cli.Watch(context.Background(), key)
	for wresp := range ch{
		for _, evt := range wresp.Events{
			fmt.Printf( "Type:%v,key:%v, value:%v",string(evt.Type), string(evt.Kv.Key), string(evt.Kv.Value))
			// 通知taillog.tskMgr
			//1.先判断操作的类型
			var newConf []*LogEntry
			if evt.Type != clientv3.EventTypeDelete {
				// 如果是删除操作，手动传递一个空的配置项
				err := json.Unmarshal(evt.Kv.Value, &newConf)
				if err != nil {
					fmt.Printf("unmarshal failed, err:%v\n", err)
					continue
				}
			}
			fmt.Printf("get new confg:%v\n", newConf)
			newConfCH <- newConf
		}
	}
}




