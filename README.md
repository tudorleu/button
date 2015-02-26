# button

## Installation

The project is set up according to the Go standards, so commands such as `go install` and `go build` should work.

The port on which the server starts is taken from an environment variable, so after building, start with: `PORT=1234 button`.

## API

 * `POST /user` - Creates a new user, returns a User object. Parameters: `email`, `firstName`, `lastName`.
 * `GET /user/:id` - Retrieves a user with the given id.
 * `POST /user/:id/transfer` - Creates a transfer for the given user. Parameters: `amount`
 * `GET /user/:id/transfers` - Lists the transfers for the given user.

## Examples:

```
$ curl -X POST https://arcane-sands-6945.herokuapp.com/user -d email=test@test.com -d firstName=Jack -d lastName=White
{"userId":1,"email":"test@test.com","firstName":"Jack","lastName":"White","points":0}

$ curl -X GET https://arcane-sands-6945.herokuapp.com/user/1
{"userId":1,"email":"test@test.com","firstName":"Jack","lastName":"White","points":0}

$ curl -X POST https://arcane-sands-6945.herokuapp.com/user/1/transfer -d amount=10
{"id":1,"userId":1,"amount":10}

$ curl -X GET https://arcane-sands-6945.herokuapp.com/user/1
{"userId":1,"email":"test@test.com","firstName":"Jack","lastName":"White","points":10}

$ curl -X GET https://arcane-sands-6945.herokuapp.com/user/1/transfers
[{"id":1,"userId":1,"amount":10}]

$ curl -X POST https://arcane-sands-6945.herokuapp.com/user/1/transfer -d amount=-20
{"error":"Insufficient points"}

$ curl -X POST https://arcane-sands-6945.herokuapp.com/user/1/transfer -d amount=-10
{"id":2,"userId":1,"amount":-10}

$ curl -X GET https://arcane-sands-6945.herokuapp.com/user/1
{"userId":1,"email":"test@test.com","firstName":"Jack","lastName":"White","points":0}

$ curl -X GET https://arcane-sands-6945.herokuapp.com/user/1/transfers
[{"id":1,"userId":1,"amount":10},{"id":2,"userId":1,"amount":-10}]
```

