package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// validasiPeminjaman mengembalikan pesan error, atau "" jika valid.
// Dipanggil saat mutex sudah terkunci, karena membaca stok sumber daya.
func (store *Store) validasiPeminjaman(input Peminjaman) string {
	if input.Peminjam == "" {
		return "peminjam wajib diisi"
	}
	if utf8.RuneCountInString(input.Peminjam) > 50 {
		return "peminjam maksimal 50 karakter"
	}
	if _, err := time.Parse("2006-01-02", input.Tanggal); err != nil {
		return "tanggal harus format YYYY-MM-DD dan tanggal yang valid"
	}
	if input.SlotMulai < 1 || input.SlotSelesai > 8 || input.SlotMulai > input.SlotSelesai {
		return "slot harus 1 sampai 8 dan slot_mulai <= slot_selesai"
	}
	sumberDaya, ditemukan := store.cariSumberDaya(input.SumberDayaID)
	if !ditemukan {
		return "sumber_daya_id tidak ditemukan"
	}
	if input.Jumlah < 1 || input.Jumlah > sumberDaya.Stok {
		return "jumlah harus antara 1 dan stok sumber daya"
	}
	return ""
}

func (store *Store) handleBuatPeminjaman(writer http.ResponseWriter, request *http.Request) {
	var input Peminjaman
	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		tulisError(writer, http.StatusBadRequest, "body JSON tidak valid")
		return
	}
	input.Peminjam = strings.TrimSpace(input.Peminjam)
	input.Tanggal = strings.TrimSpace(input.Tanggal)

	store.mutex.Lock()
	defer store.mutex.Unlock()

	if pesanError := store.validasiPeminjaman(input); pesanError != "" {
		tulisError(writer, http.StatusBadRequest, pesanError)
		return
	}

	// Nanti: cek bentrok (B5) dan aturan persona B (maks 3 slot, kuota 2) disisipkan di sini,
	// masih di dalam lock yang sama.

	input.ID = store.nextPeminjamanID
	input.Status = "aktif"
	store.nextPeminjamanID++
	store.daftarPeminjaman = append(store.daftarPeminjaman, input)
	tulisJSON(writer, http.StatusCreated, input)
}

func (store *Store) handleDaftarPeminjaman(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()

	filterSumberDayaID := 0
	adaFilterSumberDaya := false
	if teksID := query.Get("sumber_daya_id"); teksID != "" {
		angkaID, err := strconv.Atoi(teksID)
		if err != nil {
			tulisError(writer, http.StatusBadRequest, "sumber_daya_id harus angka")
			return
		}
		filterSumberDayaID = angkaID
		adaFilterSumberDaya = true
	}
	filterTanggal := query.Get("tanggal")
	filterStatus := query.Get("status")

	store.mutex.Lock()
	defer store.mutex.Unlock()

	hasil := []Peminjaman{}
	for _, peminjaman := range store.daftarPeminjaman {
		if adaFilterSumberDaya && peminjaman.SumberDayaID != filterSumberDayaID {
			continue
		}
		if filterTanggal != "" && peminjaman.Tanggal != filterTanggal {
			continue
		}
		if filterStatus != "" && peminjaman.Status != filterStatus {
			continue
		}
		hasil = append(hasil, peminjaman)
	}
	tulisJSON(writer, http.StatusOK, hasil)
}

func (store *Store) handleDetailPeminjaman(writer http.ResponseWriter, request *http.Request) {
	id, valid := ambilIDDariPath(request)
	if !valid {
		tulisError(writer, http.StatusBadRequest, "id harus angka")
		return
	}

	store.mutex.Lock()
	defer store.mutex.Unlock()

	index := store.cariIndexPeminjaman(id)
	if index == -1 {
		tulisError(writer, http.StatusNotFound, "peminjaman tidak ditemukan")
		return
	}
	tulisJSON(writer, http.StatusOK, store.daftarPeminjaman[index])
}

// handleBatalkanPeminjaman: CR-C, DELETE tidak menghapus data, hanya mengubah status menjadi "dibatalkan".
func (store *Store) handleBatalkanPeminjaman(writer http.ResponseWriter, request *http.Request) {
	id, valid := ambilIDDariPath(request)
	if !valid {
		tulisError(writer, http.StatusBadRequest, "id harus angka")
		return
	}

	store.mutex.Lock()
	defer store.mutex.Unlock()

	index := store.cariIndexPeminjaman(id)
	if index == -1 {
		tulisError(writer, http.StatusNotFound, "peminjaman tidak ditemukan")
		return
	}
	if store.daftarPeminjaman[index].Status == "dibatalkan" {
		tulisError(writer, http.StatusConflict, "peminjaman sudah dibatalkan")
		return
	}

	store.daftarPeminjaman[index].Status = "dibatalkan"
	writer.WriteHeader(http.StatusNoContent)
}
