import { createBrowserRouter, type RouteObject } from 'react-router';
import SiteLayout from './layouts/SiteLayout';
import HomePage from './pages/HomePage';
import NotFoundPage from './pages/NotFoundPage';

// Route table, like the `routes` array passed to Vue Router's createRouter.
// Child routes render inside SiteLayout's <Outlet />.
export const routes: RouteObject[] = [
  {
    path: '/',
    element: <SiteLayout />,
    children: [
      { index: true, element: <HomePage /> },
      { path: '*', element: <NotFoundPage /> },
    ],
  },
];

export const router = createBrowserRouter(routes);
