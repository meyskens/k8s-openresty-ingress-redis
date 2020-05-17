ARG ARCH
# Build go binary
FROM golang:1.14 AS gobuild

COPY ./ /go/src/github.com/meyskens/k8s-openresty-ingress-redis
WORKDIR /go/src/github.com/meyskens/k8s-openresty-ingress-redis/controller

ARG GOARCH
ARG GOARM

RUN GOARCH=${GOARCH} GOARM=${GOARM} go build ./

# Set up deinitive image
ARG ARCH
FROM maartje/openresty:1.15.8.3

# Add Dummy cert for dummy conf
RUN openssl req -new -newkey rsa:2048 -days 3650 -nodes -x509 \
       -subj '/CN=sni-support-required-for-valid-ssl' \
       -keyout /etc/ssl/resty-auto-ssl-fallback.key \
       -out /etc/ssl/resty-auto-ssl-fallback.crt

COPY --from=gobuild /go/src/github.com/meyskens/k8s-openresty-ingress-redis/controller/controller /usr/local/bin/controller

COPY ./config/default/ /etc/nginx/

EXPOSE 80
EXPOSE 443
CMD controller