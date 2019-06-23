FROM visualboxio/bootstrap:latest as bootstrap
FROM visualboxio/node:latest as node
FROM visualboxio/python3:latest as python3

FROM scratch
COPY --from=bootstrap /boot /
COPY --from=node /visualbox /node/visualbox
COPY --from=python3 /visualbox /python3/visualbox
WORKDIR /
