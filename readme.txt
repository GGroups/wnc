
本程序是简化的命令行netcat程序，方便windows下测试websocket
需求：
1.windows下没有好用的netcat
2.通过此案例熟悉golang的socket编程，以及协程。

#wnc 的用法

##服务器模式：
    wnc.exe -l int

##客户端链接模式：
    wnc.exe ip port

##额外参数：
  -v   show ips listening


#编译命令
go build -o wnc.exe main.go