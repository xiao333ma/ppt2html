package controllers

import (
	_ "bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	. "ppt2html/common"
	"strings"
)

type UploadResultInfo struct {
	ResultCode bool   `json:"resultCode"`
	Message    string `json:"message"`
}

type UploadController struct {
	MainController
}

type Sizer interface {
	Size() int64
}

var FileSize int64 = 300 * 1024 * 1024 //300M

var FileAllow = map[string]interface{}{
	"jpg":  nil, //image
	"jpeg": nil,
	"png":  nil,
	"bmp":  nil,
	"gif":  nil,
	"webp": nil,
	"swf":  nil, //video
	"flv":  nil,
	"mp4":  nil,
	"avi":  nil,
	"rmvb": nil,
	"mp3":  nil, //audio
	"m4a":  nil,
	"pcm":  nil,
	"wav":  nil,
	"wma":  nil,
	"aac":  nil,
	"wmv":  nil,
	"doc":  nil, //other
	"docx": nil,
	"xls":  nil,
	"xlsx": nil,
	"ppt":  nil,
	"pptx": nil,
	"pdf":  nil,
	"md":   nil,
	"zip":  nil,
	"rar":  nil}

//
//func (this *UploadController) Get(){
//
//	result := new(UploadResultInfo);
//	result.ResultCode = true;
//	result.Message = "你好你好,";
//	this.Data["json"] = result;
//	this.ServeJSON();
//}

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
	
	workSpace, _ := os.Getwd()
	script_path := workSpace+"/"+"ppt2html.scpt"

	fmt.Println(script_path)
	fmt.Println(full_path)
	fmt.Println(dir_path)





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
		result.Message = "http://"+ip+"/static/upload/"+str+"/index.html"
		this.Data["json"] = result
		this.ServeJSON()
	}






}

func isAllowFileSuffix(fileSuffix string) bool {
	_, flag := FileAllow[fileSuffix]
	return flag
}

type BaiduAuthInfo struct {
	AccessToken   string `json:"access_token"`
	ExpiresIn     int64  `json:"expires_in"`
	RefreshToken  string `json:"refresh_token"`
	Scope         string `json:"scope"`
	SessionKey    string `json:"session_key"`
	SessionSecret string `json:"session_secret"`
}

type BaiduASRResultInfo struct {
	CorpusNo string   `json:"corpus_no"`
	ErrMsg   string   `json:"err_msg"`
	ErrNo    int64    `json:"err_no"`
	Result   []string `json:"result"`
	SN       string   `json:"sn"`
}

var accessToken = ""
var client_id = "CdM3vOjvqn6eGyneRTBP8oyt"
var client_secret = "bFXC2fQVmBbuwkq6HnV1unbxlpRFyo8F"

func tokenBaidu() {
	//先去鉴权
	var auth = new(BaiduAuthInfo)
	accessToken = ""
	httplib.Get("https://openapi.baidu.com/oauth/2.0/token?grant_type=client_credentials&client_id=" + client_id + "&client_secret=" + client_secret).ToJSON(&auth)
	accessToken = auth.AccessToken
}

func asrBaidu(inputFileName string) string {
	result := ""
	inputFilePath := GetUloadFileBaseDir() + "/" + inputFileName
	outPutDir := GetUloadFileBaseDir() + "/asr"
	if !HasFile(outPutDir) {
		MakeFileDir("asr")
	}
	outputFilePath := outPutDir + "/" + inputFileName + ".pcm"
	cmd := exec.Command("ffmpeg", "-y", "-i", inputFilePath, "-acodec", "pcm_s16le", "-f", "s16le", "-ac", "1", "-ar", "16000", outputFilePath)
	_, err := cmd.CombinedOutput()

	fmt.Println(err)
	data, readErr := ioutil.ReadFile(outputFilePath)
	if err == nil && readErr == nil {

		if len(accessToken) < 1 {
			tokenBaidu()
		}

		//请求百度语音转文字
		url := "http://vop.baidu.com/server_api?lan=zh&cuid=liuche&token=" + accessToken
		req := httplib.Post(url)
		req.Header("Content-Type", "audio/pcm;rate=16000")
		req.Body(data)
		asrResult := new(BaiduASRResultInfo)
		req.ToJSON(&asrResult)

		if asrResult.ErrNo == 3302 {
			accessToken = ""
			asrBaidu(inputFileName)
		}
		str := ""
		for _, v := range asrResult.Result {
			v = strings.Replace(v, "，", "", -1)
			str = str + v
		}
		result = str
	}

	//删除
	os.Remove(inputFilePath)
	os.Remove(outputFilePath)
	return result
}


func (c *UploadController) Get() {
	c.TplName = "upload.html"
}

