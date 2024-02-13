# Stage  1: Build the Go app
FROM golang:alpine AS build

RUN apk update && \
  apk add --no-cache build-base

WORKDIR  /build

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=1 go build -v -o . ./cmd/htd/main.go

# Stage  2: Run the Go app
FROM alpine:latest

WORKDIR /app

# Copy the Go app binary from the builder stage
COPY --from=build /build/main .
COPY sql/ /app/sql/
COPY static/ app/static/
COPY .env /app/

CMD [ "./main" ]
