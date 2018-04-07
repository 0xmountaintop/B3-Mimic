package main

import(
    "net"
    "log"
    "fmt"
    "strconv"
    "encoding/json"
    "encoding/hex"
    "encoding/binary"
    
    "github.com/bytom/protocol/bc"
    "github.com/bytom/crypto/sha3pool"
    "github.com/bytom/consensus/difficulty"
)

type t_err struct {
    Code            int64       `json:"code"`
    Message         string      `json:"message"`
}

type t_job struct {
    Version         string      `json:"version"`
    Height          string      `json:"height"`
    PreBlckHsh      string      `json:"previous_block_hash"`
    Timestamp       string      `json:"timestamp"`
    TxMkRt          string      `json:"transactions_merkle_root"`
    TxStRt          string      `json:"transaction_status_hash"`
    Nonce           string      `json:"nonce"`
    Bits            string      `json:"bits"`
    JobId           string      `json:"job_id"` //!!!
    Seed            string      `json:"seed"`
    Target          string      `json:"target"`
}

type t_result struct {
    Id              string      `json:"id"`
    Job             t_job       `json:"job"`
    Status          string      `json:"status"`
}

type t_resp struct {
    Id              int64       `json:"id"` //!!!
    Jsonrpc         string      `json:"jsonrpc, omitempty"`
    Result          t_result    `json:"result, omitempty"`
    Error           t_err       `json:"error, omitempty"`
}

const (
    maxNonce = ^uint64(0) // 2^64 - 1 = 18446744073709551615
    poolAddr = "stratum-btm.antpool.com:6666"
    flush = "\r\n\r\n"
)

func main() {
    conn, err := net.Dial("tcp", poolAddr)
    if err != nil {
        log.Fatalln(err)
    }
    defer conn.Close()

    send_msg := `{"method": "login", "params": {"login": "antminer_1", "pass": "123", "agent": "bmminer/2.0.0"}, "id": 1}`
    conn.Write([]byte(send_msg))
    conn.Write([]byte(flush))
    log.Printf("Sent: %s", send_msg)

    buff := make([]byte, 1024)
    n, _ := conn.Read(buff)
    log.Printf("Received: %s", buff[:n])


    // body := `{
    //             "id":1,
    //             "jsonrpc":"2.0",
    //             "result":{
    //                 "id":"antminer_1",
    //                 "job":{
    //                     "version":"0100000000000000",
    //                     "height":"2315000000000000",
    //                     "previous_block_hash":"ecaeb84f7787aca9ed08199169fe8dc5a01f3ba50c90c6844452f916d645d911",
    //                     "timestamp":"239dc75a00000000",
    //                     "transactions_merkle_root":"178f3dfaf916a5f8167100602254df50f6821c243a1f6263efedde798e9271a2",
    //                     "transaction_status_hash":"c9c377e5192668bc0a367e4a4764f11e7c725ecced1d7b6a492974fab1b6d5bc",
    //                     "nonce":"4a06000001000000",
    //                     "bits":"4690890000000021",
    //                     "job_id":"16942",
    //                     "seed":"a2a62d7715ee2234e1d73c22d26a1707fb7bc0f4ee0c01d43a4c97b0328379c5",
    //                     "target":"c5a70000"
    //                 },
    //                 "status":"OK"
    //             },
    //             "error":null
    //         }`


    var resp t_resp
    json.Unmarshal([]byte(buff[:n]), &resp)
    // json.Unmarshal([]byte(body), &resp)

    mine(resp.Result.Job)
}

// Version, Height, PreviousBlockId, Timestamp, TransactionsRoot, TransactionStatusHash, Bits, Nonce
// 156 = 20+136 = 8+11+1 + 8+8+32+8+32+32+8+8
func mine(job t_job) []byte {
    inter := [156]byte{
                0x65, 0x6e, 0x74, 0x72, 0x79, 0x69, 0x64, 0x3a, //string "entryid:"
                0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, //string "blockheader"
                0x3a, //string ":"
        }


    log.Println("data parsed:")

    copy(inter[20:28], str2bytes(job.Version, 8))
    fmt.Printf("\tVersion:\t0x")
    for _, h := range inter[20:28] {
        fmt.Printf("%02x", h)
    }
    fmt.Println()

    copy(inter[28:36], str2bytes(job.Height, 8))
    fmt.Printf("\tHeight:\t\t0x")
    for _, h := range inter[28:36] {
        fmt.Printf("%02x", h)
    }
    fmt.Println()
    
    copy(inter[36:68], str2bytes(job.PreBlckHsh, 32))
    fmt.Printf("\tPreBlckHsh:\t0x")
    for _, h := range inter[36:68] {
        fmt.Printf("%02x", h)
    }
    fmt.Println()
    
    copy(inter[68:76], str2bytes(job.Timestamp, 8))
    fmt.Printf("\tTimestamp:\t0x")
    for _, h := range inter[68:76] {
        fmt.Printf("%02x", h)
    }
    fmt.Println()
    
    copy(inter[76:108], str2bytes(job.TxMkRt, 32))
    fmt.Printf("\tTxMkRt:\t\t0x")
    for _, h := range inter[76:108] {
        fmt.Printf("%02x", h)
    }
    fmt.Println()
    
    copy(inter[108:140], str2bytes(job.TxStRt, 32))
    fmt.Printf("\tTxStRt:\t\t0x")
    for _, h := range inter[108:140] {
        fmt.Printf("%02x", h)
    }
    fmt.Println()
    
    copy(inter[140:148], str2bytes(job.Bits, 8))
    fmt.Printf("\tBits:\t\t0x")
    for _, h := range inter[140:148] {
        fmt.Printf("%02x", h)
    }
    fmt.Println()

    ui64Nonce, _ := strconv.ParseUint(job.Nonce, 16, 64)
    log.Printf("Start mining from nonce: 0x%x\n", ui64Nonce)
    for ; ui64Nonce <= maxNonce; incr_nonce(&ui64Nonce) { //???
        log.Printf("Trying nonce: 0x%x\n", ui64Nonce)
        copy(inter[148:156], ui64To8Bytes(ui64Nonce))
        sha3pool.Sum256(inter[20:20+32], inter[20:20+136])
        sha3pool.Sum256(inter[20:20+32], inter[0:20+32])
        
        var header [32]byte
        copy(header[:], inter[20:20+32])
        headerHash := bc.NewHash(header)
        var seed [32]byte
        copy(seed[:], str2bytes(job.Seed, 32))
        seedHash := bc.NewHash(seed)
        bits, _ := strconv.ParseUint(job.Bits, 16, 64)
        

        log.Println("checking Pow with:")
        fmt.Printf("\theader:\t0x")
        for _, h := range header {
            fmt.Printf("%02x", h)
        }
        fmt.Printf("\n\tseed:\t0x")
        for _, s := range seed {
            fmt.Printf("%02x", s)
        }
        fmt.Printf("\n\tbits:\t0x%16x\n", bits)
        
        if difficulty.CheckProofOfWork(&headerHash, &seedHash, bits) {
            log.Printf("Valid nonce found: 0x%x\n", ui64Nonce)
            break
        }
    }

    return inter[20:20+32]
}

func incr_nonce(ui64Nonce *uint64) {
    (*ui64Nonce) += 1    
}

func str2bytes(instr string, leng uint8) []byte {
    var b [32]byte //???
    hex.Decode(b[:], []byte(instr))
    return b[0:leng]
}

// func litE2BigE(buf [32]byte) [32]byte {
//     blen := len(buf)
//     for i := 0; i < blen/2; i++ {
//         buf[i], buf[blen-1-i] = buf[blen-1-i], buf[i]
//     }
//     return buf
// }

func ui64To8Bytes(ui64 uint64) []byte {
    bs := make([]byte, 8)
    binary.LittleEndian.PutUint64(bs, ui64)
    fmt.Printf("\t\t\t\tbs:\t0x")
    for _, b := range bs {
        fmt.Printf("%02x", b)
    }
    fmt.Println()
    return bs
}
