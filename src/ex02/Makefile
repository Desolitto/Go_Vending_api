all:
	go run cmd/candy-server-server/main.go --tls-certificate localhost/cert.pem --tls-key localhost/key.pem --tls-port=3333

build_client:
	go build -o candy-client main_client.go

testSucces: build_client
	./candy-client -k AA -c 2 -m 50

testFail: build_client
	./candy-client -k AA -c 4 -m 10

clean:
	rm -rf candy-client