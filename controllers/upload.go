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

func (this *UploadController) Post() {
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
			}
			os.RemoveAll(oldPath)
			os.Remove(full_path)
		}
	}
	this.Data["json"] = result
	this.ServeJSON()
}

//func fileHandle(source string, dest string, filePath string) {
//	CopyFolder(source, dest)
//	os.RemoveAll(source)
//	os.Remove(filePath)
//}

func (c *UploadController) Get() {
	c.TplName = "upload.html"
}
