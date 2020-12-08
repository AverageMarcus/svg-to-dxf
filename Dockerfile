FROM jess/inkscape

RUN apt-get update && apt-get install -y python-lxml python-numpy

WORKDIR /app
RUN apt-get update && apt-get install -y golang
ADD main.go .
RUN go build -o main main.go

ENTRYPOINT ["/app/main"]
