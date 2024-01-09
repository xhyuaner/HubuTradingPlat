### HubuShop

---

一个使用 Go 语言编写，基于微服务和 Gin + gRPC 实现的分布式在线交易平台，支持配置中心、服务注册与发现、消息队列、分布式搜索等。

### 相关技术

---

- Web 框架：使用 Go 语言编写且应用最广泛的高性能 Web 框架——Gin
- RPC 通信：使用 Go 语言中用的最多、生态最好的 RPC 框架——gRPC
- 配置中心：可以对各个微服务的配置进行集中统一管理，支持配置的热更新，可用于构建云原生应用程序
- 服务注册/发现：可以对 service 层和 web 层的所有微服务进行统一管理，web 层可通过服务名称自动查找所依赖的 service 层服务并建立连接
- 消息队列：通过消息队列，可做到数据的异步处理，可大幅提高系统吞吐量，应对各种高并发场景
- 分布式搜索：使用分布式搜索引擎，大幅提高全文搜索的效率

### 技术难点

---

- 分布式锁
- CAP理论
- 最终一致性

- 分布式事务
- 接口幂等性
- 链路追踪
- 熔断限流

### 启动

---

1. 安装相关工具（推荐使用 Docker ），包括 MySQL、Redis、Consul、Nacos、RocketMQ、Elasticsearch、Jaeger

   ```shell
   # 1.启动MySQL
   systemctl start mysqld
   # 2.启动Mongodb
   mongod -f /home/root/apps/mongodb/mongodb.conf     
   # 3.启动Yapi
   pm2 start /home/root/apps/yapi/vendors/server/app.js --name yapi # pm2管理yapi服务  
   # url：http://192.168.88.105:3000/，账号：admin@admin.com，密码：123456
   # 4.启动docker，相关容器会自动启动，若某容器未启动，则手动启动
   systemctl start docker
   # 4.1 启动redis
   docker run -p 6379:6379 -d redis:latest redis-server
   docker update --restart=always <容器名或容器ID> 		# 设置容器自启动
   # 4.2 启动consul
   docker run -d -p 8500:8500 -p 8300:8300 -p 8301:8301 -p 8302:8302 -p 8600:8600/udp consul consul agent -dev -client=0.0.0.0   
   # 4.3 启动nacos
   docker run --name nacos-standalone -e MODE=standalone -e JVM_XMS=512m -e JVM_XMX=512m -e JVM_XMN=256m -p 8848:8848 -d nacos/nacos-server:latest
   # 4.4 启动elasticsearch
   docker run --name elasticsearch -p 9200:9200 -p 9300:9300 -e "discovery.type=single-node" -e ES_JAVA_OPTS="-Xms256m -Xmx512m" -v /data/elasticsearch/config/elasticsearch.yml:/usr/share/elasticsearch/config/elasticsearch.yml -v /data/elasticsearch/data:/usr/share/elasticsearch/data -v /data/elasticsearch/plugins:/usr/share/elasticsearch/plugins -d elasticsearch:7.10.1
   # 4.5 启动jaeger
   docker run --rm --name jaeger -p6831:6831/udp -p16686:16686 jaegertracing/all-in-one:latest
   # 5.启动rocketmq
   cd /home/root/apps/rocketmq/install && docker-compose up
   # 6.启动kong
   kong start -c /etc/kong/kong.conf
   ```

2. 启动后端 Service 层服务和 Web 层服务

3. 前端通过 `webpack.config.js` 文件 `devServer` 修改启动端口或者通过 `config/index.js` 文件 `port` 修改端口

4. 使用 `npm install` 和 `npm run dev` 命令启动前端

### 项目介绍

---

- 基于 JWT 做访问鉴权 token，Gin 做路由分发、表单验证、解决跨域等。
- 登录/注册功能：采用 service 和 web 双层架构、使用 viper 包做配置解析、web 层基于 Gin 做路由转发、使用 Redis 实现注册验证码缓存服务、使用 base64 生成验证码图片做登录验证、service  层使用 MD5 盐值加密，保证密码在数据库中加密存储。
- 商品服务功能：基于 Elasticsearch 实现商品搜索；完成如下接口：1.商品相关、2.商品品牌相关、3.商品分类类目相关、4.商品分类相关、5.商品主页轮播图相关。
- 图片文件使用 Aliyun 对象存储，使用服务端签名直传方式传输文件。
- 库存服务：库存服务的核心在于保持数据的一致性，可用性，高性能，解决在分布式高并发场景下，如何保证数据一致性，库存服务引入了 Redis 分布式锁和 RocketMQ，来实现分布式高并发场景下的数据一致性，实现扣减库存，库存超时归还，解决重复归还商品问题，保证了接口幂等性。
- 订单服务：基于 gRPC 实现订单相关服务及购物车相关服务等各类接口，使用本地 MySQL 事务保证本地数据一致性，使用 RocketMQ 从订单服务发送消息到商品服务以及库存服务（跨服务），进行商品查询和库存扣减，实现跨微服务调用，保证信息一致性。
- 用户操作接口服务： 为用户提供操作接口，实现了用户的地址，留言, 收藏管理等功能。
- 基于 Jaeger 做微服务间链路追踪，使用 Sentinel 实现服务的熔断限流。

### 优点

---

- 代码前后端分离
- 各个服务独立部署
- 更具灵活性和扩展性
- 便于单个服务的独立性能优化
- 便于自动部署和持续集成



