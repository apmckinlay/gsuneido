// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass()
		{
		return .styles $ .html $ .scripts
		}

	styles: `
<style>
	.custom-context-menu {
		display: none; /* Hidden by default */
		position: fixed; /* Use fixed for accurate positioning relative to viewport */
		z-index: 1000;
		background-color: white;
		border-spacing: 0px;
		border-radius: 0.3em;
		box-shadow: 3px 3px 3px grey;
		border: 1px solid black;
		overflow: hidden;
		list-style: none;
		padding: 0;
		margin: 0;
		font: initial;
	}
	.custom-context-menu li {
		padding: 5px 10px;
		cursor: pointer;
		min-width: 10em;
		margin: 0;
		color: black;
	}

	.custom-context-menu li:hover {
		background-color: lightblue;
	}

	/* Hides the element entirely when the 'disabled' class is applied */
	.custom-context-menu li.disabled {
		display: none !important;
	}
</style>
`
	html: `
<ul id="context-menu" class="custom-context-menu">
	<li id="menu-copy" data-action="copy">Copy</li>
	<li id="menu-copy-link" data-action="copy-link">Copy Link</li>
</ul>
`

	scripts: `<script>
	const contextMenu = document.getElementById('context-menu');
	const menuCopyLink = document.getElementById('menu-copy-link');
	let activeLinkElement = null;
	let selectedText = null;

	function handleCopy() {
		if (selectedText) {
			copyViaExecCommand(selectedText);
		}
	}

	function copyViaExecCommand(text) {
		let tempInput = document.createElement('textarea');
		tempInput.value = text;
		// Position off-screen to avoid visual disruption
		tempInput.style.position = 'absolute';
		tempInput.style.left = '-9999px';

		document.body.appendChild(tempInput);

		// Select and copy the text
		tempInput.select();
		document.execCommand('copy');

		// Clean up the temporary element
		document.body.removeChild(tempInput);
	}

	function handleCopyLink() {
		if (!activeLinkElement) {
			return;
		}

		const linkToCopy = activeLinkElement.getAttribute('data-copy-link');

		if (linkToCopy) {
			copyViaExecCommand(linkToCopy);
		}
	}

	function handleMenuClick(event) {
		const action = event.target.getAttribute('data-action');
		if (!action || event.target.classList.contains('disabled')) return;

		contextMenu.style.display = 'none';

		switch (action) {
			case 'copy':
				handleCopy();
				break;
			case 'copy-link':
				handleCopyLink();
				break;
		}
	}

	document.addEventListener('contextmenu', (event) => {
		event.preventDefault();

		const mouseX = event.clientX;
		const mouseY = event.clientY;

		const linkElement = event.target.closest('a');

		// Find the 'a' element and update state
		if (linkElement && linkElement.hasAttribute('data-copy-link')) {
			activeLinkElement = linkElement;
			menuCopyLink.classList.remove('disabled');
		} else {
			activeLinkElement = null;
			menuCopyLink.classList.add('disabled');
		}

		selectedText = window.getSelection().toString();
		// --- Position Logic ---
		// Calculate dimensions after visibility change (since hiding/showing affects dimensions)
		contextMenu.style.display = 'block';
		const menuWidth = contextMenu.offsetWidth;
		const menuHeight = contextMenu.offsetHeight;
		const viewportWidth = window.innerWidth;
		const viewportHeight = window.innerHeight;

		let left = mouseX;
		let top = mouseY;

		// Adjust horizontally if menu goes off-screen right
		if (mouseX + menuWidth > viewportWidth) {
			left = viewportWidth - menuWidth - 5;
		}
		// Adjust vertically if menu goes off-screen bottom
		if (mouseY + menuHeight > viewportHeight) {
			top = viewportHeight - menuHeight - 5;
		}

		// Ensure menu doesn't go off-screen left/top
		if (left < 0) left = 5;
		if (top < 0) top = 5;

		contextMenu.style.top = top + 'px';
		contextMenu.style.left = left + 'px';
	});

	// 2. Hide the menu on any regular click outside
	document.addEventListener('click', () => {
		contextMenu.style.display = 'none';
	});

	// 3. Prevent the menu from closing when clicking on an item
	contextMenu.addEventListener('click', handleMenuClick);

</script>
`
	}