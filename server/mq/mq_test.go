package mq

import (
	"encoding/json"
	"fmt"
	logger "gitlab.gf.com.cn/hk-common/go-tool/server/logger"
	"testing"
)

// TestSend 测试mq推送
func TestSend(t *testing.T) {
	// filer := file.NewIOFiler()
	// content, err := filer.ReadAll2ByteSlice("../../CPT汇总信息.xlsx")
	// if err != nil {
	// 	t.Fatal(err.Error())
	// }
	msg := map[string]interface{}{
		"client_id":     "222222",
		"client_name":   "libo",
		"mode":          "email",
		"monitor_index": "capital_repayment_email_notice",
		"e_mail":        "627392057@qq.com",
		// "attachments": []map[string]interface{}{
		// 	map[string]interface{}{
		// 		"filename": "testFile.xlsx",
		// 		"content":  content,
		// 	},
		// },
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
