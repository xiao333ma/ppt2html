package common

import (
	"crypto/md5"
	"fmt"
	"github.com/astaxie/beego"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

//上传文件根目录
func GetUloadFileBaseDir() string {
	root := GetStaticPath()
	upload := beego.AppConfig.String("upload_file_base_path")
	return filepath.Join(root, upload)
}

//创建上传文件夹子文件夹
func MakeFileDir(s string) (filedir string, err error) {
	filedir = GetUloadFileBaseDir() + "/" + s
	err = os.MkdirAll(filedir, 0777)
	return filedir, err
}

//判断文件或文件夹是否存在
func HasFile(s string) bool {
	f, err := os.Open(s)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	f.Close()
	return true
}

//File-File复制文件
func CopyFF(src io.Reader, dst io.Writer) error {
	_, err := io.Copy(dst, src)
	return err
}

//File-String复制文件
func CopyFS(src io.Reader, dst string) error {
	f, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, src)
	return err
}

func Md5FS(src io.Reader) string {
	h := md5.New()
	if err := CopyFF(src, h); err != nil {
		fmt.Println(err.Error())
		return ""
	}
	return fmt.Sprintf("%x", h.Sum([]byte("hunterhug")))
}

//判断是否是文件
func IsFile(filepath string) bool {
	fielinfo, err := os.Stat(filepath)
	if err != nil {
		return false
	} else {
		if fielinfo.IsDir() {
			return false
		} else {
			return true
		}
	}
}

//判断是否是文件夹
func IsDir(filepath string) bool {
	fielinfo, err := os.Stat(filepath)
	if err != nil {
		return false
	} else {
		if fielinfo.IsDir() {
			return true
		} else {
			return false
		}
	}
}

//文件状态
func FileStatus(filepath string) {
	fielinfo, err := os.Stat(filepath)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("%v", fielinfo)
	}
}

//文件夹下数量
func SizeofDir(dirPth string) int {
	if IsDir(dirPth) {
		files, _ := ioutil.ReadDir(dirPth)
		return len(files)
	}

	return 0
}

//获取文件后缀
func GetFileSuffix(f string) string {
	fa := strings.Split(f, ".")
	if len(fa) == 0 {
		return ""
	} else {
		return fa[len(fa)-1]
	}
}
