# FACEIT Challenge

## Task

Write a small microservice to manage access to `Users`, the service should be implemented in Go

## Requirements

A user must be stored using the following schema:

```json
{
    "id": "uuid",
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

* Since I decided to use MongoDB as a db, I'm using `ObjectIDs` as id. It because has been shown that using UUIDs cause performance drop for inserts.