package lib

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/lristar/go-tool/config"
	configs "github.com/lristar/go-tool/config"
	myerror "github.com/lristar/go-tool/lib/error"
	"log"
	"math"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

// Md5Sum Md5Sum
func Md5Sum(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	md5Str := hex.EncodeToString(m.Sum(nil))
	return md5Str
}

// RandMD5Str RandStr
func RandMD5Str() string {
	return Md5Sum(strconv.FormatInt(int64(math.Floor(rand.Float64()*10000))+time.Now().Unix(), 10))
}

// GetRandomString 生成随机字符串
func GetRandomString(n int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < n; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// GetPageCount 获得page count
func GetPageCount(page int, count int) (int, int) {
	if page == 0 {
		page = 1
	}
	if count == 0 {
		count = 10
	}
	limit := count
	realPage := page
	cursor := 0
	if realPage > 0 {
		cursor = (realPage - 1) * limit
	}
	return cursor, limit
}

// GetOrcalePageStartEnd 获取Oracle分页数据
func GetOrcalePageStartEnd(page, count int) (start, end int) {
	if page == 0 {
		page = 1
	}
	if count == 0 {
		count = 10
	}
	start = (page-1)*count + 1
	end = start + count - 1
	return
}

// Struct2Map struct转map
func Struct2Map(obj interface{}) map[string]interface{} {
	bytes, _ := json.Marshal(obj)
	var res map[string]interface{}
	json.Unmarshal(bytes, &res)
	return res
}

// Map2Struct map转struct data:map源数据，obj:struct指针
func Map2Struct(data interface{}, obj interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, obj)
}

// ToThousands ToThousands
func ToThousands(num float64) string {
	strNum := strconv.FormatFloat(num, 'f', -1, 64)
	numArr := strings.Split(strNum, ".")
	var result, counter = "", 0
	TransfNum := numArr[0]
	fixed := ""
	if string([]byte(TransfNum)[0]) == "-" {
		fixed = "-"
		TransfNum = string([]byte(TransfNum)[1:])
	}
	for i := len(TransfNum) - 1; i >= 0; i-- {
		counter++
		result = string([]byte(TransfNum)[i]) + result
		if counter%3 == 0 && i != 0 {
			result = "," + result
		}
	}
	if len(numArr) == 2 {
		result = result + "." + numArr[1]
	}
	return fixed + result
}

// CheckPassWord 检验密码
func CheckPassWord(pass string) bool {
	t1, t2, t3 := `[[:punct:]]`, `[a-z]`, `[A-Z]`
	m1, _ := regexp.MatchString(t1, pass)
	m2, _ := regexp.MatchString(t2, pass)
	m3, _ := regexp.MatchString(t3, pass)
	return m1 && m2 && m3
}

// SContain 判断数组是否包含value
func SContain(arr []string, value string) int {
	for i, v := range arr {
		if v == value {
			return i
		}
	}

	return -1
}

// SArrContain 判断数组是否包含其他数组
func SArrContain(arr []string, subArr []string) bool {
	for _, v := range subArr {
		if i := SContain(arr, v); i == -1 {
			return false
		}
	}
	return true
}

// SArrDelete 删除第一个值为value的元素
func SArrDelete(arr []string, value string) []string {
	i := SContain(arr, value)
	if i != -1 {
		arr = append(arr[:i], arr[i+1:]...)
	}
	return arr
}

// toFloat64 toFloat64
func toFloat64(value interface{}) float64 {
	defer func() {
		if err := recover(); err != nil {
			panic(myerror.Errorf("无法将%+v(%T)转换成float64类型", 406, value, value))
		}
	}()
	res := value.(float64)
	return res
}

// Panic 抛出model层错误
func Panic(e error) {
	if e != nil {
		panic(myerror.Errorf(e.Error(), 405))
	}
}

// IF 三元运算符替代方法
func IF(condition bool, trueValue, falseValue interface{}) interface{} {
	if condition {
		return trueValue
	}
	return falseValue
}

// Contains 在map/struct数组中查找对应条件第一次出现的下标，不存在返回-1
func Contains(source interface{}, condition configs.M) (int, error) {
	sValue := reflect.ValueOf(source)
	if sValue.Kind() == reflect.Ptr {
		sValue = sValue.Elem()
	}

	if sValue.Kind() != reflect.Slice {
		return -1, fmt.Errorf("source必须是MAP/STRUCT类型的Slice")
	}
	// soruceArr := source.([]interface{})
	for i := 0; i < sValue.Len(); i++ {
		item := sValue.Index(i)

		sMap := make(map[string]interface{})
		if item.Kind() == reflect.Map {
			sMap = (item.Interface()).(config.M)
		} else if item.Kind() == reflect.Struct {
			itemType := item.Type()
			for i := 0; i < itemType.NumField(); i++ {
				iType := itemType.Field(i)
				ikey := iType.Tag.Get("json")
				sMap[ikey] = item.Field(i).Interface()
			}
		} else if m, ok := (item.Interface()).(configs.M); ok {
			sMap = m
		} else {
			return -1, fmt.Errorf("source必须是MAP/STRUCT类型的Slice")
		}

		flag := true
		for key, value := range condition {
			// fmt.Printf("%T %v %s %t\n", sMap[key], sMap[key], value, value)
			if fmt.Sprintf("%v", sMap[key]) != fmt.Sprintf("%v", value) {
				flag = false
				break
			}
		}
		if flag {
			return i, nil
		}
	}
	return -1, nil
}

// CurrencyFormat 数字转换为货币格式显示（千分位加逗号）
func CurrencyFormat(s interface{}) string {
	var floats float64
	switch v := reflect.ValueOf(s); v.Kind() {
	case reflect.Float32, reflect.Float64:
		floats = v.Float()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		floats = float64(v.Int())
	default:
		floats = 0
	}

	str := strconv.FormatFloat(floats, 'f', -1, 64)
	numberArr := strings.Split(str, ".")
	if len(numberArr) > 2 {
		return str
	}

	length := len(numberArr[0])
	if length < 4 {
		return str
	}

	count := (length - 1) / 3
	for i := 0; i < count; i++ {
		numberArr[0] = numberArr[0][:length-(i+1)*3] + "," + numberArr[0][length-(i+1)*3:]
	}

	return strings.Join(numberArr, ".")
}

// NumberCurrencyFormat 数字转货币格式
func NumberCurrencyFormat(n float64) string {
	str := strconv.FormatFloat(n, 'f', -1, 64)
	strArr := strings.Split(str, ".")
	p := strArr[0]
	e := ""
	if len(strArr) == 2 {
		e = "." + strArr[1]
	}
	p = ReverseString(p)
	r := ""
	for i := 0; i < len(p); i++ {
		if i%3 == 0 && i != 0 {
			r = fmt.Sprintf("%s,%s", string(p[i]), r)
		} else {
			r = fmt.Sprintf("%s%s", string(p[i]), r)
		}
	}
	// println(r + e)
	return r + e
}

// ReverseString 倒置字符串
func ReverseString(s string) string {
	runes := []rune(s)
	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}
	return string(runes)
}

func JsonToField(jsonStr string) string {
	s := []rune(jsonStr)
	var resArr []string
	nextToUpper := false
	for k, v := range s {
		vStr := string(v)
		if k == 0 {
			resArr = append(resArr, strings.ToUpper(vStr))
			continue
		}
		if vStr == "_" {
			nextToUpper = true
			continue
		}
		if nextToUpper {
			vStr = strings.ToUpper(vStr)
			nextToUpper = false
		}
		resArr = append(resArr, vStr)
	}
	return strings.Join(resArr, "")
}

func Split(str, seq string) []string {
	if strings.TrimSpace(str) == "" {
		return make([]string, 0)
	}
	return strings.Split(str, seq)
}

//IsBlank 判断字符串是否为空
func IsBlank(str string) bool {
	return strings.TrimSpace(str) == ""
}

//UniqueStringSlice string数组去重
func UniqueStringSlice(languages []string) []string {
	result := make([]string, 0, len(languages))
	temp := map[string]struct{}{}
	for _, item := range languages {
		if _, ok := temp[item]; !ok { //如果字典中找不到元素，ok=false，!ok为true，就往切片中append元素。
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

//TrimEqualFold 去空格排除大小写对比
func TrimEqualFold(s, t string) bool {
	return strings.EqualFold(strings.ReplaceAll(s, " ", ""), strings.ReplaceAll(t, " ", ""))
}

// StringToByte string 转[]byte
func StringToByte(s string) []byte {
	t := (*[2]uintptr)(unsafe.Pointer(&s))
	t2 := [3]uintptr{t[0], t[1], t[1]}
	return *(*[]byte)(unsafe.Pointer(&t2))
}

// ByteToString []byte 转string
func ByteToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// FloatToFormatBalance float64 转换成千位分隔符格式string类型
func FloatToFormatBalance(v float64) string {
	buf := &bytes.Buffer{}
	if v < 0 {
		buf.Write([]byte{'-'})
		v = 0 - v
	}

	comma := []byte{','}

	parts := strings.Split(strconv.FormatFloat(v, 'f', 2, 64), ".")
	pos := 0
	if len(parts[0])%3 != 0 {
		pos += len(parts[0]) % 3
		buf.WriteString(parts[0][:pos])
		buf.Write(comma)
	}

	for ; pos < len(parts[0]); pos += 3 {
		buf.WriteString(parts[0][pos : pos+3])
		buf.Write(comma)
	}
	buf.Truncate(buf.Len() - 1)

	if len(parts) > 1 {
		buf.Write([]byte{'.'})
		buf.WriteString(parts[1])
	}
	return buf.String()
}

// Int64SliceSort int64切片快排
func Int64SliceSort(nums []int64) {
	quickSort(nums, 0, len(nums)-1)
}

func quickSort(nums []int64, l, r int) {
	if l < r {
		m := partition(nums, l, r)
		quickSort(nums, l, m-1)
		quickSort(nums, m+1, r)
	}
}

func partition(nums []int64, l int, r int) int {
	key := nums[r]
	i := l
	j := l
	for j < r {
		if nums[j] < key {
			nums[i], nums[j] = nums[j], nums[i]
			i++
		}
		j++
	}
	nums[i], nums[r] = nums[r], nums[i]
	return i
}

func GetProjectRootPath() string {
	dir := getCurrentAbPathByExecutable()
	if strings.Contains(dir, getTmpDir()) {
		return getCurrentAbPathByCaller()
	}
	return dir
}

// 获取系统临时目录，兼容go run
func getTmpDir() string {
	dir := os.Getenv("TEMP")
	if dir == "" {
		dir = os.Getenv("TMP")
	}
	res, _ := filepath.EvalSymlinks(dir)
	return res
}

// 获取当前执行文件绝对路径 go build
func getCurrentAbPathByExecutable() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	return res
}

// 获取当前执行文件绝对路径（go run）
func getCurrentAbPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}

func IsGoRun() bool {
	dir := getCurrentAbPathByExecutable()
	if strings.Contains(dir, getTmpDir()) {
		return true
	}
	return false
}

func Equal(a, b interface{}) bool {
	m1, _ := json.Marshal(a)
	m2, _ := json.Marshal(b)
	if bytes.Equal(m1, m2) {
		return true
	}
	return false
}

func IsNil(i interface{}) bool {
	of := reflect.ValueOf(i)
	if of.Kind() == reflect.Ptr {
		return of.IsNil()
	}
	return false
}

func StringToInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func IsNull(i interface{}) bool {
	if i == nil {
		return true
	}
	of := reflect.ValueOf(i)
	switch of.Kind() {
	case reflect.Array, reflect.Slice:
		return of.Len() < 1
	default:
	}
	return of.IsNil()
}

func TimeToInt(date string) int {
	v, err := strconv.Atoi(date)
	if err != nil {
		return 20991231
	}
	return v
}

//GetIntersectStringSlice 求string数据的并集
func GetIntersectStringSlice(a []string, b []string) []string {
	temp := map[string]struct{}{}
	res := make([]string, 0)
	for _, v := range a {
		temp[v] = struct{}{}
	}
	for _, v := range b {
		if _, ok := temp[v]; ok {
			res = append(res, v)
		}
	}
	return res
}

//GetDiffStringSlice 求string数据的差集
func GetDiffStringSlice(a []string, b []string) []string {
	temp := map[string]struct{}{}
	res := make([]string, 0)
	intersect := GetIntersectStringSlice(a, b)
	for _, v := range intersect {
		temp[v] = struct{}{}
	}

	for _, v := range a {
		if _, ok := temp[v]; !ok {
			res = append(res, v)
		}
	}
	for _, v := range b {
		if _, ok := temp[v]; !ok {
			res = append(res, v)
		}
	}
	return res
}
