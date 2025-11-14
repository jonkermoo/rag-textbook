import sys
import os
from database import Database
from pdf_parser import PDFParser
from chunker import TextChunker
from embeddings import EmbeddingsGenerator

def process_textbook(textbook_id, pdf_path):
    """
    Process an already-uploaded textbook by its ID
    
    textbook_id: ID of the textbook in the database
    pdf_path: Path to the PDF file on disk
    """
    print(f"\n{'='*60}")
    print(f"Processing textbook ID: {textbook_id}")
    print(f"PDF path: {pdf_path}")
    print(f"{'='*60}\n")
    
    # Initialize components
    db = Database()
    parser = PDFParser(pdf_path)
    chunker = TextChunker()
    embedder = EmbeddingsGenerator()
    
    try:
        # Extract text from PDF
        print("Extracting text from PDF...")
        pages = parser.extract_text()
        
        if not pages:
            print("ERROR: No text found in PDF")
            sys.exit(1)
        
        print(f"Extracted {len(pages)} pages")
        
        # Chunk the text
        print("\nChunking text...")
        chunks = chunker.chunk_pages(pages)
        
        if not chunks:
            print("ERROR: No chunks created")
            sys.exit(1)
        
        print(f"Created {len(chunks)} chunks")
        
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
            
            # Progress indicator every 50 chunks
            if (i + 1) % 50 == 0:
                print(f"  Stored {i + 1}/{len(chunks)} chunks")
        
        print(f"Stored all {len(chunks)} chunks")
        
        # Mark textbook as processed
        print("\nMarking textbook as processed...")
        db.mark_textbook_processed(textbook_id)
        
        print(f"\n{'='*60}")
        print(f"Textbook processed successfully!")
        print(f"{'='*60}")
        print(f"Textbook ID: {textbook_id}")
        print(f"Total chunks: {len(chunks)}")
        print(f"Total pages: {len(pages)}")
        print(f"{'='*60}\n")
        
    except Exception as e:
        print(f"\nERROR during processing: {e}")
        import traceback
        traceback.print_exc()
        sys.exit(1)
    
    finally:
        db.close()

if __name__ == "__main__":
    if len(sys.argv) < 3:
        print("Usage: python process_existing.py <textbook_id> <pdf_path>")
        print("Example: python process_existing.py 5 ../uploads/3_calculus.pdf")
        sys.exit(1)
    
    textbook_id = int(sys.argv[1])
    pdf_path = sys.argv[2]
    
    # Verify PDF exists
    if not os.path.exists(pdf_path):
        print(f"ERROR: PDF file not found: {pdf_path}")
        sys.exit(1)
    
    process_textbook(textbook_id, pdf_path)