import React, { useState, useEffect } from 'react';
import MerchItem from './MerchItem';
import './MerchShop.css';

interface Product {
    id: string;
    name: string;
    price: number;
    category: string;
    image_url: string;
    is_new: boolean;
}

const MerchShop: React.FC = () => {
    const [filter, setFilter] = useState<string>('all');
    const [products, setProducts] = useState<Product[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchProducts = async () => {
            const apiBase = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';
            try {
                const response = await fetch(`${apiBase}/merchandise`);
                if (!response.ok) throw new Error('Failed to fetch merchandise');
                const data = await response.json();
                setProducts(data.data);
            } catch (err) {
                console.error(err);
                setError('Unable to load full collection.');
            } finally {
                setLoading(false);
            }
        };

        fetchProducts();
    }, []);

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

                {loading ? (
                    <div style={{ textAlign: 'center', padding: '2rem', color: 'var(--color-gold-primary)' }}>
                        Loading collection...
                    </div>
                ) : error ? (
                    <div style={{ textAlign: 'center', padding: '2rem', color: 'red' }}>
                        {error}
                    </div>
                ) : (
                    <div className="merch-grid">
                        {filteredProducts.map(product => (
                            <MerchItem
                                key={product.id}
                                name={product.name}
                                price={product.price}
                                category={product.category}
                                image={product.image_url}
                                isNew={product.is_new}
                            />
                        ))}
                    </div>
                )}
            </div>
        </section >
    );
};

export default MerchShop;
