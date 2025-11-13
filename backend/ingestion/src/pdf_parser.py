from PyPDF2 import PdfReader

class PDFParser:
    def __init__(self, pdf_path):
        """Initialize with path to PDF file"""
        self.pdf_path = pdf_path
        self.reader = PdfReader(pdf_path)
        print(f"Loaded PDF: {pdf_path} ({len(self.reader.pages)} pages)")
    
    def extract_text(self):
        """
        Extract text from all pages
        Returns: List of tuples (page_number, text)
        """
        pages = []
        for i, page in enumerate(self.reader.pages, start=1):
            text = page.extract_text()
            if text.strip():  # Only include pages with text
                pages.append((i, text))
        
        print(f"Extracted text from {len(pages)} pages")
        return pages
    
    def get_metadata(self):
        """Get PDF metadata (title, author, etc.)"""
        return self.reader.metadata