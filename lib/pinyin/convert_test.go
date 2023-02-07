package pinyin

import (
	"testing"
)

func TestConvert(t *testing.T) {
	t.Run("中文转换", func(t *testing.T) {
		str := "北京市北京路11233号"
		t.Log(ConvertToUpper(str))
	})
}
