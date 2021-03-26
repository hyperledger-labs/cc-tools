package transactions

var txList = []Transaction{}

var basicTxs = []Transaction{
	getTx,
	GetSchema,
	GetDataTypes,
	CreateAsset,
	ReadAsset,
	ReadAssetHistory,
	Search,
	UpdateAsset,
	DeleteAsset,
	DeleteRecursive,
}

// TxList returns a copy of the txList variable
func TxList() []Transaction {
	listCopy := []Transaction{}
	for _, v := range txList {
		listCopy = append(listCopy, v)
	}
	return listCopy
}

// FetchTx returns a pointer to the Transaction object or nil if tx is not found
func FetchTx(txName string) *Transaction {
	for _, tx := range txList {
		if tx.Tag == txName {
			return &tx
		}
	}
	return nil
}

// InitTxList appends GetTx to txList to avoid initialization loop
func InitTxList(l []Transaction) {
	txList = append(l, basicTxs...)
}
