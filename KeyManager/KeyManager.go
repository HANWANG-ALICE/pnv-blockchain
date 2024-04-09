package KeyManager

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/sm3"
	"log"
	"math/big"
	"strings"
	
	"../utils"
)

var AddressLength = 32

var (
	INDEXES  []int
	bigRadix = big.NewInt(58)
	bigZero  = big.NewInt(0)
)

const (
	ALPHABET = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
)

type KeyManager struct {
	priv	 *sm2.PrivateKey
	pub      *sm2.PublicKey
}

func (km *KeyManager) Init() {
	km.priv=new(sm2.PrivateKey)
	km.pub=new(sm2.PublicKey)
}

func (km *KeyManager) GenRandomKeyPair() {
	var err error
	km.priv, err = sm2.GenerateKey() // 生成密钥对
	if err != nil {
		log.Fatal("GenRandomKeyPair Wrong:",err)
	}
	km.pub = &km.priv.PublicKey
}

func (km *KeyManager) GetPubkey() string {
	pubkey,err:=sm2.MarshalSm2PublicKey(km.pub)
	if err != nil {
		log.Fatal("MarshalSm2PublicKey Wrong:",err)
	}
	return utils.BytesToHex(pubkey)
}



func (km *KeyManager) GetPriKey() (string, error) {
	prikey,err:=sm2.MarshalSm2UnecryptedPrivateKey(km.priv)
	return utils.BytesToHex(prikey), err
}

func (km *KeyManager) SetPriKey(data string) {
	prikey, err := utils.HexToBytes(data)
	if err != nil {
		log.Fatal("SetPriKey HexToBytes Wrong:",err)
	}
	priv_recover,err:=sm2.ParsePKCS8UnecryptedPrivateKey(prikey)
	if err != nil {
		log.Fatal("ParsePKCS8UnecryptedPrivateKey:",err)
	}
	km.priv = priv_recover
	km.pub = &km.priv.PublicKey
}

func (km *KeyManager) SetPubkey(data string) {
	keyFromText, err := utils.HexToBytes(data)
	if err != nil {
		log.Fatal("SetPubkey HexToBytes Wrong:",err)
	}
	keyFromPriv,err:=sm2.MarshalSm2PublicKey(km.pub)
	if err != nil {
		log.Fatal("MarshalSm2PublicKey Wrong:",err)
	}
	if !bytes.Equal(keyFromText,keyFromPriv){
		log.Fatal("公私钥不匹配")
	}
}

func (km *KeyManager) Sign(text []byte) (string, error) {
	return sign(text, km.priv)
}

func (km *KeyManager) Verify(text []byte, signature string, pubkey string) (bool, error) {
	pubkey_bytes, _ := utils.HexToBytes(pubkey)
	key, _ := sm2.ParseSm2PublicKey(pubkey_bytes)
	return verify(text, signature, key)
}

func (km *KeyManager) VerifyWithSelfPubkey(text []byte, signature string) (bool, error) {
	return verify(text, signature, km.pub)
}

func (km *KeyManager)GetAddress() string{
	pubkey,err:=sm2.MarshalSm2PublicKey(km.pub)
	if err != nil {
		log.Fatal("MarshalSm2PublicKey Wrong:",err)
	}
	return encodeAddress(pubkey)
}

func (km *KeyManager)VerifyAddressWithPubkey(pubkey, address string)(bool, error){
	b, err :=utils.HexToBytes(pubkey)
	if err != nil{
		return false,err
	}

	return encodeAddress(b) == address, nil
}

func encodeAddress(hash []byte) string {
	tosum := make([]byte, 32)
	copy(tosum[0:15], hash[0:15])
	copy(tosum[16:],hash[len(hash)-16:])
	cksum := doubleHash(tosum)

	b := make([]byte, 25)
	copy(b[0:], hash)
	copy(b[12:], cksum[:13])

	return base58Encode(b)
}
 
/**
  对text加密，text必须是一个hash值，例如md5、sha1等
  使用私钥prk
  使用随机熵增强加密安全，安全依赖于此熵，randsign
  返回加密结果，结果为数字证书r、s的序列化后拼接，然后用hex转换为string
*/
func sign(text []byte, prk *sm2.PrivateKey) (string, error) {
	r, s, err := sm2.Sign(prk, text)
	if err != nil {
		return "", err
	}
	rt, err := r.MarshalText()
	if err != nil {
		return "", err
	}
	st, err := s.MarshalText()
	if err != nil {
		return "", err
	}
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	defer w.Close()
	_, err = w.Write([]byte(string(rt) + "+" + string(st)))
	if err != nil {
		return "", err
	}
	w.Flush()
	return hex.EncodeToString(b.Bytes()), nil
}

/**
  证书分解
  通过hex解码，分割成数字证书r，s
*/
func getSign(signature string) (rint, sint big.Int, err error) {
	byterun, err := hex.DecodeString(signature)
	if err != nil {
		err = errors.New("decrypt error, " + err.Error())
		return
	}
	r, err := gzip.NewReader(bytes.NewBuffer(byterun))
	if err != nil {
		err = errors.New("decode error," + err.Error())
		return
	}
	defer r.Close()
	buf := make([]byte, 1024)
	count, err := r.Read(buf)
	if err != nil {
		fmt.Println("decode = ", err)
		err = errors.New("decode read error," + err.Error())
		return
	}
	rs := strings.Split(string(buf[:count]), "+")
	if len(rs) != 2 {
		err = errors.New("decode fail")
		return
	}
	err = rint.UnmarshalText([]byte(rs[0]))
	if err != nil {
		err = errors.New("decrypt rint fail, " + err.Error())
		return
	}
	err = sint.UnmarshalText([]byte(rs[1]))
	if err != nil {
		err = errors.New("decrypt sint fail, " + err.Error())
		return
	}
	return

}

/**
  校验文本内容是否与签名一致
  使用公钥校验签名和文本内容
*/
func verify(text []byte, signature string, key *sm2.PublicKey) (bool, error) {

	rint, sint, err := getSign(signature)
	if err != nil {
		return false, err
	}
	result := sm2.Verify(key,text,&rint,&sint)

	return result, nil

}

func GetHash(data []byte) []byte {
	h := sm3.New()
	h.Write([]byte(data))
	sum := h.Sum(nil)
	return sum
}

func doubleHash(data []byte) []byte {
	h1 := sha256.Sum256(data)
	h2 := sha256.Sum256(h1[:])
	return h2[:]
}

// Base58Encode encodes a byte slice to a modified base58 string.
func base58Encode(b []byte) string {
	x := new(big.Int)
	x.SetBytes(b)

	answer := make([]byte, 0)
	for x.Cmp(bigZero) > 0 {
		mod := new(big.Int)
		x.DivMod(x, bigRadix, mod)
		answer = append(answer, ALPHABET[mod.Int64()])
	}

	// leading zero bytes
	for _, i := range b {
		if i != 0 {
			break
		}
		answer = append(answer, ALPHABET[0])
	}
 
	// reverse
	alen := len(answer)
	for i := 0; i < alen/2; i++ {
		answer[i], answer[alen-1-i] = answer[alen-1-i], answer[i]
	}

	return string(answer)
}