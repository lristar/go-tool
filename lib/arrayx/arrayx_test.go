package arrayx

import (
	"fmt"
	"testing"
)

func TestArrayX(t *testing.T) {
	a := ArrayX[string]{"adsf", "dafdasf", "dfafd"}
	a = a.Append("dafdasfds")
	fmt.Print(a)
}

func TestIsContain(t *testing.T) {
	a := ArrayX[int]{1, 3, 5, 7}
	fmt.Println(a.IsContain(7))
}
