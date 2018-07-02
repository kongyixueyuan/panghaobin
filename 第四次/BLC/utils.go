package BLC

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"encoding/json"
	"log"
)

func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		fmt.Println(err)
	}
	return buff.Bytes()
}

// JSON字符串转数组
func JSONToArray(jsonString string) []string {
	//json 到 []string
	var sArr []string
	if err := json.Unmarshal([]byte(jsonString), &sArr); err != nil {
		log.Panic(err)
	}
	return sArr
}