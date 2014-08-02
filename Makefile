dockersh: test
	go fmt && go build
 
test:
	go test

clean:
	rm -f dockersh
	go fmt

