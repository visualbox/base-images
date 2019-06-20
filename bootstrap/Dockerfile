FROM golang:latest as builder
WORKDIR /build
COPY . .

RUN go get github.com/sacOO7/gowebsocket
RUN CGO_ENABLED=0                              \
    GOOS=linux                                 \
    go build -a                                \
             -installsuffix cgo                \
             -ldflags '-extldflags "-static"'  \
             -o boot .

FROM scratch
COPY --from=builder /build/boot /
WORKDIR /
