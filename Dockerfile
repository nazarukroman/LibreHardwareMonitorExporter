FROM golang:1.23

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./src/*.go ./

RUN go build -o /librehardwaremonitorexporter .

CMD  ["/librehardwaremonitorexporter"]
