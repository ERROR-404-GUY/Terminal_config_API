<script lang="ts">
	import { onMount } from 'svelte';
	import { terminals } from '$lib/store/Terminal';
	import { slide } from 'svelte/transition';
	import { typewriter } from '$lib/typewriter';

	let loading = $state(true);
	let sortby = $state('tid');
	let filter = $state('refund_allowed');

	let sortedTerminals = $derived(
		[...$terminals].sort((a, b) => {
			if (sortby === 'tid') {
				if (a.tid > b.tid) return -1;
				if (a.tid < b.tid) return 1;
				return 0;
			}
			if (sortby === 'product_name') {
				if (a.product_name < b.product_name) return -1;
				if (a.product_name > b.product_name) return 1;
				return 0;
			}
			if (sortby === 'serial_number') {
				if (a.serial_number < b.serial_number) return -1;
				if (a.serial_number > b.serial_number) return 1;
				return 0;
			}
			if (sortby === 'refund_allowed') {
				if (a.refund_allowed < b.refund_allowed) return -1;
				if (a.refund_allowed > b.refund_allowed) return 1;
				return 0;
			}
			return 0;
		})
	);
	let filteredTerminals = $derived(
		sortedTerminals.filter((t) => {
			if (filter === 'all') return true;
			if (filter === 'refund_allowed') return t.refund_allowed;
			if (filter === 'refund_not_allowed') return !t.refund_allowed;
			return true;
		})
	);
	filter = "all"
	
	const toggleRefund = async (tid: string, currentStatus: boolean) => {
		try {
			const res = await fetch(`http://localhost:8081/api/terminals/${tid}?refund_allowed=${!currentStatus}`, {
				method: 'PUT'
			});
			if (res.ok) {
				const updatedTerminal = await res.json();
				terminals.update(current => 
					current.map(t => t.tid === tid ? updatedTerminal : t)
				);
			} else {
				console.error('Failed to toggle refund status');
			}
		} catch (error) {
			console.error('Error toggling refund status:', error);
		}
	};

	const deleteTerminal = async (tid: string) => {
		if (!confirm(`Are you sure you want to delete terminal ${tid}?`)) return;
		
		try {
			const res = await fetch(`http://localhost:8081/api/terminals/${tid}`, {
				method: 'DELETE'
			});
			if (res.ok) {
				terminals.update(current => current.filter(t => t.tid !== tid));
			} else {
				console.error('Failed to delete terminal');
				alert('Failed to delete terminal');
			}
		} catch (error) {
			console.error('Error deleting terminal:', error);
			alert('Error deleting terminal');
		}
	};

	onMount(async () => {
		try {
			const res = await fetch('http://localhost:8081/api/terminals');
			const data = await res.json();
			terminals.set(data);
		} catch (error) {
			console.error('Error fetching terminals:', error);
		} finally {
			loading = false;
		}
	});

</script>

<svelte:head>
	<title>Terminals :: Mainframe</title>
</svelte:head>

<div class="flex flex-col gap-6">
	<div>
		<p class="mb-4">> EXECUTING: GET /api/terminals</p>
	</div>

	{#if loading}
		<div class="flex animate-pulse flex-col gap-2" use:typewriter={{speed: 10}}>
			<p>> ESTABLISHING SECURE CONNECTION...</p>
			<p>> FETCHING AUTHORIZED TERMINALS...</p>
			<div class="mt-4 flex gap-1">
				{#each Array(20) as _}
					<span class="inline-block w-2 bg-[#4af626]">&nbsp;</span>
				{/each}
			</div>
		</div>
	{:else}
		<div use:typewriter={{speed: 10}}>
			<div class="flex items-center justify-between mb-6">
				<p class="opacity-80">> RESULTS: {$terminals.length} RECORDS FOUND</p>
				<a href="/create" class="border border-current px-4 py-1 hover:bg-[#4af626]/20 transition-colors">> CREATE TERMINAL</a>
			</div>
			<label for="terminal-select">SORT BY</label>
			<select
				id="terminal-select"
				onchange={(event) => (sortby = event.currentTarget.value)}
				bind:value={sortby}
			>
				<option value="tid">terminal id</option>
				<option value="product_name">product name</option>
				<option value="serial_number">serial number</option>
				<option value="refund_allowed">refund allowed</option>
			</select>
			<br />
			<label for="terminal-select">FILTER BY</label>
			<select
				id="terminal-select"
				onchange={(event) => (filter = event.currentTarget.value)}
				bind:value={filter}
			>
				<option value="all">all</option>
				<option value="refund_allowed">refund allowed</option>
				<option value="refund_not_allowed">refund not allowed</option>
			</select>
			<div class="overflow-x-auto">
				<div class="w-full min-w-[700px] whitespace-nowrap text-left text-sm md:text-base font-mono">
					<div class="grid grid-cols-[1.5fr_1fr_1fr_100px] border-y border-current">
						<div class="py-2 pr-4 font-normal">PRODUCT NAME</div>
						<div class="py-2 px-4 font-normal">TERMINAL ID</div>
						<div class="py-2 px-4 font-normal">SERIAL NUMBER</div>
						<div class="py-2 pl-4 font-normal">REFUND?</div>
					</div>
					<div>
						{#each filteredTerminals as t (t.tid)}
							<div transition:slide>
								<div class="grid grid-cols-[1.5fr_1fr_1fr_100px] pt-4">
									<div class="pr-4">{t.product_name}</div>
									<div class="px-4">{t.tid}</div>
									<div class="px-4">{t.serial_number}</div>
									<div class="pl-4">
										<button 
											onclick={() => toggleRefund(t.tid, t.refund_allowed)}
											class="hover:bg-[#4af626]/20 px-2 -ml-2 transition-colors border border-transparent hover:border-current"
										>
											{t.refund_allowed ? ' YES' : ' NO '}
										</button>
									</div>
								</div>
								<div class="border-b border-current pb-4 pt-1 opacity-80 flex items-center justify-between">
									<span>> ACTIVATION CODE: <span class="tracking-widest">{t.activation_code}</span></span>
									<button onclick={() => deleteTerminal(t.tid)} class="text-red-500 hover:bg-red-500/10 px-2 py-0 border border-transparent hover:border-red-500 transition-colors">[DELETE]</button>
								</div>
							</div>
						{/each}
					</div>
				</div>
			</div>

			<div class="mt-8 flex items-center gap-2">
				<span>></span>
				<span class="h-5 w-3 animate-[pulse_1s_step-end_infinite] bg-[#4af626]"></span>
			</div>
		</div>
	{/if}
</div>
