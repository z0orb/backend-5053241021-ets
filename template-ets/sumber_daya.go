package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"unicode/utf8"
)

// validasiSumberDaya mengembalikan pesan error, atau "" jika valid.
func validasiSumberDaya(input SumberDaya) string {
	if input.Nama == "" {
		return "nama wajib diisi"
	}
	if utf8.RuneCountInString(input.Nama) > 50 {
		return "nama maksimal 50 karakter"
	}
	if input.Jenis != "ruang" && input.Jenis != "alat" {
		return "jenis harus 'ruang' atau 'alat'"
	}
	if input.Stok < 1 {
		return "stok minimal 1"
	}
	return ""
}

func (store *Store) handleBuatSumberDaya(writer http.ResponseWriter, request *http.Request) {
	var input SumberDaya
	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		tulisError(writer, http.StatusBadRequest, "body JSON tidak valid")
		return
	}
	input.Nama = strings.TrimSpace(input.Nama)
	input.Jenis = strings.TrimSpace(input.Jenis)

	if pesanError := validasiSumberDaya(input); pesanError != "" {
		tulisError(writer, http.StatusBadRequest, pesanError)
		return
	}

	store.mutex.Lock()
	defer store.mutex.Unlock()

	input.ID = store.nextSumberDayaID
	store.nextSumberDayaID++
	store.daftarSumberDaya = append(store.daftarSumberDaya, input)
	tulisJSON(writer, http.StatusCreated, input)
}

func (store *Store) handleDaftarSumberDaya(writer http.ResponseWriter, request *http.Request) {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	tulisJSON(writer, http.StatusOK, store.daftarSumberDaya)
}

func (store *Store) handleDetailSumberDaya(writer http.ResponseWriter, request *http.Request) {
	id, valid := ambilIDDariPath(request)
	if !valid {
		tulisError(writer, http.StatusBadRequest, "id harus angka")
		return
	}

	store.mutex.Lock()
	defer store.mutex.Unlock()

	sumberDaya, ditemukan := store.cariSumberDaya(id)
	if !ditemukan {
		tulisError(writer, http.StatusNotFound, "sumber daya tidak ditemukan")
		return
	}
	tulisJSON(writer, http.StatusOK, sumberDaya)
}
