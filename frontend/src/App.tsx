import { useState, useEffect } from 'react';
import Navbar from './components/Navbar';
import Hero from './components/Hero';
import AccommodationCard from './components/AccommodationCard';
import MerchShop from './components/MerchShop';

// API configuration
const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';

interface Accommodation {
  id: string;
  slug: string;
  title: string;
  description: string;
  location: string;
  priceStart: number;
  capacity: number;
  category: string;
  amenities: string[];
  images: {
    hero: string;
    gallery: string[];
  };
}

function App() {
  const [accommodations, setAccommodations] = useState<Accommodation[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchAccommodations = async () => {
      try {
        const response = await fetch(`${API_URL}/accommodations`);

        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        setAccommodations(data.data || []);
      } catch (err) {
        console.error('Failed to fetch accommodations:', err);
        setError(err instanceof Error ? err.message : 'Failed to load accommodations');
      } finally {
        setLoading(false);
      }
    };

    fetchAccommodations();
  }, []);

  return (
    <div className="app">
      <Navbar />
      <Hero />

      <main id="accommodations" className="container" style={{ padding: '4rem 1rem' }}>
        <div className="section-header text-center mb-lg">
          <h2 style={{ fontSize: '2.5rem', marginBottom: '1rem' }}>Curated Sanctuaries</h2>
          <p style={{ color: 'var(--color-text-muted)', maxWidth: '600px', margin: '0 auto' }}>
            Each destination is hand-picked to offer privacy, security, and the uncompromising standards of the Wayne family.
          </p>
        </div>

        {loading && (
          <div style={{ textAlign: 'center', padding: '3rem' }}>
            <p style={{ color: 'var(--color-primary)' }}>Loading accommodations...</p>
          </div>
        )}

        {error && (
          <div style={{ textAlign: 'center', padding: '3rem' }}>
            <p style={{ color: '#e74c3c' }}>Error: {error}</p>
            <p style={{ color: 'var(--color-text-muted)', marginTop: '1rem' }}>
              Please ensure the backend API is running on port 8080
            </p>
          </div>
        )}

        {!loading && !error && (
          <div style={{
            display: 'grid',
            gridTemplateColumns: 'repeat(auto-fit, minmax(350px, 1fr))',
            gap: '2rem'
          }}>
            {accommodations.map((acc) => (
              <AccommodationCard
                key={acc.id}
                title={acc.title}
                description={acc.description}
                location={acc.location}
                priceStart={`$${acc.priceStart.toLocaleString()}`}
                image={acc.images.hero}
              />
            ))}
          </div>
        )}
      </main>

      <MerchShop />

      <footer style={{
        backgroundColor: 'var(--color-bg-secondary)',
        padding: '4rem 1rem',
        marginTop: '4rem',
        borderTop: '1px solid rgba(212, 175, 55, 0.2)'
      }}>
        <div className="container text-center">
          <h3 style={{ fontSize: '1.5rem', marginBottom: '2rem' }}>Wayne Stays</h3>
          <p style={{ color: 'var(--color-text-muted)', marginBottom: '2rem' }}>
            A subsidiary of Wayne Industries.
          </p>
          <div style={{ display: 'flex', justifyContent: 'center', gap: '2rem' }}>
            <a href="#">Privacy</a>
            <a href="#">Terms</a>
            <a href="#">Concierge</a>
          </div>
          <p style={{ marginTop: '2rem', fontSize: '0.8rem', color: 'var(--color-text-muted)' }}>
            &copy; {new Date().getFullYear()} Wayne Enterprises. All Rights Reserved. Used under license.
          </p>
        </div>
      </footer>
    </div>
  );
}

export default App;
