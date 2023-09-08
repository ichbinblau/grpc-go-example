FROM alpine:3.17

ENV HTTPS_PROXY "http://child-prc.intel.com:913"
ENV HTTP_PROXY "http://child-prc.intel.com:913"

RUN apk add --no-cache libc6-compat && mkdir -p /app/x509

WORKDIR /app

COPY ./bin .

RUN chmod +x * 

CMD ["/bin/ash", "-c", "sleep 100000000"]
