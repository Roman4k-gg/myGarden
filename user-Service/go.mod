module github.com/roman4k-gg/myGarden/user-Service

go 1.25.0

replace github.com/roman4k-gg/myGarden/pkg/user_v1 => ../pkg/user_v1

require (
	github.com/roman4k-gg/myGarden/pkg/user_v1 v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.82.0
)

require (
	github.com/golang-jwt/jwt/v5 v5.3.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.10.0 // indirect
	golang.org/x/crypto v0.54.0 // indirect
	golang.org/x/net v0.56.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.40.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260414002931-afd174a4e478 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)
