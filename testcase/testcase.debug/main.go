package main

import (
	"context"
	"errors"
	"fmt"
)

func main() {
	cli := NewClientProxy("HelloService")

	type arg struct {
		req *helloReq
		//rsp *helloRsp
	}

	args := []arg{
		{
			req: "",
		},
		{
			req: "xx",
		},
	}
	for _, v := range args {
		// set breakpoints here
		rsp, err := cli.Invoke(context.TODO(), v.req)
	}
}

func NewClientProxy(name string) Client {
	return &client{}
}

type helloReq struct {
	msg string
}

type helloRsp struct {
	code int
	msg  string
}

// Client ...
type Client interface {
	Invoke(ctx context.Context, req interface{}) (rsp interface{}, err error)
}

// client ...
type client struct {
}

func (c *client) Invoke(ctx context.Context, req interface{}) (rsp interface{}, err error) {
	v, ok := req.(*helloReq)
	if !ok || v == nil {
		return nil, errors.New("invalid req")
	}
	if len(v.msg) == 0 {
		return &helloRsp{100, "req.msg empty"}, nil
	}
	return &helloRsp{0, "ok"}, nil
}
