FROM golang:latest as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build

RUN go get github.com/sacOO7/gowebsocket
RUN CGO_ENABLED=0                              \
    GOOS=linux                                 \
    go build -a                                \
             -installsuffix cgo                \
             -ldflags '-extldflags "-static"'  \
             -o main .

FROM scratch
COPY --from=builder /build/main /
WORKDIR /
