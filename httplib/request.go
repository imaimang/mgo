package httplib

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

//Request Request
type Request struct {
	*http.Request
	AuthInfo interface{}
}

//GetJSONValue GetJSONValue
func (r *Request) GetJSONValue(entity interface{}) error {
	var datas []byte
	var err error
	if datas, err = ioutil.ReadAll(r.Body); err != nil {
		return err
	}
	r.Body.Close()
	if err = json.Unmarshal(datas, entity); err != nil {
		return err
	}
	return nil
}

func (r *Request) GetToken() string {
	return r.Header.Get("Authorization")
}

func (r *Request) getFormValue(key string) string {
	r.ParseForm()
	for keyItem, value := range r.Form {
		if strings.ToLower(keyItem) == strings.ToLower(key) {
			return value[0]
		}
	}
	return ""
}

//GetFormInt GetFormInt
func (r *Request) GetFormInt(key string) int {
	result, err := strconv.Atoi(r.getFormValue(key))
	if err == nil {
		return result
	}
	return 0
}

//GetFormInt GetFormInt
func (r *Request) GetFormBool(key string) bool {
	result, err := strconv.ParseBool(r.getFormValue(key))
	if err == nil {
		return result
	}
	return false
}

//GetFormString GetFormString
func (r *Request) GetFormString(key string) string {
	return r.getFormValue(key)
}

//GetFormIntSlice GetFormIntSlice
func (r *Request) GetFormJsonIntSlice(key string) ([]int, error) {
	var result []int
	err := json.Unmarshal([]byte(r.getFormValue(key)), &result)
	return result, err
}

//GetFormIntSlice GetFormIntSlice
func (r *Request) GetFormIntSlice(key string) []int {
	var result []int
	items := strings.Split(r.getFormValue(key), ",")
	for _, item := range items {
		itemVar, err := strconv.Atoi(item)
		if err == nil {
			result = append(result, itemVar)
		}
	}
	return result
}

//GetFormStringSlice GetFormStringSlice
func (r *Request) GetFormStringSlice(key string) []string {
	value := r.getFormValue(key)
	if value == "" {
		return nil
	}
	result := strings.Split(r.getFormValue(key), ",")
	return result
}

//GetBodyBuffer GetBodyBuffer
func (r *Request) GetBodyBuffer() []byte {
	buffer, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()
	return buffer
}
