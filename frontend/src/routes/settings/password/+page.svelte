<script lang="ts">
	import type { ActionData } from './$types';

	export let form: ActionData | null;

	let currentPassword = '';
	let newPassword = '';
	let confirmPassword = '';
	let errorMsg = '';
	let successMsg = '';

	$: if (form?.success) {
		successMsg = 'Password changed successfully';
		errorMsg = '';
		currentPassword = '';
		newPassword = '';
		confirmPassword = '';
	} else if (form?.error) {
		errorMsg = form.error;
		successMsg = '';
	}
</script>

<div class="min-h-screen flex flex-col items-center justify-center bg-gray-100 px-4">
	<h1 class="text-4xl font-extrabold mb-10">Change Password</h1>

	<div class="bg-white p-8 rounded-lg shadow-md w-full max-w-md">
		<form method="POST">
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
			>
				Change Password
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
