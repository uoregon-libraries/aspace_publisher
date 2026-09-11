FROM golang:1.25 AS build

RUN apt-get update
WORKDIR /usr/local/src/aspace_publisher

COPY go.mod ./
COPY go.sum ./
RUN go mod download && go mod verify

COPY ./handlers /usr/local/src/aspace_publisher/handlers
COPY ./views /usr/local/src/aspace_publisher/views
COPY ./alma /usr/local/src/aspace_publisher/alma
COPY ./utils /usr/local/src/aspace_publisher/utils
COPY ./oclc /usr/local/src/aspace_publisher/oclc
COPY ./marc /usr/local/src/aspace_publisher/marc
COPY ./file /usr/local/src/aspace_publisher/file
COPY ./connect /usr/local/src/aspace_publisher/connect
COPY ./aw /usr/local/src/aspace_publisher/aw
COPY ./as /usr/local/src/aspace_publisher/as
COPY ./main.go /usr/local/src/aspace_publisher/main.go

RUN go build \
    -ldflags="-s -w" \
    -o /usr/local/src/aspace_publisher/server .

FROM golang:1.25

RUN apt-get update
RUN apt-get install --no-install-recommends -y php php-dom

ARG USER_UID=1000
ARG USER_GID=1000
# Create non-root user
RUN groupadd --gid $USER_GID appuser && useradd --uid $USER_UID --gid appuser appuser

COPY --from=build /usr/local/src/aspace_publisher/server /aspace_publisher/server

USER appuser

EXPOSE 3000
ENTRYPOINT ["/aspace_publisher/server"]
