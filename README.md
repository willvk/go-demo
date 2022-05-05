# go-demo

_Goal/MVP: Build a Service in Golang to set meetup times and to track attendance with basic auth._

Extensions: 
* Deployable in AWS
* Register User flow
* SMS/Email Notification/Reminder Service
* Some cURL samples for usage

Lets build a meetup planner and attendance management service
1. Firstly we will need a crud service to manage meetups, so meetups are the initial data object
2. Users should be able to create a meeting at a date time with a list of attendee usernames
3. users will be basic auth (ext), mvp can be hardcoded

## Architecture

TODO: Pretty Picture

## Development

Run unit tests and compile code locally with.

```bash
make meetup-service
```

### Directories

<i> ** to be refactored...template for the moment ** </i>

- `$/cmd/service` - service entry point
- `$/openapi` - house OpenAPI spec document(s) for the service
- `$/mocks` - top-level mocks package for generated mocks
- `$/sam` - contain AWS infra Serverless Application Model (SAM) template
- `$/internal/api` - main API entry point, and contains API route handlers, the code that handles API requests
- `$/internal/app` - stores build-time information and information about this app
- `$/internal/flags` - initialise application config and command line arguments
- `$/internal/openapi` - (optional) contains generated code from OpenAPI spec; mainly there to not cause circular package dependency as `persistence` package also reference genrated code in this example
- `$/internal/persistence` - (optional) contains code that persists models to a storage location
- `$/internal/requestvalidation` - contains custom validation config that skips validation of security schemes defined in OpenAPI spec as it's not compatible with AWS SigV4 scheme
- `$/internal/router` - creates an instance of Echo router and register all required middlewares
- `$/internal/ssmcache` - (optional) example usage of SSM client with caching

### CI/CD

To be added

## Mocks

The service uses [gomock to generate mocks at build time.](https://github.com/golang/mock) for auto generated mocks

Generating mocks for tests, or applying updates to them is as follows:

```bash
make generate
```
