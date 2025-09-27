```markdown
# Sequence Diagram - Topic CRUD

```mermaid
sequenceDiagram
    participant Client as Frontend
    participant API as Echo API
    participant TopicH as TopicHandler
    participant DB as GORM
    participant PostgreSQL as PostgreSQL

    Client->>API: POST /api/topics
    API->>TopicH: Bind JSON -> validate
    TopicH->>DB: Insert topic
    DB->>PostgreSQL: INSERT INTO topics...
    PostgreSQL-->>DB: ID baru
    DB-->>TopicH: Return Topic object
    TopicH-->>Client: JSON 201 Created

    Client->>API: GET /api/topics
    API->>TopicH: Read all topics
    TopicH->>DB: SELECT * FROM topics
    DB->>PostgreSQL: SELECT query
    PostgreSQL-->>DB: Result
    DB-->>TopicH: Map to struct
    TopicH-->>Client: JSON array

    Client->>API: GET /api/topics/:id
    TopicH->>DB: SELECT topic WHERE id=?
    DB->>PostgreSQL: SELECT query
    PostgreSQL-->>DB: Result
    DB-->>TopicH: Topic object
    TopicH-->>Client: JSON 200 OK

    Client->>API: PUT /api/topics/:id
    TopicH->>DB: Update topic
    DB->>PostgreSQL: UPDATE topics SET name=...
    PostgreSQL-->>DB: OK
    DB-->>TopicH: Topic object updated
    TopicH-->>Client: JSON 200 OK

    Client->>API: DELETE /api/topics/:id
    TopicH->>DB: Delete topic
    DB->>PostgreSQL: DELETE FROM topics WHERE id=?
    PostgreSQL-->>DB: OK
    DB-->>TopicH: Confirmation
    TopicH-->>Client: 204 No Content