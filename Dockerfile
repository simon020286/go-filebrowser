FROM golang:alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine

WORKDIR /root

COPY --from=builder /app/main .
COPY --from=builder /app/static/ ./static

ENV BASEPATH /

EXPOSE 8080

CMD ["./main"]