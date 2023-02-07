package myredis

import (
	"fmt"
	"testing"
	"time"
)

func TestSortQueue(t *testing.T) {
	q := NewSortQueue("test_queue", 10*time.Second)
	go func() {
		for i := 0; i < 10; i++ {
			q.Publish(fmt.Sprintf("%d", i), float64(time.Now().Add(10*time.Second).Unix()))
		}
	}()
	q.Consume(func(msg string) error {
		fmt.Println(msg)
		return nil
	})
}
