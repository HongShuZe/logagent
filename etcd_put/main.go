package main

import (
	"go.etcd.io/etcd/clientv3"
	"time"
	"fmt"
	"context"
)

//etcd client put/get demo
func main()  {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"127.0.0.1:2379"},
		DialTimeout: 5 *time.Second,
	})
	if err != nil {
		// handle error!
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return
	}
	fmt.Println("connect to etcd success")
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	value := `[{"path":"/home/zwx/go/src/logagent/day2-logagent/nginx.log","topic":"web_log"},{"path":"/home/zwx/go/src/logagent/day2-logagent/redis.log","topic":"redis_log"}]`
	_, err = cli.Put(ctx, "/logagent/192.168.20.143/collect_config", value)
	cancel()
	if err != nil {
		// handle error!
		fmt.Printf("put to etcd failed, err:%v\n", err)
		return
	}
}