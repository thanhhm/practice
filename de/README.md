# Simple structure convention for able to write unit test
Go still have no standard project structure, each company, each team might have their own style but whatever the structure is, every project should have three layer proposed below

```
handlers
└── payment.go
services
└── payment.go
repositories
|── account.go 
|── transaction.go
└── user.go 
```

#### 1. repositories
* This layer will responsible for call to get/set resource from database, cache or 3rd api. Each file will be a model mapping to the resource `e.g: user.go mapping to table user in database`.
* Handle error: do not log error here, error handling belong to the upper layer `services`. Use `errors.Wrap(err, "unmarshal error here")` for additional info.
#### 2. services
* This layer will do real business logic here, `services` will call to `repositories` to get/set data and mix them match to solve business logic
* Error can be logged here if the function code is huge and complicate, otherwise just return for `handlers` to log.
#### 3. handlers
* This layer is receive request from client, validate input data. If everything is okay `handlers` will call to `services` to do the api business logic.

### Dependency
```
handlers
   ^
   |
services
   ^
   |
repository
```
As of above description, `upper` layer will depend on `lower` layer: handlers depend on services; services depend on repositories, no reverse direction.

### Step coding API
The step will be reverse direction of layer dependency: repositories -> services -> handlers.
### Rule of thumb
**Single responsibility** function should do one thing and do it well.