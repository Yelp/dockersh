FROM google/golang

WORKDIR /gopath/src/dockersh
ADD . /gopath/src/dockersh/
RUN go get
RUN make &&  chmod 755 /gopath/src/dockersh/installer.sh && ln /gopath/src/dockersh/dockersh /dockersh

CMD ["/gopath/src/dockersh/installer.sh"]

