import httpx
from typing import List, Dict

async def send_results_to_service(results: List[Dict]):
    url = "http://localhost:8080/api/ingest"
    
    async with httpx.AsyncClient(timeout=30.0) as client:
        response = await client.post(url, json=results)
        return response.status_code