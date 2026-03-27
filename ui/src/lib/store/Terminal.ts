import { writable } from 'svelte/store';

export const terminals = writable<Terminal[]>([]);
export type Terminal = {
	id: string;
	tid: string;
	serial_number: string;
	refund_allowed: boolean;
	product_name: string;
	activation_code: string;
};
