FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o provisioner ./provisioning/cmd/server

FROM scratch
COPY --from=builder /app/provisioner /provisioner
EXPOSE 8080
ENTRYPOINT ["/provisioner"]
