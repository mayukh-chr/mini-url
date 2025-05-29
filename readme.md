# URL Shortener

This project is a simple URL shortener with a Go backend and a React frontend using a SQLite3 database. The routing is done using [mux](github.com/gorilla/mux) over [net/http](https://pkg.go.dev/net/http)

## Features

### Core URL Shortening
- **Create Short URLs**: Convert long URLs into short, manageable links
  - Automatic random code generation (6 characters)
  - Custom short code support (user-defined codes)
  - Duplicate short code prevention
  - JSON API response with generated short URL

### URL Management
- **Retrieve Original URLs**: Redirect short URLs to their original destinations
  - Automatic access count tracking
  - Real-time click analytics
  - HTTP 302 redirect to original URL

- **Update Short URLs**: Modify existing URL mappings
  - Change the destination URL for existing short codes
  - Update short codes with new custom codes
  - Conflict detection for duplicate codes

- **Delete Short URLs**: Remove URL mappings from the system
  - Complete removal of short URL entries
  - Clean database management

### Analytics & Statistics
- **Access Statistics**: Track usage metrics for short URLs
  - View click counts for any short code
  - Real-time access counting
  - JSON API for statistics retrieval

### Technical Features
- **RESTful API**: Full HTTP API with proper status codes
- **SQLite Database**: Persistent storage with automatic table creation
- **Error Handling**: Comprehensive error responses and logging
- **Web Interface**: HTML form interface for easy URL management
- **React Frontend**: Modern UI with sidebar navigation for all operations
- **Cross-Platform**: Works on Windows, macOS, and Linux

### API Endpoints
- `POST /shorten` - Create new short URL
- `GET /u/{code}` - Redirect to original URL
- `PUT /u/{code}` - Update existing short URL
- `DELETE /u/{code}` - Delete short URL
- `GET /stats/{code}` - Get access statistics
- `GET /shorten` - Web interface for URL management

## Project Structure

```
/ (root)
│
├── frontend/    # React app
│
├── main.go      # Go backend entry point
```

## Prerequisites

- [Go](https://golang.org/dl/) installed
- [Node.js](https://nodejs.org/) and [npm](https://www.npmjs.com/) installed

## Getting Started

1. **Clone the repository:**

    ```bash
    git clone https://github.com/your-username/url-shortner.git
    cd url-shortner
    ```

2. Follow the steps below to run the backend and frontend.

## Running the Backend

1. Open a terminal in the project root.
2. Run:

    ```bash
    go run main.go
    ```

    The backend will start on `localhost:8080`.

## Running the Frontend

1. Open a terminal in the `frontend` folder.
2. Install dependencies:

    ```bash
    npm install
    ```

3. Start the React app:

    ```bash
    npm start
    ```

    The frontend will start, usually on `localhost:3000`.

## Usage

- Access the frontend at [http://localhost:3000](http://localhost:3000).


