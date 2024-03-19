package mqtt

import (
	"encoding/json"
	"fmt"
	"time"
)

// type MessageTwo struct {
// 	// ClientID string `json:"clientId"`
// 	// Type     string `json:"type"`
// 	// Data     string `json:"data,omitempty"`
// 	// Time     int64  `json:"time"`
// 	Msg string `json:"msg"`
// }

//	func decodeMessage(payload []byte) (*MessageTwo, error) {
//		message := new(MessageTwo)
//		decoder := json.NewDecoder(strings.NewReader(string(payload)))
//		decoder.UseNumber()
//		if err := decoder.Decode(&message); err != nil {
//			return nil, err
//		}
//		return message, nil
//	}
func test1(m Message) []byte {
	fmt.Println("m = ", m.Payload())
	p, _ := decodeMessage(m.Payload())
	fmt.Println(p)
	return m.Payload()
}

// var gtopic = "/hello"
var gtopic = "/mqtt/v1/cloud/sfl/1727884703555842048/services/get_replay"

// RegisterSubHandlers ...
func RegisterSubHandlers(adaptor *Adaptor) {
	if adaptor == nil {

	}
	consumes := []MqttRoute{
		{
			Topic:   gtopic,
			Handler: test1,
		},
	}
	adaptor.AddSubHandlers(consumes)
}

func HandlerSub(s *Adaptor) {

	RegisterSubHandlers(s)

}

func HandlerPub(s *Adaptor) {
	ticker := time.NewTicker(1 * time.Second)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				nowtime := time.Now().Format("2006-01-02 15:04:05")
				msg := "{\"msg\":\"good morning :" + nowtime + "\"}"
				fmt.Println(msg)
				s.Publish(gtopic, []byte(msg))
			}
		}
	}()
	time.Sleep(500000 * time.Second)
	ticker.Stop()
	done <- true
	fmt.Println("stop publish")

}

type data struct {
	OpsCode int32       `json:"opsCode"`
	Payload interface{} `json:"payload"`
}

type msgSubData struct {
	Data data `json:"data"`
}

type result struct {
	ErrorCode int32  `json:"errorCode"`
	ErrorMsg  string `json:"errorMsg"`
}

type msgpubData struct {
	Data   data   `json:"data"`
	Result result `json:"result"`
}

type mqttSubMsg struct {
	Tid       string     `json:"tid"`
	Bid       string     `json:"bid"`
	Sn        string     `json:"sn"`
	MsgData   msgSubData `json:"msgData"`
	Timestamp int64      `json:"timestamp"`
}

type mqttPubMsg struct {
	Tid       string     `json:"tid"`
	Bid       string     `json:"bid"`
	Sn        string     `json:"sn"`
	MsgData   msgpubData `json:"msgData"`
	Timestamp int64      `json:"timestamp"`
}

// func decodeMessage(payload []byte) (*mqttSubMsg, error) {
// 	message := new(mqttSubMsg)
// 	decoder := json.NewDecoder(strings.NewReader(string(payload)))
// 	// decoder.UseNumber()
// 	if err := decoder.Decode(&message); err != nil {
// 		return nil, err
// 	}
// 	return message, nil
// }

func decodeMessage(payload []byte) (*mqttSubMsg, error) {
	var msg mqttSubMsg
	err := json.Unmarshal(payload, &msg)
	if err != nil {
		fmt.Println("json.Unmarshal err = ", err)
		return nil, err
	}
	fmt.Println("msg = ", msg)
	return &msg, nil
}
func encodeMessage(msg mqttPubMsg) []byte {
	fmt.Println("encodeMessage = ", msg)
	bmsg, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("josn marshal error", err)
	}
	return bmsg
}
func errorMsg(errorCode, opsCode int32, tid, bid, sn, errorMsg string) mqttPubMsg {
	msg := mqttPubMsg{
		Tid: tid,
		Bid: bid,
		Sn:  sn,
		MsgData: msgpubData{
			Result: result{
				ErrorCode: errorCode,
				ErrorMsg:  errorMsg,
			},
			Data: data{
				OpsCode: opsCode,
			},
		},
		Timestamp: time.Now().UnixMilli(),
	}
	fmt.Println("msg = ", msg)
	return msg
}

const (
	fail        = -1
	success     = 0
	offline     = 10001
	jsonFailMsg = "request param invalid"
	successMsg  = "operation successful"
)

func TestCode() {
	// errMsg := errorMsg(offline, 1, "p.Tid", " p.Bid", "p.Sn", "hello")
	// fmt.Println("errMsg = ", errMsg)

	errMsg := mqttPubMsg{
		Tid: "tid",
		Bid: "bid",
		Sn:  "sn",
		MsgData: msgpubData{
			Result: result{
				ErrorCode: 1,
				ErrorMsg:  "errorMsg",
			},
			Data: data{
				OpsCode: 1,
			},
		},
		Timestamp: time.Now().UnixMilli(),
	}

	encodeMsg := encodeMessage(errMsg)
	fmt.Println("result = ", errMsg.MsgData.Result)
	fmt.Println("encodeMsg = ", encodeMsg)

	var msg mqttPubMsg
	err := json.Unmarshal(encodeMsg, &msg)
	fmt.Println("err = ", err)
	fmt.Println("msg = ", msg)
	fmt.Println("OpsCode = ", msg.MsgData.Data.OpsCode)
	fmt.Println("Result = ", msg.MsgData.Result)
	fmt.Println("r1 = ", msg.MsgData.Result.ErrorCode)
	fmt.Println("r 2 = ", msg.MsgData.Result.ErrorMsg)

}
