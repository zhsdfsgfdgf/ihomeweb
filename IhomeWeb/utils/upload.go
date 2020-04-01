package utils

import (
	"github.com/astaxie/beego"
	"github.com/weilaihui/fdfs_client"
)

/*
	上传二进制文件到dfs服务器
		参数:二进制文件,后缀
*/
func UploadByBuffer(fileBuffer []byte, fileExt string) (GroupName, RemoteFileId string, err error) {
	fdfsClient, thiserr := fdfs_client.NewFdfsClient("./conf/client.conf")
	if thiserr != nil {
		beego.Info("创建上传文件句柄失败", thiserr)
		GroupName = ""
		RemoteFileId = ""
		err = thiserr
		return
	}
	//通过句柄上传二进制的文件
	uploadResponse, thiserr := fdfsClient.UploadByBuffer(fileBuffer, fileExt)
	if thiserr != nil {
		beego.Info("上传文件失败", thiserr)
		GroupName = ""
		RemoteFileId = ""
		err = thiserr
		return
	}
	beego.Info(uploadResponse.GroupName)
	beego.Info(uploadResponse.RemoteFileId)
	//回传
	return uploadResponse.GroupName, uploadResponse.RemoteFileId, nil
}
