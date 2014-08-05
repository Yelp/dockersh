FROM google/golang

WORKDIR /gopath/src/dockersh
ADD . /gopath/src/dockersh/
RUN go get
RUN make dockersh &&  chmod 755 /gopath/src/dockersh/installer.sh && ln /gopath/src/dockersh/dockersh /dockersh && chown root:root dockersh && chmod u+s dockersh

CMD ["/gopath/src/dockersh/installer.sh"]

