package common

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"net/url"
	"strconv"
	"strings"
)

//字符串base64加密
func Base64E(urlstring string) string {
	str := []byte(urlstring)
	data := base64.StdEncoding.EncodeToString(str)
	return data
}

//字符串base64解密
func Base64D(urlxxstring string) string {
	data, err := base64.StdEncoding.DecodeString(urlxxstring)
	if err != nil {
		return ""
	}
	//s := fmt.Sprintf("%q", data)
	//s = strings.Replace(s, "\"", "", -1)
	return string(data)
}

//url转义
func UrlE(s string) string {
	return url.QueryEscape(s)
}

//url解义
func UrlD(s string) string {
	s, e := url.QueryUnescape(s)
	if e != nil {
		return e.Error()
	} else {
		return s
	}
}

//字符串是否在字符串数组中
func InArray(sa []string, a string) bool {
	for _, v := range sa {
		if a == v {
			return true
		}
	}
	return false
}

//create md5 string
func Strtomd5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	rs := hex.EncodeToString(h.Sum(nil))
	return rs
}

//password hash function
func Pwdhash(str string) string {
	return Strtomd5(str)
}

func Md5(str string) string {
	return Strtomd5(str)
}

func StringsToJson(str string) string {
	rs := []rune(str)
	jsons := ""
	for _, r := range rs {
		rint := int(r)
		if rint < 128 {
			jsons += string(r)
		} else {
			jsons += "\\u" + strconv.FormatInt(int64(rint), 16) // json
		}
	}

	return jsons
}

func Rawurlencode(str string) string {
	return strings.Replace(url.QueryEscape(str), "+", "%20", -1)
}

// 分80个文件夹
func Hashcode(asin string) string {
	dd := []byte(Md5("iloveyou"+asin+"hunterhug") + Md5(asin))
	sum := 0
	for _, i := range dd {
		sum = sum + int(i)
	}
	hashcode := sum % (80)
	s := strconv.FormatInt(int64(hashcode), 10)
	if s == "" {
		s = "xx"
	}
	return s
}

func TripAll(a string) string {
	a = strings.Replace(a, " ", "", -1)
	a = strings.Replace(a, "\n", "", -1)
	a = strings.Replace(a, "\r", "", -1)
	a = strings.Replace(a, "\t", "", -1)
	return a
}

func SubString(str string, start int, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		panic("start is wrong")
	}

	if end < 0 || end > length {
		panic("end is wrong")
	}
	return string(rs[start:end])
}
