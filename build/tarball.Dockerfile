FROM alpine:3.16
RUN apk update && apk upgrade
RUN apk add python3
RUN apk add openjdk8-jre
RUN rm -vrf /var/cache/apk/*
RUN rm -vrf /etc/ssl/*
RUN rm -vrf /etc/terminfo/*
RUN rm -vrf /usr/share/ca-certificates/*
RUN rm -vrf /usr/share/alsa/*
RUN rm -vrf /usr/share/apk/*
RUN rm -vrf /usr/share/udhcpc/*
RUN rm -vrf /lib/apk/*
# fix java libjli.so
RUN cp $(find / -name libjli.so | head -n 1) /lib
