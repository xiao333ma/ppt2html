package common

import (
	"fmt"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/utils"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"reflect"
	"strings"
)

func GetStaticPath() string {
	staticPath := ""
	workSpace, err_ws := os.Getwd()
	if err_ws != nil {
		panic(err_ws)
	}
	staticPath = filepath.Join(workSpace, "static")
	if !utils.FileExists(staticPath) {
		appPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			panic(err)
		}
		staticPath = filepath.Join(appPath, "static")
	}

	if len(staticPath) < 1 {
		staticPath = "static"
	}
	return staticPath
}

func GetUserHomePath() string {
	user, err := user.Current()
	if nil == err {
		return user.HomeDir
	}
	return ""
}

func GetViewsPath() string {
	path := ""
	workSpace, err_ws := os.Getwd()
	if err_ws != nil {
		panic(err_ws)
	}
	path = filepath.Join(workSpace, "views")
	if !utils.FileExists(path) {
		appPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			panic(err)
		}
		path = filepath.Join(appPath, "views")
	}

	if len(path) < 1 {
		path = "views"
	}
	return path
}

func GetScriptPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0])) //返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	if err != nil {
		log.Fatal(err)
	}
	strings.Replace(dir, "\\", "/", -1) //将\替换成/
	return dir
}

//获取用户IP地址
func GetClientIp(this *context.Context) string {
	s := strings.Split(this.Request.RemoteAddr, ":")
	if s[0] == "127.0.0.1" {
		if v, ok := this.Request.Header["X-Real-Ip"]; ok {
			if len(v) > 0 {
				return v[0]
			}
		}
	}
	if s[0] == "" {
		s[0] = "127.0.0.1"
	}
	return s[0]
}

func SetStructValue(obj interface{}, fieldName string, value string) {
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return
	}
	fieldNum := t.NumField()
	for i := 0; i < fieldNum; i++ {
		if t.Field(i).Name == fieldName {
			p := reflect.ValueOf(obj)
			if p.Kind() == reflect.Ptr {
				field := p.Elem().FieldByName(fieldName)
				if field.Kind() == reflect.String {
					field.SetString(value)
				}
			}
			break
		}
	}
}

func CopyFolder(source string, dest string) (err error) {
	sourceinfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dest, sourceinfo.Mode())
	if err != nil {
		return err
	}

	directory, _ := os.Open(source)

	objects, err := directory.Readdir(-1)

	for _, obj := range objects {

		sourcefilepointer := source + "/" + obj.Name()

		destinationfilepointer := dest + "/" + obj.Name()

		if obj.IsDir() {
			err = CopyFolder(sourcefilepointer, destinationfilepointer)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			err = CopyFile(sourcefilepointer, destinationfilepointer)
			if err != nil {
				fmt.Println(err)
			}
		}

	}
	return
}

func CopyFile(source string, dest string) (err error) {
	sourcefile, err := os.Open(source)
	if err != nil {
		return err
	}

	defer sourcefile.Close()

	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer destfile.Close()

	_, err = io.Copy(destfile, sourcefile)
	if err == nil {
		sourceinfo, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, sourceinfo.Mode())
		}

	}

	return
}
