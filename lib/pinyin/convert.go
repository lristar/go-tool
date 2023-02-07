package pinyin

import (
	"strings"
)

func ConvertToUpper(str string) string {
	a := NewArgs()
	a.Separator = " "
	a.Heteronym = false
	a.Fallback = func(r rune, a *Args) []string {
		return []string{string(r)}
	}
	pys := []string{}
	for _, v := range Pinyin(str, &a) {
		for _, s := range v {
			pys = append(pys, s)
		}
	}
	return strings.ToUpper(strings.Join(pys, a.Separator))
}
