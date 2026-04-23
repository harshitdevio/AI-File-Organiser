import os
from transformers import pipeline
from typing import List, Dict, Optional
from file_utils import read_file_content
from config import TARGET_TOPICS

# Initialize model once (Singleton-ish pattern)
print("Loading BART model... please wait.")
classifier_pipeline = pipeline(
    "zero-shot-classification", 
    model="facebook/bart-large-mnli", 
    device=-1
)

def classify_file(
        filepath: str, 
        topics: List[str], 
        mime_type: Optional[str] = None) -> Dict:
    if not os.path.exists(filepath):
        return {
            "filepath": filepath,
            "error": "File not found",
            "top_topic": None,
            "confidence": None,
            "mime_type": mime_type
        }
    
    try:
        content = read_file_content(filepath, mime_type)
        
        if not content.strip():
            return {
                "filepath": filepath,
                "error": "No text content extracted",
                "top_topic": None,
                "confidence": None,
                "mime_type": mime_type
            }
        
        result = classifier_pipeline(content, topics, truncation=True)
        
        return {
            "filepath": filepath,
            "top_topic": result['labels'][0],
            "confidence": float(result['scores'][0]),
            "all_scores": dict(zip(result['labels'], [float(s) for s in result['scores']])),
            "mime_type": mime_type,
            "error": None
        }
    except Exception as e:
        return {
            "filepath": filepath,
            "error": str(e),
            "top_topic": None,
            "confidence": None,
            "mime_type": mime_type
        }
    
def process_batch(filepaths: List[str]) -> List[Dict]:
    results = []
    for path in filepaths:
        result = classify_file(path, TARGET_TOPICS)
        results.append(result)
    return results