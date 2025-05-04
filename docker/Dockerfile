FROM golang:1.22-alpine as builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /mcp

FROM alpine:3.21

COPY --from=builder /mcp /mcp

EXPOSE 8080

CMD [ "/mcp" ]
