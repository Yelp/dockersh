dockersh:
	go fmt && go build
	sudo chown root:root dockersh
	sudo chmod u+s dockersh
 
clean:
	rm -f dockersh

