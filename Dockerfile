FROM kowalatech/go:1.0.1

# Get kman binary
RUN go get -u github.com/kowala-tech/kman

# Make sure kman is run by default
ENTRYPOINT ["/kowala/workspace/bin/kman"]
