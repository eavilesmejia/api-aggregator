FROM golang:1.18 as build-env

WORKDIR /go/src/app
# install dependencies first to have them cached
COPY go.mod go.mod
RUN go mod tidy
COPY . .
RUN CGO_ENABLED=0 go build -o /go/bin/app/infura

FROM gcr.io/distroless/static as run
COPY --from=build-env /go/bin/app/infura /
EXPOSE 8080
CMD ["/infura", "run", "server", "--restful"]
