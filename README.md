# News Web Crawler

This project is a **web crawler** developed in Golang that collects news articles from multiple sources, currently including **Detik** and **Liputan** news sites. It extracts and stores the collected data, making it easier to work with and analyze news content programmatically.

## Features

- **Multi-source Crawling**: Built-in crawlers for Detik and Liputan, with the option to expand to other news sites.
- **Configurable**: Set up different parameters and configurations via environment variables.
- **Modular Design**: Each source crawler is modular, allowing for flexibility and easier maintenance.
- **Database Integration**: Automatically saves crawled data to a specified database.

## Prerequisites

- **Go**: Ensure Go is installed on your system. You can download it from [go.dev/dl](https://go.dev/dl/).
- **Database**: Configure the necessary database settings in the `.env` file.

## Project Structure

- **Main Go Files**:
  - `config/config.go`: Manages configuration settings for the project.
  - `database/database.go`: Contains the logic for database connection and interactions.
  - `detik/detik.go` and `liputan/liputan.go`: Crawler modules specifically for Detik and Liputan news.
  - `models/models.go`: Contains data models used across the project.
  - `utils/utils.go`: Provides utility functions.

- **Configuration and Environment Files**:
  - `.env_template`: Contains environment variable examples.
  - `go.mod` and `go.sum`: Manage Go dependencies.

## Getting Started

### 1. Clone the Repository

    git clone <repository_url>
    cd news-crawler-main

### 2. Set up Environment Variables

Use the provided `.env_template` to create a `.env` file with the necessary configurations. Modify the values as needed.

    cp .env_template .env

### 3. Install Dependencies

Ensure all dependencies are installed by running:

    go mod tidy

### 4. Run the Crawler

Start the crawler by running:

    go run main.go

This will start the application, connecting to the specified news sites and saving data to your configured database.

## Configuration

The `.env` file holds critical configuration options:

- **Database Settings**: Configure connection details for your database.
- **Crawl Settings**: Define specific settings for each crawler if applicable.

## Extending the Crawler

To add a new news source:

1. Create a new directory under the project root.
2. Develop a custom crawler following the pattern in `detik.go` or `liputan.go`.
3. Integrate the new crawler into the main logic to start collecting data from the additional source.

## Dependencies

This project relies on Go packages listed in `go.mod` and `go.sum`. To install dependencies, run:

    go mod tidy

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository.
2. Create a new branch (`feature/your-feature`).
3. Commit your changes.
4. Push to the branch.
5. Open a pull request.

## License

This project is licensed under the MIT License.

## Contact

For any questions, feel free to reach out to the project maintainer.

---

This README provides a solid starting point for understanding, configuring, and running the web crawler project.
