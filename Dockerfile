FROM golang:1.13

WORKDIR /go/
COPY ./main.go ./

# Download all the dependencies
RUN go get -d github.com/gin-gonic/gin
RUN go get -d github.com/hoanganf/pos_domain/entity
RUN go get -d github.com/hoanganf/pos_domain/repository
RUN go get -d github.com/hoanganf/pos_domain/service
RUN go get -d github.com/hoanganf/file_upload/src

COPY ./templates ./templates
# Install the package
RUN go build -o main .

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the executable
CMD ["./main"]
