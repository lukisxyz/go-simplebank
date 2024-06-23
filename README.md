## Instalasi

1. Clone repository: `git clone https://github.com/flukis/go-simplebank.git`
2. Masuk ke direktori: `cd go-simplebank`
3. Instal dependensi: `go mod download`

## Penggunaan

1. Jalankan server: `make dev`
2. Gunakan klien RESTful API (seperti Postman) untuk berinteraksi dengan endpoint API.

## Endpoint API

API menyediakan endpoint berikut:

- **POST /api/auth/signup:** Membuat akun pengguna baru
- **POST /api/auth/login:** Mengotentikasi dan masuk sebagai pengguna
- **GET /api/accounts/:id:** Mendapatkan detail akun berdasarkan ID
- **POST /api/accounts:** Membuat akun baru
- **GET /api/accounts:** Mendapatkan daftar semua akun
- **POST /api/transfers:** Membuat transfer baru antara dua akun
- **GET /api/transfers/:id:** Mendapatkan detail transfer berdasarkan ID
- **GET /api/transfers:** Mendapatkan daftar semua transfer
