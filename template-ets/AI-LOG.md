# AI-LOG

Catat setiap penggunaan AI yang menghasilkan kode. Jika tidak memakai AI sama sekali, tulis satu baris: **tidak memakai AI**.

Log yang jujur tidak mengurangi nilai, sebanyak apa pun AI dipakai. Yang mengurangi nilai adalah log yang tidak sesuai dengan kode, atau isi log yang tidak bisa dijelaskan.

## (Nama fitur atau fungsi, misalnya: Deteksi bentrok slot, B5)
- Prompt: (ringkas, apa yang kamu minta)
- Dipakai: (bagian kode yang kamu ambil, misalnya nama fungsi dan file)
- Diubah: (apa yang kamu ubah atau perbaiki dari hasil AI, dan kenapa)

## (Persiapan workspace, mengisi konteks LLM dengan dokumen classroom)
- Prompt: JANGAN MELAKUKAN EDIT WRITE APAPUN PADA WORKSPACE INI. ACTION DIPERBOLEHKAN HANYA READ.  Baca Requirement ETS Backend - PinjamKampus untuk konteks project. Persona yang digunakan adalah PERSONA B. Jangan outputkan apa-apa dulu, simpan dalam context window.
- Dipakai: directory local ETS/ yang mengandung backend-5053241021-ets (repo submission github)
- Diubah: AI-LOG.md - Mencatat prompt

## (Penambahan aturan konvensi output LLM)
- Prompt: Output berupa codeblock dengan penjelasan unit code tersebut. Avoid penggunaan variable yang tidak jelas/magic number agar maintainable
- Dipakai: Tidak ada
- Diubah: AI-LOG.md - Mencatat prompt

## (Klarifikasi soal requirement)
- Prompt: Pertanyaan elaboration Maksud dari "siapa pengguna persona kamu dan masalah yang diselesaikan ini apa"
- Dipakai: Requirement ETS Backend - PinjamKampus.docx
- Diubah: AI-LOG.md - Mencatat prompt

## (Pilihan struktur data sebelum mulai coding)
- Prompt: sebutkan beberapa pilihan struktur data yang cocok untuk project ini, aku pilih
- Dipakai: Tidak ada
- Diubah: AI-LOG.md - Mencatat prompt

## (Pembuatan models.go project)
- Prompt: Buat models.go untuk memfasilitasi B1 & B2
- Dipakai: models.go
- Diubah: penambahan struct SumberDaya dan Peminjaman

## (Pembuatan helpers.go project)
- Prompt: Buat helpers.go untuk write JSON dan Error untuk test endpoint project ini. Sekalian function untuk mengambil ID dari path endpooint.
- Dipakai: helpers.go
- Diubah: penambahan function tulisJSON, ambilIDDariPath, dan tulisError

## (Implementasi sumber_daya.go (B1))
- Prompt: Implementasikan B1 project 
- Dipakai: sumber_daya.go
- Diubah: penambahan function validasiSumberDaya, handleBuatSumberDaya, handleDaftarSumberDaya, handleDetailSumberDaya

## (Implementasi peminjaman.go (B2))
- Prompt: Implementasikan B2 project 
- Dipakai: peminjaman.go
- Diubah: penambahan function validasiPeminjaman, handleBuatPeminjaman,

## (Implementasi Lihat peminjaman (B3))
- Prompt: Implementasikan B3 project 
- Dipakai: peminjaman.go, store.go, main.go
- Diubah: penambahan function cariIndexPeminjaman, library strconv di peminjaman.go, func handleDaftarPeminjaman, func handleDetailPeminjaman,

## (Implementasi Change Request PARSIAL)
- Prompt: Ada file baru di Folder ETS/: Lembar CR-C Pembatalan.docx. Baca, dan masukkan ke context window. Outputkan perubahan di peminjaman di chat ini.
- Dipakai: peminjaman.go, main.go
- Diubah: penambahan function handleBatalkanPeminjaman, handlefunc handleBatalkanPeminjaman. KARENA B6 BELUM ADA AS OF LOG INI.

## (Implementasi Validasi dan Error)
- Prompt: Implementasikan B4 project
- Dipakai: validasi_test.go
- Diubah: penambahan function TestValidasiB4
