# Playground Pro

Playground Pro is a sport venue reservation system designed to provide users with a convenient platform for reserving sport venues. This project aims to streamline the process of finding, booking, and managing sport venues for both users and venue owners.

# Entity Relational Diagram
![playground-pro-erd](https://github.com/playground-pro-project/playground-pro-api/assets/128290172/4b45432b-0bf1-419c-a2e6-6c8be6905c04)

# How To Run
If you want to run the app locally, you can use Docker to containerize the application and also set up Redis in Docker. This way, you can ensure consistent environments and easily manage dependencies.
You can setup redis by this way:
- Pull the Redis Docker image from the official Docker registry
```
docker pull redis
```

- Run a Redis container
```
docker run -d --name my-redis -p 6379:6379 redis
```

- To Access Redis you can run
```
docker exec -it my-redis redis-cli
```

## MVP Features

The minimum viable product (MVP) of Playground Pro includes the following key features:

1. User Management

    - Registration and login functionalities
    - OTP verification for secure account activation
    - Resending OTP for convenience
    - User profile management, including reading, updating, and deleting profiles
    - Profile picture management, allowing users to upload, delete, and change their profile pictures
    - Role upgrade from "user" to "owner" for venue management privileges
    - Password change functionality

2. Venue Management

    - Users can view a list of available venues
    - Venue owners can register their venues on the platform
    - Venue owners can edit and delete their registered venues
    - Venue image management, enabling owners to add and delete venue images

3. Reservation System

    - Users can make reservations for available sport venues

4. Review System
    - Users can add reviews for sport venues
    - Review deletion functionality for users to manage their reviews

## Technology Stack

The Playground Pro project utilizes the following technologies:

-   **Golang**: The programming language used for the backend development, providing efficiency and performance.
-   **Echo Framework**: A fast and minimalist web framework for building APIs in Golang, facilitating the development of robust and scalable endpoints.
-   **Postman**: A powerful API development environment for testing and documenting APIs during the development process.
-   **Redis**: An open-source, in-memory data structure store used for caching to improve performance and reduce database load.
-   **AWS**: Amazon Web Services (AWS) is utilized for cloud infrastructure, providing services such as file storage, hosting, and more.
-   **Docker**: Containerization technology used to package the application and its dependencies, ensuring consistency and portability across different environments.

## Authors

The Playground Pro project is developed by the following individuals:

-   **Dimas Pradana** - [GitHub](https://github.com/pradanadp)
-   **Dimas Yudhana** - [GitHub](https://github.com/dimasyudhana)

Feel free to reach out to the authors for any inquiries or contributions related to the project.

We hope that Playground Pro brings an enjoyable and seamless experience for users to discover, book, and review sport venues while providing venue owners with efficient management capabilities.
