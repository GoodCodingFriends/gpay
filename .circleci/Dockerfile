FROM golang:1.10

RUN apt update -y && \
      apt install -y mysql-client mysql-server
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh && \
      go get -u golang.org/x/lint/golint
