dockersh: test
	go fmt && go build
	sudo chown root:root dockersh
	sudo chmod u+s dockersh
 
test:
	go test

clean:
	rm -f dockersh

