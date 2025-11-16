<!--
  src/routes/signup/+page.svelte

  Purpose:
  Signup page route for creating a new user account.

  State:
  - username (string): bound to username input
  - password (string): bound to password input
  - errorMsg (string): displays backend error messages
  - isSubmitting (boolean): disables form and shows submission state

  Behavior:
  - On form submit:
      - Sends POST request to /api/signup with username and password
      - If successful, navigates to the login page ('/login')
      - If failed, displays error message below the submit button
  - Submit button is disabled while signup request is in progress
  - Provides a link to the login page for existing users

  Notes:
  - Uses SvelteKit's `goto` for navigation
  - Relies on backend endpoint /api/signup
  - Input fields use native HTML validation for required values
-->

<script lang="ts">
	import { goto } from '$app/navigation';

	const apiUrl = import.meta.env.VITE_API_URL;

	let username = '';
	let password = '';
	let errorMsg = '';
	let isSubmitting = false;

	async function handleSignup(event: SubmitEvent) {
		event.preventDefault();
		isSubmitting = true;
		errorMsg = '';

		const res = await fetch(`${apiUrl}/api/signup`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ username, password })
		});

		if (res.ok) {
			goto('/login');
		} else {
			const text = await res.text();
			errorMsg = text ?? 'Signup failed';
		}

		isSubmitting = false;
	}
</script>

<div class="min-h-screen flex flex-col items-center justify-center bg-gray-100 px-4">
	<h1 class="text-4xl font-extrabold mb-10">Sign Up</h1>

	<form on:submit={handleSignup} class="bg-white p-8 rounded-lg shadow-md w-full max-w-sm">
		<input
			type="text"
			name="username"
			bind:value={username}
			placeholder="Username"
			required
			class="w-full mb-4 px-4 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
		/>

		<input
			type="password"
			name="password"
			bind:value={password}
			placeholder="Password"
			required
			class="w-full mb-6 px-4 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
		/>

		<button
			id="submitButton"
			type="submit"
			class="w-full bg-blue-600 text-white py-2 rounded-md hover:bg-blue-700 transition"
			disabled={isSubmitting}
		>
			{isSubmitting ? 'Signing up…' : 'Sign Up'}
		</button>

		{#if errorMsg}
			<p class="mt-4 text-center text-red-600">{errorMsg}</p>
		{/if}
	</form>
	<p class="mt-4 text-sm text-center">
		Already have an account?
		<a href="/login" class="text-blue-600 hover:underline">Log in</a>
	</p>
</div>
