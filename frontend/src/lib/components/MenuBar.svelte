<script lang="ts">
	let { username }: { username: string | null } = $props();
	let menuOpen: boolean = $state(false);

	function toggleMenu() {
		menuOpen = !menuOpen;
	}

	function closeMenuOnBlur(event: FocusEvent) {
		const container = event.currentTarget as HTMLElement | null;
		const next = event.relatedTarget as Node | null;
		if (!container || !next || !container.contains(next)) menuOpen = false;
	}

	function openMenu() {
		menuOpen = true;
	}

	function closeMenu() {
		menuOpen = false;
	}

	function handleMenuKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape') {
			menuOpen = false;
		}
		if (event.key === ' ' || event.key === 'Enter') {
			event.preventDefault();
			toggleMenu();
		}
	}
</script>

<div class="flex items-center justify-between px-4 py-2 w-full">
	<div class="w-1/3"></div>
	<div class="w-1/3">
		<h1 class="text-3xl font-bold text-center">Edmonton Bike Map</h1>
	</div>
	<div class="w-1/3 flex justify-end">
		{#if username}
			<div
				class="relative inline-block"
				onfocusout={closeMenuOnBlur}
				onmouseenter={openMenu}
				onmouseleave={closeMenu}
				role="presentation"
			>
				<button
					id="usernameButton"
					class="text-black text-center px-4 py-1 hover:underline focus-visible:outline-2 focus-visible:outline-blue-600 focus-visible:outline-offset-2"
					aria-haspopup="menu"
					aria-expanded={menuOpen}
					onclick={toggleMenu}
					onkeydown={handleMenuKeydown}
				>
					{username}
				</button>
				<div
					class={`absolute left-0 mt-2 w-full bg-white text-gray-800 rounded shadow-lg z-50 transition-all duration-200 ease-out ${menuOpen ? 'opacity-100 visible translate-y-0' : 'opacity-0 invisible translate-y-1'}`}
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
