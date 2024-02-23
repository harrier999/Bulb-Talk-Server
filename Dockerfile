FROM golang:1.22

WORKDIR /usr/src/build

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY . .

RUN go mod download && go mod verify

RUN go build cmd/talk-server/main.go

RUN rm -rf /usr/src/build

WORKDIR /usr/src/app/cmd/talk-server

CMD go run main.go

