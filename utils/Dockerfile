FROM pwntus/bootstrap:latest AS bootstrap
FROM pwntus/node:latest AS node

FROM scratch
COPY --from=bootstrap /boot /
COPY --from=node /visualbox /node/visualbox
WORKDIR /
