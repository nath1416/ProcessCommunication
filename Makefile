BUILD=build/

all:
	go build -o go_server main.go
	g++ main.cpp -o cpu_server

run: all
	./go_server
clean:
	rm go_server