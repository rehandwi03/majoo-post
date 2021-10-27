## Majoo POS (Point Of Sales) <a name = "about"></a>

## Command <a name = "getting_started"></a>

### Application Lifecycle

```
$ cp .env.example .env
$ go mod download
$ go run main.go
 ┌───────────────────────────────────────────────────┐ 
 │                   Fiber v2.20.2                   │ 
 │               http://127.0.0.1:8080               │ 
 │       (bound on host 0.0.0.0 and port 8080)       │ 
 │                                                   │ 
 │ Handlers ............ 59  Processes ........... 1 │ 
 │ Prefork ....... Disabled  PID ............. 17085 │ 
 └───────────────────────────────────────────────────┘ 
```

### Docker Lifecycle

```
docker-compose up -d
```

## Endpoint <a name = "tests"></a>

| Name          | Endpoint         | Method        | With Token   | Description   |
| ------------- | -------------    | ------------- |------------- |------------- |
| Register      | */api/register*  |   *POST*      |    No        |For registering user
| Auth          | */api/login*     |   *POST*      |    No        |For login user
| User          | */api/users/:id*  |   *GET*       |    Yes       |Get detail of user
|               | */api/users*      |   *PUT*       |    Yes       |Update user
|               | */api/users/:id*  |   *DELETE*    |    Yes       |Delete user
|               | */api/users*      |   *GET*       |    Yes       |Get all user
| Merchant      | */api/merchants*  |   *POST*      |    Yes       |Create merchant
|               | */api/merchants/:id* |   *GET*    |    Yes       |Get merchant detail
|               | */api/merchants* |   *PUT*        |    Yes       |Update merchant
|               | */api/merchants/:id* |   *DELETE* |    Yes       |Delete merchant detail
|               | */api/merchants* |   *GET*        |    Yes       |Get all merchant
| Outlet        | */api/outlets*  |   *POST*      |    Yes       |Create outlet
|               | */api/outlets/:id*  |   *GET*      |    Yes       |Get outlet detail
|               | */api/outlets*  |   *PUT*      |    Yes       |Update outlet
|               | */api/outlets*  |   *GET*      |    Yes       |Get all outlet
|               | */api/outlets/:id*  |   *DELETE*      |    Yes       |Delete outlet
| Product       | */api/products*  |   *POST*      |    Yes       |Create product
|               | */api/products/:id*  |   *GET*      |    Yes       |Get product detail
|               | */api/products*  |   *PUT*      |    Yes       |Update product
|               | */api/products*  |   *GET*      |    Yes       |Get all product
|               | */api/products/:id*  |   *DELETE*      |    Yes       |Delete product
|               | */api/products/image*  |   *POST*      |    Yes       |Upload image product