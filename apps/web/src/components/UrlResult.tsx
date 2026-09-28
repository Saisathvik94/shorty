import React, { useState } from 'react';
import { ShortenedUrlResult } from '../types/url';

interface UrlResultProps {
  result: ShortenedUrlResult;
  onReset: () => void;
}

export const UrlResult: React.FC<UrlResultProps> = ({ result, onReset }) => {
  const [copied, setCopied] = useState(false);

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(result.fullShortUrl);
      setCopied(true);
      setTimeout(() => {
        setCopied(false);
      }, 2000);
    } catch {
      // Fallback for older browsers / clipboard permission issues
      const textarea = document.createElement('textarea');
      textarea.value = result.fullShortUrl;
      textarea.style.position = 'fixed';
      textarea.style.opacity = '0';
      document.body.appendChild(textarea);
      textarea.focus();
      textarea.select();
      try {
        document.execCommand('copy');
        setCopied(true);
        setTimeout(() => setCopied(false), 2000);
      } catch {
        // ignore
      }
      document.body.removeChild(textarea);
    }
  };

  return (
    <div className="result-container" aria-live="polite">
      <div className="result-header">
        <span className="result-label">Your shortened URL</span>
      </div>

      <div className="result-card">
        <div className="result-url-display">
          <span className="result-url" title={result.fullShortUrl}>
            {result.fullShortUrl}
          </span>
        </div>

        <div className="result-actions">
          <button
            type="button"
            className={`action-btn copy-btn ${copied ? 'copied' : ''}`}
            onClick={handleCopy}
            aria-label={copied ? 'URL copied to clipboard' : 'Copy shortened URL'}
          >
            {copied ? 'Copied' : 'Copy'}
          </button>

          <a
            href={result.fullShortUrl}
            target="_blank"
            rel="noopener noreferrer"
            className="action-btn open-btn"
            aria-label="Open shortened URL in new tab"
          >
            Open
          </a>
        </div>
      </div>

      <div className="result-footer">
        <button
          type="button"
          className="create-another-btn"
          onClick={onReset}
        >
          Create another
        </button>
      </div>
    </div>
  );
};
