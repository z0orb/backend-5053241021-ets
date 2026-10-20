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
### Bug 1
- Gejala: curl -i -X POST localhost:8080/sumber-daya -d '{"nama":"  Proyektor  ","jenis":"alat","stok":3}'
curl: (7) Failed to connect to localhost:8080 after 0 ms: Could not connect to server
- Cara menemukan: run go main.go di terminal
- Penyebab dan solusi: go run main.go anya mengompilasi file main.go, jadi fungsi di store.go, helpers.go, dan file lain tidak ikut terbaca. Solusi adalah jalankan 'go run .'

### Bug 2
- Gejala:
- Cara menemukan:
- Penyebab dan solusi:

## 5. Change request
(Apa yang diminta, bagian kode mana yang kamu ubah, dan satu keputusan desain yang kamu ambil.)
