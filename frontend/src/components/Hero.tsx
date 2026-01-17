import React from 'react';
import './Hero.css';

const Hero: React.FC = () => {
    return (
        <section className="hero">
            <div className="hero-wireframe"></div>
            <div className="hero-overlay"></div>
            <div className="container hero-content">
                <h1 className="hero-title">Experience the Exceptional</h1>
                <p className="hero-subtitle">
                    Unparalleled Luxury at Wayne Manor, Gotham City, and Beyond
                </p>
                <button className="cta-button">Reserve Your Stay</button>
            </div>
        </section>
    );
};

export default Hero;
