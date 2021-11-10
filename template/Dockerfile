# syntax=docker/dockerfile:1

FROM golang:1.17-bullseye as build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY cmd/*.go ./

RUN go build -o /go/bin/{{.Appname}}

# Now copy it into our base image.
FROM gcr.io/distroless/base-debian11
COPY --from=build /go/bin/{{.Appname}} /

ENTRYPOINT ["/{{.Appname}}"]