From golang:1.20.5-bullseye

RUN apt-get update
RUN apt-get install -y nano
WORKDIR /usr/local/src/aspace_publisher

COPY go.mod ./
COPY go.sum ./
RUN go mod download && go mod verify

RUN apt-get update
RUN apt-get install --no-install-recommends -y php php-dom
