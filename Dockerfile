FROM golang:1.21

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Copy the main repository files
COPY . .

# Initialize and update the submodule
RUN git submodule update --init --recursive

# Download dependencies
RUN go mod download

# Build the application
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o bankapi

EXPOSE 4100

# Run the application
CMD ["./bankapi"]