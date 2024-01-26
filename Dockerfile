########## БИЛД СТЭЙДЖ ##########
FROM golang:1.21 AS builder

# кешируем в порядке изменяемости слоев от редкого к частому
WORKDIR /src
COPY . .

# зависимости
RUN go mod download

# компилируем
RUN CGO_ENABLED=0 GOOS=linux go build -o ./binary cmd/main.go

########## РАН СТЭЙДЖ ##########
FROM alpine:latest

WORKDIR /app
COPY --from=builder /src/binary ./app
COPY --from=builder /src/config/config.yaml ./config/config.yaml

CMD ./app
