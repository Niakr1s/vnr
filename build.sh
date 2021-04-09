rm -rf dist
mkdir -p dist
(cd client && npm run build:prod)
echo "building executables"
go build -o=dist/vnr src/main.go
GOOS=windows GOARCH=amd64 go build -o dist/vnr.exe src/main.go