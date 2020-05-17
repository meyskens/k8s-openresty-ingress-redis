ARG ARCH
# Build go binary
FROM golang:1.14 AS gobuild

COPY ./ /go/src/github.com/meyskens/k8s-openresty-ingress-redis
WORKDIR /go/src/github.com/meyskens/k8s-openresty-ingress-redis/controller

ARG GOARCH
RUN GOARCH=${GOARCH} GOARM=7 go build ./

# Set up deinitive image
ARG ARCH
FROM maartje/openresty:$ARCH-1.17.8.1rc1

# Add Dummy cert for dummy conf
RUN openssl req -new -newkey rsa:2048 -days 3650 -nodes -x509 \
       -subj '/CN=sni-support-required-for-valid-ssl' \
       -keyout /etc/ssl/private/resty-auto-ssl-fallback.key \
       -out /etc/ssl/private/resty-auto-ssl-fallback.crt

COPY --from=gobuild /go/src/github.com/meyskens/k8s-openresty-ingress-redis/controller/controller /usr/local/bin/controller

COPY ./config/default/ /etc/nginx/

RUN mkdir -p /etc/ssl/private/

EXPOSE 80
EXPOSE 443
CMD controller