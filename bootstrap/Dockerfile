FROM golang:latest as builder

RUN mkdir /build
ADD /src /build/
WORKDIR /build

RUN go get github.com/mtfelian/golang-socketio
RUN CGO_ENABLED=0                              \
    GOOS=linux                                 \
    go build -a                                \
             -installsuffix cgo                \
             -ldflags '-extldflags "-static"'  \
             -o main .

FROM scratch

COPY --from=builder /build/main /app/
WORKDIR /app
