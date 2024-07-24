FROM golang:1.22-alpine

WORKDIR /app

RUN apk add --no-cache git

COPY . .

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOPRIVATE=github.com
RUN go build ./cmd/app/main.go

CMD [ "/app/main" ]