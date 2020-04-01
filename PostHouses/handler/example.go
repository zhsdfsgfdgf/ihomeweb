package handler

import (
	"context"
	"encoding/json"
	"reflect"
	"strconv"

	"sss/IhomeWeb/models"
	example "sss/PostHouses/proto/example"
	"sss/ihomeweb/utils"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/go-log/log"

	//redis缓存操作与支持驱动
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	_ "github.com/garyburd/redigo/redis"
	_ "github.com/gomodule/redigo/redis"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) PostHouses(ctx context.Context, req *example.Request, rsp *example.Response) error {
	beego.Info("用户上传房源信息")
	/* 初始化返回正确的返回值 */
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
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	// 获取用户id
	session_user_id := req.Sessionid + "user_id"
	user_id := bm.Get(session_user_id)
	beego.Info(user_id, reflect.TypeOf(user_id))
	id := int(user_id.([]uint8)[0])
	/*
		处理web端获取的数据
	*/
	var request = make(map[string]interface{})
	json.Unmarshal(req.Body, &request)
	/*
		数据准备
			1.User表
			2.Area表
			3.House表
	*/
	user := models.User{Id: id}
	area_id, _ := strconv.Atoi(request["area_id"].(string))
	area := models.Area{Id: area_id}
	facility := []*models.Facility{}
	for _, value := range request["facility"].([]interface{}) {
		//设施编号
		fid, _ := strconv.Atoi(value.(string))
		//创建临时变量,使用设施编号创建设施表对象的指针
		ftmp := &models.Facility{Id: fid}
		facility = append(facility, ftmp)
	}
	title := request["title"].(string)
	price, _ := strconv.Atoi(request["price"].(string))
	roomCount, _ := strconv.Atoi(request["room_count"].(string))
	acreage, _ := strconv.Atoi(request["acreage"].(string))
	capacity, _ := strconv.Atoi(request["capacity"].(string))
	deposit, _ := strconv.Atoi(request["deposit"].(string))
	mindays, _ := strconv.Atoi(request["min_days"].(string))
	maxdays, _ := strconv.Atoi(request["max_days"].(string))
	house := models.House{
		User:       &user,
		Title:      title,
		Price:      price * 100,
		Area:       &area,
		Address:    request["address"].(string),
		Room_count: roomCount,
		Acreage:    acreage,
		Unit:       request["unit"].(string),
		Capacity:   capacity,
		Beds:       request["beds"].(string),
		Deposit:    deposit * 100,
		Min_days:   mindays,
		Max_days:   maxdays,
	}
	/*
		数据库操作
	*/
	o := orm.NewOrm()
	_, err = o.Insert(&house)
	if err != nil {
		beego.Info("房屋表数据插入失败")
		rsp.Errno = utils.RECODE_DATAERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	//多对多表的插入
	m2m := o.QueryM2M(&house, "facilities") //创建一个m2m对象
	_, err = m2m.Add(facility)
	if err != nil {
		beego.Info("房屋设施多对多数据插入失败")
		rsp.Errno = utils.RECODE_DATAERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	/*
		给web端返回房屋id
	*/
	rsp.HouseId = strconv.Itoa(house.Id)
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
