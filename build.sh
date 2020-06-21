echo "1. build ctr ..."
cd ctr
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../bin/ctr_linux_amd64
cd ..
echo "\tbuild ctr ... ok"

echo "2. build svr ..."
cd svr
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../bin/svr_linux_amd64
cd ..
echo "\tbuild svr ... ok"
