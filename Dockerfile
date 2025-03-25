FROM golang:1.23.7
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN make build

CMD ["make", "run"]
