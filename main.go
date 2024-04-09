package main

import (
	_ "encoding/binary"
	"fmt"
	_ "unsafe"

	"./ConfigHelper"
	"./Node"
	"flag"
	_ "github.com/tinylib/msgp/msgp"
)

func Run() {
	var conf ConfigHelper.Config
	conf.ReadFile(file)
	var nodes []Node.Node
	//创建节点并设置参数
	for i := 0; i < conf.SingleServerNodeNum; i++ {
		var node Node.Node
		var parcel = conf
		parcel.MyAddress.Port += i
		parcel.ServicePort += i
		parcel.MyPubkey = parcel.PubkeyList[i]
		parcel.MyPrikey = parcel.PrikeyList[i]
		node.Init()
		node.SetConfig(parcel)
		nodes = append(nodes, node)
	}
	//初始化并启动节点
	for i := 0; i < len(nodes); i++ {
		//nodes[i].Init()
		go nodes[i].Start()
		fmt.Println("node", i, "start")
	}
	select {}
}

var file string
func init() {
	flag.StringVar(&file,"f","default","config file")
}
func main() {
	fmt.Println("欢迎使用PPoV区块链!")
	flag.Parse()
	//ConfigHelper.ConfigHelperTest()
	//ConfigHelper.CreateConfigs()

	Run()
}
