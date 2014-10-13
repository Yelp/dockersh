FROM google/golang

WORKDIR /gopath/src/github.com/Yelp/dockersh
ADD . /gopath/src/github.com/Yelp/dockersh/
RUN go get -d github.com/docker/libcontainer && cp -rv ../../docker/libcontainer/vendor/src/* ../../../ && go get
RUN make dockersh &&  chmod 755 /gopath/src/dockersh/installer.sh && ln /gopath/src/dockersh/dockersh /dockersh && chown root:root dockersh && chmod u+s dockersh

CMD ["/gopath/src/github.com/Yelp/dockersh/installer.sh"]

