
![Image of Machinable UI](images/logo.png)

Machinable gives developers the ability to store and manage their application's data in a structured, validated, RESTful way without having to write any backend code.

See the [User Documentation](https://www.machinable.io/documentation/) for more information.

#### Build

Install and build the API docker image.

```
# run glide
$ make install

# build the docker image
$ make build
```

#### Dev Environment

##### hosts

The Machinable API requires a valid hostname (with subdomain) to process requests, so you'll need to update your hosts file to include the following

```
127.0.0.1   manage.machinable.test
127.0.0.1   some-project.machinable.test
127.0.0.1   another-project.machinable.test
```

`127.0.0.1   manage.machinable.test` is required, the other lines are for any project slugs you need to test locally.

#### Testing

Run unit tests with the following command:

```
# run with make command
$ make test

# run with go test
$ go test ../... -v
```

#### CI

Github Action

### Local Development

`make install && make build && make up` will install dependencies, build images, and run the necessary containers for a local environment. This includes the posgres, redis, and api containers. The containers are defined in `docker-compose.yml`.

_NOTE: `docker-compose.yml` should only be used for your loca development environment, as it uses clear text credentials for the database_

### Web UI

_NOTE: See the UI repository for information on how to run the UI._

![Image of Machinable UI](images/ui_1.png)
