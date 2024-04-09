package MetaData

import "fmt"

//go:generate msgp
type CreatAccount struct {
	Address    string     `msg:"address"`
	Pubkey     string 	  `msg:"pubkey"`
	Timestamp  string     `msg:"timestamp"`
	Sig        string     `msg:"sig"`
}

func (ca *CreatAccount) ToByteArray() []byte {
	data, _ := ca.MarshalMsg(nil)
	return data
}

func (ca *CreatAccount) FromByteArray(data []byte) {
	_, err := ca.UnmarshalMsg(data)
	if err != nil {
		fmt.Println("err=", err)
	}
}
