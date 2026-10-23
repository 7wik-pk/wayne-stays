import React from 'react';
import './Footer.css';

const Footer: React.FC = () => {
    return (
        <footer className="site-footer">
            <div className="container text-center">
                <h3 className="site-footer-title">Wayne Stays</h3>
                <p className="site-footer-tagline">
                    A subsidiary of Wayne Enterprises.
                </p>
                <div className="site-footer-links">
                    <a href="#">Privacy</a>
                    <a href="#">Terms</a>
                    <a href="#">Concierge</a>
                </div>
                <p className="site-footer-copyright">
                    &copy; {new Date().getFullYear()} Wayne Enterprises. All Rights Reserved. Used under license.
                </p>
            </div>
        </footer>
    );
};

export default Footer;
