import React from 'react';
import { Outlet } from 'react-router';
import Navbar from '../components/Navbar';
import Footer from '../components/Footer';

// Shared frame for every page; the matched child route renders in <Outlet />
const SiteLayout: React.FC = () => {
    return (
        <div className="app">
            <Navbar />
            <Outlet />
            <Footer />
        </div>
    );
};

export default SiteLayout;
