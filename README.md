# README #

The Controller is meant to handle incoming requests and route to the appropriate business login

### What is this repository for? ###

* Version: 0.1.0

### How do I get set up? ###

You can run the serve command to give access via http. There are a few ways to start the server

#### Http Serve
1. Use cli parameters `weos-controller serve http localhost:8080 -a site-api.yml -c site-config.yml`
1. Use environment variables set environment variable `API_YAML` and `CONFIG_YAML` and then start the server `weos-controller serve http-mock`
1. Configure parameters in a config file `weos-controller serve http localhost:8080 -c weoscontroller.yml`
1. Place a config in the home folder of the service `weos-controller serve http localhost:8080`

#### Mock Responses
Getting mock responses couldn't be easier. 

##### Mock Setup
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
        '200':
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
        plugins:
          - &weosPlugin
            filename: testdata/plugins/test.so
            config:
              mysql:
                host: localhost
                user: root
                password: root
        middleware:
          - plugin: *weosPlugin
            handler: HelloWorld
      responses:
        '200':
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

### Contribution guidelines ###

This project uses [gitflow workflow](https://www.atlassian.com/git/tutorials/comparing-workflows/gitflow-workflow)

* Clone the repo to local
* Create feature branch from dev branch (e.g. feature/WEOS-1)
* Push the feature branch to the remote
* Create PR from  the ticket branch to dev branch 
* When the item is merged to master it will be deployed

To aid with this use the git flow cli (you will be able to create feature branches e.g. git flow feature start APO-1)

#### New Features ####





### Who do I talk to? ###

* Admin - Akeem Philbert <akeem.philbert@wepala.com>