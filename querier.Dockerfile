ARG GO_VERSION
ARG ALPINE_VERSION

FROM golang:${GO_VERSION} as builder
ARG VERSION
ARG GIT_COMMIT
WORKDIR /querier
COPY sdk /sdk
COPY /querier/go.mod .
COPY /querier/go.sum .
RUN go mod download
COPY ./querier /querier
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.version=${VERSION} -X main.gitCommit=${GIT_COMMIT}" -o build/querier ./cmd/querier

FROM alpine:${ALPINE_VERSION}
RUN apk add --no-cache ca-certificates && update-ca-certificates
RUN adduser -h /querier -H -D -u 1000 app
USER app
COPY --from=builder /querier/build /querier
ENTRYPOINT ["./querier/querier"]