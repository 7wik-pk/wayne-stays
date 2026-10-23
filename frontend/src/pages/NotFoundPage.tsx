import React from 'react';
import { Link } from 'react-router';
import './NotFoundPage.css';

const NotFoundPage: React.FC = () => {
    return (
        <main className="container not-found">
            <p className="not-found-code">404</p>
            <h1 className="not-found-title">This Wing Is Closed</h1>
            <p className="not-found-text">
                The page you are looking for does not exist, or has been moved elsewhere on the estate.
            </p>
            <Link to="/" className="not-found-link">Return to the Lobby</Link>
        </main>
    );
};

export default NotFoundPage;
