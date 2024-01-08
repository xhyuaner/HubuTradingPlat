package config

type OtherSrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
}

type MysqlConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Name     string `mapstructure:"db" json:"db"`
	User     string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type JaegerConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	Name string `mapstructure:"name" json:"name"`
}

type RocketMQConfig struct {
	Host           string   `mapstructure:"host" json:"host"`
	Port           int      `mapstructure:"port" json:"port"`
	GroupName      string   `mapstructure:"group_name" json:"group_name"`
	Topics         []string `mapstructure:"topics" json:"topics"`
	DelayTimeLevel int      `mapstructure:"time_level" json:"time_level"`
}

type ServerConfig struct {
	Name         string         `mapstructure:"name" json:"name"`
	Host         string         `mapstructure:"host" json:"host"`
	Tags         []string       `mapstructure:"tags" json:"tags"`
	MysqlInfo    MysqlConfig    `mapstructure:"mysql" json:"mysql"`
	ConsulInfo   ConsulConfig   `mapstructure:"consul" json:"consul"`
	RocketMQInfo RocketMQConfig `mapstructure:"rocketmq" json:"rocketmq"`
	JaegerInfo   JaegerConfig   `mapstructure:"jaeger" json:"jaeger"`

	//商品微服务的配置
	GoodsSrvInfo OtherSrvConfig `mapstructure:"goods_srv" json:"goods_srv"`
	//库存微服务的配置
	InventorySrvInfo OtherSrvConfig `mapstructure:"inventory_srv" json:"inventory_srv"`
}

type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      uint64 `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	DataId    string `mapstructure:"dataid"`
	Group     string `mapstructure:"group"`
}
