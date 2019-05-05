package common

import (
	"github.com/astaxie/beego"
	"strconv"
	"time"
)

func GetCurrentYear() string {
	return strconv.Itoa(time.Now().Year())
}

func GetCurrentMonth() string {
	t := time.Now()
	year := t.Year()
	month := int(t.Month())
	yearStr := strconv.Itoa(year)
	monthStr := strconv.Itoa(month)
	if month <= 9 {
		monthStr = "0" + monthStr
	}
	return yearStr + monthStr
}

func GetCurrentDay() string {
	t := time.Now()
	day := t.Day()
	dayStr := strconv.Itoa(day)
	return GetCurrentMonth() + dayStr
}

func GetCurrentWeek() string {
	t := time.Now()
	year, week := t.ISOWeek()
	yearStr := strconv.Itoa(year)
	weekStr := strconv.Itoa(week)
	if week <= 9 {
		weekStr = "0" + weekStr
	}
	return yearStr + weekStr
}

//得到系统时间
func GetTime() time.Time {
	timezone := float64(0)
	v := beego.AppConfig.String("timezone")
	timezone, _ = strconv.ParseFloat(v, 64)
	add := timezone * float64(time.Hour)
	return time.Now().UTC().Add(time.Duration(add))
}

/*"2006-01-02 15:04:05"*/
//得到今天日期字符串
func GetTodayString() string {
	timezone := float64(0)
	v := beego.AppConfig.String("timezone")
	timezone, _ = strconv.ParseFloat(v, 64)
	add := timezone * float64(time.Hour)
	return time.Now().UTC().Add(time.Duration(add)).Format("20060102")
}

//得到时间字符串
func GetTimeString() string {
	timezone := float64(0)
	v := beego.AppConfig.String("timezone")
	timezone, _ = strconv.ParseFloat(v, 64)
	add := timezone * float64(time.Hour)
	return time.Now().UTC().Add(time.Duration(add)).Format("20060102150405")
}
