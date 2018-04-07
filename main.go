package main

import(
    "fmt"
    "strconv"
    "encoding/json"
    "encoding/hex"
    "encoding/binary"
    
    "github.com/parnurzeal/gorequest"
    "github.com/bytom/protocol/bc"
    "github.com/bytom/crypto/sha3pool"
    "github.com/bytom/consensus/difficulty"
)

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
                                    //     0: JobId
                                    //     1: Version
                                    //     2: Height
                                    //     3: PreviousBlockHash
                                    //     4: Timestamp?
                                    //     5: TransactionsMerkleRoot
                                    //     6: TransactionStatusHash
                                    //     7: Nonce?
                                    //     8: Bits?
                                    //     9: Seed
                                    //     10: Target?
                                    // ]
    Error           Err        `json:"error, omitempty"`
}

const (
    maxNonce = ^uint64(0) // 2^64 - 1 = 18446744073709551615
)

var poolAddr = "stratum-btm.antpool.com:666/"
// var poolAddr = "221.212.212.212"



func main() {
    // fmt.Println(maxNonce)

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
            
    // body = `{
    //             "id": 10, 
    //             "result": null, 
    //             "error": { 
    //                 code: 0, 
    //                 message: "Work not ready" 
    //             }
    //         }`

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

    body = `{
                "id": 1,
                "jsonrpc": "2.0",
                "result": [
                    "1",
                    "1",
                    "0",
                    "0000000000000000000000000000000000000000000000000000000000000000", 
                    "5a685ae5", 
                    "237bf77df5c318dfa1d780043b507e00046fec7f8fdad80fc39fd8722852b27a", 
                    "53c0ab896cb7a3778cc1d35a271264d991792b7c44f5c334116bb0786dbc5635", 
                    "405510", 
                    "20000000007fffff", 
                    "e733c4b1c4ea57bc87346d9fce8c492248f1f414b9eac17faf9e9b8e0a107fa1",
                    "bdba0400"
                ]
            }`

    var jobResp JobResp
    json.Unmarshal([]byte(body), &jobResp)
    // fmt.Println(jobResp.Id)

    mine(jobResp.Result)
    // bhByte := mine(jobResp.Result)
}

// Version, Height, PreviousBlockId, Timestamp, TransactionsRoot, TransactionStatusHash, Bits, Nonce
// 156 = 20+136 = 8+11+1 + 8+8+32+8+32+32+8+8
func mine(job [11]string) []byte {
    inter := [156]byte{
                0x65, 0x6e, 0x74, 0x72, 0x79, 0x69, 0x64, 0x3a, //string "entryid:"
                0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, //string "blockheader"
                0x3a, //string ":"
        }

    // Version
    copy(inter[20:28], str2bytes(job[1], 8))
    // Height
    copy(inter[28:36], str2bytes(job[2], 8)) 
    // PreviousBlockId
    copy(inter[36:68], str2bytes(job[3], 32)) 
    // Timestamp
    copy(inter[68:76], str2bytes(job[4], 8)) 
    // TransactionsRoot
    copy(inter[76:108], str2bytes(job[5], 32)) 
    // TransactionStatusHash
    copy(inter[108:140], str2bytes(job[6], 32)) 
    // Bits
    copy(inter[140:148], str2bytes(job[8], 8)) 
    // Nonce
    ui64Nonce, _ := strconv.ParseUint(job[7], 16, 64)
    fmt.Println("Start mining from nonce:", ui64Nonce)
    for ; ui64Nonce <= maxNonce; ui64Nonce+=1 {
        fmt.Println("Trying nonce:", ui64Nonce)
        copy(inter[148:156], ui64To8Bytes(ui64Nonce))
        sha3pool.Sum256(inter[20:20+32], inter[20:20+136])
        sha3pool.Sum256(inter[20:20+32], inter[0:20+32])
        
        var header [32]byte
        copy(header[:], inter[20:20+32])
        headerHash := bc.NewHash(header)
        var seed [32]byte
        copy(seed[:], str2bytes(job[9], 32))
        seedHash := bc.NewHash(seed)
        bits, _ := strconv.ParseUint(job[8], 16, 64)
        
        // fmt.Printf("header:\t")
        // for _, h := range header {
        //     fmt.Printf("%02x", h)
        // }
        // fmt.Printf("\nseed:\t")
        // for _, s := range seed {
        //     fmt.Printf("%02x", s)
        // }
        // fmt.Printf("\nbits:\t%16x\n", bits)
        
        if difficulty.CheckProofOfWork(&headerHash, &seedHash, bits) {
            fmt.Println("Valid nonce found:", ui64Nonce)
            // break
        }
    }

    return inter[20:20+32]
}

func str2bytes(instr string, leng uint8) []byte {
    // fmt.Println([]byte(instr))
    outstr := fmt.Sprintf("%064s", instr)
    // fmt.Println(outstr)

    var b [32]byte
    hex.Decode(b[:], []byte(outstr))
    if len(instr) < 64 {
        b = litE2BigE(b)    
    }
    // fmt.Println(b)

    h := bc.NewHash(b)
    // fmt.Println(h.Bytes()[0:leng])
    return h.Bytes()[0:leng]
}

func litE2BigE(buf [32]byte) [32]byte {
    blen := len(buf)
    for i := 0; i < blen/2; i++ {
        buf[i], buf[blen-1-i] = buf[blen-1-i], buf[i]
    }
    return buf
}

func ui64To8Bytes(ui64 uint64) []byte {
    bs := make([]byte, 8)
    binary.LittleEndian.PutUint64(bs, ui64)
    // fmt.Println(bs)
    return bs
}