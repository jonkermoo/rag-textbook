import os
from dotenv import load_dotenv
from openai import OpenAI

# Load environment variables
load_dotenv('../.env')

# Initialize OpenAI client
client = OpenAI(api_key=os.getenv('OPENAI_API_KEY'))

# Test: Generate an embedding for a simple sentence
test_text = "This is a test sentence for embeddings."

try:
    response = client.embeddings.create(
        model="text-embedding-3-small",
        input=test_text
    )
    
    embedding = response.data[0].embedding
    
    print(f"OpenAI API is working!")
    print(f"Generated embedding with {len(embedding)} dimensions")
    print(f"First 5 values: {embedding[:5]}")
    
except Exception as e:
    print(f"Error connecting to OpenAI: {e}")