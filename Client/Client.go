package Client
import (
	"encoding/json"
	"fmt"
	"log"
	"net/rpc/jsonrpc"
	"../KeyManager"
	"time"
)


type BasicMessage struct {
	Code     int 		  `json:"code"`
	Message  string       `json:"message"`
	Result   interface{}  `json:"result"`
}

type Request struct {

}

/*
*12.上传一条记录
 */
type PostRecordRequest struct {
	Key    string  `json:"key"`
	Value  string  `json:"value"`
	Type   string  `json:"type"`
}

/*
*13.查询一条记录
 */
type GetRecordRequest struct {
	Key    string  `json:"key"`
	Type   string  `json:"type"`
}

/*
*14.根据哈希查询一条交易
 */
type GetTransactionByHashRequest struct {
	Hash string `json:"hash"`
}

/*
*15.获取指定区组块内的交易信息
 */
type GetTransactionsInBlockGroupRequest struct {
	Height int `json:"height"`
}

/*
*16.获取指定区组块内指定高度的交易信息
 */
type GetTransactionsInBlockRequest struct {
	Height int `json:"height"`
	BlockNums int `json:"blockNums"`
}

/*
*17.创建账户
 */
type CreatAccountMsg struct {
	Address    string     `json:"address"`
	Pubkey     string 	  `json:"pubkey"`
	Timestamp  string     `json:"timestamp"`
	Sig        string     `json:"sig"`
}

/*
*18.转账
 */
type TransferMoneyMsg struct {
	From       string     `json:"from"`
	To     	   string 	  `json:"to"`
	Pubkey	   string 	  `msg:"pubkey"`
	Amount     int 	      `json:"amount"`
	Timestamp  string     `json:"timestamp"`
	Sig        string     `json:"sig"`
}

/*
*19.查询余额
 */
type GetBalanceMsg struct {
	Address       string     `json:"address"`
}



func Client() {
	//连接远程rpc服务
	//这里使用jsonrpc.Dial
	rpc, err := jsonrpc.Dial("tcp", "127.0.0.1:8010");
	if err != nil {
		log.Fatal(err);
	}
	var choice int
	fmt.Print("请选择需要的函数: \n 1.获取当前区块高度信息\n 2.获取指定高度的区块组\n 3.获取某个区块组的所有区块\n 4.获取某个具体的区块信息\n 5.获取某个范围内所有的区块组\n 6.获取所有节点信息\n 7.根据公钥获取节点信息\n 8.获取投票节点列表\n 9.获取记账节点列表\n 10.获取候选记账节点列表\n 11.获取当前轮值记账节点\n 12.上传一条记录\n 13.查询一条记录\n 14.根据哈希查询一条交易\n 15.获取指定区组块内的交易信息\n 16.获取指定区组块内指定高度的交易信息\n 17.创建账户\n 18.转账\n 19.查询余额\n")
	fmt.Scanln(&choice)
	var basic BasicMessage
	var request Request

	switch choice{
	case 1:  //1.获取当前区块高度信息
		err2 := rpc.Call("PPoV.GetCurrentHeight",request, &basic);
		if err2 != nil {
			log.Fatal(err2);
		}
		data, _ := json.Marshal(basic)
		fmt.Println(string(data))

	case 2:  //2.获取指定高度的区块组
		err2 := rpc.Call("PPoV.GetBlockGroupByHeight",request, &basic);
		if err2 != nil {
			log.Fatal(err2);
		}
		data, _ := json.Marshal(basic)
		fmt.Println(string(data))

	case 3:  //3.获取某个区块组的所有区块
		err2 := rpc.Call("PPoV.GetBlocksInGroup",request, &basic);
		if err2 != nil {
			log.Fatal(err2);
		}
		data, _ := json.Marshal(basic)
		fmt.Println(string(data))

	case 4:  //4.获取某个具体的区块信息
		err2 := rpc.Call("PPoV.GetBlock",request, &basic);
		if err2 != nil {
			log.Fatal(err2);
		}
		data, _ := json.Marshal(basic)
		fmt.Println(string(data))

	case 5:  //5.获取某个范围内所有的区块组
		err2 := rpc.Call("PPoV.GetBlockRange",request, &basic);
		if err2 != nil {
			log.Fatal(err2);
		}
		data, _ := json.Marshal(basic)
		fmt.Println(string(data))

	case 6:  //6.获取所有节点信息
		err2 := rpc.Call("PPoV.GetNodeList",nil, &basic);
		if err2 != nil {
			log.Fatal(err2);
		}
		data, _ := json.Marshal(basic)
		fmt.Println(string(data))

	case 7:  //7.根据公钥获取节点信息
		err2 := rpc.Call("PPoV.GetNodeByPubkey",nil, &basic);
		if err2 != nil {
			log.Fatal(err2);
		}
		data, _ := json.Marshal(basic)
		fmt.Println(string(data))

	case 8:  //8.获取投票节点列表
		err2 := rpc.Call("PPoV.GetVoterList",nil, &basic);
		if err2 != nil {
			log.Fatal(err2);
		}
		data, _ := json.Marshal(basic)
		fmt.Println(string(data))

	case 9:  //9.获取记账节点列表
		err2 := rpc.Call("PPoV.GetWorkerList",nil, &basic);
		if err2 != nil {
			log.Fatal(err2);
		}
		data, _ := json.Marshal(basic)
		fmt.Println(string(data))

	case 10:  //10.获取候选记账节点列表
		err2 := rpc.Call("PPoV.GetWorkerCandidateList",nil, &basic);
		if err2 != nil {
			log.Fatal(err2);
		}
		data, _ := json.Marshal(basic)
		fmt.Println(string(data))

	case 11:  //11.获取当前轮值记账节点
		err2 := rpc.Call("PPoV.GetCurrentDutyWorker",nil, &basic);
		if err2 != nil {
			log.Fatal(err2);
		}
		data, _ := json.Marshal(basic)
		fmt.Println(string(data))

	case 12:  //12.上传一条记录
		var s PostRecordRequest
		fmt.Println("请输入Key：")
		fmt.Scanln(&s.Key)
		fmt.Println("请输入Value：")
		fmt.Scanln(&s.Value)
		fmt.Println("请输入Type：")
		fmt.Scanln(&s.Type)
		err2 := rpc.Call("PPoV.PostRecord",&s, &basic);
		if err2 != nil {
			log.Fatal(err2);
		}
		data, _ := json.Marshal(basic)
		fmt.Println(string(data))

	case 13:  //13.查询一条记录
		var s GetRecordRequest
		fmt.Println("请输入Key：")
		fmt.Scanln(&s.Key)
		fmt.Println("请输入Type：")
		fmt.Scanln(&s.Type)
		err2 := rpc.Call("PPoV.GetRecord",s, &basic);
		if err2 != nil {
			log.Fatal(err2);
		}
		data, _ := json.Marshal(basic)
		fmt.Println(string(data))

	case 14: //14.根据哈希查询一条交易
		var s GetTransactionByHashRequest
		fmt.Println("请输入Hash：")
		fmt.Scanln(&s.Hash)
		err2 := rpc.Call("PPoV.GetTransactionByHash",s, &basic);
		if err2 != nil {
			log.Fatal(err2);
		}
		data, _ := json.Marshal(basic)
		fmt.Println(string(data))

	case 15:  //15.获取指定区组块内的交易信息
		var s GetTransactionsInBlockGroupRequest
		fmt.Println("请输入Height：")
		fmt.Scanln(&s.Height)
		err2 := rpc.Call("PPoV.GetTransactionsInBlockGroup",s, &basic);
		if err2 != nil {
			log.Fatal(err2);
		}
		data, _ := json.Marshal(basic)
		fmt.Println(string(data))

	case 16: //16.获取指定区组块内指定高度的交易信息
		var s GetTransactionsInBlockRequest
		fmt.Println("请输入Height：")
		fmt.Scanln(&s.Height)
		fmt.Println("请输入BlockNums：")
		fmt.Scanln(&s.BlockNums)
		err2 := rpc.Call("PPoV.GetTransactionsInBlock",s, &basic);
		if err2 != nil {
			log.Fatal(err2);
		}
		data, _ := json.Marshal(basic)
		fmt.Println(string(data))

	case 17:  //17.创建账户
		var s CreatAccountMsg
		now := time.Now()
		s.Timestamp = fmt.Sprintf("%d", now.Unix())

		var km KeyManager.KeyManager
		km.Init()
		km.GenRandomKeyPair()
		s.Address = km.GetAddress()
		s.Pubkey = km.GetPubkey()
		s.Sig = ""
		temp2, _ := json.Marshal(s)
		s.Sig ,err = km.Sign(KeyManager.GetHash(temp2))

		err2 := rpc.Call("PPoV.CreatAccount",s, &basic);
		if err2 != nil {
			log.Fatal(err2);
		}
		data, _ := json.Marshal(basic)
		fmt.Println(string(data))

	case 18:  //18.转账
		var s TransferMoneyMsg
		fmt.Println("请输入From：")
		fmt.Scanln(&s.From)
		fmt.Println("请输入To：")
		fmt.Scanln(&s.To)

		now := time.Now()
		s.Timestamp = fmt.Sprintf("%d", now.Unix())
		fmt.Println("请输入Amount：")
		fmt.Scanln(&s.Amount)

		var km KeyManager.KeyManager
		km.Init()
		km.GenRandomKeyPair()
		s.Pubkey = km.GetPubkey()
		s.Sig = ""
		temp2, _ := json.Marshal(s)
		s.Sig ,err = km.Sign(KeyManager.GetHash(temp2))

		err2 := rpc.Call("PPoV.TransferMoney",s, &basic);
		if err2 != nil {
			log.Fatal(err2);
		}
		data, _ := json.Marshal(basic)
		fmt.Println(string(data))

	case 19:   //19.查询余额
		var s GetBalanceMsg
		var km KeyManager.KeyManager
		km.Init()
		km.GenRandomKeyPair()
		s.Address = km.GetAddress()
		err2 := rpc.Call("PPoV.GetBalance",s, &basic);
		if err2 != nil {
			log.Fatal(err2);
		}
		data, _ := json.Marshal(basic)
		fmt.Println(string(data))
	}

}
