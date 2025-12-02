FROM golang:1.24.4 as build

WORKDIR /go/src/app

COPY . .

RUN go mod download && CGO_ENABLED=0 go build -o /go/bin/app ./cmd/go-templ

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /go/bin/app /

CMD ["/app"]

