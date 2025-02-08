package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Zensey/go-archetype-project/service"
	"github.com/Zensey/go-archetype-project/service/eth/protocol"
)

func SetHttpHandlers(observer service.Parser) {
	http.HandleFunc("/current-block", func(w http.ResponseWriter, r *http.Request) {
		getCurrentBlockID(w, r, observer)
	})
	http.HandleFunc("/transactions", func(w http.ResponseWriter, r *http.Request) {
		getTransactions(w, r, observer)
	})
	http.HandleFunc("/subscribe", func(w http.ResponseWriter, r *http.Request) {
		subscribe(w, r, observer)
	})
}

func getCurrentBlockID(w http.ResponseWriter, _ *http.Request, observer service.Parser) {
	fmt.Fprintf(w, "%d", observer.GetCurrentBlock())
}

func getTransactions(w http.ResponseWriter, r *http.Request, observer service.Parser) {
	address := r.URL.Query().Get("address")
	if address == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	trx := observer.GetTransactions(address)
	json.NewEncoder(w).Encode(trx)
}

func subscribe(w http.ResponseWriter, r *http.Request, observer service.Parser) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	address := r.FormValue("address")
	if address == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// validate address value as hex
	_, err = protocol.ParseInt(address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	res := observer.Subscribe(address)
	json.NewEncoder(w).Encode(res)
}
