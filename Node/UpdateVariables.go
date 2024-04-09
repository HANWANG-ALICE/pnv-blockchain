package Node

import (
	"../KeyManager"
	"../MetaData"
	"../Network"
	"encoding/base64"
	"fmt"
	"sort"
	"time"
)

//交易排序需要
type TimePair struct {
	Key 	int
	Value 	float64
}
type TimePairList []TimePair

func (t TimePairList) Len()	int {
	return len(t)
}
func (t TimePairList) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
func (t TimePairList) Less(i, j int) bool {
	return t[i].Value < t[j].Value
}

func (node *Node) UpdateIdTransformationVaribles(transactionInterface MetaData.TransactionInterface) {
	if transaction, ok := transactionInterface.(*MetaData.IdentityTransformation); ok {
		switch transaction.Type {
		case "ApplyForVoter":
			_, ok := node.accountManager.VoterSet[transaction.Pubkey]
			if !ok {
				node.accountManager.VoterSet[transaction.Pubkey] = transaction.NodeID
			} else {
				fmt.Println("申请成为投票节点失败，已经是投票节点")
			}
			_, ok = node.network.NodeList[transaction.NodeID]
			if !ok {
				var nodelist Network.NodeInfo
				nodelist.IP = transaction.IPAddr
				nodelist.PORT = transaction.Port
				nodelist.ID = transaction.NodeID
				node.network.NodeList[transaction.NodeID] = nodelist
			}
		case "ApplyForWorkerCandidate":
			_, ok := node.accountManager.WorkerCandidateSet[transaction.Pubkey]
			if !ok {
				node.accountManager.WorkerCandidateSet[transaction.Pubkey] = transaction.NodeID
			} else {
				fmt.Println("申请成为候选记账节点失败，已经是候选记账节点")
			}
			_, ok = node.network.NodeList[transaction.NodeID]
			if !ok {
				var nodelist Network.NodeInfo
				nodelist.IP = transaction.IPAddr
				nodelist.PORT = transaction.Port
				nodelist.ID = transaction.NodeID
				node.network.NodeList[transaction.NodeID] = nodelist
			}
		case "QuitVoter":
			delete(node.accountManager.VoterSet, transaction.Pubkey)
			delete(node.network.NodeList, transaction.NodeID)
			fmt.Println("退出投票节点成功")
		case "QuitWorkerCandidate":
			delete(node.accountManager.WorkerCandidateSet, transaction.Pubkey)
			delete(node.network.NodeList, transaction.NodeID)
			fmt.Println("退出候选记账节点成功")
		}
	}
}

func (node *Node) UpdateRecordVaribles(transactionInterface MetaData.TransactionInterface) {
	if transaction, ok := transactionInterface.(*MetaData.Record); ok {
		if transaction.Command == MetaData.ADD{
			record := node.mongo.GetResultFromDatabase("Record","key",transaction.Key,"type",transaction.Type)
			_, ok := record["key"]
			if !ok{
				node.mongo.SaveRecordToDatabase("Record", *transaction)
			} else{
				record["value"] = transaction.Value
				node.mongo.UpdateRecordToDatabase("Record", record)
			}
			fmt.Println("TYPE",transaction.Type,"KEY:",transaction.Key, "VALUE",transaction.Value, "记录成功")
		}
	}
}

func (node *Node) UpdateCreatAccountTx(transactionInterface MetaData.TransactionInterface, mmap map[string]bool) bool {
	if transaction, ok := transactionInterface.(*MetaData.CreatAccount); ok {
		_, ok := mmap[transaction.Address]
		if !ok {
			_, existed := node.BalanceTable.Load(transaction.Address)
			if !existed {
				node.BalanceTable.Store(transaction.Address, 100)
				mmap[transaction.Address] = true
				return true
			} else {
				return false
			}
		} else {
			return false
		}
	}
	return false
}

func (node *Node) UpdateTransferMoneyTx(transactionInterface MetaData.TransactionInterface, mmap map[string]bool) bool {
	if transaction, ok := transactionInterface.(*MetaData.TransferMoney); ok {
		_, ok := mmap[transaction.From]
		if !ok {
			balance1, existed1 := node.BalanceTable.Load(transaction.From)
			balance2, existed2 := node.BalanceTable.Load(transaction.To)
			if existed1 && existed2 && transaction.Amount > 0 && balance1.(int) >= transaction.Amount {
				node.BalanceTable.Store(transaction.From, balance1.(int)-transaction.Amount)
				node.BalanceTable.Store(transaction.To, balance2.(int)+transaction.Amount)
				mmap[transaction.From] = true
				return true
			} else {
				return false
			}
		} else {
			return false
		}
	}
	return false
}

func (node *Node) UpdateGenesisVaribles(transactionInterface MetaData.TransactionInterface) {
	if genesisTransaction, ok := transactionInterface.(*MetaData.GenesisTransaction); ok {
		node.config.WorkerNum = genesisTransaction.WorkerNum
		node.config.VotedNum = genesisTransaction.VotedNum
		node.config.BlockGroupPerCycle = genesisTransaction.BlockGroupPerCycle
		node.config.Tcut = genesisTransaction.Tcut
		node.accountManager.WorkerSet = genesisTransaction.WorkerPubList
		node.accountManager.WorkerCandidateSet = genesisTransaction.WorkerCandidatePubList
		node.accountManager.VoterSet = genesisTransaction.VoterPubList
		var index uint32 = 0
		for _, key := range genesisTransaction.WorkerSet {
			node.accountManager.WorkerNumberSet[index] = key
			index = index + 1
		}
		index = 0
		for _, key1 := range genesisTransaction.VoterSet {
			node.accountManager.VoterNumberSet[index] = key1
			index = index + 1
		}
		for key2, _ := range genesisTransaction.WorkerCandidatePubList {
			node.accountManager.WorkerCandidateList = append(node.accountManager.WorkerCandidateList, key2)
		}

	}
}

func (node *Node) UpdateVaribles(bg *MetaData.BlockGroup) {
	if bg.Height > 0 { //normal blockgroup
		node.dutyWorkerNumber = bg.NextDutyWorker
		node.StartTime = bg.Timestamp

		//去重
		mmap := make(map[string]bool)
		
		//交易排序
		var tempTimePairList TimePairList
		for k,v := range bg.Blocks {
			pair := TimePair{
				Key:   k,
				Value: v.Timestamp,
			}
			tempTimePairList = append(tempTimePairList, pair)
		}
		sort.Sort(tempTimePairList)

		for _, v := range tempTimePairList {
			//test
			//if node.dutyWorkerNumber == node.GetMyWorkerNumber() {
			//	fmt.Println(v.Key)
			//}
			if bg.VoteResult[v.Key] != 1 {
				continue
			}
			for _, eachTransaction := range bg.Blocks[v.Key].Transactions {
				transactionHeader, transactionInterface := MetaData.DecodeTransaction(eachTransaction)
				switch transactionHeader.TXType {
				case MetaData.IdTransformation:
					node.UpdateIdTransformationVaribles(transactionInterface)
				case MetaData.Records:
					node.UpdateRecordVaribles(transactionInterface)

				case MetaData.CreatACCOUNT:
					res := node.UpdateCreatAccountTx(transactionInterface, mmap)
					if bg.ExecutionResult == nil{
						bg.ExecutionResult = make(map[string]bool)
					}
					bg.ExecutionResult[base64.StdEncoding.EncodeToString(KeyManager.GetHash(eachTransaction))] = res
				case MetaData.TransferMONEY:
					res := node.UpdateTransferMoneyTx(transactionInterface, mmap)
					if bg.ExecutionResult == nil{
						bg.ExecutionResult = make(map[string]bool)
					}
					bg.ExecutionResult[base64.StdEncoding.EncodeToString(KeyManager.GetHash(eachTransaction))] = res
				}
			}
			node.TxsAmount += uint64(len(bg.Blocks[v.Key].Transactions))
			node.TxsPeriodAmount += uint64(len(bg.Blocks[v.Key].Transactions))
		}
		node.BlockGroups.Store(bg.Height, *bg)
	}
}

func (node *Node) UpdateVariblesFromDisk(bg *MetaData.BlockGroup) {
	if bg.Height > 0 { //normal blockgroup
		node.dutyWorkerNumber = bg.NextDutyWorker
		node.StartTime = bg.Timestamp

		//去重
		mmap := make(map[string]bool)

		//交易排序
		var tempTimePairList TimePairList
		for k,v := range bg.Blocks {
			pair := TimePair{
				Key:   k,
				Value: v.Timestamp,
			}
			tempTimePairList = append(tempTimePairList, pair)
		}
		sort.Sort(tempTimePairList)

		for _, v := range tempTimePairList {
			if bg.VoteResult[v.Key] != 1 {
				continue
			}
			for _, eachTransaction := range bg.Blocks[v.Key].Transactions {
				transactionHeader, transactionInterface := MetaData.DecodeTransaction(eachTransaction)
				switch transactionHeader.TXType {
				case MetaData.IdTransformation:
					node.UpdateIdTransformationVaribles(transactionInterface)
				case MetaData.CreatACCOUNT:
					if bg.ExecutionResult[base64.StdEncoding.EncodeToString(KeyManager.GetHash(eachTransaction))] {
						node.UpdateCreatAccountTx(transactionInterface, mmap)
					}
				case MetaData.TransferMONEY:
					if bg.ExecutionResult[base64.StdEncoding.EncodeToString(KeyManager.GetHash(eachTransaction))] {
						node.UpdateTransferMoneyTx(transactionInterface, mmap)
					}
				}
			}
		}
	}
}

func (node *Node) UpdateGenesisBlockVaribles(bg *MetaData.BlockGroup) {
	if bg.Height == 0 { //genesis blockgroup
		node.dutyWorkerNumber = 0
		node.StartTime = bg.Timestamp
		if bg.Blocks[0].Height == 0 {
			transactionHeader, transactionInterface := MetaData.DecodeTransaction(bg.Blocks[0].Transactions[0])
			if transactionHeader.TXType == MetaData.Genesis {
				node.UpdateGenesisVaribles(transactionInterface)
			}
		}
		_ = node.state
		node.state <- Normal
		time.Sleep(time.Second)
	} else {
		fmt.Println("更新变量错误")
	}
}
