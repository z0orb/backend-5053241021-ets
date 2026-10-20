package main

import "net/http"

func NewHandler() http.Handler {
	store := newStore()
	mux := http.NewServeMux()

	mux.HandleFunc("POST /sumber-daya", store.handleBuatSumberDaya)
	mux.HandleFunc("GET /sumber-daya", store.handleDaftarSumberDaya)
	mux.HandleFunc("GET /sumber-daya/{id}", store.handleDetailSumberDaya)
	mux.HandleFunc("POST /peminjaman", store.handleBuatPeminjaman)
	return mux
}

func main() {
	http.ListenAndServe(":8080", NewHandler())
}
