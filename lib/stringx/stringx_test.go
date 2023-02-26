package stringx

import (
	"reflect"
	"strings"
	"testing"
)

func TestDelete(t *testing.T) {
	type args struct {
		src  []string
		dist string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "删除第一项",
			args: args{
				src:  []string{"1", "2", "3"},
				dist: "1",
			},
			want: []string{"2", "3"},
		},
		{
			name: "删除中间一项",
			args: args{
				src:  []string{"1", "2", "3"},
				dist: "2",
			},
			want: []string{"1", "3"},
		},
		{
			name: "删除最后一项",
			args: args{
				src:  []string{"1", "2", "3"},
				dist: "3",
			},
			want: []string{"1", "2"},
		},
		{
			name: "删除最后一项",
			args: args{
				src:  []string{"1", "2", "3"},
				dist: "3",
			},
			want: []string{"1", "2"},
		},
		{
			name: "删除所有1项",
			args: args{
				src:  []string{"1", "2", "1", "1", "3"},
				dist: "1",
			},
			want: []string{"2", "3"},
		},
		{
			name: "删除不存在数据",
			args: args{
				src:  []string{"1", "2", "1", "1", "3"},
				dist: "6",
			},
			want: []string{"1", "2", "1", "1", "3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Delete(tt.args.src, tt.args.dist); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRetTure(t *testing.T) {
	a := "1"
	if Includes(strings.Split(a, ","), "1") {
		t.Log("true")
	} else {
		t.Log("false")
	}

	b := "0"
	if Includes(strings.Split(b, ","), "1") {
		t.Log("true")
	} else {
		t.Log("false")
	}
}
