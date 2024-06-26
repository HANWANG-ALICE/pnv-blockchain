package MongoDB

import (
	"gopkg.in/mgo.v2/bson"
	"hash/crc32"
	"log"
	"strconv"

	"../ConfigHelper"
	"../MetaData"
	"../lib/mgo"
)

type Mongo struct {
	Pubkey string
	Height int
	Block  MetaData.BlockGroup
}

func (pl *Mongo) SetConfig(config ConfigHelper.Config) {
	pl.Pubkey = config.MyPubkey
	pl.Height=-1
	if config.DropDatabase {
		pl.deleteDB()
	}
}

func ConnecToDB() *mgo.Session {
	session, err := mgo.Dial("mongodb://localhost")
	if err != nil {
		panic(err)
	}
	//defer session.Close()
	return session
}

func (pl *Mongo) InsertToMogo(block_group MetaData.BlockGroup, index string) {
	session := ConnecToDB()
	//session.SetMode(mgo.Monotonic, true)
	defer session.Close()

	index_mongo := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(index))))
	c := session.DB("blockchain").C(index_mongo)
	err := c.Insert(&block_group)
	if err != nil {
		log.Fatal(err)
	}
}

func (pl *Mongo) InsertToMogoRecord(item MetaData.KVRecord, index string) {
	session := ConnecToDB()
	//session.SetMode(mgo.Monotonic, true)
	defer session.Close()

	//index_mongo := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(index))))
	c := session.DB("blockchain").C(index)
	err := c.Insert(&item)
	if err != nil {
		log.Fatal(err)
	}
}

func (pl *Mongo)UpdateRecordToMongo(item map[string]interface{}, index string) {
	session := ConnecToDB()
	//session.SetMode(mgo.Monotonic, true)
	defer session.Close()

	selector := bson.M{"_id":item["_id"]}
	data := bson.M{"$set":bson.M{"value":item["value"]}}

	c := session.DB("blockchain").C(index)
	err := c.Update(selector,data)
	if err != nil {
		log.Fatal(err)
	}

}

func (pl *Mongo) GetHeight() int {
	return pl.Height
}

func (pl *Mongo) deleteDB() {
	session := ConnecToDB()
	//session.SetMode(mgo.Monotonic, true)
	defer session.Close()

	tmp := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	index1 := tmp + "-" + "Record"

	_ = session.DB("blockchain").C(index1).DropCollection()
	_ = session.DB("blockchain").C(tmp).DropCollection()
}