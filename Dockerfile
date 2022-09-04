FROM golang

RUN go version
ENV GOPATH=/
git
COPY ./ ./

RUN go mod download
RUN go build -o todo-app ./cmd/main.go

CMD ["./todo-app"]
