# PROSES

## 1. Persona dan masalahnya
(2-3 kalimat dengan kata-katamu sendiri: siapa pengguna persona kamu, masalah apa yang diselesaikan fitur khasnya.)

Mahasiswa yang meminjam alat kampus mengeluh bahwa beberapa orang memborong alat sehingga yang lain tidak kebagian. Fitur khas dari persona ini menyelesaikan masalah tersebut dengan membuat batas kuota 2 peminjaman per orang per hari dan durasi maksimal 3 slot serta endpoint untuk melihat sisa kuota.

## 2. Rancangan struktur data
(Tulis struct dan slice/map apa yang kamu simpan, dan kenapa dipilih begitu. Rancang ini sebelum menulis kode, lalu catat kalau berubah di tengah jalan.)

Slice biasa dengan satu mutex, lalu bentrok dan kuota dihitung langsung dari slice peminjaman.
Karena: Urutan id sudah benar dimana data selalu di-append dengan id yang terus naik, jadi daftar otomatis terurut menurut id tanpa perlu sorting. Selain itu, engecekan bentrok, pengecekan kuota, dan penyimpanan bisa dikunci bersama yang memenuhi syarat B7. 

## 3. Satu keputusan desain
(Satu keputusan yang kamu ambil, pilihan lain yang kamu pertimbangkan, dan alasanmu memilih yang ini. Contoh: kenapa cek bentrok per slot, bukan per rentang.)

Keputusan: daftarSumberDaya []SumberDaya dan daftarPeminjaman []Peminjaman di dalam Store.
Pilihan lain: map[int]*Peminjaman dengan id sebagai key.
Alasan: Kontrak mewajibkan daftar terurut menurut id naik. Urutan map di Go acak, jadi setiap GET daftar nanti bakal harus di-sort dulu. Kalau lupa, pengujian urutan gagal. Jadi biar tidak usah uimplement sorting algorithm, keputusan tersebut saya ambil.

## 4. Bug atau kendala (minimal dua)
### Bug 1 (kendala)
- Gejala: curl -i -X POST localhost:8080/sumber-daya -d '{"nama":"  Proyektor  ","jenis":"alat","stok":3}'
curl: (7) Failed to connect to localhost:8080 after 0 ms: Could not connect to server
- Cara menemukan: run go main.go di terminal
- Penyebab dan solusi: go run main.go anya mengompilasi file main.go, jadi fungsi di store.go, helpers.go, dan file lain tidak ikut terbaca. Solusi adalah jalankan 'go run .'

### Bug 2
- Gejala: ❯ curl -i -X POST localhost:8080/peminjaman -d '{"sumber_daya_id":1,"peminjam":"  Rina ","tanggal":"2026-11-02","slot_mulai":3,"slot_selesai":4,"jumlah":1}'
HTTP/1.1 400 Bad Request
Content-Type: application/json
Date: Tue, 20 Oct 2026 09:36:08 GMT
Content-Length: 43

{"error":"sumber_daya_id tidak ditemukan"}
- Cara menemukan: Jalankan command curl -i -X POST localhost:8080/peminjaman -d '{"sumber_daya_id":1,"peminjam":"  Rina ","tanggal":"2026-11-02","slot_mulai":3,"slot_selesai":4,"jumlah":1}' di terminal langsung setelah go run . (abis terminate session sbelumnya).
- Penyebab dan solusi: Store kosong setiap kali server dijalankan ulang, sehingga harus membuat sumberdaya nya dulu baru bisa dipinjam. 

## 5. Change request
(Apa yang diminta, bagian kode mana yang kamu ubah, dan satu keputusan desain yang kamu ambil.)

5+0+5+3+2+4+1+0+2+1 = 23 mod 3 = 2. File C.

DELETE /peminjaman/{id} berubah makna menjadi pembatalan:
Status peminjaman menjadi "dibatalkan" dan datanya tetap ada. Respons tetap 204 tanpa body.
GET /peminjaman/{id} setelah dibatalkan menghasilkan 200 dengan status dibatalkan, dan peminjaman itu tetap muncul di GET /peminjaman (bisa disaring dengan ?status=dibatalkan).
Peminjaman yang dibatalkan tidak menahan slot: slot yang sama bisa langsung dipinjam orang lain.
Membatalkan peminjaman yang sudah berstatus dibatalkan menghasilkan 409.
Id bukan angka tetap 400 dan id tidak ditemukan tetap 404.

Bagian kode dirubah: peminjaman.go, nambah function handleBatalkanPeminjaman. main.go nambah store.handleBatalkanPeminjaman
