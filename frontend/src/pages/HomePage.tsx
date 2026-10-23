import React from 'react';
import Hero from '../components/Hero';
import AccommodationCard from '../components/AccommodationCard';
import MerchShop from '../components/MerchShop';
import { api } from '../api/client';
import { useApi } from '../hooks/useApi';
import './HomePage.css';

const HomePage: React.FC = () => {
    const { data: accommodations, loading, error } = useApi(api.listAccommodations);

    return (
        <>
            <Hero />

            <main id="accommodations" className="container home-accommodations">
                <div className="section-header text-center mb-lg">
                    <h2 className="section-title">Curated Sanctuaries</h2>
                    <p className="home-intro">
                        Each destination is hand-picked to offer privacy, security, and the uncompromising standards of the Wayne family.
                    </p>
                </div>

                {loading && (
                    <div className="status-message">
                        <p>Loading accommodations...</p>
                    </div>
                )}

                {error && (
                    <div className="status-message status-message--error">
                        <p>Error: {error}</p>
                        <p className="hint">Please ensure the backend API is running on port 8080</p>
                    </div>
                )}

                {accommodations && (
                    <div className="accommodation-grid">
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
        </>
    );
};

export default HomePage;
