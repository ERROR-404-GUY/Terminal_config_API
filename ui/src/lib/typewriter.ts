export function typewriter(node: HTMLElement, { speed = 5 } = {}) {
	if (typeof window === 'undefined') return;

	const texts: { node: Text; full: string }[] = [];
	const walker = document.createTreeWalker(node, NodeFilter.SHOW_TEXT, null);
	
	let n;
	while ((n = walker.nextNode())) {
		if (n.nodeValue && n.nodeValue.trim().length > 0) {
			texts.push({ node: n as Text, full: n.nodeValue });
			n.nodeValue = '';
		}
	}

	let activeNodeIdx = 0;
	let charIdx = 0;
	let frame: number;

	function tick() {
		if (activeNodeIdx >= texts.length) return;

		const current = texts[activeNodeIdx];
		charIdx += speed;
		
		if (charIdx >= current.full.length) {
			current.node.nodeValue = current.full;
			activeNodeIdx++;
			charIdx = 0;
		} else {
			current.node.nodeValue = current.full.slice(0, charIdx);
		}
		frame = requestAnimationFrame(tick);
	}

	frame = requestAnimationFrame(tick);

	return {
		destroy() {
			if (frame) cancelAnimationFrame(frame);
		}
	};
}
