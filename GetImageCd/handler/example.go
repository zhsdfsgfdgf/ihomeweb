package handler

import (
	"context"
	"encoding/json"
	"image/color"
	"sss/ihomeweb/utils"
	"time"

	"github.com/afocus/captcha"
	"github.com/astaxie/beego"
	"github.com/micro/go-log"

	//redis缓存操作与支持驱动
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	_ "github.com/garyburd/redigo/redis"
	_ "github.com/gomodule/redigo/redis"

	example "sss/GetImageCd/proto/example"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) GetImageCd(ctx context.Context, req *example.Request, rsp *example.Response) error {
	beego.Info("GET /api/v1.0/imagecode/:uuid GetImage()")
	cap := captcha.New()

	if err := cap.SetFont("comic.ttf"); err != nil {
		beego.Info("没有字体文件")
		panic(err.Error())
	}
	cap.SetSize(128, 64)
	cap.SetDisturbance(captcha.MEDIUM)
	cap.SetFrontColor(color.RGBA{255, 255, 255, 255})
	cap.SetBkgColor(color.RGBA{255, 0, 0, 255}, color.RGBA{0, 0, 255, 255}, color.RGBA{0, 153, 0, 255})
	//生成图片 返回图片和 字符串(图片内容的文本形式)
	img, str := cap.Create(4, captcha.NUM)
	beego.Info(str)
	b := *img      //解引用
	c := *(b.RGBA) //解引用
	//成功返回
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)
	//图片信息
	rsp.Pix = []byte(c.Pix)
	rsp.Stride = int64(c.Stride)
	rsp.Max = &example.Response_Point{X: int64(c.Rect.Max.X), Y: int64(c.Rect.Max.Y)}
	rsp.Min = &example.Response_Point{X: int64(c.Rect.Min.X), Y: int64(c.Rect.Min.Y)}
	/*将uuid与 随机数验证码对应的存储在redis缓存中*/
	//redis
	//准备连接信息
	redis_config_map := map[string]string{
		"key":   utils.G_server_name,
		"conn":  utils.G_redis_addr + ":" + utils.G_redis_port,
		"dbNum": utils.G_redis_dbnum,
	}
	//确定连接信息
	beego.Info(redis_config_map)
	//将连接信息由map转化json
	redis_config, _ := json.Marshal(redis_config_map)
	//连接redis
	bm, err := cache.NewCache("redis", string(redis_config))
	if err != nil {
		beego.Info("连接redis失败", err)
	} else {
		beego.Info("连接redis成功")
	}
	//验证码进行5分钟缓存
	bm.Put(req.Uuid, str, 300*time.Second)
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
