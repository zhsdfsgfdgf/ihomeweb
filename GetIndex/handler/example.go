package handler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/micro/go-log"

	example "sss/GetIndex/proto/example"
	"sss/ihomeweb/models"
	"sss/ihomeweb/utils"

	"github.com/astaxie/beego"

	//redis缓存操作与支持驱动

	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	_ "github.com/garyburd/redigo/redis"
	_ "github.com/gomodule/redigo/redis"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) GetIndex(ctx context.Context, req *example.Request, rsp *example.Response) error {
	beego.Info("获取首页图片")
	//初始化返回值
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)
	house_page_key := "home_page_data"
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
	house_page_value := bm.Get(house_page_key)
	if house_page_value != nil {
		beego.Info("获取到首页图片缓存")
		rsp.Mix = house_page_value.([]byte)
		return nil
	}
	/* 数据库查询数据,并写入缓存 */
	data := []interface{}{}
	houses := []models.House{}
	o := orm.NewOrm()
	if _, err := o.QueryTable("house").Limit(models.HOME_PAGE_MAX_HOUSES).All(&houses); err == nil {
		for _, house := range houses {
			o.LoadRelated(&house, "Area")
			o.LoadRelated(&house, "User")
			o.LoadRelated(&house, "Facilities")
			o.LoadRelated(&house, "Images")
			data = append(data, house.To_house_info())
		}
	}
	beego.Info(data, houses)
	house_page_value, _ = json.Marshal(data)
	bm.Put(house_page_key, house_page_value, 3600*time.Second)
	rsp.Mix = house_page_value.([]byte)
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
