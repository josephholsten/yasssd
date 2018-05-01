# Description

This is an API specification for a simple HTTP-based multi-user file storage
server. Your assignment is to write a server that implements this API using a
language of your choice.

There are five API endpoints to implement:

## `POST /register`

This endpoint is used to register as a new user. Usernames must be at least 3
characters and no more than 20, and may only contain alphanumeric characters.
Passwords must be at least 8 characters.

### Request
```
Content-Type: application/json

{
  "username": "the username of the user to create",
  "password": "the password of the user to create"
}
```

### Response (success)
```
Status: 204 No Content
```

### Response (failure)
```
Status: 400 Bad Request
Content-Type: application/json

{
  "error": "explanation of failure"
}
```

## `POST /login`

This endpoint is used to log in as an existing user. On success, it returns a
session token. The session token should be included in future requests to
authenticate the sender.

### Request
```
Content-Type: application/json

{
  "username": "the username of the user to log in as",
  "password": "the password of the user to log in as"
}
```

### Response (success)
```
Status: 200 OK
Content-Type: application/json

{
  "token": "an opaque session token"
}
```

### Response (failure)
```
Status: 403 Forbidden
Content-Type: application/json

{
  "error": "explanation of failure"
}
```

## `PUT /files/<filename>`

This endpoint is used to upload a file to the logged-in user's personal
storage.

### Request
```
Content-Length: <file size>
Content-Type: <content type>
X-Session: <session token>

<file bytes>
```

### Response (success)
```
Status: 201 Created
Location: /files/<filename>
```

## `GET /files/<filename>`

This endpoint is used to get a file from the logged-in user's personal storage.

### Request
```
X-Session: <session token>
```

### Response (success)
```
Status: 200 OK
Content-Length: <file size>
Content-Type: <content type>

<file bytes>
```

### Response (not logged in)
```
Status: 403 Forbidden
```

### Response (not found)
```
Status: 404 Not Found
```

## `DELETE /files/<filename>`

This endpoint is used to delete a file from the user's personal storage.

### Request
```
X-Session: <session token>
```

### Response (success)
```
Status: 204 No Content
```

### Response (not logged in)
```
Status: 403 Forbidden
```

### Response (not found)
```
Status: 404 Not Found
```

## `GET /files`

This endpoint is used to list files in this user's personal storage.

### Request
```
X-Session: <session token>
```

### Response (success)
```
Content-Type: application/json

[
  <filenames...>
]
```

### Response (not logged in)
```
Status: 403 Forbidden
```

# Submission

To submit the assignment, put the code in a repository under your GitHub account and send a link to the repository to your recruiter contact at New Alchemy.

