ARG GO_VERSION
ARG ALPINE_VERSION

FROM golang:${GO_VERSION} as builder
ARG VERSION
ARG GIT_COMMIT
WORKDIR /sender
COPY sdk /sdk
COPY /sender/go.mod .
COPY /sender/go.sum .
RUN go mod download
COPY ./sender /sender
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.version=${VERSION} -X main.gitCommit=${GIT_COMMIT}" -o build/sender ./cmd/sender

FROM alpine:${ALPINE_VERSION}
RUN apk add --no-cache ca-certificates && update-ca-certificates
RUN adduser -h /sender -H -D -u 1000 app
USER app
COPY --from=builder /sender/build /sender
ENTRYPOINT ["./sender/sender"]