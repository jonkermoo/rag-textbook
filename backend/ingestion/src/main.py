import sys
from database import Database
from pdf_parser import PDFParser
from chunker import TextChunker
from embeddings import EmbeddingsGenerator

def process_pdf(pdf_path, user_id, title):
    """
    pdf_path: Path to PDF file
    user_id: User who uploaded the textbook
    title: Title of the textbook
    """
    print(f"\n{'='*60}")
    print(f"Processing: {title}")
    print(f"{'='*60}\n")
    
    # Initialize components
    db = Database()
    parser = PDFParser(pdf_path)
    chunker = TextChunker()
    embedder = EmbeddingsGenerator()
    
    try:
        # Create textbook entry in database
        # S3 dummy key for now
        s3_key = f"textbooks/user{user_id}/{title}.pdf"
        textbook_id = db.create_textbook(user_id, title, s3_key)
        
        # Extract text from PDF
        pages = parser.extract_text()
        
        if not pages:
            print("No text found in PDF")
            return
        
        # Chunk the text
        chunks = chunker.chunk_pages(pages)
        
        if not chunks:
            print("No chunks created")
            return
        
        # Generate embeddings in batches for efficiency
        print(f"\nGenerating embeddings for {len(chunks)} chunks...")
        chunk_texts = [chunk['content'] for chunk in chunks]
        embeddings = embedder.generate_batch(chunk_texts, batch_size=50)
        
        # Insert chunks into database
        print(f"\nStoring chunks in database...")
        for i, (chunk, embedding) in enumerate(zip(chunks, embeddings)):
            db.insert_chunk(
                textbook_id=textbook_id,
                content=chunk['content'],
                page_number=chunk['page_number'],
                chunk_index=chunk['chunk_index'],
                embedding=embedding
            )
            
            # Progress indicator
            if (i + 1) % 50 == 0:
                print(f"  Stored {i + 1}/{len(chunks)} chunks")
        
        print(f"Stored all {len(chunks)} chunks")
        
        # Mark textbook as processed
        db.mark_textbook_processed(textbook_id)
        
        print(f"\n{'='*60}")
        print(f"Textbook processed successfully")
        print(f"{'='*60}")
        print(f"Textbook ID: {textbook_id}")
        print(f"Total chunks: {len(chunks)}")
        print(f"Total pages: {len(pages)}")
        print(f"\n")
        
    except Exception as e:
        print(f"\nError during processing: {e}")
        raise
    
    finally:
        db.close()

def main():
    if len(sys.argv) < 4:
        print("Usage: python main.py <pdf_path> <user_id> <title>")
        print("Example: python main.py ../test_data/calculus.pdf 1 'Calculus 101'")
        sys.exit(1)
    
    pdf_path = sys.argv[1]
    user_id = int(sys.argv[2])
    title = sys.argv[3]
    
    process_pdf(pdf_path, user_id, title)

if __name__ == "__main__":
    main()