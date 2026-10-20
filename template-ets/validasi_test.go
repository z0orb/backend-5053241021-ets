package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// kirimRequest menjalankan satu request langsung ke handler, tanpa server sungguhan.
func kirimRequest(handler http.Handler, method, path, body string) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(method, path, strings.NewReader(body))
	handler.ServeHTTP(recorder, request)
	return recorder
}

func TestValidasiB4(t *testing.T) {
	handler := NewHandler()
	kirimRequest(handler, "POST", "/sumber-daya", `{"nama":"Proyektor","jenis":"alat","stok":3}`)

	nama50 := strings.Repeat("é", 50)
	nama51 := strings.Repeat("é", 51)
	daftarKasus := []struct {
		namaKasus        string
		method           string
		path             string
		body             string
		statusDiharapkan int
	}{
		{"body bukan JSON", "POST", "/sumber-daya", `abc`, 400},
		{"body kosong", "POST", "/sumber-daya", ``, 400},
		{"nama hanya spasi", "POST", "/sumber-daya", `{"nama":"   ","jenis":"alat","stok":1}`, 400},
		{"nama 51 karakter", "POST", "/sumber-daya", `{"nama":"` + nama51 + `","jenis":"alat","stok":1}`, 400},
		{"nama 50 karakter multibyte", "POST", "/sumber-daya", `{"nama":"` + nama50 + `","jenis":"alat","stok":1}`, 201},
		{"jenis salah", "POST", "/sumber-daya", `{"nama":"A","jenis":"kursi","stok":1}`, 400},
		{"stok 0", "POST", "/sumber-daya", `{"nama":"A","jenis":"alat","stok":0}`, 400},
		{"stok tidak ada", "POST", "/sumber-daya", `{"nama":"A","jenis":"alat"}`, 400},
		{"stok teks", "POST", "/sumber-daya", `{"nama":"A","jenis":"alat","stok":"3"}`, 400},
		{"stok desimal", "POST", "/sumber-daya", `{"nama":"A","jenis":"alat","stok":1.5}`, 400},
		{"tanggal 30 Feb", "POST", "/peminjaman", `{"sumber_daya_id":1,"peminjam":"Rina","tanggal":"2026-02-30","slot_mulai":1,"slot_selesai":1,"jumlah":1}`, 400},
		{"29 Feb 2028 valid", "POST", "/peminjaman", `{"sumber_daya_id":1,"peminjam":"Rina","tanggal":"2028-02-29","slot_mulai":1,"slot_selesai":1,"jumlah":1}`, 201},
		{"slot 9", "POST", "/peminjaman", `{"sumber_daya_id":1,"peminjam":"Rina","tanggal":"2026-11-02","slot_mulai":1,"slot_selesai":9,"jumlah":1}`, 400},
		{"slot terbalik", "POST", "/peminjaman", `{"sumber_daya_id":1,"peminjam":"Rina","tanggal":"2026-11-02","slot_mulai":4,"slot_selesai":3,"jumlah":1}`, 400},
		{"sumber daya tidak ada", "POST", "/peminjaman", `{"sumber_daya_id":99,"peminjam":"Rina","tanggal":"2026-11-02","slot_mulai":1,"slot_selesai":1,"jumlah":1}`, 400},
		{"jumlah melebihi stok", "POST", "/peminjaman", `{"sumber_daya_id":1,"peminjam":"Rina","tanggal":"2026-11-02","slot_mulai":1,"slot_selesai":1,"jumlah":4}`, 400},
		{"filter id bukan angka", "GET", "/peminjaman?sumber_daya_id=abc", ``, 400},
		{"detail id bukan angka", "GET", "/peminjaman/abc", ``, 400},
		{"detail tidak ada", "GET", "/peminjaman/99", ``, 404},
	}

	for _, kasus := range daftarKasus {
		hasil := kirimRequest(handler, kasus.method, kasus.path, kasus.body)
		if hasil.Code != kasus.statusDiharapkan {
			t.Errorf("%s: status %d, diharapkan %d", kasus.namaKasus, hasil.Code, kasus.statusDiharapkan)
		}
		if hasil.Header().Get("Content-Type") != "application/json" {
			t.Errorf("%s: Content-Type bukan application/json", kasus.namaKasus)
		}
		if hasil.Code >= 400 && !strings.Contains(hasil.Body.String(), `"error"`) {
			t.Errorf("%s: body error tidak punya field error: %s", kasus.namaKasus, hasil.Body.String())
		}
	}
}

// go test ./... utk run
