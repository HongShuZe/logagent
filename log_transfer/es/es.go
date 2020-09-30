package es

import (
	"github.com/olivere/elastic"
	"strings"
	"fmt"
	"context"
	"time"
)

type LogData struct {
	Topic string `json:"topic"`
	Data string `json:"data"`
}

var (
	client *elastic.Client
	ch chan *LogData
)

func Init(address string, chanSize, nums int)(err error)  {
	if !strings.HasPrefix(address, "http://") {
		address = "http://" + address
	}
	client, err = elastic.NewClient(elastic.SetURL(address))
	if err != nil {
		return err
	}
	fmt.Println("connect to es success")
	ch = make(chan *LogData, chanSize)
	for i := 0; i<nums; i++ {
		go sendTOES()
	}
	return
}

//sendToES 发送数据到ES
func SendTOCHan(msg *LogData) {
	ch <- msg
}


func sendTOES() {
	for {
		select {
		case msg := <-ch:
			put1, err := client.Index().Index(msg.Topic).Type("XXX").BodyJson(msg).Do(context.Background())
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("Indexed student %v to index %s, type %s \n", put1.Id, put1.Index, put1.Type)
		default:
			time.Sleep(time.Second)
		}
	}
}
