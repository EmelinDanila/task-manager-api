# Task Manager API

Task Manager API is a RESTful API for managing tasks. It allows users to create, update, delete, and view tasks. The application is built with Go, using PostgreSQL for the database and JWT for authentication.

## Technologies Used

- **Go** — Main programming language.
- **PostgreSQL** — Database for storing tasks.
- **JWT** — User authentication.
- **Gin** — Web framework for handling HTTP requests.
- **Swagger** — API documentation.
- **Docker** — Containerization of the application.

---

## Installation and Running with Docker Compose

### **Prerequisites**
- Installed [Git](https://git-scm.com/)
- Installed [Docker](https://www.docker.com/)
- Installed [Docker Compose](https://docs.docker.com/compose/)

### **Steps**
1. **Clone the repository**
   ```sh
   git clone https://github.com/your-username/task-manager-api.git
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

4. **API is now available at**
   ```
   http://localhost:8080
   ```

---

## Main API Endpoints

| Method  | Endpoint        | Description                                | Authentication |
|---------|----------------|--------------------------------------------|---------------|
| `POST`  | `/register`    | Register a new user                       | ❌            |
| `POST`  | `/login`       | Authenticate user, receive JWT token      | ❌            |
| `GET`   | `/tasks`       | Get all tasks for the authenticated user  | ✅            |
| `POST`  | `/tasks`       | Create a new task                         | ✅            |
| `PUT`   | `/tasks/{id}`  | Update a task                             | ✅            |
| `DELETE`| `/tasks/{id}`  | Delete a task                             | ✅            |

🔹 **Swagger documentation is available at:**  
```
http://localhost:8080/swagger/index.html
```

---

## License

This project is licensed under the MIT License.