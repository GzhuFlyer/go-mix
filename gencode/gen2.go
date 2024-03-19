package gencode

import (
	"bytes"
	_ "embed"
	"fmt"
	goformat "go/format"
	"os"
	"text/template"
)

var teplCase = `	case ops%s:
	req := &client.%sRequest{}
	var payloadByte []byte
	payloadByte, err := json.Marshal(p.MsgData.Data.Payload)
	if err != nil {
		logger.Error("josn marshal error", err)
		return encodeMessage(errJSONTypeMsg)
	}
	err = json.Unmarshal(payloadByte, req)
	if err != nil {
		logger.Error("josn Unmarshal error", err)
		return encodeMessage(errJSONTypeMsg)
	}

	resp := &client.%sResponse{}

	err = handler.NewDeviceCenter().%s(context.Background(), req, resp)
	logger.Debug("resp = ", resp)
	if err != nil {
		errMsg := errorMsg(offline, p.MsgData.Data.OpsCode, p.Tid, p.Bid, p.Sn, err.Error())
		logger.Error("errMsg = ", errMsg)
		encodeMsg := encodeMessage(errMsg)
		logger.Error("encodeMsg = ", encodeMsg)
		return encodeMsg
	}
	smsg := successfulMsg(success, p.MsgData.Data.OpsCode, p.Tid, p.Bid, p.Sn, successMsg, resp)
	encodeOkmsg = encodeMessage(smsg)
`

var teplConst = `
	ops%s = %d;
`

//go:embed mqtt.tmpl
var codeA string

var CmdSlice = []string{
	"RadarGetBeamConfig",
	"RadarStartDetect",
	"RadarEndDetect",
	"RadarSetBeamSchedule",
	"RadarPostureCalibration",
	"RadarPostureCalibrationManual",
	"RadarGetVersionInfo",
}
var deviceName = "agx"

func Gen2() {
	// data := make(map[string]interface{})
	// var outString string
	var cs string
	var cconst string
	for k, name := range CmdSlice {
		// data := map[string]interface{}{
		// 	"CmdName": name,
		// }
		c := fmt.Sprintf(teplCase, name, name, name, name)
		cc := fmt.Sprintf(teplConst, name, k)
		// fmt.Println(c)
		cs = cs + c
		cconst = cconst + cc
		// tmpl, _ := template.New("hello").Parse(teplCase)
		// tmpl.Execute(os.Stdout, data)
		// fmt.Println()
	}
	// fmt.Println(cs)

	t := template.Must(template.New("agx").Parse(codeA))
	buffer := new(bytes.Buffer)
	data2 := map[string]interface{}{
		"CaseEvent":   cs,
		"LogicName":   "agx",
		"constDefine": cconst,
	}
	err := t.Execute(buffer, data2)
	if err != nil {
		panic(err)
	}

	code, err := goformat.Source([]byte(buffer.String()))
	if err != nil {
		// panic(err)
		code = buffer.Bytes()
	}
	fp, _ := os.Create("agx.go")
	// code := golang.FormatCode(buffer.String())
	fmt.Println(code)
	_, err = fp.WriteString(string(code))
	if err != nil {
		panic(err)
	}
	// file, _ := os.Create("agx.go")
	// tmp2, _ := template.New("hello").Parse(codeA)
	// data2 := map[string]interface{}{
	// 	"CaseEvent":   cs,
	// 	"LogicName":   "agx",
	// 	"constDefine": cconst,
	// }
	// tmp2.Execute(file, data2)

}
