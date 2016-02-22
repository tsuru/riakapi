FROM golang:1.5.3
MAINTAINER Xabier Larrakoetxea <slok69@gmail.com>

# Create the user/group for the running stuff
RUN groupadd -g 1000 dev
RUN useradd -m -u 1000 -g 1000 dev
RUN chown dev:dev -R /go

USER dev

# Install handy dependencies/tools
RUN go get github.com/Masterminds/glide
RUN go get golang.org/x/tools/cmd/cover
RUN go get github.com/axw/gocov/gocov
RUN go get github.com/mailgun/godebug


# Set environment variables
ENV RIAK_USER riakapi
ENV RIAK_PASSWORD riakapi
ENV SSH_USER riakapi
ENV SSH_PASSWORD riakapi
ENV GO15VENDOREXPERIMENT=1


WORKDIR /go/src/github.com/tsuru/riakapi
