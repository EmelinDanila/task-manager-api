# Task Manager API

Task Manager API is a RESTful API for managing tasks. It allows users to create, update, delete, and view tasks. The application is written in Go using PostgreSQL for the database and JWT for authentication.

## Technologies Used

- **Go** — main programming language.
- **PostgreSQL** — database for storing tasks.
- **JWT** — user authentication.
- **Gin** — web framework for handling HTTP requests.
- **Swagger** — API documentation.
- **Docker** — application containerization.

---

## Installation and Running with Docker Compose

### Requirements
- Installed [Git](https://git-scm.com/)
- Installed [Docker](https://www.docker.com/)
- Installed [Docker Compose](https://docs.docker.com/compose/)

### Steps
1. **Clone the repository**
   ```sh
   git clone https://github.com/EmelinDanila/task-manager-api.git
   cd task-manager-api
   ```

2. **Start the containers**
   ```sh
   docker-compose up --build -d
   ```

3. **Check running containers**
   ```sh
   docker ps
   ```

4. **The API is now available at**
   ```
   http://localhost:8080
   ```

---

## Main API Endpoints

| Method  | Endpoint       | Description                                 | Authentication |
|---------|--------------|---------------------------------------------|---------------|
| `POST`  | `/register`  | Register a new user                        | No            |
| `POST`  | `/login`     | User authentication, obtain JWT            | No            |
| `GET`   | `/tasks`     | Get all tasks for the current user         | Yes           |
| `POST`  | `/tasks`     | Create a new task                          | Yes           |
| `PUT`   | `/tasks/{id}`| Update a task                              | Yes           |
| `DELETE`| `/tasks/{id}`| Delete a task                              | Yes           |

Swagger documentation is available at:
```
http://localhost:8080/swagger/index.html
```

---

## License

This project is licensed under the MIT License.