package BLC

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		fmt.Println(err)
	}
	return buff.Bytes()
}
