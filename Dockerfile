FROM golang:1.11 AS build
COPY go.mod /usr/src/testrig/go.mod
COPY go.sum /usr/src/testrig/go.sum
WORKDIR /usr/src/testrig
RUN go mod download
COPY . /usr/src/testrig
RUN make build

FROM scratch
COPY --from=build /usr/src/testrig/bin/testrig /testrig
ENTRYPOINT ["/testrig"]
