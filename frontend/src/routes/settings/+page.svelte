<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';

	const apiUrl = import.meta.env.VITE_API_URL;

	let cyclingSpeed: number | null = null;
	let isSubmitting = false;
	let errorMsg = '';
	let successMsg = '';

	onMount(async () => {
		const token = localStorage.getItem('token');
		if (!token) {
			goto('/login');
			return;
		}

		// Try to prefill user's cycling speed
		try {
			const res = await fetch(`${apiUrl}/api/user/settings`, {
				method: 'GET',
				headers: { Authorization: `Bearer ${token}` }
			});
			if (res.ok) {
				const data = await res.json();
				cyclingSpeed = data.cyclingSpeed ?? 15;
			} else {
				cyclingSpeed = 15;
			}
		} catch (_e) {
			cyclingSpeed = 15;
		}
	});

	async function saveSpeed(e: Event) {
		e.preventDefault();
		errorMsg = '';
		successMsg = '';
		if (!cyclingSpeed || cyclingSpeed <= 0 || cyclingSpeed > 80) {
			errorMsg = 'Please enter a valid cycling speed (1-80 km/h)';
			return;
		}

		const token = localStorage.getItem('token');
		if (!token) {
			errorMsg = 'You must be logged in';
			return;
		}

		isSubmitting = true;
		const res = await fetch(`${apiUrl}/api/user/settings`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
				Authorization: `Bearer ${token}`
			},
			body: JSON.stringify({ cyclingSpeed })
		});

		if (res.ok) {
			successMsg = 'Settings saved';
		} else {
			const text = await res.text();
			errorMsg = text || 'Failed to save settings';
		}

		isSubmitting = false;
	}
</script>

<div class="min-h-screen flex flex-col items-center justify-center bg-gray-100 px-4">
	<h1 class="text-4xl font-extrabold mb-10">Settings</h1>

	<div class="bg-white p-8 rounded-lg shadow-md w-full max-w-md">
		<h2 class="text-2xl font-bold mb-6">Cycling Preferences</h2>

		<form on:submit|preventDefault={saveSpeed}>
			<label for="cyclingSpeed" class="block mb-2 font-medium">Preferred Cycling Speed (km/h)</label
			>
			<input
				type="number"
				min="1"
				max="80"
				id="cyclingSpeed"
				bind:value={cyclingSpeed}
				class="w-full mb-4 px-4 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
			/>

			<button
				type="submit"
				class="w-full bg-blue-600 text-white py-2 rounded-md hover:bg-blue-700 transition"
				disabled={isSubmitting}
			>
				{isSubmitting ? 'Saving...' : 'Save Preferences'}
			</button>

			{#if successMsg}
				<p class="mt-4 text-center text-green-600">{successMsg}</p>
			{/if}

			{#if errorMsg}
				<p class="mt-4 text-center text-red-600">{errorMsg}</p>
			{/if}
		</form>

		<div class="mt-6 text-center">
			<a href="/settings/password" class="text-blue-600 hover:underline">Change Password</a>
		</div>
		<div class="mt-2 text-center">
			<a href="/" class="text-gray-600 hover:underline">Back to Map</a>
		</div>
	</div>
</div>
