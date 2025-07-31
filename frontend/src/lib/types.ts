export interface WayFeature {
	id: number;
	tags: Record<string, string>;
}

export interface Review {
	wayID: number;
	rating: number;
	reviewText: string;
	createdAt: string;
}
