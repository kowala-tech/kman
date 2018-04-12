FROM kowalatech/go:1.0.1 as builder

WORKDIR /kowala/workspace/src/github.com/kowala-tech/kman
ADD . .

# Build kman binary
RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure
RUN make build

# Use a local version
FROM kowalatech/hugo-dev:1.0.4
COPY --from=builder /kowala/workspace/bin/kman /kowala/workspace/bin/kman
EXPOSE 8080
ADD themes themes

# Make sure kman is run by default
ENTRYPOINT ["/kowala/workspace/bin/kman"]
