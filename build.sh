mkdir -p dist
(cd client && npm run build:prod)
go build -o=dist/main src/main.go
GOOS=windows GOARCH=amd64 go build -o dist/main.exe src/main.go