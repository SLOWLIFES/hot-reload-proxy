package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"
)

type HttpDo struct {
	Url      string
	Data     map[string]string
	Raw      []byte
	FileKey  string
	FileName string
	File     io.Reader
}

func (that HttpDo) Request(method string) ([]byte, error) {
	var r http.Request
	r.ParseForm()
	for e := range that.Data {
		r.Form.Add(e, that.Data[e])
	}

	bodystr := strings.TrimSpace(r.Form.Encode())
	request, err := http.NewRequest(method, that.Url, strings.NewReader(bodystr))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Connection", "Keep-Alive")

	var resp *http.Response
	resp, err = http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (that HttpDo) PostJson() ([]byte, error) {
	body := &bytes.Buffer{}
	body.Write(that.Raw)
	request, err := http.NewRequest("POST", that.Url, body)

	request.Header.Set("Content-Type", "application/json")

	if err != nil {
		return nil, err
	}
	var resp *http.Response
	resp, err = http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	//log.Println(that.Url,string(data))
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (that HttpDo) Post() ([]byte, error) {

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if that.FileKey != "" {
		part, err := writer.CreateFormFile(that.FileKey, that.FileName)
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(part, that.File)
	}

	for e := range that.Data {
		_ = writer.WriteField(e, that.Data[e])
	}
	err := writer.Close()

	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", that.Url, body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())

	var resp *http.Response
	resp, err = http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (that HttpDo) Get() ([]byte, error) {
	return that.Request("GET")
}
