package controllers

import (
	_ "bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net"
	"os/exec"
	. "ppt2html/common"
)

type UploadResultInfo struct {
	Tip    string   `json:"tip"`
	URL    string   `json:"url"`
}

type UploadController struct {
	 MainController
}

type Sizer interface {
	Size() int64
}

var FileSize int64 = 300 * 1024 * 1024 //300M

var FileAllow = map[string]interface{}{
	"ppt":  nil,
	"pptx": nil,
	"pdf":  nil,
}


func (this *UploadController) Post() {
	this.EnableRender = false
	f,h,_ := this.GetFile("myfile")
	fileName := h.Filename
	fileSuffix := GetFileSuffix(fileName)
	fmt.Println(fileName)
	fmt.Println(fileSuffix)
	fileBytes, _ := ioutil.ReadAll(f)
	data := md5.Sum(fileBytes)

	str := hex.EncodeToString(data[:])
	file_name := fmt.Sprintf("%s.%s", str, fileSuffix)
	fmt.Println(file_name)
	dir_path := GetUloadFileBaseDir()
	dir_path = dir_path+"/"
	

	if !HasFile(GetUloadFileBaseDir()) {
		MakeFileDir("")
	}
	full_path := dir_path+file_name
	io_error := ioutil.WriteFile(full_path, fileBytes, 0777)
	script_path := GetScriptPath()+"/"+"ppt2html.scpt"

	if io_error == nil {
		cmd := exec.Command("osascript", script_path,full_path,dir_path)
		_, err := cmd.CombinedOutput()
		fmt.Println(err)
		result := new(UploadResultInfo)

		addrs, err := net.InterfaceAddrs()
		ip := ""
		for _, address := range addrs {
			// 检查ip地址判断是否回环地址
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					ip = ipnet.IP.String()+":8080"
				}
			}
		}
		result.Tip = "转换成功"
		result.URL = "http://"+ip+"/static/upload/"+str+"/index.html"
		this.Data["json"] = result
		this.ServeJSON()
	}
}

//func (c *UploadController) Get() {
//	  c.TplName = "upload.html"
//}