FROM golang

ADD . /go/src/github.com/oliread/secretshop
ADD ./cmd/conf.toml /etc/secretshop-conf.toml

RUN go get -v -d github.com/oliread/secretshop/cmd
RUN go install github.com/oliread/secretshop/cmd

ENTRYPOINT /go/bin/cmd -conf /etc/secretshop-conf.toml

EXPOSE 8080