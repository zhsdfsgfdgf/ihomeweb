# go-micro中实现限流


**在Micro中使用限流功能非常简单，只需增加一行代码就可以实现基于QPS的限流,以 uber rate limiter 插件为例**

   

 
    import (
    ...
       limiter "github.com/micro/go-plugins/wrapper/ratelimiter/uber"
    ...
    )
    
    func main() {
       const QPS = 100
       // New Service
       service := micro.NewService(
          micro.Name("com.foo.srv.hello"),
          micro.Version("latest"),
          micro.WrapHandler(limiter.NewHandlerWrapper(QPS)),
       )
    
    ...}

以上代码便为hello-srv增加了服务器限流能力， QPS上限为100。这个限制由此服务的所有handler所有method 共享。换句话说，此限制的作用域是服务级别的。