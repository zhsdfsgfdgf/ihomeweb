package handler

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/micro/go-log"

	example "sss/PostReg/proto/example"
	"sss/ihomeweb/models"
	"sss/ihomeweb/utils"

	//redis缓存操作与支持驱动
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	_ "github.com/garyburd/redigo/redis"
	_ "github.com/gomodule/redigo/redis"
)

type Example struct{}

//加密函数
func GetMd5String(s string) string {
	h := md5.New()     //创建md5对象
	h.Write([]byte(s)) //将传入的s变成二进制
	return hex.EncodeToString(h.Sum(nil))
}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) PostReg(ctx context.Context, req *example.Request, rsp *example.Response) error {
	beego.Info(" POST userreg/api/v1.0/users !!!")
	//初始化错误码
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)
	/*
		验证短信验证码
	*/

	//连接redis
	redis_config_map := map[string]string{
		"key":   utils.G_server_name,
		"conn":  utils.G_redis_addr + ":" + utils.G_redis_port,
		"dbNum": utils.G_redis_dbnum,
	}
	beego.Info(redis_config_map)
	redis_config, _ := json.Marshal(redis_config_map)
	//连接redis
	bm, err := cache.NewCache("redis", string(redis_config))
	if err != nil {
		beego.Info("连接redis失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	//验证
	value := bm.Get(req.Mobile)
	if value == nil {
		beego.Info("缓存验证码数据查询失败", value)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	/*
		将数据存入数据库
	*/
	user := models.User{Mobile: req.Mobile, Name: req.Mobile, Password_hash: req.Password}
	o := orm.NewOrm()
	id, err := o.Insert(&user)
	if err != nil {
		rsp.Errno = utils.RECODE_DBREREP
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	beego.Info("userid", id)

	/*
		返回给服务端sessionID,并在redis中写入session
	*/

	//生成sessionID 保证唯一性
	h := GetMd5String(req.Mobile + req.Password)
	rsp.SessionID = h
	//拼接key
	bm.Put(h+"name", string(user.Mobile), time.Second*3600)
	bm.Put(h+"user_id", string(user.Id), time.Second*3600)
	bm.Put(h+"mobile", string(user.Mobile), time.Second*3600)
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
