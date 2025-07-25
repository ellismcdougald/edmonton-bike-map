import { writable } from 'svelte/store';
import type { WayFeature } from '../types';

export const selectedWay = writable<WayFeature | null>(null);
