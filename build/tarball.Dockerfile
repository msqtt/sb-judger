FROM alpine:3.16
RUN apk update && apk upgrade
RUN apk add python3
RUN apk add openjdk8-jre
RUN rm -vrf /var/cache/apk/*
# fix libjli.so
RUN cp $(find / -name libjli.so | head -n 1) /lib
