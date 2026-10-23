import React, { useState } from 'react';
import MerchItem from './MerchItem';
import { api } from '../api/client';
import { useApi } from '../hooks/useApi';
import './MerchShop.css';

const CATEGORIES = ['all', 'apparel', 'accessories', 'collectibles', 'home', 'lifestyle'];

const MerchShop: React.FC = () => {
    const [filter, setFilter] = useState<string>('all');
    const { data: products, loading, error } = useApi(api.listMerchandise);

    const filteredProducts = filter === 'all'
        ? products ?? []
        : (products ?? []).filter(p => p.category === filter);

    return (
        <section id="shop" className="merch-shop">
            <div className="container">
                <div className="section-header text-center mb-lg">
                    <h2 className="section-title">The Collection</h2>
                    <p className="section-subtitle">Exclusive memorabilia from Gotham City.</p>
                </div>

                <div className="merch-filters">
                    {CATEGORIES.map(cat => (
                        <button
                            key={cat}
                            className={`filter-btn ${filter === cat ? 'active' : ''}`}
                            onClick={() => setFilter(cat)}
                        >
                            {cat}
                        </button>
                    ))}
                </div>

                {loading ? (
                    <div className="status-message">
                        Loading collection...
                    </div>
                ) : error ? (
                    <div className="status-message status-message--error">
                        Unable to load full collection.
                    </div>
                ) : (
                    <div className="merch-grid">
                        {filteredProducts.map(product => (
                            <MerchItem
                                key={product.id}
                                name={product.name}
                                price={product.price}
                                image={product.imageUrl}
                                isNew={product.isNew}
                            />
                        ))}
                    </div>
                )}
            </div>
        </section>
    );
};

export default MerchShop;
