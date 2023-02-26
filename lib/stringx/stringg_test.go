package stringx

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"
)

func TestToInteger(t *testing.T) {
	a, err := ToInteger[string, int64]("-12346513654564561")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(a)
}

func TestTrim(t *testing.T) {
	c, _ := regexp.Compile("^(0\\.?[0-9]*|[1-9][0-9]*\\.?[0-9]*|-[1-9][0-9]*\\.?[0-9]*)$")
	fmt.Println(c.MatchString("131234..000"))
}

func TestGetStrings(t *testing.T) {
	s := GetStrings()
	s.Add("afdaf", "dfafdaf", "fadfadf")
	fmt.Println(s.ToString())
}

func BenchmarkStrings_Add(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		s := []string{}
		for j := 0; j < 10; j++ {
			s = append(s, strconv.Itoa(j))
		}
		b.StartTimer()
		ss := GetStrings()
		ss.Add(s...)
	}
}

func BenchmarkStrings_Add2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		s := []string{}
		for j := 0; j < 10; j++ {
			s = append(s, strconv.Itoa(j))
		}
		b.StartTimer()
		ss := GetStrings()
		ss.Add2(s...)
	}
}
