export interface WayFeature {
	id: number;
	tags: Record<string, string>;
}

export interface Review {
	wayID: number;
	username: string;
	rating: number;
	comment: string;
	createdAt: string;
}
