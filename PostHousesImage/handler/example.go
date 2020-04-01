package handler

import (
	"context"
	"path"
	"reflect"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/micro/go-log"

	example "sss/PostHousesImage/proto/example"
	"sss/ihomeweb/models"
	"sss/ihomeweb/utils"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) PostHousesImage(ctx context.Context, req *example.Request, rsp *example.Response) error {
	beego.Info("服务端上传房源图片")
	/* 初始化返回正确的返回值 */
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)
	/* 检查图片数据是否有丢失 */
	size := len(req.Image)
	if req.Filesize != int64(size) {
		beego.Info("传输图片过程中数据有丢失")
		rsp.Errno = utils.RECODE_DATAERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
	}
	/* 调用dfs函数,上传至图片服务器 */
	//获取文件的后缀名,如jpg等
	fileext := path.Ext(req.FileName)
	Group, FileId, err := utils.UploadByBuffer(req.Image, fileext[1:])
	if err != nil {
		beego.Info("Postupavatar utils.UploadByBuffer err", err)
		rsp.Errno = utils.RECODE_IOERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	beego.Info(Group)
	beego.Info(FileId, reflect.TypeOf(FileId))
	beego.Info(req.Id, reflect.TypeOf(req.Id))
	o := orm.NewOrm()
	id, _ := strconv.Atoi(req.Id)
	house := models.House{Id: id}
	err = o.Read(&house, "id")
	if err != nil {
		beego.Info("不存在该房源", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	// 	判断index_image url 是否为空
	if house.Index_image_url == "" {
		house.Index_image_url = FileId
	}
	houseImage := models.HouseImage{House: &house, Url: FileId}
	house.Images = append(house.Images, &houseImage)

	_, err = o.Insert(&houseImage)
	if err != nil {
		beego.Info("数据插入失败：", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	_, err = o.Update(&house)
	if err != nil {
		beego.Info("房屋数据更新失败：", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	rsp.Url = FileId
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
