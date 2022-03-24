
FROM golang:1.17

LABEL manteiner="Martin Fernandez <martin@nibiru.com.uy>"

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

ENV SESSION_SECRET_KEY="nibiruuy-store-session-secret-key"
ENV ACCESS_TOKEN_DURATION="15m"
ENV REFRESH_TOKEN_DURATION="24h"

RUN go build

CMD [ "./user-service" ]