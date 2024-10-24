require('dotenv').config();
const express = require('express');
const { OpenAI } = require('openai');
const { Pool } = require('pg');
// const cron = require('node-cron');

const app = express();
app.use(express.json());

const dbConfig = {
    connectionString: 'postgresql://hackathon:hackathon@postgresql:5432/hackathon'
};
const pool = new Pool(dbConfig);

const openai = new OpenAI({ apiKey: process.env.OPENAI_API_KEY || "", baseURL: process.env.OPENAI_BASE_URL });

// setInterval(processNews, 10 * 1000); // запускаем каждые 10 секунд
// cron.schedule('*/10 * * * * *', async () => {
//     await processNews()
// });


const queueStatuses = {
    not_started: 'not_started',
    pending: 'pending',
    processing: 'processing',
    completed: 'completed', 
}


const processNews = async (newsList, prompt) => {
    console.log('Обрабатываю новости...');
    // const parsedNews = parseNews(newsList);
    const parsedNews = JSON.stringify(newsList.map(news => news.content));

    console.log({parsedNews})

    const processedContent = await sendToGPT(parsedNews, prompt);

    return processedContent;
}

const addTask = async (task) => {
    const q = 
        `INSERT INTO gpt_requests_queue(prompt, request, status)
        VALUES($1, $2, 'pending')
        RETURNING id;`;
    
    try {
        const res = await pool.query(q, [prompt, request]);
        console.log('Inserted ID:', res.rows[0].id);
        return res.rows[0].id;
    } catch (err) {
        console.error(err.stack);
    }
}

async function updateTask(id, response, status) {
    const q = 
        `UPDATE gpt_requests_queue
        SET response = $1, status = $2, updated_at = CURRENT_TIMESTAMP
        WHERE id = $3;`;

    try {
        await pool.query(q, [response, status, id]);
        console.log('Request updated');
    } catch (err) {
        console.error(err.stack);
    }
}


// async function fetchUnprocessedNews() {
//     const client = await pool.connect();

//     try {
//         const res = await client.query('SELECT id, news_content FROM news WHERE processed_content IS NULL');
//         return res.rows;
//     } finally {
//         await client.end();
//     }
// }

// const processNewsFromDB = async () => {
//     const newsList = await fetchUnprocessedNews();

//     const processedContent = processNews(newsList);

//     if (processedContent) {
//         await updateProcessedNews(newsList.map(news => news.id));
//     }

//     return processedContent;
// }


// async function updateProcessedNews(newsIds) {
//     const client = await pool.connect();

//     try {
//         const query = 'UPDATE news SET is_processed = true WHERE id = ANY($1)';
//         await client.query(query, newsIds);
//     } finally {
//         await client.end();
//     }
// }



async function sendToGPT(newsText, customPrompt) {
    try {
        const prompt = customPrompt ||
            `You are a helpful assistant for text processing.
            Compare news content for semantic similarity.
            News in json format.
            Return summary for news.
            Answer in Russian.`;

        return await getCompletionChatGPT(prompt, newsText);
    } catch (error) {
        console.error('Ошибка при обработке новости AI-сервисом:', error);
    }
}

async function getCompletionChatGPT(prompt, news, temperature) {
    const response = await openai.chat.completions.create({
        model: process.env.CHAT_MODEL || "gpt-4o",
        messages: [
            { role: "system", content: prompt },
            { role: "user", content: news },
        ],
        temperature: temperature || 0.3,
    })

    return response.choices[0].message.content;
}


app.post('/api/process-by-gpt', async (req, res) => {
    const data = req.body;
    
    if (!data) {
        return res.status(400).json({ error: 'No data provided' });
    }

    if (!data.news) {
        return res.status(400).json({ error: 'No news provided' });
    }

    try {
        const result = await processNews(data.news, data?.prompt);
        res.json({ status: 'OK', text: result });
    } catch (err) {
        res.status(500).json({ status: 'News processing error', error: err.message });
    }
});

app.get('/api/health', async (req, res) => {
    const client = await pool.connect();

    try {
        await client.query('SELECT 1');
        res.json({ status: 'Database is accessible' });
    } catch (err) {
        res.status(500).json({ status: 'Database is not accessible', error: err.message });
    } finally {
        await client.end();
    }
});

app.listen(3000, () => {
    console.log('Server is running on port 3000');
});
