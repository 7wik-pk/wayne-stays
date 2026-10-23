import { vi } from 'vitest';

// Replaces global fetch with one that answers each URL ending in a key of
// `routes` with `{ data: routes[key] }`, the backend's JSON envelope.
export function mockFetch(routes: Record<string, unknown>) {
  const fetchMock = vi.fn(async (input: RequestInfo | URL) => {
    const url = input.toString();
    const match = Object.keys(routes).find((path) => url.endsWith(path));
    if (!match) {
      return new Response(JSON.stringify({ error: 'not found' }), { status: 404 });
    }
    return new Response(JSON.stringify({ data: routes[match] }), { status: 200 });
  });

  vi.stubGlobal('fetch', fetchMock);
  return fetchMock;
}
