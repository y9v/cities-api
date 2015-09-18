FROM golang

ADD . /go/src/github.com/lebedev-yury/cities

ENV GIN_MODE release
WORKDIR /go/src/github.com/lebedev-yury/cities
RUN make build

ENTRYPOINT /go/src/github.com/lebedev-yury/cities/cities

EXPOSE 8080
