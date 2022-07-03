FROM golang:1.17-alpine

WORKDIR /data

COPY . .
RUN go build -o /data/bin-query

CMD [ "/data/bin-query" ]
