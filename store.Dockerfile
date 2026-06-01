ARG GO_VERSION
ARG ALPINE_VERSION

FROM golang:${GO_VERSION} as builder
ARG VERSION
ARG GIT_COMMIT
WORKDIR /store
COPY sdk /sdk
COPY /store/go.mod .
COPY /store/go.sum .
RUN go mod download
COPY ./store /store
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.version=${VERSION} -X main.gitCommit=${GIT_COMMIT}" -o build/api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.version=${VERSION} -X main.gitCommit=${GIT_COMMIT}" -o build/receiver ./cmd/receiver

FROM alpine:${ALPINE_VERSION}
RUN apk add --no-cache ca-certificates && update-ca-certificates
RUN adduser -h /store -H -D -u 1000 app
USER app
COPY --from=builder /store/build /store