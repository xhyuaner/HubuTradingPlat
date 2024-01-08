package config

type MysqlConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Name     string `mapstructure:"db" json:"db"`
	User     string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
}

type RedisConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type RocketMQConfig struct {
	Host      string `mapstructure:"host" json:"host"`
	Port      int    `mapstructure:"port" json:"port"`
	GroupName string `mapstructure:"group_name" json:"group_name"`
	Topic     string `mapstructure:"topics" json:"topics"`
}

type ServerConfig struct {
	Name         string         `mapstructure:"name" json:"name"`
	Host         string         `mapstructure:"host" json:"host"`
	Tags         []string       `mapstructure:"tags" json:"tags"`
	MysqlInfo    MysqlConfig    `mapstructure:"mysql" json:"mysql"`
	RedisInfo    RedisConfig    `mapstructure:"redis" json:"redis"`
	ConsulInfo   ConsulConfig   `mapstructure:"consul" json:"consul"`
	RocketMQInfo RocketMQConfig `mapstructure:"rocketmq" json:"rocketmq"`
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
