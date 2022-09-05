FROM golang

RUN go version
ENV GOPATH=/

COPY . .

RUN go mod download
RUN go build todo-app ./cmd/main.go

EXPOSE 8080

CMD ["./todo-app"]
