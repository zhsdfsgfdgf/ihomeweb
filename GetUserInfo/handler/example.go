package handler

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/micro/go-log"

	example "sss/GetUserInfo/proto/example"
	"sss/ihomeweb/models"
	"sss/ihomeweb/utils"

	//redis缓存操作与支持驱动
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	_ "github.com/garyburd/redigo/redis"
	_ "github.com/gomodule/redigo/redis"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) GetUserInfo(ctx context.Context, req *example.Request, rsp *example.Response) error {
	beego.Info("-- GET /api/v1.0/user Getuserinfo() --")
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)
	/*
		连接redis
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
		根据sessionid获取用户id,然后数据库查询相关信息
	*/
	sessioniduserid := req.Sessionid + "user_id"
	//获取到用户id值
	value_id := bm.Get(sessioniduserid)
	//查看value_id是什么类型
	beego.Info(value_id, reflect.TypeOf(value_id))
	id := int(value_id.([]uint8)[0])
	user := models.User{Id: id}
	o := orm.NewOrm()
	err = o.Read(&user)
	if err != nil {
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	//将查询到的数据依次赋值
	rsp.UserId = int64(user.Id)
	rsp.Name = user.Name
	rsp.Mobile = user.Mobile
	rsp.RealName = user.Real_name
	rsp.IdCard = user.Id_card
	rsp.AvatarUrl = user.Avatar_url

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
