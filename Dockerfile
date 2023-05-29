FROM golang:1.20

WORKDIR /app
COPY *.go go.mod go.sum ./
RUN go mod download
RUN go build -o ./stepmaniadb-backend
CMD ["./stepmaniadb-backend"]