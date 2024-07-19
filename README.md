# FACEIT Challenge

## Task

Write a small microservice to manage access to `Users`, the service should be implemented in Go

## Requirements

A user must be stored using the following schema:

```json
{
    "id": "<id>",
    "first_name": "Alice",
    "last_name": "Bob",
    "nickname": "ab234",
    "password": "secretpwd",
    "email": "alice@bob.com",
    "country": "UK",
    "created_at": "2019-10-12T07:20:50.52Z",
    "updated_at": "2020-10-12T07:20:50.52Z"
}
```

* Provide an HTTP/gRPC API
* Use a sensible storage for the Users
* Have the ability to notify other interested services of changes to User entities
* Have meaninful logs
* Be well documented
* Have an health check

## Functionalities

* Add a new User
* Modify an existing user
* Remove a user
* Return a paginated list of Users, allowing for filtering by certain criteria (All users with the country UK)


## Considerations

Most of these considerations/choices have been taken due to the lack of time and the nature of this exercise. Focusing on make things right and writing tests requires in general more time but it gives more confidence and helps writing more robust softwares, I hope it is going to take into consideration during the final evaluation.

* **ID Format:** The schema specifies UUIDs, but MongoDB's hex format for ObjectIDs is used instead. This choice improves insert performance and simplifies update/delete operations.
* **Password Security:** Passwords are hashed with bcrypt. Future improvements could include salting with an internal constant and extracting hashing into a separate module for enhanced security.
* **Pagination and Streaming:** The current implementation uses pagination. Streaming might be considered for handling larger datasets or higher limits in the future.
* **Testing:** I tried to test the most critical part of the application. The gRPC implementation lacks comprehensive testing due to time constraints. More extensive testing should be added on that part but considering the scope of this exercise I guessed that could be omitted.
* **Project Structure:** Domain-Driven Design (DDD) principles were applied for better separation of concerns. Additional field validations could be beneficial.
* Have in place more field validations.
* **Port Configuration:** The HTTP server runs on port `80`, while the gRPC server is on port `8080`. Only the HTTP server is public exposed. Typically, gRPC servers are used for internal communication and require different security considerations.
* **Missing Endpoints:** The "getUser" endpoint is omitted as it was not specified in the requirements, allowing focus on other critical functionalities.
* **Configuration and Logging:** The application uses default configuration and logging mechanisms. Future work could improve configuration management and use advanced logging libraries like [zap](https://github.com/uber-go/zap).
* **Web Server** [Gorilla/Mux](https://github.com/gorilla/mux) has been used for the web server, it is simple to use and effective for what needed.

## Testing Strategy

Testing is one of the most important part when writing software.   
Personally I like the [Test Driven Development](https://en.wikipedia.org/wiki/Test-driven_development) approach following the Boston approach (inside-out).

* In the data layer I wrote integration tests, using [Testcontainers](https://testcontainers.com/), specifically the MongoDB module. It helps to stay as close as possible to real scenarios.

* In the business level I wrote unit tests, mocking the dependencies with [testify](https://github.com/stretchr/testify).

### How to run all tests

1. Having `docker` up and running
2. `go test ./...`

![test](https://i.imgur.com/wS1jCwH.png)

## How to run the project

1. Have `docker` up and running
2. `docker-compose up --build`

![starting](https://i.imgur.com/CCEuKfp.png)

![continue](https://i.imgur.com/LbPeRde.png)

## HTTP Create user

Through the endpoint: `/api/user` using the `POST` method.

Request:
```sh
dlion@darkness ~ % curl -X POST http://localhost:80/api/user \
     -H "Content-Type: application/json" \
     -d '{
           "first_name": "John",
           "last_name": "Doe",
           "nickname": "john.doe",
           "email": "john.doe@future.com",
           "password": "supersecurepassword",
           "country": "UK"
         }'
```
Response:
```json
{
  "id": "669a5b3525ff5682bea961ba",
  "first_name": "John",
  "last_name": "Doe",
  "nickname": "john.doe",
  "email": "john.doe@future.com",
  "country": "UK",
  "created_at": "2024-07-19T12:25:25Z",
  "updated_at": "2024-07-19T12:25:25Z"
}
```

The password is not returned for security reason.

## HTTP Modify user

Through the endpoint: `/api/user/{id}` using the `PUT` method.

Request:
```sh
curl -X PUT http://localhost:80/api/user/669a5b3525ff5682bea961ba \
 -H "Content-Type: application/json" \
 -d '{ "first_name": "Paco" }'
 ```

 Response:
 ```json
 {
  "id": "669a5b3525ff5682bea961ba",
  "first_name": "Paco",
  "last_name": "Doe",
  "nickname": "john.doe",
  "email": "john.doe@future.com",
  "country": "UK",
  "created_at": "2024-07-19T12:25:25Z",
  "updated_at": "2024-07-19T12:28:52Z"
}
```

## HTTP Delete User

Through the endpoint: `/api/user/{id}` using the `DELETE` method.

```sh
curl -X DELETE http://localhost:80/api/user/669a5b3525ff5682bea961ba
```

Response: HTTP Status 200
```text
* Host localhost:80 was resolved.
* IPv6: ::1
* IPv4: 127.0.0.1
*   Trying [::1]:80...
* Connected to localhost (::1) port 80
> DELETE /api/user/669a5dd6ed1bb58da5fe8ba7 HTTP/1.1
> Host: localhost
> User-Agent: curl/8.6.0
> Accept: */*
>
< HTTP/1.1 200 OK
< Date: Fri, 19 Jul 2024 12:36:47 GMT
< Content-Length: 0
<
* Connection #0 to host localhost left intact
```

## HTTP List Users

Through the endpoint: `/api/users` using the `GET` method.

Filters available:
* `first_name`
* `last_name`
* `nickname`
* `email`
* `country`
* `limit`
* `offset`

### Request Examples


Simple request: `curl http://localhost:80/api/users`


Response:

```json
[
  {
    "id": "669a49151d6327b831fb4797",
    "first_name": "John",
    "last_name": "D11oe",
    "nickname": "11q212q1122122dssdaohqqndoe",
    "email": "121232q222dsds3222john.a@qqexample.com",
    "country": "UK",
    "created_at": "2024-07-19T11:08:05Z",
    "updated_at": "2024-07-19T11:08:05Z"
  },
  {
    "id": "669a4800ccfd366bb0843dfd",
    "first_name": "John",
    "last_name": "D11oe",
    "nickname": "11q212q1122122dssd25jaohqqndoe",
    "email": "121232q222dsds3222john.adoe@qqexample.com",
    "country": "UK",
    "created_at": "2024-07-19T11:03:28Z",
    "updated_at": "2024-07-19T11:03:28Z"
  },
  {
    "id": "669a47b7e383b231d3f1af2c",
    "first_name": "John",
    "last_name": "Doe",
    "nickname": "11212q1122122dssd25jaohqqndoe",
    "email": "121232q2dsds3222john.adoe@qqexample.com",
    "country": "UK",
    "created_at": "2024-07-19T11:02:15Z",
    "updated_at": "2024-07-19T11:02:15Z"
  },
  {
    "id": "669a46da7ce4c0c9585ed7ee",
    "first_name": "John",
    "last_name": "Doe",
    "nickname": "11212q1122122dssd25jaohndoe",
    "email": "121232q2dsds3222john.adoe@example.com",
    "country": "UK",
    "created_at": "2024-07-19T10:58:34Z",
    "updated_at": "2024-07-19T10:58:34Z"
  },
  ...
]
```

Request with a limit: `curl "http://localhost:80/api/users?limit=2"`

Response:
```json
[
  {
    "id": "669a49151d6327b831fb4797",
    "first_name": "John",
    "last_name": "D11oe",
    "nickname": "11q212q1122122dssdaohqqndoe",
    "email": "121232q222dsds3222john.a@qqexample.com",
    "country": "UK",
    "created_at": "2024-07-19T11:08:05Z",
    "updated_at": "2024-07-19T11:08:05Z"
  },
  {
    "id": "669a4800ccfd366bb0843dfd",
    "first_name": "John",
    "last_name": "D11oe",
    "nickname": "11q212q1122122dssd25jaohqqndoe",
    "email": "121232q222dsds3222john.adoe@qqexample.com",
    "country": "UK",
    "created_at": "2024-07-19T11:03:28Z",
    "updated_at": "2024-07-19T11:03:28Z"
  }
]
````

Request with offset: `curl "http://localhost:80/api/users?offset=2"`

Response:
```json
[
  {
    "id": "669a47b7e383b231d3f1af2c",
    "first_name": "John",
    "last_name": "Doe",
    "nickname": "11212q1122122dssd25jaohqqndoe",
    "email": "121232q2dsds3222john.adoe@qqexample.com",
    "country": "UK",
    "created_at": "2024-07-19T11:02:15Z",
    "updated_at": "2024-07-19T11:02:15Z"
  },
  {
    "id": "669a46da7ce4c0c9585ed7ee",
    "first_name": "John",
    "last_name": "Doe",
    "nickname": "11212q1122122dssd25jaohndoe",
    "email": "121232q2dsds3222john.adoe@example.com",
    "country": "UK",
    "created_at": "2024-07-19T10:58:34Z",
    "updated_at": "2024-07-19T10:58:34Z"
  },
  {
    "id": "669a46b47ce4c0c9585ed7ed",
    "first_name": "John",
    "last_name": "Doe",
    "nickname": "112q1122122dssd25jaohndoe",
    "email": "121121232q2dsds3222john.adoe@example.com",
    "country": "USA",
    "created_at": "2024-07-19T10:57:56Z",
    "updated_at": "2024-07-19T10:57:56Z"
  },
  ...
]
```

Request with country filter: `curl "http://localhost:80/api/users?country=UK"`

Response:
```json
[
  {
    "id": "669a49151d6327b831fb4797",
    "first_name": "John",
    "last_name": "D11oe",
    "nickname": "11q212q1122122dssdaohqqndoe",
    "email": "121232q222dsds3222john.a@qqexample.com",
    "country": "UK",
    "created_at": "2024-07-19T11:08:05Z",
    "updated_at": "2024-07-19T11:08:05Z"
  },
  {
    "id": "669a4800ccfd366bb0843dfd",
    "first_name": "John",
    "last_name": "D11oe",
    "nickname": "11q212q1122122dssd25jaohqqndoe",
    "email": "121232q222dsds3222john.adoe@qqexample.com",
    "country": "UK",
    "created_at": "2024-07-19T11:03:28Z",
    "updated_at": "2024-07-19T11:03:28Z"
  },
  {
    "id": "669a47b7e383b231d3f1af2c",
    "first_name": "John",
    "last_name": "Doe",
    "nickname": "11212q1122122dssd25jaohqqndoe",
    "email": "121232q2dsds3222john.adoe@qqexample.com",
    "country": "UK",
    "created_at": "2024-07-19T11:02:15Z",
    "updated_at": "2024-07-19T11:02:15Z"
  },
  {
    "id": "669a46da7ce4c0c9585ed7ee",
    "first_name": "John",
    "last_name": "Doe",
    "nickname": "11212q1122122dssd25jaohndoe",
    "email": "121232q2dsds3222john.adoe@example.com",
    "country": "UK",
    "created_at": "2024-07-19T10:58:34Z",
    "updated_at": "2024-07-19T10:58:34Z"
  }
]
```

Request with country filter and limit: `curl "http://localhost:80/api/users?country=UK&limit=2"`

Response:
```json
[
  {
    "id": "669a49151d6327b831fb4797",
    "first_name": "John",
    "last_name": "D11oe",
    "nickname": "11q212q1122122dssdaohqqndoe",
    "email": "121232q222dsds3222john.a@qqexample.com",
    "country": "UK",
    "created_at": "2024-07-19T11:08:05Z",
    "updated_at": "2024-07-19T11:08:05Z"
  },
  {
    "id": "669a4800ccfd366bb0843dfd",
    "first_name": "John",
    "last_name": "D11oe",
    "nickname": "11q212q1122122dssd25jaohqqndoe",
    "email": "121232q222dsds3222john.adoe@qqexample.com",
    "country": "UK",
    "created_at": "2024-07-19T11:03:28Z",
    "updated_at": "2024-07-19T11:03:28Z"
  }
]
```

## gRPC Functions

* `GetUsers(GetUsersRequest) returns (GetUsersResponse);`
* `CreateUser(CreateUserRequest) returns (User);`
* `UpdateUser(UpdateUserRequest) returns (User);`
* `DeleteUser (DeleteUserRequest) returns (Empty);`
* `Watch(google.protobuf.Empty) returns (stream WatchResponse);`