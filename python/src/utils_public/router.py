# coding: utf-8
# -*- coding: utf-8 -*-

import pickle
import time
import traceback
import os
import pandas as pd
import httpx
import openai
from fastapi import APIRouter, HTTPException, Depends

from utils_modules.utils_logger import logging
from .utils_model import UtilsRequestModel, UtilsResponseModel
from httpx_socks import AsyncProxyTransport

from dotenv import load_dotenv

load_dotenv()


# os.environ['PROXY_URL'] = "http://api.baseai.ru/openai/v1"
# os.environ['PROXY_TOKEN'] = "sk-bMscQ1wZCpFzCMu9BBA0BaseAI5NNM6dCBEmDb2YOpHDdl"


proxy_transport = AsyncProxyTransport.from_url(os.getenv("PROXY_URL")) if os.getenv("PROXY_URL") else None
http_client = httpx.AsyncClient(transport=proxy_transport)
client_ai = openai.Client(api_key=os.getenv("PROXY_TOKEN"), base_url=os.getenv("PROXY_URL"))

logger = logging.getLogger('LOGGER_NAME')

# initialize router
router_app = APIRouter()


@router_app.post("/summary",
            tags=["summary text data"],
            summary="...",
            response_description="Return HTTP Response JSON file")
async def get_summary_data(data: UtilsRequestModel ):
    
    
    try:
        df = pd.DataFrame([s.__dict__ for s in data.request])
        
        print(df.head())
        
        message = [{'role': 'system', 'content': str(data.prompt)}, 
                    {'role': 'user', 'content': str(df.text.to_list())}]
        
        response = client_ai.chat.completions.create(model="gpt-4o",
                                                messages=message,
                                                max_tokens=3000,
                                                temperature=0.4,
                                                )
        
        # logging response
        logger.info(f"Response: {data.request} : {response.choices[0].message.content}")
        return {"response": response.choices[0].message.content}
    
    except Exception as e:
        print(traceback.format_exc())
        logger.error(f"Error: {e}")
        raise HTTPException(status_code=500, detail=str(e))
