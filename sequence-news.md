```markdown
# Sequence Diagram - News CRUD

## GET /api/news dengan filter status & topic

```mermaid
sequenceDiagram
    participant Client as User / Frontend
    participant API as News API (Golang Echo)
    participant Handler as NewsHandler
    participant DB as GORM
    participant PostgreSQL as PostgreSQL
    participant Swagger as Swagger

    Client->>API: GET /api/news?status=published&topic_id=2
    API->>Handler: Terima request
    Handler->>DB: Build query dengan filter status & topic
    DB->>PostgreSQL: SELECT * FROM news
        JOIN news_topics ON news.id = news_topics.news_id
        WHERE status='published' AND topic_id=2
    PostgreSQL-->>DB: Hasil query
    DB-->>Handler: Map hasil ke struct News + Topics
    Handler-->>Client: Response JSON
    API->>Swagger: Endpoint terdokumentasi

sequenceDiagram
    participant Client as Frontend
    participant API as Echo API
    participant NewsH as NewsHandler
    participant DB as GORM
    participant PostgreSQL as PostgreSQL

    Client->>API: POST /api/news
    API->>NewsH: Bind JSON -> validate
    NewsH->>DB: Insert news + topic relation
    DB->>PostgreSQL: INSERT INTO news...
    PostgreSQL-->>DB: ID baru
    DB-->>NewsH: Return News object
    NewsH-->>Client: JSON response 201 Created

    Client->>API: PUT /api/news/:id
    API->>NewsH: Bind JSON -> validate
    NewsH->>DB: Update news
    DB->>PostgreSQL: UPDATE news SET ...
    PostgreSQL-->>DB: OK
    DB-->>NewsH: Return updated News
    NewsH-->>Client: JSON 200 OK

    Client->>API: DELETE /api/news/:id
    API->>NewsH: Get ID
    NewsH->>DB: Mark news as deleted
    DB->>PostgreSQL: UPDATE news SET status='deleted'
    PostgreSQL-->>DB: OK
    DB-->>NewsH: Confirmation
    NewsH-->>Client: 204 No Content