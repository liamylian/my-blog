FROM golang:1.18-alpine AS builder

ENV CGO_ENABLED=0
ENV GOPRIVATE=""
ENV GOPROXY="https://goproxy.cn,direct"
ENV GOSUMDB="sum.golang.google.cn"

WORKDIR /root/my-blog/

ADD . .
RUN go mod download \
    && go test --cover $(go list ./... | grep -v /vendor/) \
    && go build -o main main.go

FROM alpine
ENV TZ Asia/Shanghai
WORKDIR /root/

COPY --from=builder /root/my-blog/main main
COPY --from=builder /root/my-blog/statics/ ./statics
RUN chmod +x main

ENTRYPOINT ["/root/main"]
