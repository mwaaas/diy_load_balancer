run:
	rm -f main
	go build src/main/main.go
	./main

arpc_client:
	rm -f get_arp
	go build src/main/get_arp.go
	./get_arp -i tap0 -ip 10.10.1.4