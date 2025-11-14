// User types
export interface User {
  id: number;
  email: string;
  verfied: boolean;
  created_at: string;
  last_login: string;
}

export interface loginRequest{
  email: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  user: User;
}

// Textbook types
export interface Textbook {
  id: number;
  user_id: number;
  title: string;
  s3_key: string;
  uploaded_at: string;
  process: boolean;
}

export interface TextbookStatus {
  textbook_id: number;
  title: string;
  processed: boolean;
  chunk_count: number;
  uploaded_at: string;
}

// Query types
export interface QueryRequest{
  textbook_id: number;
  question: string;
}

export interface QueryResponse {
  answer: string;
  sources: Source[];
}

export interface Source {
  page_number: number;
  content: string;
}