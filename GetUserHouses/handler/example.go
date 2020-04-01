package handler

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"

	example "sss/GetUserHouses/proto/example"
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
func (e *Example) GetUserHouses(ctx context.Context, req *example.Request, rsp *example.Response) error {
	beego.Info("获取用户房源信息")
	//初始化返回值
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
	/* 通过sessionid获得用户id，进而获得用户房源信息 */
	user_id_str := bm.Get(req.SessionId + "user_id")
	id := int(user_id_str.([]uint8)[0])
	beego.Info(id, reflect.TypeOf(id))
	o := orm.NewOrm()
	qs := o.QueryTable("house")
	houses_list := []models.House{}
	_, err = qs.Filter("user_id", id).All(&houses_list)
	if err != nil {
		beego.Info("查询房屋信息失败")
		rsp.Errno = utils.RECODE_DATAERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	/* 返回给web端,注意编码成二进制数据 */
	house, _ := json.Marshal(houses_list)
	rsp.Mix = house
	return nil
}
