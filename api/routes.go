package api

import (
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/add-record", AddRecord).Methods("POST")
	router.HandleFunc("/api/v1/blockchain", GetBlockchain).Methods("GET")
	router.HandleFunc("/api/v1/get-record/{cid}", GetFileFromIPFS).Methods("GET")
	return router
}
