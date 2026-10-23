import React, { useState, useEffect } from 'react';
import { Link } from 'react-router';
import './Navbar.css';

const Navbar: React.FC = () => {
    const [scrolled, setScrolled] = useState(false);

    useEffect(() => {
        const handleScroll = () => {
            const offset = window.scrollY;
            if (offset > 50) {
                setScrolled(true);
            } else {
                setScrolled(false);
            }
        };

        window.addEventListener('scroll', handleScroll);
        return () => {
            window.removeEventListener('scroll', handleScroll);
        };
    }, []);

    return (
        <nav className={`navbar ${scrolled ? 'scrolled' : ''}`}>
            <div className="container navbar-content">
                <Link to="/" className="logo">
                    <img src="/logo.png?v=2" alt="Wayne Stays" />
                </Link>
                {/* Section links point at the home page so they also work from other routes */}
                <ul className="nav-links">
                    <li><a href="/#accommodations">Accommodations</a></li>
                    <li><a href="/#shop">Shop</a></li>
                </ul>
            </div>
        </nav>
    );
};

export default Navbar;
