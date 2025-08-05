FROM golang:1.24

ENV PROJECT_DIR=/app \
    GO111MODULE=on \
    CGO_ENABLED=0

WORKDIR /app
RUN mkdir "/build"
COPY . .

# copy migrations file to image
COPY ./migrations /app/migrations


RUN go get github.com/githubnemo/CompileDaemon
RUN go install github.com/githubnemo/CompileDaemon
ENTRYPOINT CompileDaemon -build="go build -o /build/app -buildvcs=false ./cmd/api_server" -command="/build/app"
