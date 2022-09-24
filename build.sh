while getopts ah flag
do
    case "${flag}" in
        a) angular=true;;
        h) help=true
    esac
done

if [ $help ]
then
    echo "Usage:
    -h : Help
    -a : Build angular client"
    exit 0
fi

rm -rf dist
mkdir -p dist

if [ $angular ]
then
    echo "building client"
    (cd client && npm run build)
fi

echo "building executables"

GOOS=windows GOARCH=amd64 go build -ldflags "-H=windowsgui" -o dist/vnr.exe src/main.go
rcedit dist/vnr.exe --set-icon icon.ico