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

* The id in the provided schema is clearly an uuid but since I decided to use MongoDB as a db, I'm using the hex format of the MongoDB's `ObjectIDs`. Because has been shown that using UUIDs cause performance drop for inserts in MongoDB and in my opinion it's easier for an user providing it during the update/delete operation.
* The password is hashed using the bcrypt function and saving it into the db. As future improvement it could be salted as well with some internal constant. For this exercise I let the `user-repository` taking care of it, but as a future improvement it could be extracted into a separate module, injecting it properly.
* Since we are paginating I'm thinking that streaming back to the user the content doesn't have much sense in this scenario, for larger dataset or bigger limits, it can have sense implement it.
* In the gRPC there is an evident lack of testing, due of the lack of time available to complete this exercise.
* I tried to use DDD to have a better separation in the project structure.
* Have in place more field validations.
* I exposed the http server, it is okay to be public available, I didn't expose the gRPC server, normally gRPC are made for internal communications and it should be treated differently.
* I didn't add a "getUser" endpoint, it wasn't specified in the requirements, so I could focus on more important functionalities.

## Testing Strategy

Testing is one of the most important part when writing software.   
Personally I like the Test Driven Development approach.

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

## gRPC functions:

* `GetUsers(GetUsersRequest) returns (GetUsersResponse);`
* `CreateUser(CreateUserRequest) returns (User);`
* `UpdateUser(UpdateUserRequest) returns (User);`
* `DeleteUser (DeleteUserRequest) returns (Empty);`
* `Watch(google.protobuf.Empty) returns (stream WatchResponse);`