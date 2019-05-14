FROM golang

# Fetch dependencies
RUN go get -u github.com/golang/dep/cmd/dep

# Add project directory to Docker image.
ADD . /go/src/github.com/punocracy/punocracy

ENV USER alvaro
ENV HTTP_ADDR :8888
ENV HTTP_DRAIN_INTERVAL 1s
ENV COOKIE_SECRET zu7HZy1Da2abXWPP

# Replace this with actual PostgreSQL DSN.
ENV DSN $GO_BOOTSTRAP_MYSQL_DSN

WORKDIR /go/src/github.com/punocracy/punocracy

RUN dep init
RUN dep ensure

RUN go build

EXPOSE 8888

CMD ./punocracy
