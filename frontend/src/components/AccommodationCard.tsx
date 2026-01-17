import React from 'react';
import './AccommodationCard.css';

interface AccommodationCardProps {
    title: string;
    description: string;
    location: string;
    priceStart: string;
}

const AccommodationCard: React.FC<AccommodationCardProps> = ({ title, description, location, priceStart }) => {
    return (
        <div className="card">
            <div className="card-wireframe">
                <div className="wireframe-placeholder">Visual Placeholder</div>
            </div>
            <div className="card-content">
                <h3 className="card-title">{title}</h3>
                <p className="card-location"><span className="icon">📍</span> {location}</p>
                <p className="card-desc">{description}</p>
                <div className="card-footer">
                    <span className="card-price">From {priceStart} / night</span>
                    <button className="card-button">View Details</button>
                </div>
            </div>
        </div>
    );
};

export default AccommodationCard;
