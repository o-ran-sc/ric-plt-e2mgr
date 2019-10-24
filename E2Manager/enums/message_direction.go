package enums

import (
	"encoding/json"
	"strconv"
)

type MessageDirection int32

var messageDirectionEnumName = map[int32]string{
	0: "UNKNOWN_MESSAGE_DIRECTION",
	1: "RAN_TO_RIC",
	2: "RIC_TO_RAN",
}

const (
	UNKNOWN_MESSAGE_DIRECTION MessageDirection = 0
	RAN_TO_RIC                MessageDirection = 1
	RIC_TO_RAN                MessageDirection = 2
)

func (md MessageDirection) String() string {
	s, ok := messageDirectionEnumName[int32(md)]
	if ok {
		return s
	}
	return strconv.Itoa(int(md))
}

func (md MessageDirection) MarshalJSON() ([]byte, error) {
	_, ok := messageDirectionEnumName[int32(md)]

	if !ok {
		return nil,&json.UnsupportedValueError{}
	}

	v:= int32(md)
	return json.Marshal(v)
}
