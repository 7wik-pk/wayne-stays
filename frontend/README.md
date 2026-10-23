# Wayne Stays: Frontend

React 19 + TypeScript + Vite single-page app for Wayne Stays.

```bash
npm install
npm run dev          # http://localhost:5173 (expects the backend on :8080)
npm test             # Vitest, run once
npm run test:watch   # Vitest, watch mode
npm run lint         # ESLint
npm run build        # Type-check + production build into dist/
```

Requires Node.js 20.19+. The API base URL comes from `VITE_API_URL`
(see `.env.example`), defaulting to `http://localhost:8080/api/v1`.

## Layout

```
src/
├── main.tsx          # Entry point (createRoot + StrictMode)
├── App.tsx           # Mounts the router
├── router.tsx        # Route table (React Router 7)
├── layouts/          # Page frames: SiteLayout = Navbar + <Outlet /> + Footer
├── pages/            # One component per route
├── components/       # Reusable components, each with its own .css file
├── api/              # Typed fetch client (client.ts) + response types (types.ts)
├── hooks/            # Custom hooks (useApi: loading/error state + abort on unmount)
└── test/             # Vitest setup + fetch mock
```

Tests live next to the code they cover (`*.test.ts(x)`).

## Learning notes

New React concepts are logged with their Vue 3 equivalents in
[react_vue_map.md](react_vue_map.md).
