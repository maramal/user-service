FROM golang:1.18 as builder

LABEL manteiner="Martin Fernandez <maramal@outlook.com>"

# Variables de entorno
ENV SESSION_SECRET_KEY="maramal-store-session-secret-key"
ENV ACCESS_TOKEN_DURATION="1h"
ENV REFRESH_TOKEN_DURATION="24h"


WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Construir aplicación
RUN go build -a -installsuffix cgo -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiar el binario de la pre-construcción del paso anterior
COPY --from=builder /app/main .

CMD [ "./main" ]