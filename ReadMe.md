
# NFT Payment

Proyek nft payment adalah sebuah aplikasi web yang menyediakan API untuk pembelian nft

## Pencegahan Keamanan

Langkah pencegahan
1. **Validasi Input**: melakukan validasi parameter

2. **Perlindungan Terhadap Serangan Injection**: menggunakan parameterized queries atau prepared statements saat berinteraksi dengan basis data untuk mencegah SQL Injection.
3. **Autentikasi**: menerapkan autentikasi yang kuat menggunakan token JWT (JSON Web Tokens)

## Running System

change app.yaml.example in folder config to app.yaml and setting the configuration based on your machine

Using Docker
```bash
  docker-compose up -d
```

install dependencies

```bash
  go mod tidy
```

running database migration (create table and seed data into the table)

```bash
  go run main.go db:migrate up
```


reset database (delete database and existing data)

```bash
  go run main.go db:migrate reset
```

running api server

```bash
  go run main.go api
```




## Tech Stack

**Database:** PostgresSQL

**Framework:** Echo golang

**Migration:** Goose
## API Reference



#### register user (example api contract was attached)

```http
  GET /v1/register
  
  param :
  username: {{username}}
  password: {{password}}
```

### login user
```http
  GET /v1/login
  
  param :
  username: {{username}}
  password: {{password}}
```

### Authorization

To Access endpoint always using bearer Authorization

```
Bearer {{token from login}}
```


#### example operation

postman file already attached in repo


