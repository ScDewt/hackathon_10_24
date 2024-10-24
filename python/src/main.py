# coding: utf-8
# -*- coding: utf-8 -*-


from fastapi import FastAPI
from databases import Database

import time

from utils_modules.utils_logger import logging
from utils_public.router import router_app


DATABASE_URL = "postgresql://hackathon:hackathon@postgresql:5432/hackathon"
database = Database(DATABASE_URL)
logger = logging.getLogger(__name__)


app = FastAPI(title="Deduplicated", 
                description="API to summary text", 
                version="0.1")


@app.middleware("http")
async def track_request_stats(request, call_next):

    last_response=time.time()


    # get IP address from request
    client_host = request.client.host

    # done processing request
    response = await call_next(request)
    # latency request processing
    process_time = time.time() - last_response

    # log request
    logger.info(f"Request: {request.method} {request.url} - IP: {client_host} Completed in "
                                        f"{process_time:.2f}s Status code: {response.status_code}")

    return response



@app.get("/api/health")
async def health_check():
    try:
        await database.connect()
        await database.execute("SELECT 1")
        await database.disconnect()
        return {"status": "Database is accessible"}
    except Exception as e:
        return {"status": "Database is not accessible", "error": str(e)}, 500
    


# router for check if API is working
app.include_router(router=router_app, prefix="/api")







