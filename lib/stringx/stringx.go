package stringx

import (
	"encoding/json"
	"fmt"
	"gitlab.gf.com.cn/hk-common/go-tool/server/logger"
	"strconv"
	"strings"
)

const (
	COMMA     = ","
	SEMICOLON = ";"
)

type TempParam map[string]string

func Template(s string, param TempParam) string {
	for k, v := range param {
		s = strings.ReplaceAll(s, fmt.Sprintf("{{%s}}", k), v)
	}
	return s
}

func ParseInt64(s string) int64 {
	n, _ := strconv.ParseInt(s, 10, 64)
	return n
}

func ParseFloat(s string) float64 {
	n, _ := strconv.ParseFloat(s, 64)
	return n
}

func Atoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

func Judge(b bool, v1, v2 string) string {
	if b {
		return v1
	}
	return v2
}

func ToStrings(s, split string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(s, split)
}

func Snake(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		// or通过ASCII码进行大小写的转化
		// 65-90（A-Z），97-122（a-z）
		//判断如果字母为大写的A-Z就在前面拼接一个_
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	//ToLower把大写字母统一转小写
	return strings.ToLower(string(data[:]))
}

func Includes(src []string, dist string) bool {
	for _, s := range src {
		if s == dist {
			return true
		}
	}
	return false
}

func Delete(src []string, dist string) []string {
	j := 0
	for _, v := range src {
		if v != dist {
			src[j] = v
			j++
		}
	}
	return src[:j]
}

func Marshal(src interface{}) string {
	bt, err := json.Marshal(src)
	if err != nil {
		logger.Error(err)
	}
	return string(bt)
}

func ToArray1(src string) []string {
	return strings.Split(src, COMMA)
}

func ToArray2(src string) []string {
	return strings.Split(src, SEMICOLON)
}
