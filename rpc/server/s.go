package main

import (
	"fmt"
	"net"
)

/*
这个过程中最重要的就是序列化和反序列化了，因为数据传输的数据包必须是二进制的，你直接丢一个Java对象过去，人家可不认识，你必须把Java对象序列化为二进制格式，传给Server端，Server端接收到之后，再反序列化为Java对象。

要实现一个RPC不算难，难的是实现一个高性能高可靠的RPC框架。

比如，既然是分布式了，那么一个服务可能有多个实例，你在调用时，要如何获取这些实例的地址呢？

这时候就需要一个服务注册中心，比如在Dubbo里头，就可以使用Zookeeper作为注册中心，在调用时，从Zookeeper获取服务的实例列表，再从中选择一个进行调用。

那么选哪个调用好呢？这时候就需要负载均衡了，于是你又得考虑如何实现复杂均衡，比如Dubbo就提供了好几种负载均衡策略。

这还没完，总不能每次调用时都去注册中心查询实例列表吧，这样效率多低呀，于是又有了缓存，有了缓存，就要考虑缓存的更新问题，blablabla......

你以为就这样结束了，没呢，还有这些：

客户端总不能每次调用完都干等着服务端返回数据吧，于是就要支持异步调用；
服务端的接口修改了，老的接口还有人在用，怎么办？总不能让他们都改了吧？这就需要版本控制了；
服务端总不能每次接到请求都马上启动一个线程去处理吧？于是就需要线程池；
服务端关闭时，还没处理完的请求怎么办？是直接结束呢，还是等全部请求处理完再关闭呢？
......
如此种种，都是一个优秀的RPC框架需要考虑的问题。
*/

func cal(a int, b int) int {
	return a + b
}

// cal(1,3)

func main() {
	fmt.Println("server")
	ip := net.IP{127, 0, 0, 1}
	add := net.TCPAddr{
		IP:   ip,
		Port: 3333,
		Zone: "0",
	}
	b, _ := net.ListenTCP("tcp4", &add)
	fmt.Println("tcp listened", b)
	defer b.Close()
	//b.Accept()
	for {
		con, _ := b.Accept()
		bf := make([]byte, 2048)
		x, _ := con.Read(bf)
		fmt.Println("tcp received", string(bf[:x]))
		con.Write([]byte("pong"))
	}
}
