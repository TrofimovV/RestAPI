FROM golang

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN go mod download
RUN go build ./cmd/main.go

CMD ["./main"]
