package pgs

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

const (
	sep = ","
)

type Strings []string

func (s Strings) MarshalText() (text []byte, err error) {
	return []byte(s.ToString()), nil
}

func (s Strings) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	return strings.Join(s, sep), nil
}

func (s *Strings) Scan(src interface{}) error {
	if bytes, ok := src.(string); ok {
		ss := make([]string, 0)
		if bytes != "" {
			ss = strings.Split(bytes, sep)
		}
		s.setValue(ss)
	} else {
		return fmt.Errorf("数据错误data=%s", string(bytes))
	}
	return nil
}

func (s *Strings) setValue(data []string) {
	*s = data
}

func (s Strings) ToString() string {
	return strings.Join(s, sep)
}
