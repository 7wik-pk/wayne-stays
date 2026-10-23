import React from 'react';
import './MerchItem.css';

interface MerchItemProps {
    name: string;
    price: number;
    image: string;
    isNew?: boolean;
}

const MerchItem: React.FC<MerchItemProps> = ({ name, price, image, isNew }) => {
    const handleAddToCart = () => {
        // Placeholder for cart logic
        alert(`Added ${name} to cart`);
    };

    return (
        <div className="merch-item">
            <div className="merch-image-container">
                <img src={image} alt={name} className="merch-image" />
            </div>
            <div className="merch-details">
                {isNew && <span className="merch-badge">New Arrival</span>}
                <h4 className="merch-name">{name}</h4>
                <div className="merch-footer">
                    <span className="merch-price">${price.toFixed(2)}</span>
                    <button onClick={handleAddToCart} className="merch-add-btn">
                        +
                    </button>
                </div>
            </div>
        </div>
    );
};

export default MerchItem;
