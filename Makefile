dockersh: dockersh_nosudo
	sudo chown root:root dockersh
	sudo chmod u+s dockersh
 
dockersh_nosudo: test
	go fmt && go build
 
test:
	go test

clean:
	rm -f dockersh

