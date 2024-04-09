package MetaData

import "fmt"

//go:generate msgp
type TransferMoney struct {
	From       string     `msg:"from"`
	To     	   string 	  `msg:"to"`
	Pubkey	   string 	  `msg:"pubkey"`
	Amount     int 	      `msg:"amount"`
	Timestamp  string     `msg:"timestamp"`
	Sig        string     `msg:"sig"`
}

func (tm *TransferMoney) ToByteArray() []byte {
	data, _ := tm.MarshalMsg(nil)
	return data
}

func (tm *TransferMoney) FromByteArray(data []byte) {
	_, err := tm.UnmarshalMsg(data)
	if err != nil {
		fmt.Println("err=", err)
	}
}
