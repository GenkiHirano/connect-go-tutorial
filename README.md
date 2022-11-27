# 目的

- connect-go を公式チュートリアルを通して、理解を深める。

# コマンド

## サーバー実行コマンド

```
make server
```

## 実行コマンド

```
make client
```

または

```
curl \
    --header "Content-Type: application/json" \
    --data '{"name": "Gopherくん"}' \
    http://localhost:8080/greet.v1.GreetService/Greet
```

または

```
grpcurl \
    -protoset <(buf build -o -) -plaintext \
    -d '{"name": "Gopherくん"}' \
    localhost:8080 greet.v1.GreetService/Greet
```
