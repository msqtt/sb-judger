FROM alpine:3.16
RUN apk update && apk upgrade
RUN apk add python3
RUN apk add openjdk8
RUN rm -vrf /var/cache/apk/*
