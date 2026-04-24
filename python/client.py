import httpx
from typing import List, Dict

async def send_results_to_service(results: List[Dict]):
    url = "https://"
    
    async with httpx.AsyncClient(timeout=30.0) as client:
        response = await client.post(url, json=results)
        return response.status_code