
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

# run container
$ make up
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

##### configuration

The application config has the following structure:

```config.json
{
    "Version": "0.0.0",
    "AppSecret": "",
    "ReCaptchaSecret": "",
    "IPStackKey": ""
}
```

|Key|Description|Required|
|**Version**|The version of the API|`False`|
|**AppSecret**|The secret string used to salt passwords|`True`|
|**ReCaptchaSecret**|The Google reCaptcha secret used for user registration|`True`|
|**IPStackKey**|The API Key for IP Stack|`False`|

The secret config values can also be provided as environment variables in `docker-compose.yml`:

```yml
    - APP_SECRET
    - RECAPTCHA_SECRET
    - IPSTACK_KEY
```

#### Testing

Run unit tests with the following command:

```
# run with make command
$ make test

# run with go test
$ go test ../... -v
```

#### CI

Github Tag Action - [https://github.com/anothrNick/github-tag-action](https://github.com/anothrNick/github-tag-action)

See [./github/workflows/main.yml](./github/workflow/main.yml) for the full Github workflow.

#### Packages

Docker Image: [machinable](https://github.com/anothrNick/machinable/packages/54301)

### Local Development

`make install build up` will install dependencies, build images, and run the necessary containers for a local environment. This includes the posgres, redis, and api containers. The containers are defined in `docker-compose.yml`.

_NOTE: `docker-compose.yml` should only be used for your local development environment, as it uses clear text credentials for the database_

### Web UI

_NOTE: See the [UI repository](https://github.com/anothrNick/machinable-ui) for information on how to run the UI._

![Image of Machinable UI](images/ui_1.png)
