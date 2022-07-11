package db

type KV struct {
	//DB *badger.DB
}

func (kv *KV) Init() {

	//opt := badger.DefaultOptions("").WithInMemory(true)

	// thats way too much logic for init, shouldnt this be a constructor?
	//KVStoreConfig := config.New().Service.DNA
	//fmt.Println(KVStoreConfig)

	//// It will be created if it doesn't exist.
	//db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
	//
	//kv.DB = db
	//
	//if err != nil {
	//fmt.Println(err)
	//}
	//err = db.Update(func(txn *badger.Txn) error {
	//	txn.Set([]byte("author"), []byte("relix"))
	//	return nil
	//})
	//if err != nil {
	//fmt.Println(err)
	//}
	//defer db.Close()
	//// Your code here…
	//db.View(func(txn *badger.Txn) error {
	//	name, err := txn.Get([]byte("author"))
	//	fmt.Println(err)
	//	fmt.Println(name)
	//	return nil
	//})
}

func (kv *KV) Set(key string, value string) {

}

func (kv *KV) Get(key string) {
	//return kv.DB.ge
}
