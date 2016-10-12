package requests

//import "fmt"
import "time"
import "strings"
import "net/url"
import "net/http"
import "io/ioutil"
import "encoding/json"

type M map[string]interface{}

type Request struct {
	Url  string
	Args M
}

type Response struct {
	Content string
}

func (r *Request) initArgs() {
	if r.Args == nil {
		r.Args = M{}
	}
}

func (r *Request) setTimeout(timeout int) error {
	r.Args["timeout"] = int(timeout)
	return nil
}

func Timeout(timeout int) func(*Request) error {
	return func(r *Request) error {
		return r.setTimeout(timeout)
	}
}

func (r *Request) setProxies(proxy string) error {
	r.Args["proxy"] = proxy
	return nil
}

func Proxies(proxy string) func(*Request) error {
	return func(r *Request) error {
		return r.setProxies(proxy)
	}
}

func (r *Request) setCookies(cookies map[string]string) error {
	r.Args["cookies"] = cookies
	return nil
}

func Cookies(cookies map[string]string) func(*Request) error {
	return func(r *Request) error {
		return r.setCookies(cookies)
	}
}

func (r *Request) setHeaders(headers map[string]string) error {
	r.Args["headers"] = headers
	return nil
}

func Headers(headers map[string]string) func(*Request) error {
	return func(r *Request) error {
		return r.setHeaders(headers)
	}
}

func (r *Request) setParams(params map[string]string) error {
	r.Args["params"] = params
	return nil
}

func Params(params map[string]string) func(*Request) error {
	return func(r *Request) error {
		return r.setParams(params)
	}
}

func (r *Request) setForm(form map[string]string) error {
	r.Args["form"] = form
	return nil
}

func Form(form map[string]string) func(*Request) error {
	return func(r *Request) error {
		return r.setForm(form)
	}
}

func (r *Request) setData(data map[string]string) error {
	r.Args["data"] = data
	return nil
}

func Data(data map[string]string) func(*Request) error {
	return func(r *Request) error {
		return r.setData(data)
	}
}

func (r *Request) setBin(bin map[string]string) error {
	r.Args["bin"] = bin
	return nil
}

func Bin(bin map[string]string) func(*Request) error {
	return func(r *Request) error {
		return r.setBin(bin)
	}
}

func (r *Request) setJson(json map[string]string) error {
	r.Args["json"] = json
	return nil
}

func Json(json map[string]string) func(*Request) error {
	return func(r *Request) error {
		return r.setJson(json)
	}
}

func (r *Request) setOptions(options M) error {
	for k, v := range options {
		r.Args[k] = v
	}
	return nil
}

func Options(options M) func(*Request) error {
	return func(r *Request) error {
		return r.setOptions(options)
	}
}

func (r *Request) Options(uri string, options ...func(*Request) error) (*Response, error) {
	return r.MakeRequest("Options", uri)

}

func (r *Request) MakeRequest(method string, uri string, options ...func(*Request) error) (*Response, error) {

	r.initArgs()

	for _, option := range options {
		err := option(r)
		if err != nil {
			panic(err)
		}
	}

	payload := ""
	if data, ok := r.Args["form"].(map[string]string); ok {
		var req http.Request
		req.ParseForm()
		for key, val := range data {
			req.Form.Add(key, val)
		}
		payload = strings.TrimSpace(req.Form.Encode())
	}

	if data, ok := r.Args["data"].(map[string]string); ok {
		body, err := json.Marshal(data)
		if err == nil {
			payload = string(body)
		}
	}

	req, err := http.NewRequest(method, uri, strings.NewReader(payload))

	if params, ok := r.Args["params"].(map[string]string); ok {
		q := req.URL.Query()
		for key, val := range params {
			q.Add(key, val)
		}
		req.URL.RawQuery = q.Encode()
	}

	if headers, ok := r.Args["headers"].(map[string]string); ok {
		for key, val := range headers {
			req.Header.Set(key, val)
		}
	}

	if cookies, ok := r.Args["cookies"].(map[string]string); ok {
		for key, val := range cookies {
			cookie := &http.Cookie{Name: key, Value: val, HttpOnly: false}
			req.AddCookie(cookie)
		}
	}

	transport := &http.Transport{}
	if proxy, ok := r.Args["proxy"].(string); ok {
		if proxy != "" {
			proxyUrl, err := url.Parse(proxy)
			if err == nil {
				transport.Proxy = http.ProxyURL(proxyUrl)
			}
		}
	}

	client := &http.Client{
		Transport: transport,
	}

	timeoutSeconds := r.Args["timeout"].(int)
	timeout := time.Duration(0) * time.Second
	if timeoutSeconds > 0 {
		timeout = time.Duration(timeoutSeconds) * time.Second
	}
	client.Timeout = timeout

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	response := &Response{}
	if err != nil {
		panic(err)
	} else {
		response.Content = string(body)
	}
	return response, nil
}
