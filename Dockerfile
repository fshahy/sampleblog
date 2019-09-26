FROM golang:1.12.7

ENV SOURCES /src/github.com/fshahy/sampleblog
WORKDIR ${SOURCES}

COPY . ${SOURCES}

RUN go build -o server

ENTRYPOINT [ "./server" ]