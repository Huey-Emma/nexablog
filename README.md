# Project Support

### Introduction

NexablogAPI is an API that allows user create blog articles and perform CRUD - Create - Read - Update - Delete operations

### Project Support Features

-   Users can register
-   Users can login
-   Users can create a blog post
-   Users can fetch a blog post
-   Users can update a blog post
-   Users can delete a blog post

### Installation Guide

-   Clone the repository [here](https://github.com/Huey-Emma/nexablog.git)
-   The main branch is the most stable branch at any given time. Ensure you are working from it

-   Run `go mod download` to install all dependencies
-   You can either work with PGAdmin or PSQL to access your datasbase
-   Create a `.env` file in the project root folder and add your variables. See `.env.sample` for assistance

### API Endpoints

| HTTP Verbs | Endpoints                 | Action               |
| ---------- | ------------------------- | -------------------- |
| POST       | `api/users`               | Register a user      |
| GET        | `api/users/me`            | Fetch user's profile |
| POST       | `api/tokens/authenticate` | Get auth token       |
| GET        | `api/posts`               | Fetch all posts      |
| POST       | `api/posts`               | Create a post        |
| GET        | `api/posts/1`             | Fetch a post by id   |
| PUT        | `api/post/1`              | Update a post        |
| DELETE     | `api/post/1`              | Delete a post        |

### Technologies Used

-   [Golang](https://go.dev)
    An open source programming language supported by Google. It allows for installation and management of dependencies and communication with databases.
-   [go-chi](https://go-chi.io/#/)
    A lightweight, idiomatic and composable router for building Go HTTP services.
-   [Postgres](https://www.postgresql.org)
    A powerful open source object relational database that has a strong reputation for reliability, feature robustness and performance

### License

This project is available for use under the MIT License
