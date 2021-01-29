FROM golang:1.15 AS build
WORKDIR /go/src

COPY pkg ./pkg
COPY go.mod .
COPY go.sum .
COPY specification.yaml .
COPY main.go .

ENV CGO_ENABLED=0
RUN go get -d -v ./...

RUN go build -a -installsuffix cgo -o paged .

FROM scratch AS runtime
ENV GIN_MODE=release
COPY --from=build /go/src/paged ./
EXPOSE 3000/tcp
ENTRYPOINT ["./paged"]
