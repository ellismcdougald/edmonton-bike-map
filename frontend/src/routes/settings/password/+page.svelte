<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';

	const apiUrl = import.meta.env.VITE_API_URL;

	let currentPassword: string = '';
	let newPassword: string = '';
	let confirmPassword: string = '';
	let errorMsg: string = '';
	let successMsg: string = '';
	let isSubmitting: boolean = false;

	onMount(() => {
		const token = localStorage.getItem('token');
		if (!token) {
			goto('/login');
		}
	});

	async function handleChangePassword(event: SubmitEvent) {
		event.preventDefault();
		isSubmitting = true;
		errorMsg = '';
		successMsg = '';

		if (newPassword !== confirmPassword) {
			errorMsg = 'New passwords do not match';
			isSubmitting = false;
			return;
		}

		const token = localStorage.getItem('token');
		if (!token) {
			errorMsg = 'You must be logged in to change your password';
			isSubmitting = false;
			return;
		}

		const res = await fetch(`${apiUrl}/api/change-password`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
				Authorization: `Bearer ${token}`
			},
			body: JSON.stringify({
				currentPassword,
				newPassword
			})
		});

		if (res.ok) {
			successMsg = 'Password changed successfully';
			currentPassword = '';
			newPassword = '';
			confirmPassword = '';
		} else {
			const text = await res.text();
			errorMsg = text || 'Failed to change password';
		}

		isSubmitting = false;
	}
</script>

<div class="min-h-screen flex flex-col items-center justify-center bg-gray-100 px-4">
	<h1 class="text-4xl font-extrabold mb-10">Change Password</h1>

	<div class="bg-white p-8 rounded-lg shadow-md w-full max-w-md">
		<form on:submit={handleChangePassword}>
			<input
				type="password"
				name="currentPassword"
				bind:value={currentPassword}
				placeholder="Current Password"
				required
				class="w-full mb-4 px-4 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
			/>

			<input
				type="password"
				name="newPassword"
				bind:value={newPassword}
				placeholder="New Password"
				required
				class="w-full mb-4 px-4 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
			/>

			<input
				type="password"
				name="confirmPassword"
				bind:value={confirmPassword}
				placeholder="Confirm New Password"
				required
				class="w-full mb-6 px-4 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
			/>

			<button
				id="submitButton"
				type="submit"
				class="w-full bg-blue-600 text-white py-2 rounded-md hover:bg-blue-700 transition"
				disabled={isSubmitting}
			>
				{isSubmitting ? 'Changing password...' : 'Change Password'}
			</button>

			{#if successMsg}
				<p class="mt-4 text-center text-green-600">{successMsg}</p>
			{/if}

			{#if errorMsg}
				<p class="mt-4 text-center text-red-600">{errorMsg}</p>
			{/if}
		</form>

		<div class="mt-6 text-center">
			<a href="/settings" class="text-blue-600 hover:underline">Back to Settings</a>
		</div>
	</div>
</div>
