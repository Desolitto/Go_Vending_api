all:
	cd ./cmd/candy-server-server && go run main.go --port 3333

testMore:
	curl -XPOST -i -H "Content-Type: application/json" -d '{"money": 20, "candyType": "AA", "candyCount": 1}' http://127.0.0.1:3333/buy_candy

testLess:
	curl -XPOST -i -H "Content-Type: application/json" -d '{"money": 20, "candyType": "YR", "candyCount": 2}' http://127.0.0.1:3333/buy_candy

testWrong1:
	curl -XPOST -i -H "Content-Type: application/json" -d '{"money": 20, "candyType": "YR", "candyCount": -1}' http://127.0.0.1:3333/buy_candy

testWrong2:
	curl -XPOST -i -H "Content-Type: application/json" -d '{"money": -1, "candyType": "YR", "candyCount": 1}' http://127.0.0.1:3333/buy_candy

testWrong3:
	curl -XPOST -i -H "Content-Type: application/json" -d '{"money": 20, "candyType": "XX", "candyCount": 1}' http://127.0.0.1:3333/buy_candy

clean:
	rm -rf candy-client