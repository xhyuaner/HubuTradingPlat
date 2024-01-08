package config

type OtherSrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"key" json:"key"`
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

type AlipayConfig struct {
	AppID        string `mapstructure:"app_id" json:"app_id"`
	PrivateKey   string `mapstructure:"private_key" json:"private_key"`
	AliPublicKey string `mapstructure:"ali_public_key" json:"ali_public_key"`
	NotifyURL    string `mapstructure:"notify_url" json:"notify_url"`
	ReturnURL    string `mapstructure:"return_url" json:"return_url"`
}

type ServerConfig struct {
	Name             string         `mapstructure:"name" json:"name"`
	Host             string         `mapstructure:"host" json:"host"`
	Port             int            `mapstructure:"port" json:"port"`
	Tags             []string       `mapstructure:"tags" json:"tags"`
	GoodsSrvInfo     OtherSrvConfig `mapstructure:"goods_srv" json:"goods_srv"`
	OrderSrvInfo     OtherSrvConfig `mapstructure:"order_srv" json:"order_srv"`
	InventorySrvInfo OtherSrvConfig `mapstructure:"inventory_srv" json:"inventory_srv"`
	JWTInfo          JWTConfig      `mapstructure:"jwt" json:"jwt"`
	ConsulInfo       ConsulConfig   `mapstructure:"consul" json:"consul"`
	AliPayInfo       AlipayConfig   `mapstructure:"alipay" json:"alipay"`
	JaegerInfo       JaegerConfig   `mapstructure:"jaeger" json:"jaeger"`
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
