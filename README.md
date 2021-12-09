# README

The Controller is meant to handle incoming requests and route to the appropriate business login

### What is this repository for?

- Version: 0.1.0

### How do I get set up?

This module should be imported into an api and initialized.

```shell
go get github.com/wepala/weos-controller
```

Then setup a configuration file `api.yaml` for example

```yaml
openapi: 3.0.2
info:
  title: WeOS REST API
  version: 1.0.0
  description: REST API for passing information into WeOS

x-weos-config:
  session:
    key: "${SESSION_KEY}"
    path: ""
  logger:
    level: ${LOG_LEVEL}
    report-caller: true
    formatter: ${LOG_FORMAT}
  applicationId: ${APPLICATION_ID}
  applicationTitle: ${APPLICATION_TITLE}
  accountId: ${ACCOUNT_ID}
  database:
    host: ${POSTGRES_HOST}
    database: ${POSTGRES_DB}
    username: ${POSTGRES_USER}
    password: ${POSTGRES_PASSWORD}
    port: ${POSTGRES_PORT}
  middleware:
    - RequestID
    - Recover
    - Static
  jwtConfig:
    key: ${JWT_KEY}
    tokenLookup: ${JWT_LOOKUP}
    claims:
      email: string
      real: bool
    authScheme: ${JWT_SCHEME}
    contextKey: ${JWT_CONTEXT}
    signingMethod: ${JWT_SIGNING_METHOD}

paths:
  /health:
    summary: Health Check
    get:
      x-weos-config:
        handler: HealthCheck
      responses:
        "200":
          description: Health Response
        "500":
          description: API Internal Error
  /admin:
    summary: Admin Endpoint
    get:
      x-weos-config:
        group: True
        middleware:
          - Static
      responses:
        200:
          description: Admin Endpoint
```

NB: The content of the yaml can also be passed into the `Initialize` method

##### Config Explained

The configuration uses the [OAS 3.0 api specification](https://swagger.io/specification/). Using a standard format allows
us to take advantage of the rich tooling that is available. WeOS specific configuration can be found under the weos
specific configuration extension `x-weos-config`. When this extension is applied to a path, the configuration details
is available as a path config.

#### Mock Responses

Getting mock responses couldn't be easier.

##### Mock Setup

**THIS IS CURRENTLY NOT SETUP FOR THE NEW DIRECTION**

In your api yaml there are a couple ways you can setup mock responses;

1. Don't configure the path with a `x-weos-config` (it will automatically return example responses you have defined on the path or on the component schema of the response)
1. Setup a sub property `mock` to true on the `x-weos-config`

Example Mock Config

```yaml
openapi: "3.0.0"
info:
  version: 1.0.0
  title: Basic Site
paths:
  /:
    get:
      summary: Landing page of the site
      responses:
        "200":
          description: Landing page
          content:
            text/html:
              example: test
              schema:
                type: string
  /about:
    get:
      summary: About Page
      x-weos-config:
        mock: true
        handler: HelloWorld
      responses:
        "200":
          description: About Page
          content:
            text/html:
              example: test
              schema:
                type: string
```

You can define mocks in a few ways

1. `example` under content in the response
1. `examples` these are named examples
1. `example` on the component

Read more about examples in swagger - https://swagger.io/docs/specification/adding-examples/

##### Mock Request

There are a few options you can use when making a mock request

1. `X-Mock-Status-Code` - Use this to specify the response code you'd like to receive (this is if the path has multiple responses defined)
1. `X-Mock-Example` - Use this to specify which example should be used (this is when the `examples` option is used on the response)
1. `X-Mock-Example-Length` - Use this to specify the amount of items to return (this is useful when the response is an array)
1. `X-Mock-Content-Type` - Use this to specify the content type desired. (This is required if there are multiple content types defined)

### Contribution guidelines

This project uses [gitflow workflow](https://www.atlassian.com/git/tutorials/comparing-workflows/gitflow-workflow)

- Clone the repo to local
- Create feature branch from dev branch (e.g. feature/WEOS-1)
- Push the feature branch to the remote
- Create PR from the ticket branch to dev branch
- When the item is merged to master it will be deployed

To aid with this use the git flow cli (you will be able to create feature branches e.g. git flow feature start APO-1)

#### New Features

### Who do I talk to?

- Admin - Akeem Philbert <akeem.philbert@wepala.com>
