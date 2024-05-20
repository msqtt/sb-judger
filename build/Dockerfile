FROM golang:1.18.7-alpine3.16 AS build
RUN apk add make
WORKDIR /app
COPY . .
RUN make


FROM alpine:3.16 AS prod
RUN apk update && apk upgrade
RUN apk add gcc             # 11.2.1_git20220219-r2
RUN apk add g++             # 11.2.1_git20220219-r2
RUN apk add go              # 1.18.7-r0
RUN apk add python3         # 3.10.14-r1
RUN apk add openjdk8        # 8.392.08-r0
RUN apk add rust            # 1.60.0-r2
RUN rm -vrf /var/cache/apk/*
WORKDIR /app
COPY --from=build /app/rootfs ./rootfs
COPY --from=build /app/web ./web
COPY --from=build /app/configs ./configs
COPY --from=build /app/sandbox .
COPY --from=build /app/sb-judger .
EXPOSE 8080
EXPOSE 9090
CMD [ "/app/sb-judger" ]
