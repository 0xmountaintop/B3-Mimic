package main

import(
    "net"
    "log"
    "fmt"
    "strconv"
    "encoding/json"
    // "time"
    "math/big"
    "encoding/hex"
    
    "github.com/bytom/testutil"
    "github.com/bytom/protocol/bc/types"
    "github.com/bytom/consensus/difficulty"
    // "github.com/bytom/consensus"
)

type t_err struct {
    Code            uint64      `json:"code"`
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
    JobId           string      `json:"job_id"`
    Seed            string      `json:"seed"`
    Target          string      `json:"target"`
}

type t_result struct {
    Id              string      `json:"id"`
    Job             t_job       `json:"job"`
    Status          string      `json:"status"`
}

type t_resp struct {
    Id              int64       `json:"id"`
    Jsonrpc         string      `json:"jsonrpc, omitempty"`
    Result          t_result    `json:"result, omitempty"`
    Error           t_err       `json:"error, omitempty"`
}

type t_jobntf struct {
    Jsonrpc         string      `json:"jsonrpc, omitempty"`
    Method          string      `json:"method, omitempty"`
    Params          t_job       `json:"params, omitempty"`
}

const (
    maxNonce = ^uint64(0) // 2^64 - 1 = 18446744073709551615
    poolAddr = "stratum-btm.antpool.com:6666" //39.107.125.245
    login = `poisoned.1`

    flush = "\r\n\r\n"
    MOCK = false
    DEBUG = false
    esHR  = uint64(166) //estimated Hashrate. 1 for Go, 10 for simd, 166 for gpu, 900 for B3
)

var (
    MsgId = uint64(0)
    Diff1 = StringToBig("0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
    newestJob = ""
)

func main() {
    MsgId += 1
    conn, err := net.Dial("tcp", poolAddr)
    if err != nil {
        log.Fatalln(err)
    }
    defer conn.Close()

    send_msg := `{"method": "login", "params": {"login": "`
    send_msg += login 
    send_msg += `", "pass": "123", "agent": "bmminer/2.0.0"}, "id": `
    send_msg += strconv.FormatUint(MsgId, 10) 
    send_msg += `}`
    MsgId += 1
    conn.Write([]byte(send_msg))
    conn.Write([]byte(flush))
    log.Printf("Sent: %s", send_msg)
    buff := make([]byte, 1024)
    n, _ := conn.Read(buff)
    log.Printf("----login job received----\n%s\n", buff[:n])
    var resp t_resp
    json.Unmarshal([]byte(buff[:n]), &resp)

        
    if DEBUG && MOCK {
        mock_input(&resp)
    }

    newestJob = resp.Result.Job.JobId
    go func(job t_job, conn net.Conn){
        mine(job, conn)
    }(resp.Result.Job, conn)
    
    for true {
        buff = make([]byte, 1024)
        n, _ = conn.Read(buff)
        var jobntf t_jobntf
        json.Unmarshal([]byte(buff[:n]), &jobntf)
        
        if jobntf.Method == "job" {
            log.Printf("----new job received----\n%s\n", buff[:n])
            newestJob = jobntf.Params.JobId
            go func(job t_job, conn net.Conn){
                mine(job, conn)
            }(jobntf.Params, conn)
        } else {
            log.Printf("Received: %s\n", buff[:n])
        }
    }
}

/*
type BlockHeader struct {
    Version             uint64  // The version of the block.
    Height              uint64  // The height of the block.
    PreviousBlockHash   bc.Hash // The hash of the previous block.
    Timestamp           uint64  // The time of the block in seconds.
    Nonce               uint64  // Nonce used to generate the block.
    Bits                uint64  // Difficulty target for the block.
    BlockCommitment     types.BlockCommitment{
                            TransactionsMerkleRoot: node.transactionsMerkleRoot,
                            TransactionStatusHash:  node.transactionStatusHash,
                        },
}
*/

func mine(job t_job, conn net.Conn) bool {
    bh := &types.BlockHeader{
                Version:            str2ui64Bg(job.Version),
                Height:             str2ui64Bg(job.Height),
                PreviousBlockHash:  testutil.MustDecodeHash(job.PreBlckHsh),
                Timestamp:          str2ui64Bg(job.Timestamp),
                Bits:               str2ui64Bg(job.Bits),
                BlockCommitment:    types.BlockCommitment{
                                        TransactionsMerkleRoot: testutil.MustDecodeHash(job.TxMkRt),
                                        TransactionStatusHash:  testutil.MustDecodeHash(job.TxStRt),
                                    },
    }
    if DEBUG {
        view_parsing(bh, job)
    }

    log.Printf("Job %s: Mining at height: %d\n", job.JobId, bh.Height)
    seedHash := testutil.MustDecodeHash(job.Seed)
    padded := make([]byte, 32)
    targetHex := job.Target
    decoded, _ := hex.DecodeString(targetHex)
    decoded = reverse(decoded)
    copy(padded[:len(decoded)], decoded)
    newDiff := new(big.Int).SetBytes(padded)
    newDiff = new(big.Int).Div(Diff1, newDiff)
    log.Printf("Job %s: Old target: %v\n", job.JobId, difficulty.CompactToBig(bh.Bits))
    newDiff = new(big.Int).Mul(difficulty.CompactToBig(bh.Bits), newDiff)
    log.Printf("Job %s: New target: %v\n", job.JobId, newDiff)

    nonce := str2ui64Li(job.Nonce)
    log.Printf("Job %s: Start from nonce:\t0x%016x = %d\n", job.JobId, nonce, nonce)
    // for i := nonce; i <= nonce+consensus.TargetSecondsPerBlock*esHR && i <= maxNonce; i++ {
    for i := nonce; i <= maxNonce; i++ {
        if job.JobId != newestJob {
            log.Printf("Job %s: Expired", job.JobId)
            return false
        } else {
            // log.Printf("Checking PoW with nonce: 0x%016x = %d\n", i, i)
            bh.Nonce = i
            headerHash := bh.Hash()
            if DEBUG {
                fmt.Printf("Job %s: HeaderHash: %v\n", job.JobId, headerHash.String())
            }

            // if difficulty.CheckProofOfWork(&headerHash, &seedHash, bh.Bits) {
            if difficulty.CheckProofOfWork(&headerHash, &seedHash, difficulty.BigToCompact(newDiff)) {
                log.Printf("Job %s: Target found! Proof hash: 0x%v\n", job.JobId, headerHash.String())

                nonceStr := strconv.FormatUint(i, 16)
                // nonceStr = strSwitchEndian(fmt.Sprintf("%016s", nonceStr))
                nonceStr = fmt.Sprintf("%016s", nonceStr)
                if DEBUG {
                    log.Printf("Job %s: Sending back nonce as string: %s", job.JobId, nonceStr)
                }

                send_msg := `{"method": "submit", "params": {"id": "`
                send_msg += login 
                send_msg += `", "job_id": "`
                send_msg += job.JobId
                send_msg += `", "nonce": "`
                send_msg += nonceStr
                send_msg += `"}, "id":`
                send_msg += strconv.FormatUint(MsgId, 10)
                send_msg += `}`
                MsgId += 1
                conn.Write([]byte(send_msg))
                conn.Write([]byte(flush))
                log.Printf("Job %s: Sent: %s", job.JobId, send_msg)
                // return true
            }        
        }
    }
    log.Printf("Job %s: Stop at nonce:\t\t0x%016x = %d\n", job.JobId, bh.Nonce, bh.Nonce)
    return false
}

func mock_input(presp *t_resp) {
    body := `{
                "id":1,
                "jsonrpc":"2.0",
                "result":{
                    "id":"antminer_1",
                    "job":{
                        "version":"0100000000000000",
                        "height":"0000000000000000",
                        "previous_block_hash":"0000000000000000000000000000000000000000000000000000000000000000",
                        "timestamp":"e55a685a00000000",
                        "transactions_merkle_root":"237bf77df5c318dfa1d780043b507e00046fec7f8fdad80fc39fd8722852b27a",
                        "transaction_status_hash":"53c0ab896cb7a3778cc1d35a271264d991792b7c44f5c334116bb0786dbc5635",
                        "nonce":"1055400000000000",
                        "bits":"ffff7f0000000020",
                        "job_id":"16942",
                        "seed":"8636e94c0f1143df98f80c53afbadad4fc3946e1cc597041d7d3f96aebacda07",
                        "target":"c5a70000"
                    },
                    "status":"OK"
                },
                "error":null
            }`
    json.Unmarshal([]byte(body), &(*presp))
}

func view_parsing(bh *types.BlockHeader, job t_job) {
    log.Println("Printing parsing result:")
    fmt.Println("\tVersion:", bh.Version)
    fmt.Println("\tHeight:", bh.Height)
    fmt.Println("\tPreviousBlockHash:", bh.PreviousBlockHash.String())
    fmt.Println("\tTimestamp:", bh.Timestamp)
    fmt.Println("\tbits_str:", job.Bits)
    fmt.Println("\tBits:", bh.Bits)
    fmt.Println("\tTransactionsMerkleRoot:", bh.BlockCommitment.TransactionsMerkleRoot.String())
    fmt.Println("\tTransactionStatusHash:", bh.BlockCommitment.TransactionStatusHash.String())
    fmt.Println("\ttarget_str:", job.Target)
    fmt.Println("\ttarget_ui64Bg:", str2ui64Bg(job.Target))
}

func str2ui64Bg(str string) uint64 {
    ui64, _ := strconv.ParseUint(strSwitchEndian(str), 16, 64)
    return ui64
}

func str2ui64Li(str string) uint64 {
    ui64, _ := strconv.ParseUint(str, 16, 64)
    return ui64
}

func strSwitchEndian(oldstr string) string {
    // fmt.Println("old str:", oldstr)
    slen := len(oldstr)
    if slen%2 != 0 {
        panic("hex string format error")
    }

    newstr := ""
    for i := 0; i < slen; i+=2 {
        newstr += oldstr[slen-i-2:slen-i]
    }
    // fmt.Println("new str:", newstr)
    return newstr
}

func StringToBig(h string) *big.Int {
    n := new(big.Int)
    n.SetString(h, 0)
    return n
}

func reverse(src []byte) []byte {
    dst := make([]byte, len(src))
    for i := len(src); i > 0; i-- {
        dst[len(src)-i] = src[i-1]
    }
    return dst
}
