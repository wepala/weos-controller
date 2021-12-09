# LOGGING

The previous system made use of logrus for logging however it has now transitioned into maintainance mode. Therefore we have switched to Zap logger, which was setup to implement the same interface as before. Custom log levels are also now supported and can be set on a per request basis.

### How do i get set up?

Use the `Loglevel` middleware in `api.yaml` file.

**NB: Must include `ZapLogger` middleware in `api.yaml` file since the LogLevel is set in the context/per request and `ZapLogger` is responsible for setting the logger itself on the context.**

Example configurations:

#### 1.Added globally to the middleware

```yaml
x-weos-config:
  session:
    key: "${SESSION_KEY}"
    path: ""
  applicationId: ${APPLICATION_ID}
  applicationTitle: ${APPLICATION_TITLE}
  accountId: ${ACCOUNT_ID}
  database:
    host: ${POSTGRES_HOST}
    database: ${POSTGRES_DB}
    username: ${POSTGRES_USER}
    password: ${POSTGRES_PASSWORD}
    port: ${POSTGRES_PORT}
  pre-middleware:
    - PreGlobalMiddleware
  middleware:
    - LogLevel
    - ZapLogger
    - GlobalMiddleware
  jwtConfig:
    key: ${JWT_KEY}
    tokenLookup: ${JWT_LOOKUP}
    claims:
      email: "email@mail"
      real: false
```

#### 2.Added to the middleware on a specfic path

```yaml
/endpoint:
  summary: Some user endpoint
  get:
    x-weos-config:
      handler: FooBar
      middleware:
        - LogLevel
        - ZapLogger
        - Middleware
        - PreMiddleware
    responses:
      200:
        description: Admin Endpoint
```

### Log Levels

Log levels are now passed in using the `X-LOG-LEVEL` header in requests. (The constant `weoscontroller.HeaderXLogLevel` can also be used when testing.)

Available headers are:

1. `debug`
2. `info`
3. `warn`
4. `error`

**NB: The above list represents the hierarchy of levels, therefore whatever the current level is set at, the levels below it will also be outputted.
`I.E.` If the level is set at `info`, you will also see `warn` and `error` logs.**

If no header is passed in, the default log level will be set to `error`

### Writing Log statements

Ensure that you are calling the echo logger and `not` the standard golang logger when using log outputs.

This follows this format `e.Logger().x`, `x` representing the type of log you require `(Debug, Info, Warn, Error)`.

**NB: `e` represents the `echo.context`**

### Zap Logger

Use the `ZapLogger` middleware in `api.yaml` file.

Example configurations:

#### 1.Added globally to the middleware

```yaml
x-weos-config:
  session:
    key: "${SESSION_KEY}"
    path: ""
  applicationId: ${APPLICATION_ID}
  applicationTitle: ${APPLICATION_TITLE}
  accountId: ${ACCOUNT_ID}
  database:
    host: ${POSTGRES_HOST}
    database: ${POSTGRES_DB}
    username: ${POSTGRES_USER}
    password: ${POSTGRES_PASSWORD}
    port: ${POSTGRES_PORT}
  pre-middleware:
    - PreGlobalMiddleware
  middleware:
    - ZapLogger
    - GlobalMiddleware
  jwtConfig:
    key: ${JWT_KEY}
    tokenLookup: ${JWT_LOOKUP}
    claims:
      email: "email@mail"
      real: false
```

**NB: This middleware is used to set the logger on the context. If the middleware is not specified, the logger on the echo framework `only` will be set.**
