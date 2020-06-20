s0: generate a gorpc server

```bash
gorpc create -protofile helloworld.proto
cd helloworld
```
    
s1: start a gorpc server

```bash
cd helloworld
go build -o helloworld
./helloworld
```

s2: run `gorpc <rpc>` to test

```bash
gorpc rpc -protofile=helloworld/helloworld/helloworld.proto \
        -reqbody=helloworld.HelloReq \
        -body='{"msg":"my name is zhangjie"}' \
        -rspbody=helloworld.HelloRsp \
        -target=ip://127.0.0.1:8000 \
        -func='/helloworld.helloworld_svr/Hello' \
        -times=10 \
        2> /dev/null
```
