FROM golang:1.18.7-alpine3.16 AS build
RUN apk add make
WORKDIR /app
COPY . .
RUN make
RUN mkdir prod-data
COPY ./rootfs prod-data/
COPY ./web prod-data/
COPY ./configs prod-data/
COPY ./sandbox prod-data/
COPY ./sb-judger prod-data/


FROM alpine:3.16 AS prod
RUN apk update && apk upgrade && \
    apk add gcc g++ go python3 openjdk8 rust && \
    # gcc: 11.2.1_git20220219-r2
    # g++: 11.2.1_git20220219-r2
    # go: 1.18.7-r0
    # python: 3.10.14-r1
    # java: 8.392.08-r0
    # rust: 1.60.0-r2
    rm -vrf /var/cache/apk/*
WORKDIR /app
COPY --from=build /app/prod-data/* .
EXPOSE 8080
EXPOSE 9090
CMD [ "/app/sb-judger" ]
