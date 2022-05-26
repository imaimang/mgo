package core

import (
	"encoding/json"
	"reflect"
	"strconv"
)

//Byte2Bits Byte2Bits
func Byte2Bits(data byte) []int {
	dst := make([]int, 0)
	for i := 0; i < 8; i++ {
		move := uint(7 - i)
		dst = append(dst, int((data>>move)&1))
	}
	return dst
}

//ConvertToBufer ConvertToBufer
func ConvertToBufer(content interface{}) []byte {
	var contentBuffer []byte
	switch content.(type) {
	case int:
		value, _ := content.(int)
		contentBuffer = []byte(strconv.FormatInt(int64(value), 10))
	case *int:
		value, _ := content.(*int)
		contentBuffer = []byte(strconv.FormatInt(int64(*value), 10))
	case int8:
		value, _ := content.(int8)
		contentBuffer = []byte(strconv.FormatInt(int64(value), 10))
	case *int8:
		value, _ := content.(*int8)
		contentBuffer = []byte(strconv.FormatInt(int64(*value), 10))
	case int16:
		value, _ := content.(int16)
		contentBuffer = []byte(strconv.FormatInt(int64(value), 10))
	case *int16:
		value, _ := content.(*int16)
		contentBuffer = []byte(strconv.FormatInt(int64(*value), 10))
	case int32:
		value, _ := content.(int32)
		contentBuffer = []byte(strconv.FormatInt(int64(value), 10))
	case *int32:
		value, _ := content.(*int32)
		contentBuffer = []byte(strconv.FormatInt(int64(*value), 10))
	case int64:
		value, _ := content.(int64)
		contentBuffer = []byte(strconv.FormatInt(int64(value), 10))
	case *int64:
		value, _ := content.(*int64)
		contentBuffer = []byte(strconv.FormatInt(int64(*value), 10))
	case uint:
		value, _ := content.(uint)
		contentBuffer = []byte(strconv.FormatInt(int64(value), 10))
	case *uint:
		value, _ := content.(*uint)
		contentBuffer = []byte(strconv.FormatInt(int64(*value), 10))
	case uint8:
		value, _ := content.(uint8)
		contentBuffer = []byte(strconv.FormatInt(int64(value), 10))
	case *uint8:
		value, _ := content.(*uint8)
		contentBuffer = []byte(strconv.FormatInt(int64(*value), 10))
	case uint16:
		value, _ := content.(uint16)
		contentBuffer = []byte(strconv.FormatInt(int64(value), 10))
	case *uint16:
		value, _ := content.(*uint16)
		contentBuffer = []byte(strconv.FormatInt(int64(*value), 10))
	case uint32:
		value, _ := content.(uint32)
		contentBuffer = []byte(strconv.FormatInt(int64(value), 10))
	case *uint32:
		value, _ := content.(*uint32)
		contentBuffer = []byte(strconv.FormatInt(int64(*value), 10))
	case uint64:
		value, _ := content.(uint64)
		contentBuffer = []byte(strconv.FormatInt(int64(value), 10))
	case *uint64:
		value, _ := content.(*uint64)
		contentBuffer = []byte(strconv.FormatInt(int64(*value), 10))
	case float32:
		value, _ := content.(float32)
		contentBuffer = []byte(strconv.FormatFloat(float64(value), 'f', -1, 32))
	case *float32:
		value, _ := content.(*float32)
		contentBuffer = []byte(strconv.FormatFloat(float64(*value), 'f', -1, 32))
	case float64:
		value, _ := content.(float64)
		contentBuffer = []byte(strconv.FormatFloat(float64(value), 'f', -1, 64))
	case *float64:
		value, _ := content.(*float64)
		contentBuffer = []byte(strconv.FormatFloat(float64(*value), 'f', -1, 64))
	case string:
		value, _ := content.(string)
		contentBuffer = []byte(value)
	case *string:
		value, _ := content.(*string)
		contentBuffer = []byte(*value)
	case []byte:
		value, _ := content.([]byte)
		contentBuffer = value
	case nil:
	default:
		if reflect.TypeOf(content).Kind() == reflect.Ptr {
			if !reflect.ValueOf(content).IsNil() {
				contentBuffer, _ = json.Marshal(content)
			}
		} else if reflect.TypeOf(content).Kind() == reflect.Slice {
			if reflect.ValueOf(content).Len() == 0 {
				contentBuffer = []byte("[]")
			} else {
				contentBuffer, _ = json.Marshal(content)
			}
		} else {
			contentBuffer, _ = json.Marshal(content)
		}
	}
	return contentBuffer
}
