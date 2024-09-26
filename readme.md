# GoLang Server

This project is a GoLang server that includes Sentry integration for error tracking, performance tracing, and profiling.

## Prerequisites

Before you begin, ensure you have the following installed on your machine:

- [Go](https://golang.org/doc/install) (version 1.16 or later)
- [Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)

## Installation

1. **Clone the repository**:
    ```sh
    git clone https://github.com/yourusername/your-repo.git
    cd your-repo
    ```

2. **Install dependencies**:
    ```sh
    go mod tidy
    ```

3. **Add a `.env` file**:
    Create a file named `.env` in the root directory of your project and add the following content:
    ```sh
    SENTRY_DSN=YOUR_DSN_KEY
    SENTRY_RELEASE=your_package@your_version
    ```



## Running the Server

To start the server, run the following command:

```sh
go run main.go