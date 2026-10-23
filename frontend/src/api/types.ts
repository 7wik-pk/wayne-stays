// Response shapes returned by the Go backend (see API.md)

export interface Accommodation {
  id: string;
  slug: string;
  title: string;
  description: string;
  location: string;
  priceStart: number; // dollars
  capacity: number;
  category: string;
  amenities: string[];
  images: {
    hero: string;
    gallery: string[];
  };
}

export interface Merchandise {
  id: string;
  name: string;
  description?: string;
  price: number; // dollars
  category: string;
  imageUrl: string;
  isNew: boolean;
}
