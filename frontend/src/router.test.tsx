import { describe, expect, it } from 'vitest';
import { render, screen } from '@testing-library/react';
import { createMemoryRouter, RouterProvider } from 'react-router';
import { routes } from './router';
import { mockFetch } from './test/mockFetch';

// A memory router runs the real route table without touching the browser URL
function renderAt(path: string) {
  const router = createMemoryRouter(routes, { initialEntries: [path] });
  render(<RouterProvider router={router} />);
}

describe('routes', () => {
  it('renders the home page at /', async () => {
    mockFetch({
      '/accommodations': [{
        id: '1', slug: 'clocktower-penthouse', title: 'The Clocktower Penthouse',
        description: '', location: 'Burnley, Old Gotham', priceStart: 1400, capacity: 2,
        category: 'historic', amenities: [], images: { hero: '/hero.jpg', gallery: [] },
      }],
      '/merchandise': [],
    });
    renderAt('/');

    expect(await screen.findByText('The Clocktower Penthouse')).toBeInTheDocument();
    expect(screen.getByText('From $1,400 / night')).toBeInTheDocument();
  });

  it('renders the not found page for unknown paths, inside the site layout', () => {
    renderAt('/batcave');

    expect(screen.getByRole('heading', { name: 'This Wing Is Closed' })).toBeInTheDocument();
    expect(screen.getByRole('link', { name: 'Return to the Lobby' })).toHaveAttribute('href', '/');
    expect(screen.getByAltText('Wayne Stays')).toBeInTheDocument(); // navbar logo
  });
});
