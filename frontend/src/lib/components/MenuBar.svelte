<script lang="ts">
	let { username }: { username: string | null } = $props();
</script>

<div class="flex items-center justify-between px-4 py-2 w-full">
	<div class="w-1/3"></div>
	<div class="w-1/3">
		<h1 class="text-3xl font-bold text-center">Edmonton Bike Map</h1>
	</div>
	<div class="w-1/3 flex justify-end">
		{#if username}
			<div class="relative inline-block group">
				<button
					id="usernameButton"
					class="text-black text-center px-4 py-1 hover:underline focus:outline-none"
				>
					{username}
				</button>
				<div
					class="absolute left-0 mt-2 w-full bg-white text-gray-800 rounded shadow-lg z-50
	             opacity-0 invisible translate-y-1 group-hover:visible group-hover:opacity-100 group-hover:translate-y-0
	             transition-all duration-200 ease-out"
				>
					<a
						id="settingsLink"
						href="/settings"
						data-sveltekit-reload
						class="block w-full px-4 py-2 text-center hover:bg-gray-100"
					>
						Settings
					</a>
					<form method="POST" action="/logout" class="w-full">
						<button
							id="logoutButton"
							type="submit"
							class="block w-full px-4 py-2 text-center hover:bg-gray-100 border-none bg-transparent cursor-pointer text-gray-800"
						>
							Log out
						</button>
					</form>
				</div>
			</div>
		{:else}
			<a id="loginLink" href="/login" class="text-blue-600 hover:underline" data-sveltekit-reload>
				Log in
			</a>
		{/if}
	</div>
</div>
