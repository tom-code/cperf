FROM alpine
RUN apk update && apk add curl
ADD cperf /cperf

