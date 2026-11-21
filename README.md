# Servery Reviews frontend

## Run the frontend

You need [Node.js](https://nodejs.org) 20.9 or newer.

```bash
npm install
npm run dev
```

Then open http://localhost:3000.

The frontend sends every request that starts with `/api/` to your backend at
`http://localhost:8080` (see `next.config.ts`). If the page shows
`500 Internal Server Error`, your backend probably isn't running on port 8080.

All the frontend code is in `app/page.tsx`.
