export TLS_PORT=3333
export TLS_CERTIFICATE=server/127.0.0.1/cert.pem
export TLS_PRIVATE_KEY=server/127.0.0.1/key.pem
export TLS_CA_CERTIFICATE=ca/minica.pem

go run server/cmd/serv-server/main.go