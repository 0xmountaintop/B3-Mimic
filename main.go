package main

import(
    "fmt"

    "github.com/parnurzeal/gorequest"
)

var poolAddr = "stratum-btm.antpool.com:6666/"

func main() {
    fmt.Println(poolAddr)

    request := gorequest.New()
    // resp, body, errs := request.Get(req_url).End()

    // fmt.Println(string(json_dat))

    // _, body, _ := request.Post("https://httpbin.org/post").
    //     Send(string(json_dat)).
    //     // Send(`{"data":1}`).
    //     End()

    _, body, _ := request.Post(poolAddr).
        Send(`{
                  "id": 1,
                  "jsonrpc": "2.0",
                  "method": "login",
                  "params": [
                        "0xb85150eb365e7df0941f0cf08235f987ba91506a",//login
                        "",//Pass
                        "agent"//Agent
                    ]
            }`).
        End()

    fmt.Println(body)
}