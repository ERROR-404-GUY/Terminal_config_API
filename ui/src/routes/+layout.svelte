<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import { typewriter } from '$lib/typewriter';

	let { children } = $props();
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>

<div class="flex min-h-screen flex-col p-4 leading-relaxed tracking-widest uppercase md:p-8">
	<header class="mb-8 border-b-2 border-[#4af626] pb-4 opacity-80" use:typewriter={{speed: 4}}>
		<pre class="font-bold whitespace-pre-wrap select-none">
====================================================
           [ CONFIG API TERMINAL ] v4.0.4           
   UNAUTHORIZED ACCESS IS STRICTLY PROHIBITED       
====================================================
		</pre>
	</header>

	<main class="flex-grow">
		{@render children()}
	</main>

	<footer class="mt-8 flex justify-between border-t-2 border-[#4af626] pt-4 text-xs opacity-60" use:typewriter={{speed: 2}}>
		<span>SYSTEM STATUS: ONLINE</span>
		<span>USER: ROOT</span>
	</footer>
</div>

