FROM golang:1.23-alpine

RUN apk add --no-cache git
RUN go install github.com/air-verse/air@v1.52.3

RUN addgroup -S apiGroup \
&& adduser -S -G apiGroup apiUser

WORKDIR /home/apiUser/app

COPY go.mod go.sum* ./
RUN go mod download

COPY --chown=apiUser:apiGroup . .

USER apiUser

EXPOSE 5000
CMD ["air", "-c", ".air.toml"]