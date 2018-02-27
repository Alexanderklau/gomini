package main

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"

	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
)

const (
	upload_path string = "/files/"
)

var (
	ACCESS_KEY = "********************"
	SECRET_KEY = "********************"
	BUCKET     = "blog"
)

func helloHandle(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello world!")
}

//上传
func uploadHandle(w http.ResponseWriter, r *http.Request) {
	//从请求当中判断方法
	if r.Method == "GET" {
		tmp, err := template.ParseFiles("templates/upload.html")
		if err != nil {
			fmt.Println("模版渲染失败")
		}
		tmp.Execute(w, nil)
	} else {
		//获取文件内容 要这样获取
		file, head, err := r.FormFile("file")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		//创建文件
		pwd, _ := os.Getwd()
		err = os.Mkdir(pwd+upload_path, os.ModePerm)
		if err != nil {
			fmt.Println("dir is create Error")
		}
		fW, err := os.Create(pwd + upload_path + head.Filename)
		if err != nil {
			fmt.Println("文件创建失败")
			return
		}
		fmt.Println(*fW)
		defer fW.Close()
		_, err = io.Copy(fW, file)
		if err != nil {
			fmt.Println("文件保存失败")
			return
		}
		upload_qiniu(pwd + upload_path + head.Filename)
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

//七牛上传
func upload_qiniu(filePath string) {
	key := "github-x.png"
	putPolicy := storage.PutPolicy{
		Scope: BUCKET,
	}
	mac := qbox.NewMac(ACCESS_KEY, SECRET_KEY)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	//空间对应机房
	cfg.Zone = &storage.ZoneHuanan
	//不启用HTTPS域名
	cfg.UseHTTPS = false
	//不使用CND加速
	cfg.UseCdnDomains = false
	//构建上传表单对象
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	// 可选
	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": "github logo",
		},
	}
	err := formUploader.PutFile(context.Background(), &ret, upToken, key, filePath, &putExtra)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ret.Key, ret.Hash)
}

func main() {
	//启动一个http 服务器
	http.HandleFunc("/", helloHandle)
	//上传
	http.HandleFunc("/upload", uploadHandle)
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		fmt.Println("服务器启动失败")
		return
	}
	fmt.Println("服务器启动成功")
}
