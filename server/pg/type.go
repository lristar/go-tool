package pgs

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

type Strings []string

func (s Strings) Value() (driver.Value, error) {
	//ss := *s
	return strings.Join(s, ","), nil
}

func (s Strings) Scan(src interface{}) error {
	var bytes []byte
	var err error
	switch v := src.(type) {
	case []byte:
		ss := strings.Split(string(v), ",")
		s.setValue(ss)
	case string:
		ss := strings.Split(v, ",")
		s.setValue(ss)
	default:
		bytes, err = json.Marshal(src)
		if err != nil {
			return fmt.Errorf("数据错误data=%s", string(bytes))
		}
	}
	return nil
}

func (s Strings) setValue(data []string) {
	bt, _ := json.Marshal(data)
	json.Unmarshal(bt, s)
}

func (s Strings) ToString() string {
	return strings.Join(s, ",")
}
