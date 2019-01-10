package httpbatchlib

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
)

// 实现http batch request：目的client请求server建立的connection数
// 由于每次进行http连接都会带来一定的开销，通过batch http request在单个http connection放入多个API 调用
// 常见实例
// 1.更新元数据(metadata)：在多个对象上更新权限等
// 2.删除很多对象
// ...
// 在每一种情况下，都可以将它们组合到单个HTTP请求中，而不是分别发送每个调用.
// 注意：主请求的Content-Type必须指明：multipart/mixed；而主http请求会包括多个内嵌的http请求，每个内嵌http请求
//       对应的http头：Content-Type: application/http并且包含一个Content-ID选项；这些头部信息标记每个内嵌http请求的开始
//       当服务端拿到主请求后就会分别提取出每个内嵌http，并将对应头部标记部分忽略
// 同样响应也是标准的http response对应的content type = multipart/mixed，主响应中的每一个块对应着一个内嵌http请求的响应

// 注意格式中boundary的内容就是用来分割每个request的
// Batch 请求格式
//	POST /batch HTTP/1.1
//	Host: localhost:8080
//	Content-Length: 457
//	Content-Type: multipart/mixed; boundary="===============7330845974216740156=="
//	--===============7330845974216740156==
//	Content-Type: application/http
//	Content-Transfer-Encoding: binary
//	Content-ID: <b29c5de2-0db4-490b-b421-6a51b598bd22+1>
//
//	GET /ie=utf-8&mod=1&isbd=1&isid=411B89C151082533&ie=utf-8&f=8&rsv_bp=0&rsv_idx=1&tn=baidu&wd=hello&rsv_pq=8a36c22d00024742&rsv_t=7a87AqER8FLHRaU37tvDhD4RZNSII0U9dut0gSpMKj44SyTMyqhcBWLcHns&rqlang=cn&rsv_enter=0&rsv_sug3=6&rsv_sug1=6&rsv_sug7=101&inputT=5442&rsv_sug4=6169&rsv_sid=1447_25809_21126_28131_26350_28266_28140_22073&_ss=1&clist=&hsug=&f4s=1&csor=5&_cr1=28436 HTTP/1.1
//	Host: www.baidu.com
//	Accept: */*
//
//
//	--===============7330845974216740156==
//	Content-Type: application/http
//	Content-Transfer-Encoding: binary
//	Content-ID: <b29c5de2-0db4-490b-b421-6a51b598bd22+2>
//
//	GET https://	www.jd.com HTTP/1.1
//	Accept: */*
//
//
//	--===============7330845974216740156==--


var (
	ErrMalformedMediaType = fmt.Errorf("Malformed media type")
	ErrInvalidMediaType = fmt.Errorf("Invalid media type, it should be multipart/mixed with boundary parameter")
	ErrMalformedRequest = fmt.Errorf("Invalid request part format")
)

// ========================batch start========================
// TODO: no implement
func Batch(requests ...*http.Request) *BeegoHTTPRequest{
	return NewBeegoRequestBatch("BATCH", requests...)
}

// ========================= batch request =======================
// VaildRequest: 验证Batch是不是采用的POST方式、对应的mediaType必须设定内容
// mediaType对应的内容：Content-Type: multipart/mixed; boundary="===============分割符=="
func ValidRequest(r *http.Request) (status int){
	if r.Method != http.MethodPost{
		status = http.StatusMethodNotAllowed
        log.Printf("Error http status method not allowed!!! statuscode = %v",status)
      	return
	}

	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		log.Printf("Error parsing the media type: %s", err)
		status = http.StatusBadRequest
		return
	}

	if mediaType != "multipart/mixed" {
		status = http.StatusUnsupportedMediaType
		log.Printf("Error http status unsupported media type: statuscode=%v", status)
		return
	}
	return 0
}
func UnpackRequests(r *http.Request) ([]*http.Request, error) {
	delimiter, err := extractDelimiter(r)
	if err != nil {
		return nil, err
	}

	var requests []*http.Request
	multiReader := multipart.NewReader(r.Body, delimiter)
	for {
		part, err := multiReader.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		partReq, err := http.ReadRequest(bufio.NewReader(part))
		if err != nil {
			log.Printf(err.Error())
			return nil, ErrMalformedRequest
		}
		fixRequestForClientCall(partReq)

		requests = append(requests, partReq)
	}

	return requests, nil
}
// 提取batch request中的分隔符内容
func extractDelimiter(r *http.Request) (string, error) {
	mediaType, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))

	if err != nil {
		return "", ErrMalformedMediaType
	}
	if mediaType != "multipart/mixed" {
		return "", ErrInvalidMediaType
	}

	boundary, ok := params["boundary"]
	if !ok {
		return "", ErrInvalidMediaType
	}
	return boundary, nil
}
// 修正请求响应客户端
func fixRequestForClientCall(r *http.Request) {
	r.RequestURI = ""
	r.URL.Host = r.Host

	if r.URL.Scheme == "" {
		r.URL.Scheme = "http"
	}
}
// ========================= batch request =======================
// ========================= batch response ======================
func BuildResponses(requests []*http.Request) []*http.Response{
	var responses = make([]*http.Response,0,len(requests))
	for _, req := range requests{
		resp, err := http.DefaultClient.Do(req)
		if err != nil{
			log.Printf("Error during the fanout requests: %s", err)
			return nil
		}
		responses = append(responses, resp)
	}
	return responses
}

func WriteResponses(w http.ResponseWriter, responses []*http.Response) error {
	var buf bytes.Buffer
	multipartWriter := multipart.NewWriter(&buf)

	mimeHeaders := textproto.MIMEHeader(make(map[string][]string))
	mimeHeaders.Set("Content-Type", "application/http")

	for _, resp := range responses {
		part, err := multipartWriter.CreatePart(mimeHeaders)
		if err != nil {
			return err
		}
		resp.Write(part)
	}

	multipartWriter.Close()

	w.Header().Set("Content-Type", mime.FormatMediaType("multipart/mixed", map[string]string{"boundary": multipartWriter.Boundary()}))
	w.WriteHeader(http.StatusOK)
	buf.WriteTo(w)
	return nil
}
// ========================= batch response ======================

const Default_HTTP_BATCH_URL = "/batch"
func NewBeegoRequestBatch(method string, requests ...*http.Request) *BeegoHTTPRequest{
	var body bytes.Buffer
	var resp http.Response

	req , err := http.NewRequest("POST",Default_HTTP_BATCH_URL, ioutil.NopCloser(&body))
	if err != nil{
		log.Print(fmt.Errorf("Impossible to create the wrapper request: %v\n", err))
		return nil
	}
	// 设置属性
	req.Header = make(http.Header)
	req.Proto = "HTTP/1.1"
	req.ProtoMajor = 1
	req.ProtoMinor = 1

	multiWriter := multipart.NewWriter(&body)
	partHeaders := textproto.MIMEHeader{}
	partHeaders.Set("Content-Type", "application/http")
	for _, r := range requests {
		partBody, errFor := multiWriter.CreatePart(partHeaders)
		if errFor != nil {
			log.Printf("Error while creating request part: %v\n", errFor)
			return nil
		}
		r.WriteProxy(partBody)
	}
	multiWriter.Close()

	mediaType := mime.FormatMediaType("multipart/mixed", map[string]string{"boundary": multiWriter.Boundary()})
	req.Header.Set("Content-Type", mediaType)

	return &BeegoHTTPRequest{
		url:     Default_HTTP_BATCH_URL,
		req:     req,
		params:  map[string][]string{},
		files:   map[string]string{},
		setting: defaultSetting,
		resp:    &resp,
	}
}
// ========================batch end========================
