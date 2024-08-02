.PHONY: server client

server :
	cd server; go run ./main 5005

client :
	cd client; go run ./main localhost 5005 rich.bitch