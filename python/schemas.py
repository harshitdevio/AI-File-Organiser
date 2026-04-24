from pydantic import BaseModel
from typing import List, Optional, Dict

class FileRequest(BaseModel):
    filepath: str
    topics: Optional[List[str]] = None
    mime_type: Optional[str] = None

class BatchFileRequest(BaseModel):
    files: List[FileRequest]

class FolderBatchRequest(BaseModel):
    filepaths: List[str] 
class ClassificationResult(BaseModel):
    filepath: str
    top_topic: Optional[str] = None
    confidence: Optional[float] = None
    all_scores: Optional[Dict[str, float]] = None
    mime_type: Optional[str] = None
    error: Optional[str] = None