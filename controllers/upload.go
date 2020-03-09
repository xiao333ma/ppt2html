package controllers

import (
	_ "bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	. "ppt2html/common"
)

type UploadResultInfo struct {
	Code     int    `json:"code"`
	HTMLFolderName string `json:"HTMLFolderName"`
	IMAGEFolderName string `json:"IMAGEFolderName"`
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
	inputFileName :=  this.GetString("fileName")
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
		var folderName string
		if len(inputFileName)>0 {
			folderName = inputFileName
		}else{
			data := md5.Sum(fileBytes)
			folderName = hex.EncodeToString(data[:])
		}
		file_name := fmt.Sprintf("%s.%s", folderName, fileSuffix)
		dir_path := GetUloadFileBaseDir()
		dir_path = dir_path + "/"

		if !HasFile(GetUloadFileBaseDir()) {
			MakeFileDir("")
		}
		full_path := dir_path + file_name
		io_error := ioutil.WriteFile(full_path, fileBytes, 0777)
		if io_error == nil && len(mountPath) > 0 {
			exportHTML_ERR,html_folderName  :=  exportHTML(full_path,dir_path,mountPath,folderName)
			exportImage_ERR,image_folderName  :=  exportImage(full_path,dir_path,mountPath,folderName)
			if exportHTML_ERR == nil && exportImage_ERR == nil{
				result.Code = 0
				result.HTMLFolderName = html_folderName
				result.IMAGEFolderName = image_folderName
			}else {
				result.Code = 2
			}
			os.Remove(full_path)
		}
	}
	this.Data["json"] = result
	this.ServeJSON()
}

func (this *UploadController) Get() {
	this.TplName = "upload.html"
}

func exportHTML(input_file_path string,input_dir_path string,mountPath string,folderName string) (error,string) {
	var err error
	new_folderName := folderName+"_html"
	if len(input_file_path)<1||len(input_dir_path)<1  {
		err = errors.New("file or dir is not found")
	}else {
		if len(input_file_path)<1||len(input_dir_path)<1  {
			err = errors.New("fileName is nil")
		}else{
			script_path := GetScriptPath() + "/" + "ppt2html.scpt"
			cmd := exec.Command("osascript", script_path, input_file_path, input_dir_path)
			_, cmdErr := cmd.CombinedOutput()
			fmt.Println(cmdErr)
			oldPath := input_dir_path  + folderName
			newPath := mountPath + "/" + new_folderName
			err = CopyFolder(oldPath, newPath)
			os.RemoveAll(oldPath)
			commond := exec.Command("killall", "Keynote")
			commond.CombinedOutput()
			if err != nil {
				new_folderName = ""
			}
		}
	}
	return err,new_folderName
}
func exportImage(input_file_path string,input_dir_path string,mountPath string,folderName string) (error,string) {
	var err error
	new_folderName := folderName+"_image"
	if len(input_file_path)<1||len(input_dir_path)<1  {
		err = errors.New("file or dir is not found")
	}else {
		if len(input_file_path)<1||len(input_dir_path)<1  {
			err = errors.New("fileName is nil")
		}else{
			script_path := GetScriptPath() + "/" + "ppt2image.scpt"
			cmd := exec.Command("osascript", script_path, input_file_path, input_dir_path)
			_, cmdErr := cmd.CombinedOutput()
			fmt.Println(cmdErr)
			oldPath := input_dir_path  + folderName
			newPath := mountPath + "/" + new_folderName
			err = CopyFolder(oldPath, newPath)
			os.RemoveAll(oldPath)
			commond := exec.Command("killall", "Keynote")
			commond.CombinedOutput()
			if err != nil {
				new_folderName = ""
			}
		}
	}
	return err,new_folderName
}