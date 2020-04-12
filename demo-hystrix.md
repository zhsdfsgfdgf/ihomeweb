# go-micro中实现熔断


**在创建service实例的时候指定hystrix插件，系统便具备了自动熔断能力，所有从此节点发出的Micro服务调用都会受到熔断插件的限制和保护。 当请求超时或并发数超限，调用方会立即接收到熔断错误**



    import (
       hystrixGo "github.com/afex/hystrix-go/hystrix"
      "github.com/micro/go-plugins/wrapper/breaker/hystrix"
    )
    func main(){
    ...
      service := micro.NewService(
         micro.Name("com.foo.breaker.example"),
         micro.WrapClient(hystrix.NewClientWrapper()),
      )
      service.Init()
      //默认的超时时间是1000毫秒， 默认最大并发数是10,我们也可以自行修改
      hystrix.DefaultMaxConcurrent = 3
      hystrix.DefaultTimeout = 200 
    ...
    }

