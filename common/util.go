package common

import (
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/utils"
	"os"
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

func SetStructValue(obj interface{},fieldName string,value string)  {
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return
	}
	fieldNum := t.NumField()
	for i := 0; i < fieldNum; i++ {
		if t.Field(i).Name == fieldName  {
			p := reflect.ValueOf(obj)
			if p.Kind() == reflect.Ptr{
				field := p.Elem().FieldByName(fieldName)
				if field.Kind() == reflect.String {
					field.SetString(value)
				}
			}
			break
		}
	}
}
