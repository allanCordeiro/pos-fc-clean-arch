FROM golang:1.22.0 as builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /ordersystem cmd/api/main.go


FROM scratch
COPY --from=builder /ordersystem /ordersystem
COPY --from=builder usr/src/app/.env /.env


EXPOSE 8000
EXPOSE 8080
EXPOSE 50051
ENTRYPOINT [ "/ordersystem" ]