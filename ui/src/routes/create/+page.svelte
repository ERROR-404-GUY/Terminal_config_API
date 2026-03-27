<script lang="ts">
	import { goto } from '$app/navigation';
	import { typewriter } from '$lib/typewriter';

	let loading = $state(false);
	let error = $state('');

	let tid = $state('');
	let product_name = $state('');
	let serial_number = $state('');
	let refund_allowed = $state(false);

	const handleSubmit = async (e: Event) => {
		e.preventDefault();
		loading = true;
		error = '';

		try {
			const res = await fetch('http://localhost:8081/api/terminals', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({
					product_name,
					refund_allowed
				})
			});

			if (res.ok) {
				goto('/');
			} else {
				const data = await res.json().catch(() => null);
				error = data?.error || 'Failed to create terminal';
			}
		} catch (err) {
			console.error('Error creating terminal:', err);
			error = 'Error connecting to server';
		} finally {
			loading = false;
		}
	};
</script>

<svelte:head>
	<title>Create Terminal :: Mainframe</title>
</svelte:head>

<div class="flex flex-col gap-6 font-mono" use:typewriter={{speed: 10}}>
	<div class="flex items-center justify-between">
		<p>> EXECUTING: INIT_TERMINAL_PROTOCOL</p>
		<a href="/" class="opacity-80 hover:opacity-100 hover:text-red-500 transition-colors">[ABORT]</a>
	</div>

	<div class="border border-current p-6">
		<h2 class="mb-6 opacity-80">> ENTER TERMINAL PARAMETERS:</h2>

		{#if error}
			<div class="mb-6 border border-red-500 p-2 text-red-500">
				> ERROR: {error}
			</div>
		{/if}

		<form onsubmit={handleSubmit} class="flex flex-col gap-4 max-w-lg">
			<div class="flex flex-col gap-1">
				<label for="product_name" class="opacity-80 mt-2">PRODUCT NAME</label>
				<select
					id="product_name"
					bind:value={product_name}
					required
					class="border border-current bg-black px-3 py-2 text-[#4af626] focus:outline-none focus:ring-1 focus:ring-[#4af626]"
				>
					<option value="" disabled selected hidden>Select a product...</option>
					<option value="POS Terminal">POS Terminal</option>
					<option value="Model X-90">Model X-90</option>
					<option value="Handheld Reader">Handheld Reader</option>
					<option value="Register Pro">Register Pro</option>
				</select>
			</div>

			<div class="flex items-center gap-3 mt-4">
				<input
					type="checkbox"
					id="refund_allowed"
					bind:checked={refund_allowed}
					class="h-4 w-4 appearance-none border border-current checked:bg-[#4af626] focus:outline-none cursor-pointer"
				/>
				<label for="refund_allowed" class="opacity-80 cursor-pointer">ALLOW REFUNDS</label>
			</div>

			<div class="mt-8 flex items-center gap-4">
				<button
					type="submit"
					disabled={loading}
					class="border border-[#4af626] px-6 py-2 text-[#4af626] hover:bg-[#4af626]/20 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
				>
					{loading ? '> TRANSMITTING...' : '> INITIALIZE REGISTRATION'}
				</button>
			</div>
		</form>
	</div>
</div>
