FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o provisioner ./provisioning/cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o admin ./provisioning/cmd/admin
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o worker ./provisioning/cmd/worker

FROM scratch
COPY --from=builder /app/provisioner /provisioner
COPY --from=builder /app/admin /admin
COPY --from=builder /app/worker /worker
EXPOSE 8080
CMD ["/provisioner"]
