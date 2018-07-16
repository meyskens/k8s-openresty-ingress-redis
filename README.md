Kubernetes OpenResty Ingress Controller with redis
=======================================

## Why OpenResty?
I work in an high traffic envoirement where we have tested different reverse proxies. Our tests showed NGINX as the best one on resource usage. But NGINX didn't fit our needs ins customisability, that is why we chose OpenResty as a solution. While this setup is quite minimal we have more stuff going on in our production configuration. 

## Why Redis?
This is a fork of my [k8s-openresty-ingress](https://github.com/meyskens/k8s-openresty-ingress) repo to use Redis instead. I noticed reloading on config changes caused very high load on Nginx when running for some time. That's why this plugin will query redis for which host to proxy to to be able to re-route without loading in new configuration.


## Thank you to
- [traefik](https://github.com/containous/traefik/) For having understandable code on the Kubernetes backend