from pydantic import BaseModel
from typing import List, Optional

class FileRequest(BaseModel):
    filepath: str
    topics: List[str]
    mime_type: Optional[str] = None

class BatchFileRequest(BaseModel):
    files: List[FileRequest]

class FolderBatchRequest(BaseModel):
    filepaths: List[str] 
class ClassificationResult(BaseModel):
    filepath: str
    top_topic: Optional[str]
    confidence: Optional[float]
    error: Optional[str] = None