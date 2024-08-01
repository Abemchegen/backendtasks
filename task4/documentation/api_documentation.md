# Task Manager API Documentation

The link for postman documentation is : https://documenter.getpostman.com/view/37350482/2sA3kdAck8

This project provides a suite of endpoints for managing tasks within a task management system. The collection includes functionality for creating, retrieving, updating, and deleting tasks. Below is a summary of each endpoint:

## Endpoints

### GET /tasks

**Description**: This endpoint retrieves a list of tasks.

**URL**: `http://localhost:8080/tasks`

**Response**: The response of this request is documented below as a JSON schema.

```json
{
  "1": {
    "id": 1,
    "title": "Task 1",
    "description": "first task",
    "duedate": "0001-01-01T00:00:00Z",
    "status": "Pending"
  },
  "2": {
    "id": 2,
    "title": "Task 2",
    "description": "second task",
    "duedate": "0001-01-01T00:00:00Z",
    "status": "Pending"
  },
  "3": {
    "id": 3,
    "title": "Task 3",
    "description": "third task",
    "duedate": "0001-01-01T00:00:00Z",
    "status": "Pending"
  }
}
```

### GET /tasks/:id

**Description**: This endpoint retrieves the details of a specific task identified by the provided ID.

**URL**: http://localhost:8080/tasks/:id

**Response**: The response of this request is a JSON object representing the details of the task. The JSON schema for the response can be documented as follows:

```json
{
  "id": 1,
  "title": "Task 4",
  "description": "fourth task",
  "duedate": "0001-01-01T00:00:00Z",
  "status": "Pending"
}
```

**Path Variables**:
**id**: The ID of the task.

### PUT /tasks

**Description**: This endpoint allows the user to update a specific task identified by the provided ID.

**URL**: http://localhost:8080/tasks/:id

**Request**: The request should be made using an HTTP PUT method to the specified URL, with the ID of the task included in the URL path. The request body should contain the updated details of the task including the title, description, and status.

**Request Body**:

title (string, required): The updated title of the task.
description (string, required)**: The updated description of the task.
status (string, required)**: The updated status of the task.
Response: The response of this request is a JSON schema representing the structure of the response data. This schema will define the properties and their types that will be included in the response when the task is successfully updated.

```json
{
  "id": 1,
  "title": "Task 4",
  "description": "fourth task",
  "duedate": "0001-01-01T00:00:00Z",
  "status": "Pending"
}
```

**Path Variables**:
**id**: The ID of the task.

### DELETE /tasks/:id

**Description**: This endpoint is used to delete a specific task identified by its ID.

**URL**: http://localhost:8080/tasks/:id

**Request**: Method: DELETE

**Response**: The response of this request is a JSON schema representing the structure of the response data. The schema will define the properties and their data types that will be returned upon successful deletion of the task.

**Path Variables**:
**id**: The ID of the task.

### POST /tasks

**Description**: The endpoint allows the creation of a new task by sending an HTTP POST request to the specified URL.

**URL**: http://localhost:8080/tasks

**Request Body**:

Title (string, required): The title of the task.
Description (string, required): Description of the task.
Status (string, required): Current status of the task.
Response: The response of this request is a JSON schema representing the structure of the response data.

Example Request Body:

```json
{
  "Title": "Task 4",
  "Description": "fourth task",
  "Status": "Pending"
}
```

This Markdown file includes the full documentation for all the endpoints provided, including their descriptions, URLs, request bodies, and response schemas.
