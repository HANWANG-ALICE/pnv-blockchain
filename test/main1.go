package main

import (
	"../KeyManager"
	"encoding/json"
	"fmt"
	"log"
	"net/rpc/jsonrpc"
)

type BasicMessage struct {
	Code     int 		  `json:"code"`
	Message  string       `json:"message"`
	Result   interface{}  `json:"result"`
}

type PostRecordRequest struct {
	Key    string  `json:"key"`
	Value  string  `json:"value"`
	Type   string  `json:"type"`
}

type Request struct {
	Height int `json:"height"`
}

type CreatAccountMsg struct {
	Address    string     `json:"address"`
	Pubkey     string 	  `json:"pubkey"`
	Timestamp  string     `json:"timestamp"`
	Sig        string     `json:"sig"`
}

type GetTransactionByHashRequest struct {
	Hash string `json:"hash"`
}
/*
func main(){
	rpc, err := jsonrpc.Dial("tcp", "127.0.0.1:8010");
	if err != nil {
		log.Fatal(err);
	}
	var basic BasicMessage

	var request Request
	request.Height = 82

	//调用远程方法
	//注意第三个参数是指针类型
	//GetTransactionByHash  PostRecord  GetTransactionsInBlockGroup GetTransactionsInBlock
	err2 := rpc.Call("PPoV.GetBlockGroupByHeight",request, &basic);
	if err2 != nil {
		log.Fatal(err2);
	}
	data, _ := json.Marshal(basic)
	fmt.Println(string(data))
}

*/
type GetTransactionsInBlockRequest struct {
	Height int `json:"height"`
	BlockNums int `json:"blockNums"`
}/*
func main(){
	rpc, err := jsonrpc.Dial("tcp", "127.0.0.1:8010");
	if err != nil {
		log.Fatal(err);
	}
	var basic BasicMessage

	var request GetTransactionsInBlockRequest
	request.Height = 35
	request.BlockNums =1

	//调用远程方法
	//注意第三个参数是指针类型
	//GetTransactionByHash  PostRecord  GetTransactionsInBlockGroup GetTransactionsInBlock
	err2 := rpc.Call("PPoV.GetTransactionsInBlock",request, &basic);
	if err2 != nil {
		log.Fatal(err2);
	}
	data, _ := json.Marshal(basic)
	fmt.Println(string(data))
}
*/
/*
func main(){
	//连接远程rpc服务
	//这里使用jsonrpc.Dial
	rpc, err := jsonrpc.Dial("tcp", "127.0.0.1:8010");
	if err != nil {
		log.Fatal(err);
	}
	var basic BasicMessage

	var request GetTransactionByHashRequest
	request.Hash = "UDwDNyi82OTesNGhJQLqLKjinH1EF83auWJXR0vcPZc="

	//调用远程方法
	//注意第三个参数是指针类型
	//GetTransactionByHash  PostRecord  GetTransactionsInBlockGroup GetTransactionsInBlock
	err2 := rpc.Call("PPoV.GetTransactionByHash",request, &basic);
	if err2 != nil {
		log.Fatal(err2);
	}
	data, _ := json.Marshal(basic)
	fmt.Println(string(data))
}*/

func main() {
	//连接远程rpc服务
	//这里使用jsonrpc.Dial
	rpc, err := jsonrpc.Dial("tcp", "127.0.0.1:8010");
	if err != nil {
		log.Fatal(err);
	}
	var basic BasicMessage


	var request CreatAccountMsg
	var km KeyManager.KeyManager
	km.Init()
	km.GenRandomKeyPair()
	request.Pubkey = km.GetPubkey()
	fmt.Println(len(request.Pubkey))
	request.Timestamp = "1111"
	request.Address = km.GetAddress()
	fmt.Println(request.Address)
	request.Sig = ""
	temp2, _ := json.Marshal(request)
	request.Sig ,err = km.Sign(KeyManager.GetHash(temp2))

	//调用远程方法
	//注意第三个参数是指针类型
	//GetTransactionByHash  PostRecord  GetTransactionsInBlockGroup GetTransactionsInBlock
	err2 := rpc.Call("PPoV.CreatAccount",request, &basic);
	if err2 != nil {
		log.Fatal(err2);
	}
	data, _ := json.Marshal(basic)
	fmt.Println(string(data))
}
