# Use an official Golang runtime as a parent image
FROM golang:1.17

# Set the working directory to /go/src/app
WORKDIR /go/src/app

# Copy the current directory contents into the container at /go/src/app
COPY . .

# Initialize the module, the name assumed here: loan-ms-go
RUN go mod init loan-ms-go

# Install any needed packages specified in go.mod
RUN go mod download

# Add dependencies explicitly
RUN go get -u github.com/go-redis/redis
RUN go get -u github.com/go-sql-driver/mysql
RUN go get -u github.com/gorilla/mux

# Build the Go app
RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
