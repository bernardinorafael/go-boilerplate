FROM golang:1.24 AS builder

# Add Maintainer info
LABEL maintainer="Rafael Bernardino <rafaelferreirab2@gmail.com>"

# Set the working directory in the container
WORKDIR /app

# Copy the go mod and sum files
COPY go.mod go.sum ./
# Download all the dependencies
RUN go mod download

# Copy the project files
COPY . /app

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o main cmd/api/main.go

FROM scratch AS prod
COPY --from=builder /app/main .

EXPOSE 8080

ENTRYPOINT ["./main"]