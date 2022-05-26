package httplib

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type HttpProxy struct {
	headers map[string]string
	token   string
	baseUrl string
}

func (h *HttpProxy) GetBaseUrl() string {
	return h.baseUrl
}

func (h *HttpProxy) SetBaseUrl(url string) {
	if !strings.HasSuffix(url, "/") {
		h.baseUrl = url + "/"
	} else {
		h.baseUrl = url
	}
	h.headers = make(map[string]string)
}

func (h *HttpProxy) SetAuthorization(key string) {
	h.token = key
}

func (h *HttpProxy) SetHeader(key, value string) {
	h.headers[key] = value
}

//GetMessage GetMessage
func (h *HttpProxy) GetMessageWidthTimeout(url string, params url.Values, timeout time.Duration) ([]byte, error) {
	address := h.joinUrl(url)
	if params != nil {
		address += "?" + params.Encode()
	}

	req, err := http.NewRequest("GET", address, nil)
	if err == nil {
		for k, v := range h.headers {
			req.Header.Add(k, v)
		}
		if h.token != "" {
			req.Header.Add("Authorization", h.token)
		}

		client := http.Client{
			Timeout: timeout,
		}
		resp, err := client.Do(req)
		if err == nil {
			defer resp.Body.Close()
			datas, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			if resp.StatusCode == http.StatusOK {
				return datas, err
			}
			return nil, errors.New(string(datas))
		}
	}
	return nil, err
}

//GetMessage GetMessage
func (h *HttpProxy) GetMessage(url string, params url.Values) ([]byte, error) {
	address := h.joinUrl(url)
	if params != nil {
		address += "?" + params.Encode()
	}
	var buffer []byte
	req, err := http.NewRequest("GET", address, nil)
	if err == nil {
		for k, v := range h.headers {
			req.Header.Add(k, v)
		}
		if h.token != "" {
			req.Header.Add("Authorization", h.token)
		}

		var resp *http.Response
		resp, err = http.DefaultClient.Do(req)
		if err == nil {
			defer resp.Body.Close()
			buffer, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				err = errors.New(string(buffer))
			}
		}
	}
	return buffer, err
}

//GetMessage GetMessage
func (h *HttpProxy) DeleteMessage(url string, params url.Values) ([]byte, error) {
	address := h.joinUrl(url)
	if params != nil {
		address += "?" + params.Encode()
	}
	var buffer []byte
	req, err := http.NewRequest("DELETE", address, nil)
	if err == nil {
		for k, v := range h.headers {
			req.Header.Add(k, v)
		}
		if h.token != "" {
			req.Header.Add("Authorization", h.token)
		}

		var resp *http.Response
		resp, err = http.DefaultClient.Do(req)
		if err == nil {
			defer resp.Body.Close()
			buffer, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				err = errors.New(string(buffer))
			}
		}
	}
	return buffer, err
}

func (h *HttpProxy) PutMessage(url string, content interface{}) ([]byte, error) {
	var buffer []byte
	req, err := http.NewRequest("PUT", h.joinUrl(url), nil)
	if err == nil {
		req.Header.Add("Authorization", h.token)
		for k, v := range h.headers {
			req.Header.Add(k, v)
		}
		var resp *http.Response
		resp, err = http.DefaultClient.Do(req)
		if err == nil {
			defer resp.Body.Close()
			buffer, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				err = errors.New(string(buffer))
			}
		}
	}
	return buffer, err
}

//PostMessage PostMessage
func (h *HttpProxy) PostMessage(requestURL string, content interface{}) ([]byte, error) {
	return h.PostMessages(requestURL, content, 10*time.Second)
}

//PostMessage PostMessage
func (h *HttpProxy) PostMessages(requestURL string, content interface{}, timeout time.Duration) ([]byte, error) {
	buffer := new(bytes.Buffer)
	contentType := "application/x-www-form-urlencoded"

	switch content := content.(type) {
	case string:
		contentType = "text/html"
		buffer.Write([]byte(content))
	case []byte:
		contentType = "application/octet-stream"
		buffer.Write(content)
	default:
		params, _ := json.Marshal(content)
		contentType = "application/json"
		buffer.Write([]byte(params))
	}
	var bufferResult []byte
	req, err := http.NewRequest("POST", h.joinUrl(requestURL), buffer)
	if err == nil {

		req.Header.Add("Content-Type", contentType)
		req.Header.Add("Authorization", h.token)
		for k, v := range h.headers {
			req.Header.Add(k, v)
		}
		client := http.Client{
			Timeout: timeout,
		}
		var resp *http.Response
		resp, err = client.Do(req)
		if err == nil {
			defer resp.Body.Close()
			bufferResult, err = ioutil.ReadAll(resp.Body)
			if resp.StatusCode != http.StatusOK {
				err = errors.New(string(bufferResult))
			}
		}
	}
	return bufferResult, err
}

func (h *HttpProxy) joinUrl(url string) string {
	if strings.HasPrefix(url, "/") {
		return h.baseUrl + strings.TrimLeft(url, "/")
	}
	return h.baseUrl + url
}

func (h *HttpProxy) DownloadFile(remoteFilename, localFilename string) error {
	resp, err := http.Get(h.baseUrl + remoteFilename)
	fmt.Println("下载文件", h.baseUrl+remoteFilename)
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusOK {
		file, err := os.Create(localFilename)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(file, resp.Body)
		return err
	}
	return errors.New(resp.Status)
}

func (h *HttpProxy) UploadFile(url string, fileName string, datas []byte) error {
	boundary := "ASSDFWDFBFWEFWWDF" //可以自己设定，需要比较复杂的字符串作

	picData := "--" + boundary + "\n"
	picData = picData + "Content-Disposition: form-data; name=\"userfile\"; filename=" + fileName + "\n"
	picData = picData + "Content-Type: application/octet-stream\n\n"
	picData = picData + string(datas) + "\n"
	picData = picData + "--" + boundary + "\n"

	address := h.joinUrl(url)
	req, err := http.NewRequest("POST", address, strings.NewReader(picData))
	req.Header.Set("Content-Type", "multipart/form-data; boundary="+boundary)
	if err == nil {
		_, err := http.DefaultClient.Do(req)
		return err
	}
	return err
}
