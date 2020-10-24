FROM golang:1.14-stretch

RUN go get github.com/cespare/reflex && \
    go get github.com/go-delve/delve/cmd/dlv

COPY . /app
WORKDIR /app

CMD ["reflex", "-c", "reflex.conf"]
