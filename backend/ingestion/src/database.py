import os
import psycopg
from pgvector.psycopg import register_vector
from dotenv import load_dotenv

load_dotenv('../.env')

class Database:
    def __init__(self):
        """Initialize database connection"""
        self.conn = psycopg.connect(
            host=os.getenv('DB_HOST'),
            port=os.getenv('DB_PORT'),
            dbname=os.getenv('DB_NAME'),
            user=os.getenv('DB_USER'),
            password=os.getenv('DB_PASSWORD')
        )
        # Register pgvector types
        register_vector(self.conn)
        print("Connected to database")
    
    def create_textbook(self, user_id, title, s3_key):
        """
        Create a textbook entry
        Returns: textbook_id
        """
        with self.conn.cursor() as cur:
            cur.execute("""
                INSERT INTO textbooks (user_id, title, s3_key, processed)
                VALUES (%s, %s, %s, FALSE)
                RETURNING id
            """, (user_id, title, s3_key))
            textbook_id = cur.fetchone()[0]
            self.conn.commit()
            print(f"Created textbook with ID: {textbook_id}")
            return textbook_id
    
    def insert_chunk(self, textbook_id, content, page_number, chunk_index, embedding):
        """
        Insert a text chunk with its embedding
        """
        with self.conn.cursor() as cur:
            cur.execute("""
                INSERT INTO chunks (textbook_id, content, page_number, chunk_index, embedding)
                VALUES (%s, %s, %s, %s, %s)
            """, (textbook_id, content, page_number, chunk_index, embedding))
            self.conn.commit()
    
    def mark_textbook_processed(self, textbook_id):
        """Mark a textbook as fully processed"""
        with self.conn.cursor() as cur:
            cur.execute("""
                UPDATE textbooks SET processed = TRUE WHERE id = %s
            """, (textbook_id,))
            self.conn.commit()
            print(f"Marked textbook {textbook_id} as processed")
    
    def close(self):
        """Close database connection"""
        self.conn.close()
        print("Database connection closed")