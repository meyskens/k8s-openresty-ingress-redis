Kubernetes OpenResty Ingress Controller with redis
=======================================

This is a Kubernetes Ingress controller which uses OpenResty with a Redis configuration backend.
Routes are stored inside Redis, and TLS certificates on disk. On each request OpenResty will reach out to Redis to get
the proxy configuration. This allows for quick routing updates without the memory overhead of reloading NGINX.

## Why OpenResty?
I work in a high traffic (streaming) environment where we have tested different reverse proxies. 
Our tests showed NGINX as the best one on resource usage. 
But NGINX didn't fit our needs ins customizability, that is why we chose OpenResty as a solution. 
While this setup is quite minimal we have more stuff going on in our production configuration. 

## Why Redis?
This is a fork of my [k8s-openresty-ingress](https://github.com/meyskens/k8s-openresty-ingress) repo to use Redis instead. 
We noticed reloading on config changes caused very high load on NGINX when running for some time as it kept the old configuration in memory
till all clients were disconnected, which in streaming is almost never the case. 
That's why this plugin will query redis for which host to proxy to be able to re-route without loading in new configuration.

## Using this ingress controller
The code is written very specifically for our usecase, I reccomend only to use this if you have experience with the stack we use.
The code also isn't fully tested against all cluster edge cases. Buyer be aware.

## Thank you to
- [traefik](https://github.com/containous/traefik/) For having understandable code on the Kubernetes backend