export interface WayFeature {
	id: number;
	tags: Record<string, string>;
}

export interface Review {
	wayId: number;
	username: string;
	rating: number;
	comment: string;
	createdAt: string;
}
