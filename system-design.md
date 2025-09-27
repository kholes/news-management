# System Design - News Management

## 1. High-Level Architecture

```mermaid
flowchart LR
    A[User / Frontend] --> B[API Gateway (Echo)]
    B --> C[Handlers]
    C --> D[Service / Business Logic]
    D --> E[(PostgreSQL)]
    C --> F[Swagger / OpenAPI]
