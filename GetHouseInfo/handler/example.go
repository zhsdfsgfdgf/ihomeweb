package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/micro/go-log"

	example "sss/GetHouseInfo/proto/example"
	"sss/ihomeweb/models"
	"sss/ihomeweb/utils"

	//redis缓存操作与支持驱动
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/astaxie/beego/orm"
	_ "github.com/garyburd/redigo/redis"
	_ "github.com/gomodule/redigo/redis"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) GetHouseInfo(ctx context.Context, req *example.Request, rsp *example.Response) error {
	beego.Info("获取房屋信息")
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
		根据房屋id查找房屋缓存信息
	*/
	houseId, _ := strconv.Atoi(req.Id)
	houseInfoKey := fmt.Sprintf("house_info_%s", houseId)
	houseIfoValue := bm.Get(houseInfoKey)
	if houseIfoValue != nil {
		beego.Info("房屋信息在缓存中")
		rsp.Housedata = houseIfoValue.([]byte)
	}
	/*
		查找数据库,并写入缓存
	*/
	o := orm.NewOrm()
	house := models.House{Id: houseId}
	o.Read(&house)

	// 关联查询
	o.LoadRelated(&house, "Area")
	o.LoadRelated(&house, "User")
	o.LoadRelated(&house, "Images")
	o.LoadRelated(&house, "Facilities")
	beego.Info(house)
	// 将查询的结果存入redis中
	houseMix, err := json.Marshal(house)
	bm.Put(houseInfoKey, houseMix, time.Second*3600)

	rsp.Housedata = houseMix

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
