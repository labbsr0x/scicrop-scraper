FROM golang:1.12.5-stretch as builder

RUN mkdir /scicrop-scraper
WORKDIR /scicrop-scraper

ADD go.mod .
#ADD go.sum .

RUN go mod download

ADD . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /go/bin/scicrop-scraper .

FROM alpine

COPY --from=builder /go/bin/scicrop-scraper /app/
COPY checkIsUp.sh /app/

WORKDIR /app
#CMD ["./scicrop-scraper"]
CMD ["sh","checkIsUp.sh"]
