module users-api

go 1.20

require (
	github.com/bradfitz/gomemcache v0.0.0-20230905024940-24af94b03874
	github.com/go-sql-driver/mysql v1.8.1
	github.com/golang-jwt/jwt/v4 v4.5.0 // Para manejar los tokens JWT
	github.com/gorilla/mux v1.8.0 // Para manejar las rutas HTTP
	golang.org/x/crypto v0.12.0 // Para manejar el hashing de contraseñas
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
)
