FROM jess/inkscape

RUN apt-get update && apt-get install -y python-lxml python-numpy wget

RUN wget https://golang.org/dl/go1.16.2.linux-amd64.tar.gz && rm -rf /usr/local/go && tar -C /usr/local -xzf go1.16.2.linux-amd64.tar.gz && ln -s /usr/local/go/bin/go /usr/local/bin/go && go version

WORKDIR /app
ADD main.go index.html ./
RUN go build -o main main.go

ENTRYPOINT ["/app/main"]
