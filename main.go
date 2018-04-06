package main

import(
    "fmt"

    "github.com/parnurzeal/gorequest"
)

var poolAddr = "stratum-btm.antpool.com:6666/"
// var poolAddr = "https://api.typeform.com/forms"


func main() {
    request := gorequest.New()

    // resp, body, _ := request.Post(poolAddr).
    //     Send(`{
    //             "id": 1,
    //             "jsonrpc": "2.0",
    //             "method": "login",
    //             "params": [
    //                 "0xb85150eb365e7df0941f0cf08235f987ba91506a",//login
    //                 "",//Pass
    //                 "agent"//Agent
    //             ]
    //         }`).
    //     End()
    // fmt.Println(resp)


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