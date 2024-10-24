from fastapi import FastAPI
from databases import Database

app = FastAPI()

DATABASE_URL = "postgresql://hackathon:hackathon@postgresql:5432/hackathon"
database = Database(DATABASE_URL)


@app.get("/api/health")
async def health_check():
    try:
        await database.connect()
        await database.execute("SELECT 1")
        await database.disconnect()
        return {"status": "Database is accessible"}
    except Exception as e:
        return {"status": "Database is not accessible", "error": str(e)}, 500
