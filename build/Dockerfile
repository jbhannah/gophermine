FROM golang:1.13-alpine as builder

WORKDIR /go/src/github.com/jbhannah/gophermine
ENV CGO_ENABLED=0
ENV GO111MODULE=on

RUN apk add bash build-base curl git

COPY go.mod go.sum ./
RUN go mod download

COPY cmd cmd/
COPY internal internal/
COPY pkg pkg/
RUN go test ./... \
    && go install ./cmd/mc

FROM gcr.io/distroless/base-debian10
WORKDIR /opt/gophermine
COPY --from=builder /go/bin/mc ./bin/mc
ENTRYPOINT [ "./bin/mc" ]
