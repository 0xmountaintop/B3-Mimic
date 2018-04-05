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
        Send(`  {
                    "command":"test",
                }`).
        End()

    fmt.Println(body)
}