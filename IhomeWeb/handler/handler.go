package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"net/http"
	"regexp"
	"sss/ihomeweb/utils"
	"time"

	"github.com/afocus/captcha"
	"github.com/astaxie/beego"
	"github.com/julienschmidt/httprouter"
	example "github.com/micro/examples/template/srv/proto/example"
	"github.com/micro/go-grpc"
	"github.com/micro/go-micro/client"

	//调用area的proto
	DELETESESSION "sss/DeleteSession/proto/example"
	GETAREA "sss/GetArea/proto/example"
	GETHOUSEINFO "sss/GetHouseInfo/proto/example"
	GETIMAGECD "sss/GetImageCd/proto/example"
	GETINDEX "sss/GetIndex/proto/example"
	GETSESSION "sss/GetSession/proto/example"
	GETSMSCD "sss/GetSmscd/proto/example"
	GETUSERHOUSES "sss/GetUserHouses/proto/example"
	GETUSERINFO "sss/GetUserInfo/proto/example"
	"sss/IhomeWeb/models"
	POSTAUTHUSER "sss/PostAuthUser/proto/example"
	POSTAVATAR "sss/PostAvatar/proto/example"
	POSTHOUSES "sss/PostHouses/proto/example"
	POSTHOUSESIMAGE "sss/PostHousesImage/proto/example"
	POSTLOGIN "sss/PostLogin/proto/example"
	POSTREG "sss/PostReg/proto/example"
)

func ExampleCall(w http.ResponseWriter, r *http.Request) {
	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// call the backend service
	exampleClient := example.NewExampleService("go.micro.srv.template", client.DefaultClient)
	rsp, err := exampleClient.Call(context.TODO(), &example.Request{
		Name: request["name"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"msg": rsp.Msg,
		"ref": time.Now().UnixNano(),
	}

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

//需要参数,但我又不要传参(要匹配url后面值的处理才要)
func GetArea(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	beego.Info(" GetArea api/v1.0/areas !!!")
	//创建grpc服务
	server := grpc.NewService()
	//服务初始化
	server.Init()
	//调用服务,返回句柄
	exampleClient := GETAREA.NewExampleService("go.micro.srv.GetArea", server.Client())
	//调用服务,返回数据
	//没传数据
	rsp, err := exampleClient.GetArea(context.TODO(), &GETAREA.Request{})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	//接收数据
	//准备接收切片
	////循环读取服务返回的数据
	area_list := []models.Area{}
	for _, value := range rsp.Data {
		tmp := models.Area{Id: int(value.Aid), Name: value.Aname}
		area_list = append(area_list, tmp)
	}
	//创建返回数据map
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":   area_list,
	}
	//回传数据的时候不能直接发送过去,需要设置数据格式
	w.Header().Set("Content-Type", "application/json")
	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func GetIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	/*
		调用后端服务
	*/
	server := grpc.NewService()
	server.Init()
	exampleClient := GETINDEX.NewExampleService("go.micro.srv.GetIndex", server.Client())
	rsp, err := exampleClient.GetIndex(context.TODO(), &GETINDEX.Request{})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	data := []interface{}{}
	json.Unmarshal(rsp.Mix, &data)
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":   data,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func GetSession(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	/*
		创建服务及句柄
	*/
	service := grpc.NewService()
	service.Init()
	exampleClient := GETSESSION.NewExampleService("go.micro.srv.GetSession", service.Client())
	/*
		cookie操作
	*/
	//获取cookie
	userlogin, err := r.Cookie("userlogin")

	//如果不存在就返回
	if err != nil {
		//创建返回数据map
		response := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 503)
			beego.Info(err)
			return
		}
		return
	}

	//存在就发送数据给服务
	rsp, err := exampleClient.GetSession(context.TODO(), &GETSESSION.Request{
		Sessionid: userlogin.Value,
	})
	if err != nil {
		http.Error(w, err.Error(), 502)
		beego.Info(err)
		return
	}
	/*
		服务响应,给前端返回数据
	*/
	//将获取到的用户名返回给前端
	data := make(map[string]string)
	data["name"] = rsp.Data
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":   data,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func GetImageCd(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	beego.Info("获取图片验证码 url：api/v1.0/imagecode/:uuid")
	//创建服务
	server := grpc.NewService()
	//服务初始化
	server.Init()
	//连接服务
	exampleClient := GETIMAGECD.NewExampleService("go.micro.srv.GetImageCd",
		server.Client())
	//获取前端发送过来的唯一uuid
	beego.Info(ps.ByName("uuid"))
	//通过句柄调用我们proto协议中准备好的函数
	//第一个参数为默认,第二个参数 proto协议中准备好的请求包
	rsp, err := exampleClient.GetImageCd(context.TODO(), &GETIMAGECD.Request{
		Uuid: ps.ByName("uuid"),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	//处理发送过来的图片信息
	var img image.RGBA
	img.Stride = int(rsp.Stride)
	img.Rect.Min.X = int(rsp.Min.X)
	img.Rect.Min.Y = int(rsp.Min.Y)
	img.Rect.Max.X = int(rsp.Max.X)
	img.Rect.Max.Y = int(rsp.Max.Y)
	img.Pix = []uint8(rsp.Pix)
	var image captcha.Image
	image.RGBA = &img
	//将图片发送给前端
	png.Encode(w, image)
}

func GetSmscd(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	beego.Info("获取短信验证码")
	//创建服务
	server := grpc.NewService()
	server.Init()
	//获取 前端发送过来的手机号
	mobile := ps.ByName("mobile")
	beego.Info(mobile)
	//后端进行正则匹配
	//创建正则句柄
	myreg := regexp.MustCompile(`0?(13|14|15|17|18|19)[0-9]{9}`)
	//进行正则匹配
	bo := myreg.MatchString(mobile)
	//如果手机号错误则
	if bo == false {
		resp := map[string]interface{}{
			"errno":  utils.RECODE_NODATA,
			"errmsg": "手机号错误",
		}
		//设置返回数据格式
		w.Header().Set("Content-Type", "application/json")
		//将错误发送给前端
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			beego.Info(err)
			return
		}
		beego.Info("手机号错误返回")
	}

	//获取url携带的参数
	text := r.URL.Query()["text"][0]
	id := r.URL.Query()["id"][0]
	exampleClient := GETSMSCD.NewExampleService("go.micro.srv.GetSmscd", server.Client())
	rsp, err := exampleClient.GetSmscd(context.TODO(), &GETSMSCD.Request{
		Mobile: mobile,
		//uuid
		Id: id,
		//验证码
		Text: text,
	})
	if err != nil {
		http.Error(w, err.Error(), 502)
		beego.Info(err) //beego.Debug(err)
		return
	}
	//创建返回map
	resp := map[string]interface{}{
		"errno":  rsp.Error,
		"errmsg": rsp.Errmsg,
	}
	//设置返回格式
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
func Postreg(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	beego.Info(" 注册请求 /api/v1.0/users ")

	/*获取前端发送过来的json数据*/

	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	//由于前端没做所以后端进行下操作
	if request["mobile"] == "" || request["password"] == "" || request["sms_code"] == "" {
		resp := map[string]interface{}{
			"errno":  utils.RECODE_NODATA,
			"errmsg": "信息有误请重新输入",
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			beego.Info(err)
			return
		}
		beego.Info("所发送数据为空")
		return
	}
	server := grpc.NewService()
	server.Init()
	// call the backend service
	exampleClient := POSTREG.NewExampleService("go.micro.srv.PostReg", server.Client())
	rsp, err := exampleClient.PostReg(context.TODO(), &POSTREG.Request{
		Mobile:   request["mobile"].(string),
		Password: request["password"].(string),
		SmsCode:  request["sms_code"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	//读取cookie
	cookie, err := r.Cookie("userlogin")
	//如果读取失败或者cookie的value中不存在则创建cookie
	if err != nil || "" == cookie.Value {
		cookie := http.Cookie{Name: "userlogin", Value: rsp.SessionID, Path: "/", MaxAge: 3600}
		http.SetCookie(w, &cookie)
	}
	//准备回传数据
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
	}
	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func PostLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	beego.Info("登陆 api/v1.0/sessions")
	/*
		获取前端post请求发送的内容,并进行基本判断
	*/
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	//判断账号密码是否为空
	if request["mobile"] == "" || request["password"] == "" {
		resp := map[string]interface{}{
			"errno":  utils.RECODE_NODATA,
			"errmsg": "信息有误请从新输入",
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			beego.Info(err)
			return
		}
		beego.Info("账号和密码输入不能为空")
		return
	}
	/*
		调用后端服务
	*/
	server := grpc.NewService()
	server.Init()
	exampleClient := POSTLOGIN.NewExampleService("go.micro.srv.PostLogin", server.Client())
	rsp, err := exampleClient.PostLogin(context.TODO(), &POSTLOGIN.Request{
		Password: request["password"].(string),
		Mobile:   request["mobile"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	cookie, err := r.Cookie("userlogin")
	if err != nil || "" == cookie.Value {
		cookie := http.Cookie{Name: "userlogin", Value: rsp.SessionID, Path: "/", MaxAge: 3600}
		http.SetCookie(w, &cookie)
	}
	resp := map[string]interface{}{
		"errno": rsp.Errno, "errmsg": rsp.Errmsg,
	}
	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
func DeleteSession(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	/*
		获取传来的cookie
	*/
	userlogin, err := r.Cookie("userlogin")
	//如果没有数据说明没有的登陆直接返回错误
	if err != nil {
		resp := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			beego.Info(err)
			return
		}
		return
	}
	/*
		调用后端服务
	*/
	server := grpc.NewService()
	server.Init()
	exampleClient := DELETESESSION.NewExampleService("go.micro.srv.DeleteSession", server.Client())
	rsp, err := exampleClient.DeleteSession(context.TODO(), &DELETESESSION.Request{
		Sessionid: userlogin.Value,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	/*
		删除cookie
			cookie没有删除的写法,可以把时间设置为负的
	*/
	cookie, err := r.Cookie("userlogin")
	if err != nil || "" == cookie.Value {
		return
	} else {
		cookie := http.Cookie{Name: "userlogin", Path: "/", MaxAge: -1}
		http.SetCookie(w, &cookie)
	}
	resp := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), 503)
		beego.Info(err)
		return
	}
	return
}
func GetUserInfo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	beego.Info("GetUserInfo 获取用户信息 /api/v1.0/user")
	/*
		获取传来的cookie
	*/
	userlogin, err := r.Cookie("userlogin")
	if err != nil {
		resp := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			beego.Info(err)
			return
		}
		return
	}
	/*
		调用后端服务
	*/
	server := grpc.NewService()
	server.Init()
	exampleClient := GETUSERINFO.NewExampleService("go.micro.srv.GetUserInfo", server.Client())
	rsp, err := exampleClient.GetUserInfo(context.TODO(), &GETUSERINFO.Request{
		Sessionid: userlogin.Value,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	/*
		将服务返回的数据发送到前端
	*/
	data := make(map[string]interface{})
	data["user_id"] = int(rsp.UserId)
	data["name"] = rsp.Name
	data["mobile"] = rsp.Mobile
	data["real_name"] = rsp.RealName
	data["id_card"] = rsp.IdCard
	data["avatar_url"] = utils.AddDomain2Url(rsp.AvatarUrl)
	resp := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":   data,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
func PostAvatar(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	beego.Info("上传用户头像 PostAvatar /api/v1.0/user/avatar")
	/* 查看用户是否登录 */
	userlogin, err := r.Cookie("userlogin")
	if err != nil {
		resp := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}
	/* 获取前端发送过来的文件 */
	file, hander, err := r.FormFile("avatar")
	if err != nil {
		beego.Info("图片发送失败", err)
		resp := map[string]interface{}{
			"errno":  utils.RECODE_IOERR,
			"errmsg": utils.RecodeText(utils.RECODE_IOERR),
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}
	/* 打印基本信息 */
	beego.Info(file, hander)
	beego.Info("文件大小", hander.Size)
	beego.Info("文件名", hander.Filename)
	/*
		将读取的文件处理,为发送给服务端做准备
			1.二进制的空间用来存储文件
			2.将文件读取到准备空间里
	*/
	filebuffer := make([]byte, hander.Size)
	_, err = file.Read(filebuffer)
	if err != nil {
		beego.Info("Postupavatarfile.Read(filebuffer) err", err)
		resp := map[string]interface{}{
			"errno":  utils.RECODE_IOERR,
			"errmsg": utils.RecodeText(utils.RECODE_IOERR),
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}
	/* 调用后端服务 */
	server := grpc.NewService()
	server.Init()
	exampleClient := POSTAVATAR.NewExampleService("go.micro.srv.PostAvatar", server.Client())
	rsp, err := exampleClient.PostAvatar(context.TODO(), &POSTAVATAR.Request{
		Sessionid: userlogin.Value,
		Filename:  hander.Filename,
		Filesize:  hander.Size,
		Avatar:    filebuffer,
	})
	if err != nil {
		http.Error(w, err.Error(), 502)
		beego.Info("接受后端服务出错", err)
		return
	}
	/*
		接受服务端数据,处理后发送给前端
			1.准备回传数据空间
			2.url拼接回传数据(加上图片服务器地址和端口)

	*/
	data := make(map[string]interface{})
	data["avatar_url"] = utils.AddDomain2Url(rsp.AvatarUrl)
	resp := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":   data,
	}
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
func GetUserAuth(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	beego.Info("GetUserInfo 获取用户信息 /api/v1.0/user")
	/*
		获取传来的cookie
	*/
	userlogin, err := r.Cookie("userlogin")
	if err != nil {
		resp := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			beego.Info(err)
			return
		}
		return
	}
	/*
		调用后端服务
	*/
	server := grpc.NewService()
	server.Init()
	exampleClient := GETUSERINFO.NewExampleService("go.micro.srv.GetUserInfo", server.Client())
	rsp, err := exampleClient.GetUserInfo(context.TODO(), &GETUSERINFO.Request{
		Sessionid: userlogin.Value,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	/*
		将服务返回的数据发送到前端
	*/
	data := make(map[string]interface{})
	data["user_id"] = int(rsp.UserId)
	data["name"] = rsp.Name
	data["mobile"] = rsp.Mobile
	data["real_name"] = rsp.RealName
	data["id_card"] = rsp.IdCard
	data["avatar_url"] = utils.AddDomain2Url(rsp.AvatarUrl)
	resp := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":   data,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
func PostAuthUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	/*
		获取传来的cookie
	*/
	userlogin, err := r.Cookie("userlogin")
	if err != nil {
		resp := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			beego.Info(err)
			return
		}
		return
	}
	/* 获取前端发送的信息 */
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	/*
		调用后端服务
	*/
	server := grpc.NewService()
	server.Init()
	exampleClient := POSTAUTHUSER.NewExampleService("go.micro.srv.PostAuthUser", server.Client())
	rsp, err := exampleClient.PostAuthUser(context.TODO(), &POSTAUTHUSER.Request{
		RealName:  request["real_name"].(string),
		IDCard:    request["id_card"].(string),
		SessionID: userlogin.Value,
	})
	if err != nil {
		beego.Info("669错误")
		http.Error(w, err.Error(), 500)
		return
	}

	/* 获取后端服务返回数据,发送给前端 */
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
func GetUserHouses(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	beego.Info("GetUserInfo 获取用户房源信息 /api/v1.0/user")
	/*
		获取传来的cookie
	*/
	userlogin, err := r.Cookie("userlogin")
	if err != nil {
		resp := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			beego.Info(err)
			return
		}
		return
	}
	/*
		调用后端服务
	*/
	server := grpc.NewService()
	server.Init()
	exampleClient := GETUSERHOUSES.NewExampleService("go.micro.srv.GetUserHouses", server.Client())
	rsp, err := exampleClient.GetUserHouses(context.TODO(), &GETUSERHOUSES.Request{
		SessionId: userlogin.Value,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	/*
		将服务返回的数据发送到前端
	*/
	houses_list := []models.House{}
	json.Unmarshal(rsp.Mix, &houses_list)
	var houses []interface{}
	for _, value := range houses_list {
		houses = append(houses, value.To_house_info())
	}
	data := make(map[string]interface{})
	data["houses"] = houses
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":   data,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
func PostHouses(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	beego.Info("用户发布房源信息 /api/v1.0/houses")
	/*
		获取传来的cookie
	*/
	userlogin, err := r.Cookie("userlogin")
	if err != nil {
		resp := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			beego.Info(err)
			return
		}
		return
	}
	/* 获取前端发送的信息 */
	//将前端发送过来的数据整体读取
	//body就是json的二进制流
	body, _ := ioutil.ReadAll(r.Body)
	/*	调用后端服务 */
	server := grpc.NewService()
	server.Init()
	exampleClient := POSTHOUSES.NewExampleService("go.micro.srv.PostHouses", server.Client())
	rsp, err := exampleClient.PostHouses(context.TODO(), &POSTHOUSES.Request{
		Sessionid: userlogin.Value,
		Body:      body,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	/* 获取后端服务返回数据,发送给前端 */
	data := make(map[string]interface{})
	data["house_id"] = rsp.HouseId
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":   data,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
func PostHousesImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	beego.Info("用户上传房源图片 /api/v1.0/houses/:id/images")
	id := ps.ByName("id")
	/*
		获取传来的cookie
	*/
	_, err := r.Cookie("userlogin")
	if err != nil {
		resp := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			beego.Info(err)
			return
		}
		return
	}
	/* 获取前端发送过来的文件 */
	file, hander, err := r.FormFile("house_image")
	if err != nil {
		beego.Info("图片发送失败", err)
		resp := map[string]interface{}{
			"errno":  utils.RECODE_IOERR,
			"errmsg": utils.RecodeText(utils.RECODE_IOERR),
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}
	/* 打印基本信息 */
	beego.Info(file, hander)
	beego.Info("文件大小", hander.Size)
	beego.Info("文件名", hander.Filename)
	/*
		将读取的文件处理,为发送给服务端做准备
			1.二进制的空间用来存储文件
			2.将文件读取到准备空间里
	*/
	filebuffer := make([]byte, hander.Size)
	_, err = file.Read(filebuffer)
	if err != nil {
		beego.Info("图片读取出现问题", err)
		resp := map[string]interface{}{
			"errno":  utils.RECODE_IOERR,
			"errmsg": utils.RecodeText(utils.RECODE_IOERR),
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}
	/* 调用后端服务 */
	server := grpc.NewService()
	server.Init()
	exampleClient := POSTHOUSESIMAGE.NewExampleService("go.micro.srv.PostHousesImage", server.Client())
	rsp, err := exampleClient.PostHousesImage(context.TODO(), &POSTHOUSESIMAGE.Request{
		Id:       id,
		Image:    filebuffer,
		FileName: hander.Filename,
		Filesize: hander.Size,
	})
	if err != nil {
		http.Error(w, err.Error(), 502)
		beego.Info("接受后端服务出错", err)
		return
	}
	/*
		接受服务端数据,处理后发送给前端
			1.准备回传数据空间
			2.url拼接回传数据(加上图片服务器地址和端口)

	*/
	data := make(map[string]interface{})
	data["url"] = utils.AddDomain2Url(rsp.Url)
	resp := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":   data,
	}
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
func GetHouseInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	beego.Info("GetUserInfo 获取房源信息")
	/*
		调用后端服务
	*/
	server := grpc.NewService()
	server.Init()
	exampleClient := GETHOUSEINFO.NewExampleService("go.micro.srv.GetHouseInfo", server.Client())
	rsp, err := exampleClient.GetHouseInfo(context.TODO(), &GETHOUSEINFO.Request{
		Id: ps.ByName("id"),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	/*
		将服务返回的数据发送到前端
	*/
	house := models.House{}
	json.Unmarshal(rsp.Housedata, &house)
	dataMap := make(map[string]interface{})
	dataMap["house"] = house.To_one_house_desc()

	resp := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":   dataMap,
	}
	fmt.Println(dataMap["house"])
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
