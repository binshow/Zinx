# Zinx
go实现tcp服务器框架:https://www.yuque.com/aceld/npyr8s/bgftov

1. 简单的一个 server 的实现
2. 链接 connection 的实现，保证了 tcp 原生的 conn
3. 简单的路由实现，
4. 封装request
5. 全局配置文件的加载
6. 消息封装，解决TCP的粘包问题
7. 读写分离：一个专门负责从客户端读取数据，一个专门负责向客户端写数据
