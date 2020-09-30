package main

import (
	"go.etcd.io/etcd/clientv3"
	"time"
	"fmt"
	"context"
)

func main()  {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Printf("connect to etcd failed , err:%v", err)
		return
	}
	fmt.Println("connect to eetcd success")
	defer cli.Close()

	//watch
	//派一个哨兵一直监视着 hsz这个key的变化(新增、修改、删除)
	ch := cli.Watch(context.Background(), "hsz")
	for wresp := range ch{
		for _, evt := range wresp.Events{
			fmt.Printf("Type: %v, key: %v, value:%v \n", evt.Type, string(evt.Kv.Key), string(evt.Kv.Value))
		}

	}
}
