FROM alpine:3.16
RUN apk update && \
    apk upgrade && \
    apk add python3 openjdk8-jre && \
    # python: 3.10.14-r1
    # java: 8.392.08-r0
    rm -vrf /var/cache/apk/* && \
    rm -vrf /etc/ssl/* && \
    rm -vrf /etc/terminfo/* && \
    rm -vrf /usr/share/ca-certificates/* && \
    rm -vrf /usr/share/alsa/* && \
    rm -vrf /usr/share/apk/* && \
    rm -vrf /usr/share/udhcpc/* && \
    rm -vrf /lib/apk/* && \
    # fix java libjli.so
    cp $(find / -name libjli.so | head -n 1) /lib
