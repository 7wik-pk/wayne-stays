import { describe, expect, it, vi } from 'vitest';
import { API_URL, ApiError, apiGet } from './client';
import { mockFetch } from '../test/mockFetch';

describe('apiGet', () => {
  it('requests the path under API_URL and unwraps the data envelope', async () => {
    const fetchMock = mockFetch({ '/accommodations': [{ id: '1' }] });

    await expect(apiGet('/accommodations')).resolves.toEqual([{ id: '1' }]);
    expect(fetchMock).toHaveBeenCalledWith(`${API_URL}/accommodations`, { signal: undefined });
  });

  it('throws an ApiError carrying the status and server message', async () => {
    vi.stubGlobal('fetch', vi.fn(async () =>
      new Response(JSON.stringify({ error: 'Accommodation not found' }), { status: 404 }),
    ));

    const err = await apiGet('/accommodations/slug/batcave').catch((e) => e);
    expect(err).toBeInstanceOf(ApiError);
    expect(err).toMatchObject({ status: 404, message: 'Accommodation not found' });
  });
});
