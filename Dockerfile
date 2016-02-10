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


ENV GO15VENDOREXPERIMENT=1


WORKDIR /go/src/gitlab.qdqmedia.com/shared-projects/riakapi
