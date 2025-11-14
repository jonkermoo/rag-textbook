CREATE EXTENSION IF NOT EXISTS vector;

-- Users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    verified BOOLEAN DEFAULT FALSE,
    verification_token VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP
);

-- Textbooks table
CREATE TABLE textbooks (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(500) NOT NULL,
    s3_key VARCHAR(1000) NOT NULL,
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    processed BOOLEAN DEFAULT FALSE
);

-- Chunks table (stores text chunks with embeddings)
CREATE TABLE chunks (
    id SERIAL PRIMARY KEY,
    textbook_id INTEGER REFERENCES textbooks(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    page_number INTEGER,
    chunk_index INTEGER,
    embedding vector(1536),  -- OpenAI embeddings are 1536 dimensions
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index for fast vector similarity search
CREATE INDEX ON chunks USING ivfflat (embedding vector_cosine_ops) WITH (lists = 100);

-- Create index for faster lookups
CREATE INDEX idx_chunks_textbook_id ON chunks(textbook_id);
CREATE INDEX idx_textbooks_user_id ON textbooks(user_id);