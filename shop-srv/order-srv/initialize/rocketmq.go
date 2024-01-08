package initialize

import (
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"shop-srv/order-srv/global"
)

func InitMQ() {
	rocketMQInfo := global.ServerConfig.RocketMQInfo
	var err error
	//初始化producer
	global.MQProducerClient, err = rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{fmt.Sprintf("%s:%d", rocketMQInfo.Host, rocketMQInfo.Port)})),
		// TODO:1-添加了producer组名
		producer.WithGroupName("order-general-producer"),
	)
	if err != nil {
		panic("生成MQ producer失败")
	}

	if err = global.MQProducerClient.Start(); err != nil {
		panic("启动MQ producer失败")
	}
}
