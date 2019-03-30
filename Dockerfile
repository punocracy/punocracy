FROM golang

# Fetch dependencies
RUN go get github.com/tools/godep

# Add project directory to Docker image.
ADD . /go/src/github.com/alvarosness/sample

ENV USER alvaro
ENV HTTP_ADDR :8888
ENV HTTP_DRAIN_INTERVAL 1s
ENV COOKIE_SECRET LwxL6R3E0iayWXnm

# Replace this with actual PostgreSQL DSN.
ENV DSN postgres://alvaro@localhost:5432/sample?sslmode=disable

WORKDIR /go/src/github.com/alvarosness/sample

RUN godep go build

EXPOSE 8888
CMD ./sample