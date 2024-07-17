# Cancer Cell Detector Backend

## Overview

This project is a backend API for a Cancer Cell Detector application, built using [Golang](https://go.dev/learn/) with the Gin framework. The API processes medical imaging data to detect cancer cells using machine learning models.

## Features

- High-performance API using Go and Gin framework.
- Supports image upload for cancer cell detection.
- Integrates with machine learning models for prediction.
- Uses GORM for seamless PostgreSQL database operations.
- Automatic interactive API documentation.

## Installation

1. **Clone the repository:**

    ```sh
    git clone git@github.com/anojaryal/Cancer-Cell-Detector.git
    ```
    or
    ```sh
    git clone https://github.com/anojaryal/Cancer-Cell-Detector.git
    ```

2. **Navigate to the project directory:**

    ```sh
    cd Cancer-Cell-Detector
    ```

3. **Install the dependencies:**

    ```sh
    go mod download
    ```

4. **Set up environment variables:**

    Create a `.env` file in the project root with the following content:

    ```env
    PORT=''
    DB_URL="host=localhost user='' password='' dbname='' port=5432 sslmode=disable"
    SECRET_KEY
    SMTP_SERVER
    SMTP_PORT
    SENDER_EMAIL
    EMAIL_PASSWORD
    ```

5. **Run the application:**

    ```sh
    go run main.go
    ```

## Contributing

If you'd like to contribute, please fork the repository and use a feature branch. Pull requests are warmly welcome.

## License

This project is licensed under the terms of the *MIT license*.

## Contact

For any questions or suggestions, feel free to reach out to [anoj1810@gmail.com](mailto:anoj1810@gmail.com).
