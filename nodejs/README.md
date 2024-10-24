
```
-- add enum
CREATE TYPE request_status AS ENUM ('not_started', 'pending', 'completed', 'failed');

-- add table
CREATE TABLE gpt_requests_queue (
    id SERIAL PRIMARY KEY,
    prompt TEXT NOT NULL,
    request TEXT NOT NULL,
    response TEXT,
    status request_status NOT NULL DEFAULT 'not_started',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```
