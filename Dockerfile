FROM golang:1.23.0-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN mkdir -p /output

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /output/describe ./cmd/describe
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /output/histogram ./cmd/histogram
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /output/logregpredict ./cmd/logregpredict
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /output/logregtrain ./cmd/logregtrain
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /output/pairplot ./cmd/pairplot
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /output/scatterplot ./cmd/scatterplot

RUN chmod +x /output/*

CMD cp -r /output/* /binaries/
