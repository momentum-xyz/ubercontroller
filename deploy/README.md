# Deployment

There are three parts that need to be deployed:

- [PostgreSQL database](https://www.postgresql.org/)
- Controller
- Frontend web application

The PostgreSQL is commonly used, open source, relational database.

The Controller is the backend application of the platform that manages the Odyssey universe state with the worlds and provides a websocket and JSON API.
It is a single, statically compiled binary. Each backend controller can serve one or more specific worlds (managed on the blockchain) and requires file storage space for user uploaded content.
There should only be a single one running, it does not scale horizontally (it manages world state in-memory). To scale add more CPUs and memory or 'migrate' worlds to another (separately running) controller.

The frontend web application is a compiled HTML + javascript application (SPA, Single Page Application) and can be statically hosted anywhere as long as the backend controller is configured to allow its domain through [CORS](https://enable-cors.org/).
It could be deployed on a CDN for performance/scalability. Optionally the controller can also host the web application files itself if there is no other static file hosting option available.


## Hosting

Recommended hosting method to run this platform is using a reverse proxy/loadbalancer in front of it (e.g. nginx) which terminates SSL, proxies websocket and API requests to the controller and serves the static frontend web application. PostgreSQL is often available as managed service but can also run it yourself (see PostgreSQL own documentation for that).

### Minimal requirements

- PostgreSQL database as a externally managed service or running it directly.
- Need to be able to run (docker) container images or a compiled binary for the backend.
- Access to read&write storage space.
- Backend needs outgoing network connections to 3rd party services (e.g. blockchain RPC services)
- Serve static HTML, JS and other asset files.
- Resources: depends on the amount of worlds hosted and the usage (concurrent users) of these worlds. But minimum, with only a couple of worlds and a handful of users, starts at a around 2 CPU, 256MB for the backend part (so excluding the database).


## Docker

Running the platform with plain [docker](https://www.docker.com/).

Using the docker compose method described below is easiest and serves as an example showing the separate parts and their (minimal) configuration.

The container images to use:

- Official [PostgreSQL image](https://hub.docker.com/_/postgres)
- Controller image on [GitHub Container registry](https://github.com/momentum-xyz/ubercontroller/pkgs/container/ubercontroller)
- Frontend image on [Github Container registry](https://github.com/momentum-xyz/ui-client/pkgs/container/ui-client)


### Docker compose

Running the platform with docker [compose](https://docs.docker.com/compose/).

Example docker compose configuration is inside the `deploy/compose` directory.
Executing `docker compose up` inside that directory will start up the platform and make it accessible on http://localhost:8080

```console
cd deploy/compose
docker compose up
```

This will run the platform in 'development' mode, connected to a private blockchain testnet.
To run against the Goerli ethereum testnet use the `compose.testnet.yml` override file.

```console
cd deploy/compose
docker compose -f compose.yml -f compose.goerli.yml
```

Or use `compose.nova.yml` to connect to the real Arbitrum Nova blockchain on Ethereum.


Remove/cleanup the docker compose environment:

```
docker compose down -v --rmi local
```


## Kubernetes

The `deploy/k8s` contains a [Kustomize](https://kubectl.docs.kubernetes.io/) configuration to deploy the platform to a cluster.

```console
kubectl kustomize deploy/k8s/overlays/dev/ | kubectl apply -f -
```

The example is for a development environment that uses a standalone postgresql database.
The `FRONTEND_URL` should match how you expose the cluster ingress.

For production environment it is advised to use a [postgres operator](https://operatorhub.io/?keyword=postgres) to manage your database and configure SSL certificates for https access.


## Running from source

Checkout this git repository (https://github.com/momentum-xyz/ubercontroller) and https://github.com/momentum-xyz/ui-client

- Start a PostgreSQL database, see its [docs](https://www.postgresql.org/docs/current/admin.html)
- Build ui-client, see its README.md (`yarn install; yarn build`)
- Build controller, run `make build`.
- Configure a config.yaml (see config.example.yaml)
- Run it: `FRONTEND_SERVE_DIR=../ui-client/packages/app/build/ ./bin/ubercontroller`
- Open http://localhost:4000 in web browser.
