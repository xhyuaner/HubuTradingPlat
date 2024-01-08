package main

import (
	"context"
	"fmt"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

func main() {
	p, err := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{"192.168.88.105:9876"})),
	)
	if err != nil {
		panic("生成producer失败")
	}

	if err = p.Start(); err != nil {
		panic("启动producer失败")
	}

	msg := primitive.NewMessage("hellomp", []byte("this is delay message"))
	msg.WithDelayTimeLevel(3)
	res, err := p.SendSync(context.Background(), msg)
	if err != nil {
		fmt.Printf("发送失败: %s\n", err)
	} else {
		fmt.Printf("发送成功: %s\n", res.String())
	}

	if err = p.Shutdown(); err != nil {
		panic("关闭producer失败")
	}

}
