<script lang="ts">
	export let form;
</script>

<div class="min-h-screen flex flex-col items-center justify-center bg-gray-100 px-4">
	<h1 class="text-4xl font-extrabold mb-10">Login</h1>

	<form method="POST" class="bg-white p-8 rounded-lg shadow-md w-full max-w-sm">
		<input
			type="text"
			name="username"
			placeholder="Username"
			required
			class="w-full mb-4 px-4 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
		/>

		<input
			type="password"
			name="password"
			placeholder="Password"
			required
			class="w-full mb-6 px-4 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
		/>

		<button
			type="submit"
			class="w-full bg-blue-600 text-white py-2 rounded-md hover:bg-blue-700 transition"
		>
			Log In
		</button>

		{#if form?.error}
			<p class="mt-4 text-center text-red-600">{form.error}</p>
		{/if}
	</form>

	<p class="mt-4 text-sm text-center">
		Don't have an account?
		<a href="/signup" class="text-blue-600 hover:underline">Sign up</a>
	</p>
</div>
