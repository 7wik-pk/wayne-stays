import React, { useState } from 'react';
import MerchItem from './MerchItem';
import './MerchShop.css';

const MerchShop: React.FC = () => {
    const [filter, setFilter] = useState<string>('all');

    const products = [
        { id: 1, name: 'Wayne Enterprises Leather Portfolio', price: 320.00, category: 'Accessories', isNew: true },
        { id: 2, name: 'Wayne Manor Heritage Cashmere Scarf', price: 450.00, category: 'Apparel', isNew: true },
        { id: 3, name: 'Martha\'s Pearls (Replica)', price: 450.00, category: 'Collectibles' },
        { id: 4, name: 'Wayne Manor: Architectural Legacy', price: 125.00, category: 'Home' },
        { id: 5, name: 'Martha Wayne Memorial Rose Set', price: 95.00, category: 'Home' },
        { id: 6, name: 'Wayne Foundation Gala Program', price: 350.00, category: 'Collectibles' },
        { id: 7, name: 'Kane County Reserve Wine', price: 420.00, category: 'Lifestyle' },
        { id: 8, name: 'Gotham Skyline Silk Tie', price: 185.00, category: 'Apparel' },
        { id: 9, name: 'Wayne Stays Luxury Linen Set', price: 380.00, category: 'Home' },
        { id: 10, name: 'Gotham Opera House Poster', price: 165.00, category: 'Collectibles' },
    ];

    const filteredProducts = filter === 'all'
        ? products
        : products.filter(p => p.category.toLowerCase() === filter);

    return (
        <section id="shop" className="merch-shop">
            <div className="container">
                <div className="section-header text-center mb-lg">
                    <h2 className="section-title">The Collection</h2>
                    <p className="section-subtitle">Exclusive memorabilia from Gotham City.</p>
                </div>

                <div className="merch-filters">
                    {['all', 'apparel', 'accessories', 'collectibles', 'home', 'lifestyle'].map(cat => (
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
        </section >
    );
};

export default MerchShop;
