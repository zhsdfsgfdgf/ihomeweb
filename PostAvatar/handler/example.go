package handler

import (
	"context"
	"encoding/json"
	"path"
	"reflect"

	"sss/IhomeWeb/models"
	example "sss/PostAvatar/proto/example"
	"sss/ihomeweb/utils"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/micro/go-log"

	//redis缓存操作与支持驱动
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	_ "github.com/garyburd/redigo/redis"
	_ "github.com/gomodule/redigo/redis"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) PostAvatar(ctx context.Context, req *example.Request, rsp *example.Response) error {
	beego.Info("上传用户头像 PostAvatar /api/v1.0/user/avatar")
	/* 初始化返回正确的返回值 */
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)
	/* 检查图片数据是否有丢失 */
	size := len(req.Avatar)
	if req.Filesize != int64(size) {
		beego.Info("传输图片过程中数据有丢失")
		rsp.Errno = utils.RECODE_DATAERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
	}
	/* 调用dfs函数,上传至图片服务器 */
	//获取文件的后缀名,如jpg等
	fileext := path.Ext(req.Filename)
	Group, FileId, err := utils.UploadByBuffer(req.Avatar, fileext[1:])
	if err != nil {
		beego.Info("Postupavatar utils.UploadByBuffer err", err)
		rsp.Errno = utils.RECODE_IOERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	beego.Info(Group)
	beego.Info(FileId)
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
	/* 获取用户id */
	session_user_id := req.Sessionid + "user_id"
	user_id := bm.Get(session_user_id)
	beego.Info(user_id, reflect.TypeOf(user_id))
	id := int(user_id.([]uint8)[0])
	/* 更新表中数据 */
	beego.Info(FileId, reflect.TypeOf(FileId))
	beego.Info(id)
	o := orm.NewOrm()
	user := models.User{Id: id, Avatar_url: FileId}
	beego.Info(user)
	_, err = o.Update(&user, "avatar_url")
	if err != nil {
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
	}
	rsp.AvatarUrl = FileId
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
