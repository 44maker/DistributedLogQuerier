# DistributedLogQuerier

## Run

#在所有compute节点中
```bash
go run server.go
```
#在controller中
```bash
go build client.go
./client [query] [log file name]
```
If running for demo purpose, [log_file_name] may be vm1.log or /var/log/nova/nova-compute.log

## Unit Test:
Generate unit test:
```bash
go run generate_testfiles.go
```
#在所有compute节点中
```bash
go run server.go
```
#在controller中
```bash
go test -v client_test.go
```

## Result
![result1.png](https://cdn.ipfsscan.io/ipfs/QmeMx26j5nSqrzqYfcZeLWuorUzg34hzaNV2XxPaHGGZSD?filename=image.png)
![result2.png](https://cdn.ipfsscan.io/ipfs/QmX8TSsf83v5WvFAgNpUifyztkHKDrM6V7PhdXjq8AFeoT?filename=Results.png)
