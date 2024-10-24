const express = require('express');
const { Client } = require('pg');

const app = express();

app.get('/api/health', async (req, res) => {
    const client = new Client({
        connectionString: 'postgresql://hackathon:hackathon@postgresql:5432/hackathon'
    });

    try {
        await client.connect();
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
