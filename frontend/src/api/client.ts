import type { Accommodation, Merchandise } from './types';

export const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';

// Thrown for non-2xx responses so callers can tell HTTP failures from network ones
export class ApiError extends Error {
  readonly status: number;

  constructor(status: number, message: string) {
    super(message);
    this.name = 'ApiError';
    this.status = status;
  }
}

// GETs a path under API_URL and returns the `data` field of the JSON envelope
export async function apiGet<T>(path: string, signal?: AbortSignal): Promise<T> {
  const response = await fetch(`${API_URL}${path}`, { signal });

  if (!response.ok) {
    const body = await response.json().catch(() => null);
    throw new ApiError(response.status, body?.error ?? `HTTP error! status: ${response.status}`);
  }

  const body = await response.json();
  return body.data;
}

// One function per endpoint. They are module-level, so their identity never
// changes, which makes them safe to pass straight into useApi.
export const api = {
  listAccommodations: (signal?: AbortSignal) =>
    apiGet<Accommodation[]>('/accommodations', signal),

  listMerchandise: (signal?: AbortSignal) =>
    apiGet<Merchandise[]>('/merchandise', signal),
};
