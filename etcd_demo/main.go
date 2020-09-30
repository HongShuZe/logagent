package main

import (
	"go.etcd.io/etcd/clientv3"
	"time"
	"fmt"
	"context"
)

func main()  {
	cli,err := clientv3.New(clientv3.Config{
		Endpoints: []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Printf("connect to etcd failed , err:%v\n", err)
		return
	}
	fmt.Println("connect to etcd success")
	defer cli.Close()

	//put
	ctx, cancle := context.WithTimeout(context.Background(), time.Second)
	_, err = cli.Put(ctx, "hsz", "hszz")
	cancle()
	if err != nil {
		fmt.Printf("put to etcd failed err :%v\n", err)
		return
	}

	//get
	ctx, cancle = context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, "hsz")
	cancle()
	if err != nil {
		fmt.Printf("get from etcd failed err :%v\n", err)
		return
	}

	for _, ev := range resp.Kvs{
		fmt.Printf("%s:%s\n", ev.Key, ev.Value)
	}
}
