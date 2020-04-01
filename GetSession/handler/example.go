package handler

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/micro/go-log"

	example "sss/GetSession/proto/example"
	"sss/IhomeWeb/utils"
	//redis缓存操作与支持驱动
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/garyburd/redigo/redis"
	_ "github.com/garyburd/redigo/redis"
	_ "github.com/gomodule/redigo/redis"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) GetSession(ctx context.Context, req *example.Request, rsp *example.Response) error {
	//初始化返回值
	rsp.Errno = utils.RECODE_SESSIONERR
	rsp.Errmsg = utils.RecodeText(rsp.Errno)
	/*
		redis连接
	*/
	redis_config_map := map[string]string{
		"key":   utils.G_server_name,
		"conn":  utils.G_redis_addr + ":" + utils.G_redis_port,
		"dbNum": utils.G_redis_dbnum,
	}
	beego.Info(redis_config_map)
	redis_config, _ := json.Marshal(redis_config_map)
	bm, err := cache.NewCache("redis", string(redis_config))
	if err != nil {
		beego.Info("连接redis失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	/*
		获取session
	*/
	sessionidname := req.Sessionid + "name"
	//从缓存中获取session 那么使用唯一识别码 通过key查询用户名
	user_info_value := bm.Get(sessionidname)
	//通过redis方法进行转换
	name, err := redis.String(user_info_value, nil)
	if err != nil {
		rsp.Errno = utils.RECODE_DATAERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
	}
	beego.Info(name, reflect.TypeOf(name))
	//获取到了session
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)
	rsp.Data = name
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *Example) Stream(ctx context.Context, req *example.StreamingRequest, stream example.Example_StreamStream) error {
	log.Logf("Received Example.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Logf("Responding: %d", i)
		if err := stream.Send(&example.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *Example) PingPong(ctx context.Context, stream example.Example_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Logf("Got ping %v", req.Stroke)
		if err := stream.Send(&example.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
