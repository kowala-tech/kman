FROM kowalatech/go:1.0.1

WORKDIR /kowala/workspace/src/github.com/kowala-tech/kman
ADD . .

# Build kman binary
RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure
RUN go install github.com/kowala-tech/kman/cmd/kman

# Make sure kman is run by default
ENTRYPOINT ["/kowala/workspace/bin/kman"]
