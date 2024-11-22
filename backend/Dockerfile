FROM golang:1.22.2-alpine

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /opt
ADD .  /opt

RUN go build -o main ./main.go

EXPOSE 8080

CMD ["/opt/main"]