package mq

import (
	"encoding/json"
	"fmt"
	"gitlab.gf.com.cn/hk-common/go-tool/server/logger"
	"testing"
)

// TestSend 测试mq推送
func TestSend(t *testing.T) {
	msg := map[string]interface{}{
		"client_id":     "222222",
		"client_name":   "libo",
		"mode":          "email",
		"monitor_index": "capital_repayment_email_notice",
		"e_mail":        "627392057@qq.com",
	}
	res, err := json.Marshal(msg)
	if err != nil {
		logger.Panic(err.Error())
	}
	err = Send([]string{string(res)}, "public_queue")
	if err != nil {
		t.Fail()
	}
}

func TestSendExchange(t *testing.T) {
	bodys := map[string]interface{}{
		"client_id":     "",
		"monitor_index": "",
		"data_expired":  "",
	}
	err := SendExchange(bodys, "work_order", "#.order.#")
	fmt.Println(err)
}
