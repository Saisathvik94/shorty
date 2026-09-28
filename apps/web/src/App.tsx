import React, { useState, useEffect } from 'react';
import { UrlForm } from './components/UrlForm';
import { UrlResult } from './components/UrlResult';
import { ErrorMessage } from './components/ErrorMessage';
import Grainient from './components/Gradient';
import { createShortUrl, getFullShortUrl, getBackendRedirectUrl } from './api/urls';
import { ShortenedUrlResult } from './types/url';

export const App: React.FC = () => {
  const [isLoading, setIsLoading] = useState(false);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);
  const [result, setResult] = useState<ShortenedUrlResult | null>(null);
  const [isRedirecting, setIsRedirecting] = useState(false);

  // Check if current URL pathname is a short code (e.g. /nQYAYb)
  useEffect(() => {
    const pathname = window.location.pathname.replace(/^\/+/, '').trim();
    if (pathname && pathname !== 'index.html') {
      setIsRedirecting(true);
      const targetUrl = getBackendRedirectUrl(pathname);
      window.location.replace(targetUrl);
    }
  }, []);

  const handleSubmit = async (url: string, expiresAt: string) => {
    setIsLoading(true);
    setErrorMessage(null);

    try {
      const response = await createShortUrl({
        url,
        expires_at: expiresAt,
      });

      const fullShortUrl = getFullShortUrl(response.short_code);
      setResult({
        shortCode: response.short_code,
        fullShortUrl,
        originalUrl: url,
      });
    } catch (err: unknown) {
      if (err instanceof Error) {
        setErrorMessage(err.message);
      } else {
        setErrorMessage('Something went wrong. Please try again.');
      }
    } finally {
      setIsLoading(false);
    }
  };

  const handleReset = () => {
    setResult(null);
    setErrorMessage(null);
  };

  if (isRedirecting) {
    return (
      <div className="redirect-screen">
        <div className="spinner-large" aria-hidden="true" />
        <p className="redirect-text">Redirecting to destination...</p>
      </div>
    );
  }

  return (
    <div className="app-layout">
      <div className="background-grainient">
        <Grainient
          color1="#be185d"
          color2="#000000"
          color3="#52525b"
          timeSpeed={0.15}
          colorBalance={0.0}
          warpStrength={0.5}
          warpFrequency={3.0}
          warpSpeed={1.0}
          warpAmplitude={25.0}
          blendAngle={90.0}
          blendSoftness={0.25}
          rotationAmount={0.0}
          noiseScale={1.5}
          grainAmount={0.03}
          grainScale={2.0}
          grainAnimated={false}
          contrast={0.9}
          gamma={1.1}
          saturation={1}
          centerX={0.0}
          centerY={0.0}
          zoom={1.1}
          lightMode
        />
      </div>

      <main className="main-content">
        <div className="hero-section">
          <h1 className="hero-title">Shorten your links.</h1>
          <p className="hero-subtitle">Simple, fast URL shortening.</p>
        </div>

        <div className="card-wrapper">
          <ErrorMessage
            message={errorMessage || ''}
            onDismiss={() => setErrorMessage(null)}
          />

          {!result ? (
            <UrlForm
              isLoading={isLoading}
              onSubmit={handleSubmit}
              onError={(msg) => setErrorMessage(msg)}
            />
          ) : (
            <UrlResult result={result} onReset={handleReset} />
          )}
        </div>
      </main>

      <footer className="footer">
        <p className="footer-text">Shorty &bull; Focused link utility</p>
      </footer>
    </div>
  );
};

export default App;
