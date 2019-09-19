FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git

WORKDIR /app/serviceB

COPY go.mod .
COPY go.sum .
RUN  go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o serviceB


# Our final image
FROM scratch

WORKDIR /app

COPY --from=builder /app/serviceB .

ENTRYPOINT [ "./serviceB" ]