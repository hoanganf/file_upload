FROM golang:1.13

WORKDIR /go/src/app
COPY ./src ./github.com/file-upload/
COPY ./templates ./github.com/file-upload/
COPY ./main.go ./github.com/file-upload/

# Download all the dependencies
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the executable
CMD ["main"]
