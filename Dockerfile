FROM golang:1.14-stretch

RUN go get github.com/cespare/reflex

COPY . /app
WORKDIR /app

CMD ["reflex", "-c", "reflex.conf"]
