import { CreateUrlPayload, CreateUrlResponse, ApiErrorResponse } from '../types/url';

const API_BASE_URL = import.meta.env.VITE_API_URL?.trim() || 'http://localhost:3000';

/**
 * Returns the frontend URL for a given short code (e.g. http://localhost:5173/:shortCode)
 */
export function getFullShortUrl(shortCode: string): string {
  const base = typeof window !== 'undefined' ? window.location.origin : 'http://localhost:5173';
  return `${base.replace(/\/+$/, '')}/${encodeURIComponent(shortCode)}`;
}

/**
 * Returns the backend redirect URL for a given short code (e.g. http://localhost:3000/:shortCode)
 */
export function getBackendRedirectUrl(shortCode: string): string {
  return `${API_BASE_URL.replace(/\/+$/, '')}/${encodeURIComponent(shortCode)}`;
}

/**
 * Sends a request to the backend to create a short URL.
 */
export async function createShortUrl(payload: CreateUrlPayload): Promise<CreateUrlResponse> {
  const endpoint = API_BASE_URL ? `${API_BASE_URL.replace(/\/+$/, '')}/api/urls` : '/api/urls';

  let response: Response;
  try {
    response = await fetch(endpoint, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Accept: 'application/json',
      },
      body: JSON.stringify(payload),
    });
  } catch {
    throw new Error('Unable to reach Shorty. Make sure the API is running.');
  }

  if (!response.ok) {
    let errorData: ApiErrorResponse | null = null;
    try {
      errorData = (await response.json()) as ApiErrorResponse;
    } catch {
      // Non-JSON response body
    }

    if (response.status === 429) {
      throw new Error(errorData?.message || errorData?.error || 'Too many requests. Please try again shortly.');
    }

    if (response.status >= 500) {
      throw new Error('Something went wrong. Please try again.');
    }

    // 400 or other client error: use backend message if available
    const errorMessage = errorData?.error || errorData?.message || 'Please check your input and try again.';
    throw new Error(errorMessage);
  }

  const data = (await response.json()) as CreateUrlResponse;
  if (!data.short_code) {
    throw new Error('Invalid response received from server.');
  }

  return data;
}
