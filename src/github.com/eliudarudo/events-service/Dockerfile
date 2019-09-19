FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git

WORKDIR /app/serviceA

COPY go.mod .
COPY go.sum .
RUN  go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o serviceA


# Our final image
FROM scratch

WORKDIR /app

COPY --from=builder /app/serviceA .

ENTRYPOINT [ "./serviceA" ]