import React, { useState } from 'react';
import MerchItem from './MerchItem';
import './MerchShop.css';

const MerchShop: React.FC = () => {
    const [filter, setFilter] = useState<string>('all');

    const products = [
        { id: 1, name: 'GCPD Tactical Fleece', price: 85.00, category: 'Apparel', isNew: true },
        { id: 2, name: 'Wayne Ent. Journal', price: 45.00, category: 'Accessories' },
        { id: 3, name: 'Iceberg Lounge Coasters', price: 35.00, category: 'Home' },
        { id: 4, name: 'Vote Dent Pin', price: 15.00, category: 'Collectibles' },
        { id: 5, name: 'Martha\'s Pearls (Replica)', price: 450.00, category: 'Collectibles' },
        { id: 6, name: 'Arkham Asylum Keytag', price: 25.00, category: 'Accessories' },
    ];

    const filteredProducts = filter === 'all'
        ? products
        : products.filter(p => p.category.toLowerCase() === filter);

    return (
        <section id="shop" className="merch-shop">
            <div className="container">
                <div className="section-header text-center mb-lg">
                    <h2 className="section-title">The Collection</h2>
                    <p className="section-subtitle">Exlcusive memorabilia from Gotham City.</p>
                </div>

                <div className="merch-filters">
                    {['all', 'apparel', 'collectibles', 'accessories'].map(cat => (
                        <button
                            key={cat}
                            className={`filter-btn ${filter === cat ? 'active' : ''}`}
                            onClick={() => setFilter(cat)}
                        >
                            {cat}
                        </button>
                    ))}
                </div>

                <div className="merch-grid">
                    {filteredProducts.map(product => (
                        <MerchItem
                            key={product.id}
                            name={product.name}
                            price={product.price}
                            category={product.category}
                            isNew={product.isNew}
                        />
                    ))}
                </div>
            </div>
        </section>
    );
};

export default MerchShop;
