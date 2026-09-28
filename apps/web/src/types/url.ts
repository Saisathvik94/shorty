export type ExpirationPreset = '1d' | '7d' | '30d';

export interface CreateUrlPayload {
  url: string;
  expires_at: string;
}

export interface CreateUrlResponse {
  short_code: string;
}

export interface ApiErrorResponse {
  error?: string;
  message?: string;
}

export interface ShortenedUrlResult {
  shortCode: string;
  fullShortUrl: string;
  originalUrl: string;
}
