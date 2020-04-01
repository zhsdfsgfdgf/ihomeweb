package handler

import (
	"context"
	"encoding/json"
	"reflect"

	example "sss/PostAuthUser/proto/example"
	"sss/ihomeweb/models"
	"sss/ihomeweb/utils"

	"github.com/astaxie/beego"
	"github.com/micro/go-log"

	//redis缓存操作与支持驱动
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/astaxie/beego/orm"
	_ "github.com/garyburd/redigo/redis"
	_ "github.com/gomodule/redigo/redis"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) PostAuthUser(ctx context.Context, req *example.Request, rsp *example.Response) error {
	beego.Info("用户实名认证")
	rsp.Errno = utils.RECODE_OK
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
		rsp.Errno = utils.RECODE_DBREDISER
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	/* 从session中获取用户id,进行数据库更新操作 */
	user_id_str := bm.Get(req.SessionID + "user_id")
	beego.Info(user_id_str, reflect.TypeOf(user_id_str))
	id := int(user_id_str.([]uint8)[0])
	beego.Info(id, reflect.TypeOf(id))
	o := orm.NewOrm()
	user := models.User{Id: id, Real_name: req.RealName, Id_card: req.IDCard}
	_, err = o.Update(&user, "real_name", "id_card")
	if err != nil {
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
	}
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
