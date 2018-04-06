package main

import(
    "fmt"
    "encoding/json"
    
    "github.com/parnurzeal/gorequest"
)

var poolAddr = "stratum-btm.antpool.com:6666/"
// var poolAddr = "221.212.212.212"


type Err struct {
    Code            int64       `json:"code"`
    Message         string      `json:"message"`
}

type AuthResp struct {
    Id              int64       `json:"id"`
    Jsonrpc         string      `json:"jsonrpc"`
}

type JobResp struct {
    Id              int64       `json:"id"`
    Jsonrpc         string      `json:"jsonrpc, omitempty"`
    Result          [11]string  `json:"result, omitempty"`
                                    // [
                                    //     JobId
                                    //     Version
                                    //     Height
                                    //     PreviousBlockHash
                                    //     Timestamp
                                    //     TransactionsMerkleRoot
                                    //     TransactionStatusHash
                                    //     Nonce
                                    //     Bits
                                    //     Seed
                                    //     Target
                                    // ]
    Error           Err        `json:"error, omitempty"`
}





func main() {
    request := gorequest.New()

    // resp, body, _ := request.Post(poolAddr).
    _, body, _ := request.Post(poolAddr).
        Send(`{
                  "id": 1,
                  "jsonrpc": "2.0",
                  "method": "login",
                  "params": [
                     "antminer",//login
                     "001",//Pass
                     "agent"//Agent
                  ]
                }`).
        End()
    // fmt.Println(resp)
    // fmt.Println(body)


    body = `{
                "id": 1,
                "jsonrpc": "2.0",
                "result": [
                    "1",
                    "1",
                    "1", 
                    "e733c4b1c4ea57bc87346d9fce8c492248f1f414b9eac17faf9e9b8e0a107fa1", 
                    "5aa39c6e", 
                    "15bd7762b3ee8057ecb83b792e2168c6b6bddaf10163d110f7e63db387e6aacf", 
                    "53c0ab896cb7a3778cc1d35a271264d991792b7c44f5c334116bb0786dbc5635", 
                    "8000000000000000", 
                    "20000000007fffff", 
                    "e733c4b1c4ea57bc87346d9fce8c492248f1f414b9eac17faf9e9b8e0a107fa1",
                    "bdba0400"
                ]
            }`


            
    // body = `{ 
    //             "id": 10, 
    //             "result": null, 
    //             "error": { 
    //                 code: 0, 
    //                 message: "Work not ready" 
    //             } 
    //         }`



    var jobResp JobResp
    json.Unmarshal([]byte(body), &jobResp)

    fmt.Println(jobResp.Id)
}