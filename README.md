# backend_xyz
Making backend about finance using Go and PostgreSQL

## Setup Project
- Install [Go](https://go.dev/doc/install)
- Konfigurasi file ```src/.env``` untuk mengakses database dan endpoint

## Setup Database
- Install [PostgreSQL](https://www.postgresql.org/download/)
- Make database
- Import ```database/xyz.pgsql``` ke database

## Compiles and hot-reloads for development
```
cd ./src
air
```

## Endpoint List
- `POST /api/auth/customer` : 
  Mendaftarkan customer baru (fullname, email, pass) dan menghasilkan token untuk authorization
- `POST /api/auth/login` : 
  Melakukan login oleh customer (email, pass) dan menghasilkan token untuk authorization
- `PUT  /api/customer` : 
  Update data customer kecuali password (fullname, email, pass), password hanya digunakan untuk authorization
- `PUT  /api/customer` : 
  Update password customer (oldPass, newPass)
- `POST /api/customer/details` : 
  Create dan/atau Update detail data customer (nik, legalname, place_of_birth, date_of_birth, salary, ktp_img, selfie_img) sekaligus Create dan/atau Update tenor limit yang dimiliki customer
- `GET  /api/tenorLimit?monthTenor=` : 
  Menampilkan data tenor limit yang dimiliki customer dengan parameter (monthTenor)
- `POST /api/transaction` : 
  Create transaksi yang dilakukan oleh customer (otr, assetname, admin_fee, installment_amount, interest_amount)
- `GET  /api/transaction?otrMin=&otrMax=&month=&year=` : 
  Menampilkan data transaksi oleh pengguna dengan parameter (otrMin, otrMax, month, year)

## Others
- Anda dapat melihat keterangan pembuatan endpoint [disini](./doc/keterangan.txt)
- Anda dapat melihat arsitektur pembuatan database [disini](./doc/relasi_xyz_db.png)
- Anda dapat melihat rencana arsitektur aplikasi [disini](./doc/architectural_diagram.png)
