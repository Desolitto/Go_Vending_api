```
swagger generate server -f swagger.yaml
swagger generate client -f swagger.yaml
go mod tidy


1st terminal:
go run cmd/candy-server-server/main.go --port=3333

2d terminal:
curl -XPOST -H "Content-Type: application/json" -d '{"money": 20, "candyType": "AA", "candyCount": 1}' http://127.0.0.1:3333/buy_candy

```