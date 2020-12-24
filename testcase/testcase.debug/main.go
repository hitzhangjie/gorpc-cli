package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	go func() {
		for {
			time.Sleep(time.Second)
		}
	}()
	cli := NewXXXClientProxy("HelloService")

	type arg struct {
		req *helloReq
		//rsp *helloRsp
	}

	args := []arg{
		{
			req: &helloReq{msg: ""},
		},
		{
			req: &helloReq{msg: "xx"},
		},
	}
	for _, v := range args {
		// set breakpoints here
		rsp, err := cli.Invoke(context.TODO(), v.req)
		if err != nil {
			panic("invalid rsp")
		}
		fmt.Println(rsp, err)
	}
}

func NewXXXClientProxy(name string) *Client {
	return &Client{}
}

type helloReq struct {
	msg string
}

type helloRsp struct {
	code int
	msg  string
}

// client ...
type Client struct {
}

func (c *Client) Invoke(ctx context.Context, req *helloReq) (rsp *helloRsp, err error) {
	if len(req.msg) == 0 {
		return &helloRsp{100, "req.msg empty"}, nil
	}
	return &helloRsp{0, "ok"}, nil
}
