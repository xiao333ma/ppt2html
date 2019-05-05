package common

import (
	"strconv"
)

const (
	MaxScoreScope uint64 = 20
	Connector     string = "@@@@@"
)

type ScoreTimeType uint64

const (
	//_ ScoreTimeType = iota
	ScoreTimeDay ScoreTimeType = iota
	ScoreTimeWeek
	ScoreTimeMonth
	ScoreTimeYear
	ScoreTimeAll
)

//分表用 根据特定的微信openid 算出应该划分到哪个表里
func GetOpenidIndexNum(str string) int64 {
	var num int64
	var index int64 = -1
	if len(str) > 0 {
		for _, value := range []byte(str) {
			num += int64(value)
		}

		if num >= 10 {
			index = num % 10
		} else {
			index = num
		}
	}
	return index
}

func GeneratorScoreCacheKeysByScope(scope uint64) []string {

	if scope > MaxScoreScope {
		scope = MaxScoreScope
	}

	blockStr := strconv.Itoa(int(scope))
	return []string{
		strconv.Itoa(int(ScoreTimeDay)) + Connector + GetCurrentDay() + Connector + blockStr,
		strconv.Itoa(int(ScoreTimeWeek)) + Connector + GetCurrentWeek() + Connector + blockStr,
		strconv.Itoa(int(ScoreTimeMonth)) + Connector + GetCurrentMonth() + Connector + blockStr,
		strconv.Itoa(int(ScoreTimeYear)) + Connector + GetCurrentYear() + Connector + blockStr,
		strconv.Itoa(int(ScoreTimeAll)) + Connector + "all" + Connector + blockStr,
	}
}
