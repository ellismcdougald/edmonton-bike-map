/**
 * utils/route.ts
 *
 * Purpose:
 * Provides helper functions to interpret route metadata from OpenStreetMap tags
 * for display in the sidebar or map UI.
 *
 * Functions:
 * - determineRouteType(tags: Record<string, string>): string
 *     - Determines the road type based on `highway`, `cycleway`, and `foot` tags
 *     - Examples: 'Residential Street', 'Collector Road with Bike Lane', 'Shared Use Path', etc.
 *     - Returns 'Unknown' if no known classification matches
 * - determineBicycleRoute(tags: Record<string, string>): 'Yes' | 'No'
 *     - Determines if the route is part of the local bicycle network
 *     - Checks for `bicycle: designated` or `lcn: yes` tags
 *
 * Behavior:
 * - Pure functions; do not modify input
 * - Used by Sidebar and other components to display human-readable route info
 *
 * Notes:
 * - Logic in determineRouteType is specific to Edmonton and may need updates for other areas
 * - TODO: simplify and improve classification logic for common and edge cases
 */

// determineRouteType uses the tags for a route to determine its road type (arterial, collector, local, etc) and bike infastructure (if possible)
// TODO: simplify and improve logic, continue to ensure that all common types are classified
export function determineRouteType(tags: Record<string, string>): string {
	if (tags.highway === 'residential') return 'Residential Street';
	if (tags.highway === 'tertiary') {
		return tags.cycleway === 'lane' ? 'Collector Road with Bike Lane' : 'Collector Road';
	}
	if (tags.highway === 'path') {
		if (tags.bicycle === 'designated' && tags.foot === 'designated') return 'Shared Use Path';
		if (tags.bicycle === 'designated') return 'Bike Path';
		if (tags.foot === 'designated') return 'Footpath';
	}
	if (tags.highway === 'footway')
		return tags.bicycle === 'yes' || tags.bicycle === 'designated' ? 'Shared Use Path' : 'Footpath';
	if (tags.highway === 'cycleway')
		return tags.foot === 'designated' || tags.foot === 'yes' ? 'Shared Use Path' : 'Cycleway';
	if (tags.highway === 'secondary') {
		if (
			['share_busway'].includes(tags.cycleway) ||
			['share_busway'].includes(tags['cycleway:left']) ||
			['share_busway'].includes(tags['cycleway:right'])
		)
			return 'Arterial Road With Bus-Bike Lane';
		return 'Arterial Road';
	}
	if (tags.highway === 'primary') return 'Major Arterial Road';
	if (tags.highway === 'unclassified') {
		if (
			['separate'].includes(tags.cycleway) ||
			['separate'].includes(tags['cycleway:left']) ||
			['separate'].includes(tags['cycleway:right'])
		)
			return 'Local Road with Separated Bike Lane';
		return 'Local Road';
	}
	return 'Unknown';
}

// determineBicycleRoute uses the tags to determine if a route is part of the local bicycle network
export function determineBicycleRoute(tags: Record<string, string>): 'Yes' | 'No' {
	return tags.bicycle === 'designated' || tags.lcn === 'yes' ? 'Yes' : 'No';
}
