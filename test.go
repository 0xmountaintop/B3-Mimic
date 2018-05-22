package main

import (
    "fmt"
    "math/big"
    "strconv"
    "encoding/hex"

    // "github.com/bytom/consensus/difficulty"
)

var Diff1 = StringToBig("0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")

func StringToBig(h string) *big.Int {
    n := new(big.Int)
    n.SetString(h, 0)
    return n
}

func main() {
    // fmt.Println("Diff1:", Diff1)
    padded := make([]byte, 32)
    // padded := make([]byte, 32)

    // bitsStr := "540b02000000001d"

    // bitsStr = strSwitchEndian(bitsStr)
    // bitsUint := str2ui64(bitsStr)
    // diff := difficulty.CompactToBig(bitsUint)
    // fmt.Println("Diff:", diff)



    // // diff,_ := new(big.Int).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF",16)
    // // diff.Rsh(diff, uint(1))
    // // diffBuff := new(big.Int).Div(Diff1, diff).Bytes()
    // // fmt.Println(diffBuff)
    // // copy(padded[32-len(diffBuff):], diffBuff)
    // copy(padded[:], diff.Bytes())
    // fmt.Println(padded)
    // // buff := padded[0:4]
    // buff := padded[:]
    // fmt.Println(buff)
    // // targetHex := hex.EncodeToString(reverse(buff))
    // targetHex := hex.EncodeToString(buff)
    // fmt.Println(targetHex)

    // i := big.NewInt(32000)
    // fmt.Println(i)
    // diffBuff := i.Bytes()
    // fmt.Println(diffBuff)
    // copy(padded[32-len(diffBuff):], diffBuff)
    // fmt.Println(padded)
    // buff := padded[0:4]
    // fmt.Println(buff)
    // targetHex := hex.EncodeToString(reverse(buff))
    // fmt.Println(targetHex)
    
    // targetHex := "ffff0300"
    // targetHex := "ffff3f00"
    targetHex := "c5a70000"
    // targetHex := "ffffffff"
    fmt.Println(targetHex)
    decoded, _ := hex.DecodeString(targetHex)
    fmt.Println(decoded)
    decoded = reverse(decoded)
    fmt.Println(decoded)
    copy(padded[:len(decoded)], decoded)
    fmt.Println(padded)
    newDiff := new(big.Int).SetBytes(padded)
    fmt.Println(newDiff)
    newDiff = new(big.Int).Div(Diff1, newDiff)
    fmt.Println(newDiff)
}

func reverse(src []byte) []byte {
    dst := make([]byte, len(src))
    for i := len(src); i > 0; i-- {
        dst[len(src)-i] = src[i-1]
    }
    return dst
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


func str2ui64(str string) uint64 {
    ui64, _ := strconv.ParseUint(str, 16, 64)
    return ui64
}