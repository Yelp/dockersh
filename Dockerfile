FROM google/golang

ENV GOPATH $GOPATH:/gopath/src/github.com/docker/libcontainer/vendor
WORKDIR /gopath/src/github.com/Yelp/dockersh
ADD . /gopath/src/github.com/Yelp/dockersh/
RUN go get
RUN make dockersh && chmod 755 /gopath/src/github.com/Yelp/dockersh/installer.sh && ln /gopath/src/github.com/Yelp/dockersh/dockersh /dockersh && chown root:root dockersh && chmod u+s dockersh

CMD ["/gopath/src/github.com/Yelp/dockersh/installer.sh"]

