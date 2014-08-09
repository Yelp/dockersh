dockersh: test
	go fmt && go build -ldflags "-linkmode external -extldflags -static"
	strip dockersh
 
test:
	go test

clean:
	rm -f dockersh
	go fmt

