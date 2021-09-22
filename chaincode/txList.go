package main

import (
	tx "github.com/goledgerdev/cc-tools/transactions"
)

var txList = []tx.Transaction{
	tx.CreateAsset,
	tx.UpdateAsset,
	tx.DeleteAsset,
	//	txdefs.CreateNewLibrary,
	//	txdefs.GetNumberOfBooksFromLibrary,
	//	txdefs.UpdateBookTenant,
}
