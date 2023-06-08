From golang:1.20.5-bullseye

WORKDIR /usr/local/src/aspace_publisher

COPY go.mod ./
COPY go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o /aspace_publisher/aspace_publisher

