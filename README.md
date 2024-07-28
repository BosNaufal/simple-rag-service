# Simple RAG Service Go Lang

A simple service for handling retrieval-augmented generation (RAG) with Go.

## Table of Contents
- [Benefit of using RAG](#benefit-of-using-rag)
- [Installation](#installation)
- [Usage](#usage)
- [Makefile Commands](#makefile-commands)
- [License](#license)

## Benefit of using RAG
Incorporating Retrieval-Augmented Generation (RAG) into our application offers several benefits. It enhances the accuracy and relevance of responses by grounding them in retrieved documents, which ensures better handling of specific and rare information. This leads to more contextually appropriate and coherent outputs, reducing the risk of generating incorrect or nonsensical answers (hallucinations). Additionally, RAG allows for dynamic updating of knowledge without needing to retrain the entire model, making it easier to keep the application up-to-date. The scalability of RAG enables efficient processing of large datasets, and its versatility means it can improve various functionalities such as customer support, information retrieval, and content generation within our application.

## Installation

1. **Clone the repository:**
```sh
git clone https://gitlab.com/yourusername/simple-rag-service.git
cd simple-rag-service
```

2. **Install dependencies:**
Ensure you have Go installed. Then run:
```sh
go mod download
```

3. **Install the migrate binary:**
Download and install migrate from [here](https://github.com/golang-migrate/migrate).

## Configuration
### Database Configuration
Set the DATABASE_URL environment variable to point to your PostgreSQL database. Example:
```sh
# don't forget to activate the sslmode when on production
export DATABASE_URL=postgresql://postgres:root@localhost:5431/personal_ai?sslmode=disable
```

## Usage
You can use the provided Makefile commands to run the service and manage database migrations.

## Makefile Commands
Make sure to read the makefile

## License
This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details