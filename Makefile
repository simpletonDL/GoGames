.PHONY: server client

server :
	cd server; go run ./main 5005

client :
	cd client; go run ./main 5005