import { describe, expect, it, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import MerchShop from './MerchShop';
import { mockFetch } from '../test/mockFetch';
import type { Merchandise } from '../api/types';

const products: Merchandise[] = [
  { id: '1', name: 'Gotham Skyline Silk Tie', price: 185, category: 'apparel', imageUrl: '/tie.png', isNew: false },
  { id: '2', name: 'Kane County Reserve Wine', price: 420, category: 'lifestyle', imageUrl: '/wine.png', isNew: true },
];

describe('MerchShop', () => {
  it('lists products from the API and filters them by category', async () => {
    mockFetch({ '/merchandise': products });
    render(<MerchShop />);

    expect(await screen.findByText('Gotham Skyline Silk Tie')).toBeInTheDocument();
    expect(screen.getByText('Kane County Reserve Wine')).toBeInTheDocument();
    expect(screen.getByText('New Arrival')).toBeInTheDocument();

    await userEvent.click(screen.getByRole('button', { name: 'apparel' }));

    expect(screen.getByText('Gotham Skyline Silk Tie')).toBeInTheDocument();
    expect(screen.queryByText('Kane County Reserve Wine')).not.toBeInTheDocument();
  });

  it('shows an error message when the API fails', async () => {
    vi.spyOn(console, 'error').mockImplementation(() => {});
    mockFetch({}); // every request 404s
    render(<MerchShop />);

    expect(await screen.findByText('Unable to load full collection.')).toBeInTheDocument();
  });
});
