package utils

import (
	"encoding/binary"
	"encoding/hex"
	"strconv"
	"time"
)

//

func Uint64ToString(i uint64) string {
	return strconv.FormatUint(i, 10)
}

func Uint64ToBytes(i uint64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, i)
	return buf
}

func IntToBytes(i int) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint32(buf, uint32(i))
	return buf
}

func BytesToInt(buf []byte) int {
	return int(binary.BigEndian.Uint32(buf))
}
func BytesToUint64(buf []byte) uint64 {
	return binary.BigEndian.Uint64(buf)
}

func BytesToHex(data []byte) (dst string) {
	return hex.EncodeToString(data)
}
func HexToBytes(data string) (dst []byte, err error) {
	return hex.DecodeString(data)
}

//获取时间戳,单位为秒
func GetCurrentTime() float64 {
	return float64(time.Now().UTC().UnixNano()) / 1e9
}

//获取时间戳,单位为毫秒
func GetCurrentTimeMilli() float64 {
	return float64(time.Now().UTC().UnixNano()) / 1e6
}