# REST API Server with Golang and Redis

This project is centered around the implementation of a low-latency REST API server in Go. It demonstrates the fundamental CRUD (Create, Read, Update, Delete) operations by leveraging Redis as the data store. The combination of Go's performance characteristics and Redis's in-memory data handling makes this API server highly efficient and responsive, ideal for applications requiring real-time data access.

## Features
- User Authentication using BCrypt
- Basic CRUD operations for items
- RESTful API Endpoints
- Integration with Redis as a datastore

## Architecture 

This project employs Redis due to its performance, low latency, and simplicity, which are essential for the applications built with this codebase. Here are the reasons why Redis is chosen:

- **In-Memory Nature**: Redis primarily stores data in the server's RAM, which provides blazing fast read and write operations as compared to disk storage.

- **Low Latency**: Redis's in-memory nature enables it to offer extremely low data retrieval latency, making it an excellent choice for applications requiring real-time data access.

- **Optional Data Persistence**: Redis, though in-memory, offers the option to persist data to disk. This feature allows for finding the right balance between performance and data durability.

- **Simplicity and Adequacy**: Redis is simple and easy to set up. For applications that don’t demand extensive scalability and where a single-node configuration is sufficient, Redis is an excellent choice.

- **Efficiency in Memory and Transactions**: When the application isn't transaction-heavy and doesn't have extensive memory usage, Redis is advantageous due to its memory-efficient nature and lightweight transaction model.

- **Flexibility Over ACID**: In cases where ACID properties are not crucial for data consistency in the application, Redis is preferable because it offers a more flexible data consistency model than traditional relational databases.


## Getting Started

### Prerequisites

- Go (version >= 1.16)
- Docker and Docker Compose
- Redis

### Installation

1. Clone the repository

```sh
git clone https://github.com/yourusername/your-repo-name.git
```

2. Navigate to the project directory

```sh
cd your-repo-name
```

3. Start the Redis server using Docker Compose

```sh
docker-compose up -d
```

### Running the Server

- The server will start on `http://localhost:8000`

## API Endpoints

| Method | Endpoint           | Description                |
|--------|--------------------|----------------------------|
| POST   | /login             | Authenticate or create a user |
| POST   | /items             | Create a new item          |
| GET    | /items/{id}        | Retrieve an item by ID     |
| PUT    | /items/{id}        | Update an item by ID       |
| DELETE | /items/{id}        | Delete an item by ID       |

## Usage
### Prerequisites
Ensure you have Docker installed on your machine. You can download Docker from here.
Ensure you have docker-compose installed. It normally comes with Docker Desktop, but you can check the installation guide here.

### Getting Started
- Clone the repository to your local machine.
```sh
git clone <repository_url>
cd <repository_directory>
```
Note: Replace <repository_url> with the URL of the Github repository, and <repository_directory> with the name of the directory where the repository is cloned.
- Build and Run the Containers.
  Inside the repository directory, run the following command to build and run the API and Redis containers using docker-compose.
```sh
docker-compose up --build -d
```

#### Authenticate or Create a User

- User Authentication.
  To authenticate a user, send a POST request to /login endpoint with JSON payload containing loginID and password.

```sh
  curl -X POST -H "Content-Type: application/json" -d '{"loginID": "user1", "password": "pass123"}' http://localhost:8000/login
```

Payload structure:
```json
{
  "loginID": "<login_id>",
  "password": "<password>"
}
```
#### Interacting with the API.
The API server will be running on port 8000. Below are the commands for CRUD operations:

- Create an New Item:
Send a POST request to /items endpoint with JSON payload containing id and value.

curl -X POST -H "Content-Type: application/json" -d '{"id": "1", "value": "sample_item"}' http://localhost:8000/items

Payload structure:

```json
{
  "id": "<item_id>",
  "value": "<item_value>"
}
```
- Read an Item:
Send a GET request to /items/{id} endpoint to retrieve an item.

```sh
curl -X GET http://localhost:8000/items/1
```

- Update an Item:
Send a PUT request to /items/{id} endpoint with JSON payload containing value.

```sh
curl -X PUT -H "Content-Type: application/json" -d '{"value": "updated_item"}' http://localhost:8000/items/1
```
Payload structure:
```json
{
  "value": "<new_item_value>"
}
```
- Delete an Item:
  Send a DELETE request to /items/{id} endpoint to delete an item.

```sh
  curl -X DELETE http://localhost:8000/items/1
```

## Security Considerations

This project incorporates a few fundamental security practices but is mainly for educational purposes. When deploying this in a production environment, additional security measures should be taken into consideration.

### Included Security Features:
- Passwords are hashed using BCrypt before storing them in Redis. This ensures that even if there is unauthorized access to the database, the actual passwords are not compromised.

### Not Included and Recommended for Production:
- Securing the Redis Instance: This project doesn't enforce the security of the Redis instance. It's crucial to ensure that your Redis instance is not exposed to the public internet and is properly secured.
- HTTPS: The project does not include HTTPS, which is essential for encrypting data in transit. Implementing HTTPS is recommended to prevent man-in-the-middle attacks.
- Authentication & Authorization: The project lacks a complete authentication and authorization mechanism. Implementing something like OAuth, JWT, or API tokens and role-based access control is highly recommended.
- Rate Limiting: Implementing rate limiting for API requests is important to prevent abuse and is not included in this project.
- Input Validation: The project does not perform extensive input validation. Sanitizing and validating input data is crucial for preventing injection attacks and ensuring data integrity.

These additional security measures are essential for maintaining the security and integrity of the data and should be considered for any production deployment.

## Future Improvements

- Implement role-based authorization.
- Enable the configuration of Redis connection through environment variables.
- Implement rate limiting for API requests.
- Implement HTTPS and stronger authentication mechanisms.
- Enhance input validation and sanitization.


