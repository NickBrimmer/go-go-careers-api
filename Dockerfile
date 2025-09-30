FROM golang:1.23-alpine

RUN apk add --no-cache git
RUN go install github.com/air-verse/air@latest

RUN addgroup -S apiGroup \
&& adduser -S -G apiGroup apiUser

WORKDIR /home/apiUser/app
COPY --chown=apiUser:apiGroup . .
RUN chown -R apiUser:apiGroup /home/apiUser/app

USER apiUser
RUN go mod download

EXPOSE 5000
CMD ["air", "-c", ".air.toml"]