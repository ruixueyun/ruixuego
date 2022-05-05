// Package ruixuego httpclient HTTP 客户端
package ruixuego

import (
	"crypto/tls"
	"time"

	"github.com/valyala/fasthttp"
)

const defaultContentType = "application/json"

// NewHTTPClient create http client
func NewHTTPClient(timeout time.Duration, concurrency int) *HTTPClient {
	return &HTTPClient{
		client: fasthttp.Client{
			TLSConfig:          &tls.Config{InsecureSkipVerify: true},
			ReadTimeout:        timeout,
			WriteTimeout:       timeout,
			MaxConnWaitTimeout: timeout,
		},
		timeout:     timeout,
		concurrency: make(chan struct{}, concurrency),
	}
}

// HTTPClient a fasthttp client
type HTTPClient struct {
	client      fasthttp.Client
	timeout     time.Duration
	concurrency chan struct{}
}

func (c *HTTPClient) prepareDo(
	url string, req *fasthttp.Request) *fasthttp.Response {

	c.concurrency <- struct{}{}
	req.SetRequestURI(url)
	return fasthttp.AcquireResponse()
}

func (c *HTTPClient) getRequestWithArgs(args *fasthttp.Args) *fasthttp.Request {

	req := GetRequest()
	req.SetBody(args.QueryString())
	// args need to be get by fashhttp.AcquireArgs()
	fasthttp.ReleaseArgs(args)
	return req
}

func (c *HTTPClient) afterDo(req *fasthttp.Request) {

	fasthttp.ReleaseRequest(req)
	<-c.concurrency
}

// Do 发起接口请求
func (c *HTTPClient) Do(url string, args *fasthttp.Args) (*fasthttp.Response, error) {

	return c.DoWithTimeout(url, args, c.timeout)
}

// DoWithTimeout 发起一个带有超时时间的请求
func (c *HTTPClient) DoWithTimeout(
	url string,
	args *fasthttp.Args,
	timeout time.Duration) (*fasthttp.Response, error) {

	req := c.getRequestWithArgs(args)
	req.Header.SetContentType(defaultContentType)
	return c.DoRequestWithTimeout(url, req, timeout)
}

// DoWithoutTimeout 发起一个没有超时时间的请求
func (c *HTTPClient) DoWithoutTimeout(
	url string, args *fasthttp.Args) (*fasthttp.Response, error) {

	req := c.getRequestWithArgs(args)
	req.Header.SetContentType(defaultContentType)
	return c.DoRequestWithoutTimeout(url, req)
}

// DoContentTypeWithTimeout 发起一个带有超时时间的请求
func (c *HTTPClient) DoContentTypeWithTimeout(
	contentType, url string,
	args *fasthttp.Args, timeout time.Duration) (*fasthttp.Response, error) {

	req := c.getRequestWithArgs(args)
	req.Header.SetContentType(contentType)
	return c.DoRequestWithTimeout(url, req, timeout)
}

// DoContentTypeWithoutTimeout 发起一个没有超时时间的请求
func (c *HTTPClient) DoContentTypeWithoutTimeout(
	contentType, url string, args *fasthttp.Args) (*fasthttp.Response, error) {

	req := c.getRequestWithArgs(args)
	req.Header.SetContentType(contentType)
	return c.DoRequestWithoutTimeout(url, req)
}

// DoRequest 指定请求头内容类型发起一个带有超时时间的请求
func (c *HTTPClient) DoRequest(
	url string,
	req *fasthttp.Request) (*fasthttp.Response, error) {

	resp := c.prepareDo(url, req)
	defer c.afterDo(req)
	if err := c.client.DoTimeout(req, resp, c.timeout); err != nil {
		return nil, err
	}
	return resp, nil
}

// DoRequestWithTimeout 指定请求头内容类型发起一个带有超时时间的请求直接返回 *fasthttp.Response
func (c *HTTPClient) DoRequestWithTimeout(
	url string,
	req *fasthttp.Request,
	timeout time.Duration) (*fasthttp.Response, error) {

	resp := c.prepareDo(url, req)
	defer c.afterDo(req)
	if err := c.client.DoTimeout(req, resp, timeout); err != nil {
		return nil, err
	}
	return resp, nil
}

// DoRequestWithoutTimeout 指定请求头内容类型发起一个没有超时时间的请求
func (c *HTTPClient) DoRequestWithoutTimeout(
	url string, req *fasthttp.Request) (*fasthttp.Response, error) {

	resp := c.prepareDo(url, req)
	defer c.afterDo(req)
	if err := c.client.Do(req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetArgs 获取 fasthttp 参数对象
func GetArgs() *fasthttp.Args {
	return fasthttp.AcquireArgs()
}

// GetRequest 获取 fasthttp 请求参数, 用于需要自定义请求头的场景
func GetRequest() *fasthttp.Request {
	return fasthttp.AcquireRequest()
}

// PutResponse 将 *fasthttp.Response 返回对象池
func PutResponse(resp *fasthttp.Response) {
	fasthttp.ReleaseResponse(resp)
}
