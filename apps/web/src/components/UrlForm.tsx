import React, { useState } from 'react';
import { ExpirationPreset } from '../types/url';

interface UrlFormProps {
  isLoading: boolean;
  onSubmit: (url: string, expiresAt: string) => Promise<void>;
  onError: (msg: string) => void;
}

export const UrlForm: React.FC<UrlFormProps> = ({ isLoading, onSubmit, onError }) => {
  const [url, setUrl] = useState('');
  const [preset, setPreset] = useState<ExpirationPreset>('1d');

  // Calculate RFC3339 timestamp based on current selection
  const calculateExpiresAt = (): string => {
    const now = new Date();

    if (preset === '1d') {
      const target = new Date(now.getTime() + 1 * 24 * 60 * 60 * 1000);
      return target.toISOString();
    }

    if (preset === '7d') {
      const target = new Date(now.getTime() + 7 * 24 * 60 * 60 * 1000);
      return target.toISOString();
    }

    // Default 30 days (max backend limit)
    const target = new Date(now.getTime() + 30 * 24 * 60 * 60 * 1000);
    return target.toISOString();
  };

  const isValidUrlFormat = (inputUrl: string): boolean => {
    try {
      const parsed = new URL(inputUrl);
      return parsed.protocol === 'http:' || parsed.protocol === 'https:';
    } catch {
      return false;
    }
  };

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (isLoading) return;

    const trimmedUrl = url.trim();
    if (!isValidUrlFormat(trimmedUrl)) {
      onError('Please enter a valid URL.');
      return;
    }

    const expiresAt = calculateExpiresAt();
    await onSubmit(trimmedUrl, expiresAt);
  };

  return (
    <form className="url-form" onSubmit={handleSubmit} noValidate>
      <div className="form-group">
        <label htmlFor="url-input" className="form-label">
          Destination URL
        </label>
        <input
          id="url-input"
          type="url"
          className="form-input"
          placeholder="https://example.com/very-long-url"
          value={url}
          onChange={(e) => setUrl(e.target.value)}
          disabled={isLoading}
          required
          autoFocus
          autoComplete="off"
          spellCheck={false}
        />
      </div>

      <div className="form-row">
        <div className="expiration-container">
          <label id="expiration-label" className="form-label">
            Expiration
          </label>
          <div
            className="preset-group"
            role="radiogroup"
            aria-labelledby="expiration-label"
          >
            {(['1d', '7d', '30d'] as const).map((opt) => (
              <button
                key={opt}
                type="button"
                role="radio"
                aria-checked={preset === opt}
                className={`preset-btn ${preset === opt ? 'active' : ''}`}
                onClick={() => setPreset(opt)}
                disabled={isLoading}
              >
                {opt === '1d' && '1 day'}
                {opt === '7d' && '7 days'}
                {opt === '30d' && '30 days'}
              </button>
            ))}
          </div>
        </div>

        <div className="submit-container">
          <button
            type="submit"
            className="submit-btn"
            disabled={isLoading || !url.trim()}
          >
            {isLoading ? (
              <span className="loading-state">
                <span className="spinner" aria-hidden="true" />
                Shortening...
              </span>
            ) : (
              'Shorten'
            )}
          </button>
        </div>
      </div>
    </form>
  );
};
