package httplib

import (
	"net/http"

	"github.com/imaimang/mgo/core"
)

//Response Response
type Response struct {
	http.ResponseWriter
}

//ResponseResult ResponseResult
func (r *Response) ResponseResult(content interface{}, err error) {
	if err == nil {

		switch content.(type) {
		case int, *int, int8, *int8, int16, *int16, int32, *int32, int64, *int64, uint, *uint,
			uint8, *uint8, uint16, *uint16, uint32, *uint32, uint64, *uint64, float32, *float32, float64, *float64, string, *string, nil:
			r.Header().Set("Content-Type", "text/plain; charset=utf-8")
		case []byte:
		default:
			r.Header().Set("Content-Type", "application/json; charset=utf-8")
		}
		r.WriteHeader(200)
		contentBuffer := core.ConvertToBufer(content)
		r.Write(contentBuffer)
	} else {
		r.WriteHeader(400)
		r.Write([]byte(err.Error()))
	}
}
