package handler

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/micro/go-log"

	example "sss/GetSmscd/proto/example"
	"sss/IhomeWeb/models"

	"fmt"
	"math/rand"
	"sss/ihomeweb/utils"
	sms "submail_go_sdk/submail/sms"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"

	//redis缓存操作与支持驱动
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/garyburd/redigo/redis"
	_ "github.com/garyburd/redigo/redis"
	_ "github.com/gomodule/redigo/redis"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) GetSmscd(ctx context.Context, req *example.Request, rsp *example.Response) error {
	beego.Info(" GET smscd api/v1.0/smscode/:id ")
	//初始化返回正确的返回值
	rsp.Error = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Error)

	/*
		验证手机号
	*/

	o := orm.NewOrm()
	user := models.User{Mobile: req.Mobile}
	err := o.Read(&user)
	if err == nil {
		beego.Info("用户已经存在")
		rsp.Error = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Error)
		return nil
	}
	beego.Info(err)

	/*
		验证图片验证码
	*/

	//连接redis数据库
	redis_config_map := map[string]string{
		"key":   utils.G_server_name,
		"conn":  utils.G_redis_addr + ":" + utils.G_redis_port,
		"dbNum": utils.G_redis_dbnum,
	}
	beego.Info(redis_config_map)
	redis_config, _ := json.Marshal(redis_config_map)
	bm, err := cache.NewCache("redis", string(redis_config))
	if err != nil {
		beego.Info("缓存创建失败", err)
		rsp.Error = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Error)
		return nil
	}
	beego.Info(req.Id, reflect.TypeOf(req.Id))
	//查询相关数据
	value := bm.Get(req.Id)
	if value == nil {
		beego.Info("获取到缓存图片验证码失败", value)
		rsp.Error = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Error)
		return nil
	}
	//打印数据类型,proto中是string类型,redis取出来是int类型
	// beego.Info(value, reflect.TypeOf(value))
	value_str, _ := redis.String(value, nil)
	// beego.Info(value_str, reflect.TypeOf(value_str))
	//数据比对
	if req.Text != value_str {
		beego.Info("图片验证码 错误 ")
		rsp.Error = utils.RECODE_YANZHENG
		rsp.Errmsg = utils.RecodeText(rsp.Error)
		return nil
	}

	/*
		发送手机验证码
	*/

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	size := r.Intn(9999) + 1001
	beego.Info(size)
	// SMS 短信服务配置 appid & appkey 请前往：https://www.mysubmail.com/chs/sms/apps 获取
	config := make(map[string]string)
	config["appid"] = "46739"
	config["appkey"] = "07cd7704750bb58cfd36335a06d1f413"
	// SMS 数字签名模式 normal or md5 or sha1 ,normal = 明文appkey鉴权 ，md5 和 sha1 为数字签名鉴权模式
	config["signType"] = "normal"
	//创建 短信 Send 接口
	submail := sms.CreateSend(config)
	//设置联系人 手机号码
	submail.SetTo(req.Mobile)
	//设置短信正文，请注意：国内短信需要强制添加短信签名，并且需要使用全角大括号 “【签名】”标识，并放在短信正文的最前面
	str := fmt.Sprintf("【golang租房项目】您的验证码是：%d,请在5分钟输入", size)
	submail.SetContent(str)
	//执行 Send 方法发送短信
	submail.Send()

	/*通过手机号将验证短信进行缓存*/
	err = bm.Put(req.Mobile, size, time.Second*300)
	if err != nil {
		beego.Info("缓存出现问题")
		rsp.Error = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Error)
		return nil
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
