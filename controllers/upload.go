package controllers

import (
	_ "bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	. "ppt2html/common"
)

type UploadResultInfo struct {
	Code     int    `json:"code"`
	FileName string `json:"fileName"`
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
	"key":  nil,
}

// code 0 成功
// code -1 未知功能
// code 1000 非法后缀名
// code 2 拷贝失败
func (this *UploadController) Post() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("捕获到了从panic产生的异常：", err) // 这里的err就是panic传入的内容
		}
	}()
	this.EnableRender = false
	result := new(UploadResultInfo)
	result.Code = -1
	f, h, _ := this.GetFile("myfile")
	fileName := h.Filename
	fileSuffix := GetFileSuffix(fileName)

	if _, ok := FileAllow[fileSuffix]; !ok {
		result.Code = 1000 //代表非法的后缀名
	} else {
		mountPath := ""
		homeDir := GetUserHomePath()
		if len(homeDir) > 0 {
			mountPath = homeDir + "/data/ppt"
		}

		fileBytes, _ := ioutil.ReadAll(f)
		data := md5.Sum(fileBytes)

		str := hex.EncodeToString(data[:])
		file_name := fmt.Sprintf("%s.%s", str, fileSuffix)
		dir_path := GetUloadFileBaseDir()
		dir_path = dir_path + "/"

		if !HasFile(GetUloadFileBaseDir()) {
			MakeFileDir("")
		}
		full_path := dir_path + file_name
		io_error := ioutil.WriteFile(full_path, fileBytes, 0777)
		script_path := GetScriptPath() + "/" + "ppt2html.scpt"
		fmt.Println(full_path)
		if io_error == nil && len(mountPath) > 0 {
			cmd := exec.Command("osascript", script_path, full_path, dir_path)
			_, cmdErr := cmd.CombinedOutput()
			fmt.Println(cmdErr)
			oldPath := dir_path + "/" + str
			newPath := mountPath + "/" + str
			copyErr  := CopyFolder(oldPath, newPath)
			if copyErr == nil {
				result.Code = 0
				result.FileName = str
			}else {
				result.Code = 2
			}
			os.RemoveAll(oldPath)
			os.Remove(full_path)
		}
	}

	cmd := exec.Command("killall", "Keynote")
	cmd.CombinedOutput()

	this.Data["json"] = result
	this.ServeJSON()
}

func (this *UploadController) Get() {
	this.TplName = "upload.html"
}
