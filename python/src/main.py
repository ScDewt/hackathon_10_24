# coding: utf-8
# -*- coding: utf-8 -*-

import uvicorn


if __name__ == '__main__':
    uvicorn.run('app:app',
                host="localhost",
                port=2894,
                reload=True,
                log_level="info",
                workers=1,
                )