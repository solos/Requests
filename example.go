package main

import r "../requests/requests"
import "fmt"

func main() {
	kwargs := r.M{}
	options := r.M{
		"timeout": 10,
	}
	cookies := map[string]string{
		"user": "solos",
	}
	headers := map[string]string{
		"content-Type": "application/json",
		//"content-Type": "application/x-www-form-urlencoded",
	}

	data := map[string]string{
		"hello": "world",
	}

	req := &r.Request{Args: kwargs}
	resp, _ := req.MakeRequest("GET", "http://www.example.com", r.Timeout(10), r.Headers(headers), r.Cookies(cookies), r.Options(options), r.Data(data))
	fmt.Println(resp.Content)
	resp, _ = req.MakeRequest("POST", "http://www.example.com", r.Timeout(10), r.Headers(headers), r.Cookies(cookies), r.Options(options), r.Data(data))
	fmt.Println(resp.Content)
}
