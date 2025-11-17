# Vexra

## Overview

Vexra is a cloud-deployed, serverless Retrieval-Augmented Generation (RAG) system that transforms students' uploaded textbooks and lecture notes into personalized study guides and explanations. Students can upload their course materials and ask questions to receive AI-generated answers grounded in their specific textbooks, complete with page citations.

## Tech Stack

### Backend

- Go: API server for handling queries and RAG orchestration
- Python: PDF ingestion and processing service
- PostgreSQL + pgvector: Vector database for semantic search
- Docker - Containerization for local development
- JWT: Stateless authentication tokens
- bcrypt: Secure password Hashing

### AI/ML

- OpenAI Embeddings API: Converts text to 1536-dimensional vectors
- OpenAI GPT-4: Generates contextual answers from retrieved content
- pgvector: PostgreSQL extension for fast vector similarity search

###Frontend

- React + TypeScript: User interface
- Vite: Build tool and dev server
- Tailwind CSS: Styling

### Email Service

- Resend: Email verification and notifications

### Cloud Infrastructure

- AWS EC2: Backend deployment
- Amazon RDS: Managed database
- Amazon S3: Textbook storage

### What am I working on right now?
i hate aws so much security dude

