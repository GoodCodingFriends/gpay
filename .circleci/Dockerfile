FROM golang:1.10

RUN apt update -y && \
      apt install -y mysql-client mysql-server
RUN curl https://glide.sh/get | sh && \
      go get -u golang.org/x/lint/golint
