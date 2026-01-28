import Navbar from './components/Navbar';
import Hero from './components/Hero';
import AccommodationCard from './components/AccommodationCard';
import MerchShop from './components/MerchShop';

function App() {
  const accommodations = [
    {
      title: "Wayne Tower Sky Suites",
      description: "Executive suites in the iconic Wayne Tower. Direct access to Wayne Enterprises board rooms, private viewing galleries, and the legendary penthouse lounge with 360° views from Gotham's tallest building.",
      location: "Financial District, Gotham City",
      priceStart: "$1,850"
    },
    {
      title: "Miagani Island Sanctuary",
      description: "Secluded island retreat accessible only by private ferry. Complete privacy, lighthouse views, and untouched natural beauty. An entire island experience with dedicated staff.",
      location: "Miagani Island, Gotham Harbor",
      priceStart: "$3,200"
    },
    {
      title: "The Clocktower Penthouse",
      description: "Luxury apartment within Gotham's famous clocktower. Vaulted ceilings, original clockwork machinery, and stunning Gothic architecture. Live inside a functioning historic landmark.",
      location: "Burnley, Old Gotham",
      priceStart: "$1,400"
    },
    {
      title: "Arkham District Gentleman's Club",
      description: "Restored Victorian gentleman's club with leather libraries, billiard rooms, and private dining. A step back to Gotham's gilded age with members-only access to exclusive lounges.",
      location: "Arkham, North Gotham",
      priceStart: "$950"
    },
    {
      title: "Kane County Vineyard Estate",
      description: "Rolling vineyard estate producing Gotham's finest wines. Private wine cellar, tasting room, and farm-to-table cuisine by Michelin-starred chefs with exclusive vintage tastings.",
      location: "Kane County, Gotham Outskirts",
      priceStart: "$2,100"
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
