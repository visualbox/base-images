FROM node:12-alpine as modules
WORKDIR /build
COPY . .
RUN yarn

FROM scratch
COPY --from=modules /build /visualbox
WORKDIR /
