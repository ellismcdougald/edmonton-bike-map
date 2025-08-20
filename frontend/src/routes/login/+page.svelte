<script lang="ts">
	import { goto } from '$app/navigation';

	let username: string = '';
	let password: string = '';
	let errorMsg: string = '';
	let isSubmitting: boolean = false;

	async function handleLogin(event: SubmitEvent) {
		event.preventDefault();
		isSubmitting = true;
		errorMsg = '';

		const res = await fetch('http://localhost:8080/api/login', {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ username, password })
		});

		if (res.ok) {
			const data = await res.json();
			localStorage.setItem('token', data.token);
			goto('/');
		} else {
			const text = await res.text();
			errorMsg = text ?? 'Login failed.';
		}

		isSubmitting = false;
	}
</script>

<div class="min-h-screen flex flex-col items-center justify-center bg-gray-100 px-4">
	<h1 class="text-4xl font-extrabold mb-10">Login</h1>

	<form on:submit={handleLogin} class="bg-white p-8 rounded-lg shadow-md w-full max-w-sm">
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
			{isSubmitting ? 'Logging in...' : 'Log In'}
		</button>

		{#if errorMsg}
			<p class="mt-4 text-center text-red-600">{errorMsg}</p>
		{/if}
	</form>
	<p class="mt-4 text-sm text-center">
		Don't have an account?
		<a href="/signup" class="text-blue-600 hover:underline">Sign up</a>
	</p>
</div>
