FROM golang:1.26.10-alpine@sha256:5cb67f19c5e7994dcb80c074f7cfca8aeeada45bbdca035f73df269cefdcce61 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o provisioner ./provisioning/cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o admin ./provisioning/cmd/admin
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o worker ./provisioning/cmd/worker

FROM scratch
USER 65532:65532
COPY --from=builder /app/provisioner /provisioner
COPY --from=builder /app/admin /admin
COPY --from=builder /app/worker /worker
EXPOSE 8080
CMD ["/provisioner"]
