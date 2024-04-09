package MongoDB

import (
	"hash/crc32"
	"log"
	"strconv"

	"../MetaData"
	"../lib/mgo/bson"
)

func (pl *Mongo) QueryHeight() int {
	session := ConnecToDB()

	//session.SetMode(mgo.Monotonic, true)
	defer session.Close()

	var blocks []MetaData.BlockGroup
	var height int = -1
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	c := session.DB("blockchain").C(index)
	err := c.Find(nil).Sort("-height").Limit(1).All(&blocks)
	//err = c.Find(nil).All(&blocks)
	if err != nil {
		log.Println(err)
	}
	for _, x := range blocks {
		if x.Height > height {
			height = x.Height
		}
	}
	return height
}

func (pl *Mongo) GetAmount() int {
	return pl.QueryHeight() + 1
}
func (pl *Mongo) PushbackBlockToDatabase(block MetaData.BlockGroup) {
	pl.InsertToMogo(block, pl.Pubkey)
	pl.Block = block
	pl.Height = block.Height

}

func (pl *Mongo) GetBlockFromDatabase(height int) MetaData.BlockGroup {
	session := ConnecToDB()
	//session.SetMode(mgo.Monotonic, true)
	defer session.Close()

	var block MetaData.BlockGroup
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	c := session.DB("blockchain").C(index)
	err := c.Find(bson.M{"height": height}).One(&block)
	if err != nil {
		log.Println(err)
	}
	return block
}

func (pl *Mongo) GetBlockByTxHashFromDatabase(hash []byte) MetaData.BlockGroup{
	session := ConnecToDB()
	//session.SetMode(mgo.Monotonic, true)
	defer session.Close()

	var block MetaData.BlockGroup
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	c := session.DB("blockchain").C(index)
	err := c.Find(bson.M{"blocks":bson.M{"$elemMatch":bson.M{"transactionshash":bson.M{"$elemMatch":bson.M{"$eq":hash}}}}}).One(&block)

	if err != nil {
		log.Println(err)
	}
	return block
}
