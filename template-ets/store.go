package main

import "sync"

type Store struct {
	mutex            sync.Mutex
	daftarSumberDaya []SumberDaya
	daftarPeminjaman []Peminjaman
	nextSumberDayaID int
	nextPeminjamanID int
}

func newStore() *Store {
	return &Store{
		daftarSumberDaya: []SumberDaya{},
		daftarPeminjaman: []Peminjaman{},
		nextSumberDayaID: 1,
		nextPeminjamanID: 1,
	}
}

// cariSumberDaya dipanggil saat mutex sudah terkunci.
func (store *Store) cariSumberDaya(id int) (SumberDaya, bool) {
	for _, sumberDaya := range store.daftarSumberDaya {
		if sumberDaya.ID == id {
			return sumberDaya, true
		}
	}
	return SumberDaya{}, false
}
