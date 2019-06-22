# Scaled service on Docker via SSL 

This repository is an example of using a Go service behind
an NGINX webserver that is accessible via HTTPS/SSL on
<https://localhost:443>. It also illustrates how to scale
a service in local development with Docker Compose.

## Getting started

First, you need to create an SSL certificate as described
on [Let's Encrypt](https://letsencrypt.org/docs/certificates-for-localhost/)
for `localhost:443`.

```sh
openssl req -x509 -out etc/nginx/localhost.crt -keyout etc/nginx/localhost.key \
  -newkey rsa:2048 -nodes -sha256 \
  -subj '/CN=localhost' -extensions EXT -config <( \
   printf "[dn]\nCN=localhost\n[req]\ndistinguished_name = dn\n[EXT]\nsubjectAltName=DNS:localhost\nkeyUsage=digitalSignature\nextendedKeyUsage=serverAuth")
```

Next, follow these steps:

1. `make build` will create the Go binary `./service`.
2. `docker-compose up --build -d` will download, build, and
   start the application.
3. Open <https://localhost:443> and accept the warning the
   browser will give you. You should read the current time
   every time the browser refreshes the page.

By default, only one instance of the service is running.
So if you `curl https://localhost:443` multiple times, it
will also print the same hostname.

But you can scale multiple instances of the service by
invoking e.g.:

```sh
docker-compose up --scale service=3
```

That will create 3 identical instances of the service
which NGINX will proxy to all of them.

When you now `curl https://localhost:443` multiple times, and
it should print different hostnames. Notice that it may take
up to 30 seconds (as defined in the NGINX configuration file)
for the DNS resolver to purge its cache.

This works by using the default DNS resolver that Docker creates
by default on `127.0.0.11`.

## License

MIT.
