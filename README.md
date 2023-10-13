# GovTech-Technical

Student portal for teachers

## Prerequisites

Before getting started, make sure you have the following software installed on your system:

- [Docker](https://www.docker.com/get-started)

## Quick Start with Docker Compose

1. **Copy the docker-compose.yml file**:

   ```yml
   version: "3.8"
   name: govtech-technical
   services:
     db:
       image: postgres:alpine3.18
       volumes:
         - ./databases/db:/var/lib/postgresql/data
       environment:
         - POSTGRES_NAME=pg
         - POSTGRES_USER=pg
         - POSTGRES_PASSWORD=pg
       ports:
         - "5432:5432"
       healthcheck:
         test: ["CMD-SHELL", "pg_isready -U pg -d pg"]
         interval: 1s
         timeout: 1s
     govtech-technical:
       image: mingyuanc/govtech-technical:latest
       depends_on:
         db:
           condition: service_healthy
       ports:
         - "8282:8282"
       environment:
         - DATABASE_URL=postgres://pg:pg@db:5432/pg
       command: ["sh", "-c", "./backend"]
   ```

1. **Then run `docker-compose up` to spin up the server and the database**:

   The server will be listening on `localhost:8282`

## Endpoints

### `GET /api/ping` - Health route

Checks if the server is up and running

Example of usage:

`/api/ping`

Expected outcome:

```
pong
```

### `GET api/commonstudents/` - Finds common student among list of teachers

Retrieve a list of students common to a given list of teachers, in other words, retrieve students who are register to all of the given teachers. As per user story, suspended students will be shown.

Constraints:

- The teacher query parameter cannot be empty.
- The teacher must already be registered.
- Email must be valid.

Example of usage:

`/api/commonstudents?teacher=teacherken%40gmail.com`

Expected outcome:

```json
{
  "students": [
    "commonstudent1@gmail.com",
    "commonstudent2@gmail.com",
    "student_only_under_teacher_ken@gmail.com"
  ]
}
```

Example Error outcome:

1. Invalid teacher email

```json
{
  "error": "Teacher parameter at index 0 is an invalid email: testNotFound@gmail"
}
```

2. Teacher is not registered

```json
{
  "error": "Teacher parameter at index 0 is not found: testNotFound@gmail.com"
}
```

### `POST /api/register` - Registers a teacher to a list of students

A teacher can be registered to multiple students, students can also be register to multiple teachers.

If the teacher / students are not found in the database, they will be automatically registered and their association will be captured as well.

Constraints:

- As per user guide, the students array in the request body cannot be empty.
- Email must be valid.

Example of request body:

```json
{
  "teacher": "teacherken@gmail.com"
  "students":
    [
      "studentjon@gmail.com",
      "studenthon@gmail.com"
    ]
}
```

Expected outcome:

`HTTP 204`

Example Error outcome:

1. Invalid student email

```json
{
  "error": "Key: 'RegisterBody.Students[0]' Error:Field validation for 'Students[0]' failed on the 'email' tag"
}
```

2. No student provided

```json
{
  "error": "Must provide at least one student"
}
```

### `POST api/suspend/` - Suspends a student

Suspend a given student by his email

Constraints:

- The student query parameter cannot be empty.
- The student must already be registered.

Example of request body:

```json
{
  "student": "suspend1@gmail.com"
}
```

Expected outcome:

`HTTP 204`

Example Error outcome:

1. Invalid student email

```json
{
  "error": "Key: 'SuspendBody.Student' Error:Field validation for 'Student' failed on the 'email' tag"
}
```

2. Student is not registered

```json
{
  "error": "Student not found: studentmary@gmail.com"
}
```

### `POST api/retrievefornotifications/` - Retrieve a list of student who can receive a notification

The request body should contain:

- The teacher who is sending the notification.
- The text of the notification.

To receive notifications, the student:

- MUST NOT be suspended,
- Be registered with the teacher OR is mentioned

Constraints:

- Student mentioned must have valid email and is registered.
- Student email must have the prefix '@' to be recognized.
- List of students should not have any duplicate/repetitions

Example of request body:

```json
{
  "teacher": "teacherken@gmail.com",
  "notification": "Hello students! @studentagnes@gmail.com @studentmiche@gmail.com"
}
```

Expected outcome:

```json
{
  "recipients": [
    "studentbob@gmail.com",
    "studentagnes@gmail.com",
    "studentmiche@gmail.com"
  ]
}
```

Example Error outcome:

1. Invalid student email

```json
{
  "error": "Student mentioned has an invalid email: notsuspend1.com"
}
```

2. Student is not registered

```json
{
  "error": "Student parameter at index 0 is not found: notsuspended@gmail.com"
}
```

3. Teacher is not registered

```json
{
  "error": "Teacher not found: teacherkenner@gmail.com"
}
```
