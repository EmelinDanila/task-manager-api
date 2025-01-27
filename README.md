
 Task Manager API
================

Task Manager API is a RESTful API for managing tasks. The application will allow users to create, update, delete, and view tasks. It will use Go for the backend, PostgreSQL for the database, and JWT for authentication.

Technologies
-------------

- Go — main programming language.
- PostgreSQL — database for storing task data.
- JWT (JSON Web Tokens) — for user authentication and authorization.
- Gin — web framework for handling HTTP requests.
- Swagger — API documentation.
- Docker — containerization of the application.

Description
-----------

The goal of this project is to build a simple Task Manager API with the following features:
- User registration and login (JWT authentication).
- CRUD operations for managing tasks.
- Data validation for all inputs.
- API documentation using Swagger.

Installation
------------

### Prerequisites

- Go (1.18 or higher)
- PostgreSQL
- Docker (for containerization)

### Steps to Run Locally

1. Clone the repository:
```bash
git clone https://github.com/your-username/task-manager-api.git   
```

2.Navigate to the project directory:
```bash
cd task-manager-api
```

3.Set up environment variables for the database connection and JWT secret.

4.Run the application using Docker or directly with Go.

Endpoints (Planned)
```bash
POST /api/register — register a new user.
POST /api/login — authenticate a user and receive a JWT.
GET /api/tasks — get all tasks (authentication required).
POST /api/tasks — create a new task (authentication required).
PUT /api/tasks/{id} — update a task (authentication required).
DELETE /api/tasks/{id} — delete a task (authentication required).
```

License
---------
This project is licensed under the MIT License. 
