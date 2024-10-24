from aiogram import Bot, Dispatcher, types
from aiogram.dispatcher.router import Router
from aiogram import F
import asyncio
import sqlite3
import random
import string

API_TOKEN = '%API_TOKEN%'

async def set_bot_commands_and_description():
    bot = Bot(token=API_TOKEN)

    await bot.set_my_description(description=(
        "FinamНовости — это ваш надежный помощник для получения уникального и актуального контента из множества телеграм-каналов. "
        "Укажите боту каналы, которые вы читаете, и он будет сканировать их, отбирая только оригинальные новости и посты. "
        "Если текст пересекается на 40% и более."
    ))

bot = Bot(token=API_TOKEN)
dp = Dispatcher()
router = Router()
dp.include_router(router)

def init_db():
    conn = sqlite3.connect('bot_data.db')
    cursor = conn.cursor()
    cursor.execute('''CREATE TABLE IF NOT EXISTS users (
        chat_id INTEGER PRIMARY KEY,
        token TEXT,
        has_sent_command INTEGER DEFAULT 0
    )''')
    conn.commit()
    conn.close()

def generate_token():
    """Генерирует случайный токен из 16 символов."""
    return ''.join(random.choices(string.ascii_letters + string.digits, k=16))

@router.message(F.text == '/start')
async def start_command(message: types.Message):
    chat_id = message.chat.id

    conn = sqlite3.connect('bot_data.db')
    cursor = conn.cursor()

    cursor.execute('SELECT has_sent_command FROM users WHERE chat_id = ?', (chat_id,))
    result = cursor.fetchone()

    if result is None:
        cursor.execute('INSERT INTO users (chat_id) VALUES (?)', (chat_id,))
        conn.commit()

    if result is None or result[0] == 0:
        description = (
            "Начнём? "

        )
        await message.answer(description)

    conn.close()

@router.message(F.text == '/token')
async def token_command(message: types.Message):
    chat_id = message.chat.id
    token = generate_token()

    conn = sqlite3.connect('bot_data.db')
    cursor = conn.cursor()

    cursor.execute('SELECT token FROM users WHERE chat_id = ?', (chat_id,))
    result = cursor.fetchone()

    if result is None:
        cursor.execute('INSERT INTO users (chat_id, token, has_sent_command) VALUES (?, ?, 1)', (chat_id, token))
        conn.commit()
        await message.answer(f"Регистрация успешна! Ваш токен: {token}")
    else:
        await message.answer("Вы уже зарегистрированы и у вас есть токен.")

    cursor.execute('UPDATE users SET has_sent_command = 1 WHERE chat_id = ?', (chat_id,))
    conn.commit()
    conn.close()

async def main():
    init_db()

    await dp.start_polling(bot)

if __name__ == '__main__':
    asyncio.run(main())
