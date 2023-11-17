FROM golang:1.18.7-alpine3.16 AS build
RUN apk add make
WORKDIR /app
COPY . .
RUN make


FROM alpine:3.16 AS prod
RUN apk update && apk upgrade
RUN apk add gcc
RUN apk add g++
RUN apk add go
RUN apk add python3
RUN apk add openjdk8
RUN apk add rust
RUN rm -vrf /var/cache/apk/*
WORKDIR /app
COPY --from=build /app/rootfs ./rootfs
COPY --from=build /app/configs ./configs
COPY --from=build /app/sandbox .
COPY --from=build /app/sb-judger .
EXPOSE 8080
EXPOSE 9090
CMD [ "/app/sb-judger" ]
