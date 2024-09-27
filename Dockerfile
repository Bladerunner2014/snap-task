# Latest golang image on alpine linux
FROM golang:1.23-alpine

# Work directory
WORKDIR /app

# Installing dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copying all the files
COPY . .

# Starting our application
CMD ["sh", "-c", "cd cmd && go run main.go"]

# Exposing server port
EXPOSE 8080

