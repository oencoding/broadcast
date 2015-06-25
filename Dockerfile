# Start with golang base image
FROM golang
MAINTAINER Omar Qazi (omar@predpol.com)
ENV HOME /root

# Compile latest source
ADD . /go/src/github.com/omarqazi/broadcast
RUN go get github.com/omarqazi/broadcast
RUN go install github.com/omarqazi/broadcast

CMD /go/bin/broadcast
EXPOSE 8080
