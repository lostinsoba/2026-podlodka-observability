ARG GO_VERSION
ARG ALPINE_VERSION

FROM golang:${GO_VERSION} as builder
ARG VERSION
ARG GIT_COMMIT
WORKDIR /tenant-registry
COPY sdk /sdk
COPY /tenant-registry/go.mod .
COPY /tenant-registry/go.sum .
RUN go mod download
COPY ./tenant-registry /tenant-registry
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.version=${VERSION} -X main.gitCommit=${GIT_COMMIT}" -o build/api ./cmd/api

FROM alpine:${ALPINE_VERSION}
RUN apk add --no-cache ca-certificates && update-ca-certificates
RUN adduser -h /tenant-registry -H -D -u 1000 app
USER app
COPY --from=builder /tenant-registry/build /tenant-registry
ENTRYPOINT ["./tenant-registry/api"]