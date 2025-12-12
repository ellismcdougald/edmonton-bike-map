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

export interface WayFeatureProperties {
	id?: number;
	[key: string]: string | number | boolean | undefined;
}

export type WayFeatureGeoJSON = GeoJSON.Feature<GeoJSON.Geometry, WayFeatureProperties>;

export type WayFeatureCollection = GeoJSON.FeatureCollection<
	GeoJSON.Geometry,
	WayFeatureProperties
>;

export type FeatureCollection = GeoJSON.FeatureCollection<GeoJSON.Geometry, WayFeatureProperties>;
