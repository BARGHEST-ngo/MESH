FROM golang:1.26-alpine AS build

RUN apk add --no-cache git patch

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download

COPY analyst/src ./analyst
RUN go generate ./...
COPY .git ./.git
RUN ./analyst/build.sh

FROM alpine:latest

WORKDIR /usr/bin

COPY --from=build /usr/src/app/analyst/mesh /usr/bin/mesh
COPY analyst/meshcli /usr/bin/meshcli
COPY analyst/entrypoint.sh /usr/bin/entrypoint.sh
RUN chmod +x /usr/bin/entrypoint.sh

ENTRYPOINT ["/usr/bin/entrypoint.sh"]
CMD ["/usr/bin/meshcli"]