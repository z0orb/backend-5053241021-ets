package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func tulisJSON(writer http.ResponseWriter, statusCode int, data any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	json.NewEncoder(writer).Encode(data)
}

func tulisError(writer http.ResponseWriter, statusCode int, pesan string) {
	tulisJSON(writer, statusCode, map[string]string{"error": pesan})
}

// ambilIDDariPath mengembalikan false jika {id} bukan angka (dibalas 400).
func ambilIDDariPath(request *http.Request) (int, bool) {
	id, err := strconv.Atoi(request.PathValue("id"))
	return id, err == nil
}
