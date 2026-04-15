from pydantic import BaseModel
from typing import List, Optional

class FileRequest(BaseModel):
    filepath: str
    topics: List[str]
    mime_type: Optional[str] = None

class BatchFileRequest(BaseModel):
    files: List[FileRequest]