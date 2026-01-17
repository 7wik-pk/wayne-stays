import Navbar from './components/Navbar';
import Hero from './components/Hero';
import AccommodationCard from './components/AccommodationCard';
import MerchShop from './components/MerchShop';

function App() {
  const accommodations = [
    {
      title: "The Manor West Wing",
      description: "Experience the history of the Wayne legacy. Stay in the meticulously restored West Wing of the Wayne Manor, featuring original Victorian architecture blended with modern amenities.",
      location: "Bristol Township, Gotham",
      priceStart: "$2,500"
    },
    {
      title: "Gotham Heights Penthouse",
      description: "A private sanctuary above the clouds. This high-rise complex offers panoramic views of the Gotham skyline, private elevator access, and helipad privileges.",
      location: "Diamond District, Gotham",
      priceStart: "$1,800"
    },
    {
      title: "Martha's Vineyard Retreat",
      description: "A serene escape dedicated to the memory of Martha Wayne. This coastal villa offers private beach access, organic farm-to-table dining, and complete seclusion.",
      location: "Martha's Vineyard, MA",
      priceStart: "$3,200"
    }
  ];

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

        <div style={{
          display: 'grid',
          gridTemplateColumns: 'repeat(auto-fit, minmax(350px, 1fr))',
          gap: '2rem'
        }}>
          {accommodations.map((acc, index) => (
            <AccommodationCard key={index} {...acc} />
          ))}
        </div>
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
