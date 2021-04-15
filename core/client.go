package core

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"
)

// Get 向微信支付发送一个http get请求
func Get(ctx context.Context, hc *http.Client, credential Credential, validator Validator, requestURL string) ([]byte, error) {
	return DoRequest(ctx, hc, credential, validator, http.MethodGet, requestURL, ApplicationJSON, "", "")
}

// Post 向微信支付发送一个http post请求
func Post(ctx context.Context, hc *http.Client, credential Credential, validator Validator, requestURL string, requestBody interface{}) ([]byte, error) {
	return do(ctx, hc, credential, validator, http.MethodPost, requestURL, requestBody)
}

// Patch 向微信支付发送一个http patch请求
func Patch(ctx context.Context, hc *http.Client, credential Credential, validator Validator, requestURL string, requestBody interface{}) ([]byte, error) {
	return do(ctx, hc, credential, validator, http.MethodPatch, requestURL, requestBody)
}

// Put 向微信支付发送一个http put请求
func Put(ctx context.Context, hc *http.Client, credential Credential, validator Validator, requestURL string, requestBody interface{}) ([]byte, error) {
	return do(ctx, hc, credential, validator, http.MethodPut, requestURL, requestBody)
}

// Delete 向微信支付发送一个http delete请求
func Delete(ctx context.Context, hc *http.Client, credential Credential, validator Validator, requestURL string, requestBody interface{}) ([]byte, error) {
	return do(ctx, hc, credential, validator, http.MethodDelete, requestURL, requestBody)
}

func do(ctx context.Context, hc *http.Client, credential Credential, validator Validator, method, requestURL string, body interface{}) ([]byte, error) {
	var reqBody string
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("json marshal body err:%v", err)
		}
		reqBody = string(bodyBytes)
	}
	return DoRequest(ctx, hc, credential, validator, method, requestURL, ApplicationJSON, reqBody, reqBody)
}

func DoRequest(ctx context.Context, hc *http.Client, credential Credential, validator Validator,
	method, requestURL, contentType, reqBody, signBody string) ([]byte, error) {
	var err error
	var authorization string
	request, err := http.NewRequestWithContext(ctx, method, requestURL,
		strings.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	request.Header.Set(Accept, "*/*")
	request.Header.Set(ContentType, contentType)
	request.Header.Set(UserAgent, UserAgentContent)
	// 生产授权信息
	authorization, err = credential.GenerateAuthorizationHeader(ctx, method,
		request.URL.RequestURI(), signBody)
	if err != nil {
		return nil, fmt.Errorf("generate authorization err:%s", err.Error())
	}
	request.Header.Set(Authorization, authorization)
	response, err := hc.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := CheckResponse(response)
	if err != nil {
		return body, err
	}
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return body, fmt.Errorf("read response body err:[%s]", err.Error())
	}
	if err = validator.Validate(ctx, body, response.Header); err != nil {
		return body, err
	}
	return body, nil
}

// CheckResponse 校验回包是否有错误
//
// 当http回包的状态码的范围不是200-299之间的时候，会返回相应的错误信息，主要包括http状态码、回包错误码、回包错误信息提示
func CheckResponse(res *http.Response) ([]byte, error) {
	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		return nil, nil
	}
	slurp, err := ioutil.ReadAll(res.Body)
	if err == nil {
		jerr := &Error{StatusCode: res.StatusCode}
		err = json.Unmarshal(slurp, jerr)
		if err == nil {
			return slurp, jerr
		}
	}
	return slurp, &Error{
		StatusCode: res.StatusCode,
		Body:       string(slurp),
		Header:     res.Header,
	}
}

// CreateFormField 设置form-data 中的普通属性
//
//示例内容
//	Content-Disposition: form-data; name="meta";
//	Content-Type: application/json
//
//	{ "filename": "file_test.mp4", "sha256": " hjkahkjsjkfsjk78687dhjahdajhk " }
//
// 如果要设置上述内容
//	CreateFormField(w, "meta", "application/json", meta)
func CreateFormField(w *multipart.Writer, fieldName, contentType string, fieldValue []byte) error {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s";`, fieldName))
	h.Set("Content-Type", contentType)
	part, err := w.CreatePart(h)
	if err != nil {
		return err
	}
	_, err = part.Write(fieldValue)
	return err
}

// CreateFormFile 设置form-data中的文件
//
// 示例内容：
//	Content-Disposition: form-data; name="file"; filename="file_test.mp4";
//	Content-Type: video/mp4
//
//	pic1  //pic1即为媒体视频的二进制内容
//
// 如果要设置上述内容，则CreateFormFile(w, "file_test.mp4", "video/mp4", pic1)
func CreateFormFile(w *multipart.Writer, filename, contentType string, file []byte) error {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "file", filename))
	h.Set("Content-Type", contentType)
	part, err := w.CreatePart(h)
	if err != nil {
		return err
	}
	_, err = part.Write(file)
	return err
}
