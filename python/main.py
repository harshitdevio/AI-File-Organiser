from fastapi import FastAPI, HTTPException
from schemas import FileRequest, BatchFileRequest, FolderBatchRequest
from classifier import classify_file, process_batch
from typing import List
from config import TARGET_TOPICS
from client import send_results_to_service

app = FastAPI()

@app.get("/health")
async def health():
    return {"status": "ok", "model": "facebook/bart-large-mnli"}

@app.post("/detect-topic")
async def detect_topic(request: FileRequest):
    topics = request.topics if request.topics else TARGET_TOPICS
    result = classify_file(request.filepath, topics, request.mime_type)
    if result.get("error"):
        raise HTTPException(status_code=400, detail=result["error"])
    return result

@app.post("/detect-topics-batch")
async def detect_topics_batch(request: BatchFileRequest):
    results = [
        classify_file(f.filepath, f.topics if f.topics else TARGET_TOPICS, f.mime_type) 
        for f in request.files
    ]
    
    return {
        "total": len(results),
        "successful": sum(1 for r in results if not r.get("error")),
        "failed": sum(1 for r in results if r.get("error")),
        "results": results
    }

@app.post("/process-unit-2")
async def process_unit_2(request: FolderBatchRequest):
    results = process_batch(request.filepaths)
    return {
        "count": len(results),
        "results": results
    }

@app.post("/classify-and-sync")
async def run_pipeline(paths: List[str]):
    final_results = process_batch(paths) 
    
    status = await send_results_to_service(final_results)
    return {"sent": len(final_results), "remote_status": status}