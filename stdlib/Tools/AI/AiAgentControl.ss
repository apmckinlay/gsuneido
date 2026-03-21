// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
// BuiltDate > 20260219
Controller
	{
	Xmin: 800
	Ymin: 800
	Title: "AI Chat"
	url: "https://openrouter.ai/api/v1"
	modelSettingKey: "AiAgentControl_model"
	validModels: #(
		"anthropic/claude-haiku-4.5",
		"deepseek/deepseek-v3.2",
		"minimax/minimax-m2.7",
		"moonshotai/kimi-k2.5",
		"openai/gpt-5.3-codex",
		"qwen/qwen3-coder-next",
		"nvidia/nemotron-3-super-120b-a12b:free",
		"xiaomi/mimo-v2-pro",
		"x-ai/grok-code-fast-1",
		"z-ai/glm-5",
		"z-ai/glm-5-turbo"
		)
	defaultModel: "minimax/minimax-m2.7"

	New()
		{
		DeleteOldFiles('.ai/', -7) /*= one week */
		prompt = Query1("suneidoc", path: "/res", name: "AiPrompt").text
		.model = .FindControl("model")
		model = UserSettings.Get(.modelSettingKey, .defaultModel)
		if not .validModels.Has?(model)
			model = .defaultModel
		.model.Set(model)

		.agent = AiAgent(.url, .getApiKey(), model, .output, prompt)
		.vert = .VertSplit.Vert
		.editor = .FindControl("editor")
		Defer({ .editor.SetFocus() })
		}
	getApiKey()
		{
		s = Https.Get("https://appserver.internal.axonsoft.com:8088/Wiki?OpenRouterApiKeyForAiAgentControl")
		return s.Extract(`sk-or-v1-[0-9a-f]+`)
		}
	Controls()
		{
		["VertSplit",
			[#Mshtml, .page, name: "webView", xstretch: 1, ystretch: 4],
			[#Vert,
				#Skip,
				.normalButtons(),
				#Skip,
				#(ScintillaAddons, name: "editor", wrap:, xstretch: 1)
				]
			]
		}
	normalButtons(model = false)
		{
		if model is false
			model = .defaultModel
		return [#Horz,
			#Fill,
			#(Button, "Send", tip: "send a message to the AI"),
			#Fill,
			#(Button, "Stop", tip: "interrupt the AI"),
			#Fill,
			#(Button, "New", tip: "start a new conversation"),
			#Fill,
			#(Button, "Load", tip: "load a previous conversation"),
			#Fill,
			[#ChooseButton model, list: .validModels, name: "model"],
			#Fill
			]
		}
	approveButtons: #(Horz,
		Fill,
		(EnhancedButton, "Allow", tip: "let the action go ahead"
			mouseEffect:, buttonStyle:, pad: 20, weight: bold, textColor: 0x007700)
		Fill,
		(EnhancedButton, "Deny", tip: "block the action"
			mouseEffect:, buttonStyle:, pad: 20, weight: bold, textColor: 0x0000ff)
		Fill
		)
	sending: false
	On_Send()
		{
		text = .userText()
		if text is ""
			return
		.FindControl("Send").SetEnabled(false)
		.sending = true
		.agent.Input(text)
		}
	userText()
		{
		text = .editor.Get().Trim()
		if text isnt ""
			.AppendMd(text $ "\n\n", "user")
		.editor.Set("")
		.editor.SetFocus()
		return text
		}
	Enter_Pressed(pressed = false)
		{
		if KeyPressed?(VK.SHIFT, :pressed)
			return 0 // allow default (newline)
		.On_Send()
		return false
		}

	output(what, data, update = false)
		{
		switch what
			{
		case "user":
			.Defer({ .AppendMd("**You:** " $ data, what) })
		case "think", "tool", "output":
			.Defer({ .AppendMd(data, what) })
		case "complete":
			.Defer()
				{
				//.AppendMd("\n\n---\n")
				.AppendMd(.endMarker)
				.FindControl("Send").SetEnabled(true)
				.sending = false
				}
		default:
			}
		if update isnt false
			.Defer({ .showAllowDeny(update) })
		}
	endMarker: "[END_OF_MESSAGE]"

	buttonsRowIndex: 1
	replaceBottomRow(row)
		{
		.vert.Remove(.buttonsRowIndex)
		.vert.Insert(.buttonsRowIndex, row)
		}

	showAllowDeny(update)
		{
		.pendingUpdate = update
		.selectedModel = .model.Get()
		.replaceBottomRow(.approveButtons)
		}

	restoreNormalButtons()
		{
		model = .selectedModel
		.replaceBottomRow(.normalButtons(model))
		.model = .FindControl("model")
		.model.Set(model)
		.FindControl("Send").SetEnabled(false)
		}

	On_Allow()
		{
		if .pendingUpdate is false
			return
		update = .pendingUpdate
		.pendingUpdate = false
		update.Allow(.userText())
		.restoreNormalButtons()
		}

	On_Deny()
		{
		if .pendingUpdate is false
			return
		update = .pendingUpdate
		.pendingUpdate = false
		update.Deny(.userText())
		.restoreNormalButtons()
		}

	On_Stop()
		{
		.agent.Interrupt()
		.AppendMd("\n\n*Response interrupted by user.*\n\n")
		.AppendMd(.endMarker)
		.FindControl("Send").SetEnabled(true)
		.sending = false
		}

	On_New()
		{
		.agent.ClearHistory()
		.Defer({ .FindControl("webView").Set(.page) })
		}
	On_Load()
		{
		filename = OpenFileName(filter: "Log Files (*.md)|*.md|All Files (*.*)|*.*")
		if filename isnt ""
			{
			.FindControl("webView").Set(.page)
			.agent.LoadConversation(filename)
			}
		}

	agent: false
	model: false
	NewValue(value, source)
		{
		if source is .model and .agent isnt false
			{
			.agent.SetModel(value)
			.selectedModel = value
			}
		}

	AppendMd(chunk, type = "output")
		{
		// base64 encode to avoid unicode issues
		b64 = Base64.Encode(chunk)
		html = `<i data-b64="` $ b64 $ `" data-type="` $ type $ `"></i>`
		.FindControl("webView").InsertAdjacentHTML("base64-sink", "beforeend", html)
		}

page: `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<style>
		body {
			font-family: sans-serif; line-height: 1.5;
			margin: 0 15px; padding-bottom: 20px; }
		.message-bubble { border-bottom: 1px solid #0056b3; padding: 10px 0; }

		/* Kill the extra gap at the top & bottom of sections */
		[class^="msg-"] > :not(.copy-btn):first-of-type { margin-top: 0 !important; }
		[class^="msg-"] > :not(.copy-btn):last-of-type { margin-bottom: 0 !important; }
		[class^="msg-"] {
			position: relative;
			padding-right: calc(1.5em + 10px) !important;
		}

		.msg-think {
			color: #666; font-style: italic; background: #f9f9f9;
			border-left: 3px solid #ccc; padding: 8px 12px; margin: 10px 0;
			font-size: 0.95em;
		}
		.msg-tool {
			font-family: monospace; background: #fffbe6;
			border: 1px solid #ffe58f; border-radius: 4px;
			color: #856404; padding: 8px 12px; margin: 10px 0; font-size: 0.85em;
		}
		.msg-output { color: #000; }
		.msg-user {
			color: #0056b3;
			margin-bottom: 5px;
		}
		.msg-user > p:first-child::before {
			content: "You: ";
			font-weight: 700;
		}

		pre {
			background: #f4f4f4;
			padding: 10px;
			overflow-x: auto;
			border-radius: 5px;
			position: relative;
		}
		code { background: #eee; padding: 2px 4px; border-radius: 3px; }
		pre code { background: transparent; padding: 0; border-radius: 0; }
		#bottom-anchor { height: 1px; margin-top: -1px; }
		/* Copy button styles */
		.copy-btn {
			position: absolute;
			top: 0px;
			right: 0px;
			background: transparent;
			border-radius: 4px;
			border: 1px solid transparent;
			cursor: pointer;
			padding: 4px;
			opacity: 0;
			transition: opacity 0.2s, background 0.2s;
			z-index: 10;
		}
		.copy-btn:hover {
			opacity: 1;
			border-color: #ddd;
			background: #f0f0f0;
		}
		.copy-btn.copied {
			background: #4caf50;
			border-color: #4caf50;
			color: white;
		}
		/* Section copy button - show on hover */
		[class^="msg-"]:hover > .copy-btn {
			opacity: 0.5;
		}
		/* Code block copy button - slightly different styling */
		pre .copy-btn {
			top: 5px;
			right: 5px;
			font-size: 0.85em;
		}
		pre:hover .copy-btn {
			opacity: 0.5;
		}
		pre:hover .copy-btn:hover {
			opacity: 1;
			border-color: #ddd;
			background: #f0f0f0;
		}
		[class^="msg-"]:hover > .copy-btn:hover {
			opacity: 1;
			border-color: #ddd;
			background: #f0f0f0;
		}
	</style>
</head>
<body>
	<div id="chat-history"></div>
	<div id="active-display"></div>
	<div id="base64-sink" style="display:none;"></div>
	<div id="bottom-anchor"></div>

	<script src="https://cdn.jsdelivr.net/npm/marked/marked.min.js"></script>
	<script src="https://cdn.jsdelivr.net/npm/dompurify/dist/purify.min.js"></script>

	<script>
		const sink = document.getElementById('base64-sink');
		const display = document.getElementById('active-display');
		const history = document.getElementById('chat-history');
		const anchor = document.getElementById('bottom-anchor');

		const END_MARKER = "W0VORF9PRl9NRVNTQUdFXQ==";
		let isAtBottom = true;
		let sections = [];

		const visibilityObserver = new IntersectionObserver(([entry]) => {
			isAtBottom = entry.isIntersecting;
		}, { threshold: 0, rootMargin: "0px 0px 50px 0px" });
		visibilityObserver.observe(anchor);

		function scrollToBottom() {
			if (isAtBottom) window.scrollTo(0, document.body.scrollHeight);
		}

		// Copy to clipboard function
		async function copyToClipboard(text, button) {
			try {
				if (navigator.clipboard) {
					await navigator.clipboard.writeText(text);
				} else { // fallback. navigator.clipboard is not supported in WebView2
					const textarea = document.createElement('textarea');
					textarea.value = text;
					textarea.style.position = 'fixed';
					textarea.style.opacity = '0';
					document.body.appendChild(textarea);
					textarea.select();
					const success = document.execCommand('copy');
					document.body.removeChild(textarea);
					if (!success) {
						throw "execCommand copy is not supported";
					}
				}
				button.classList.add('copied');
				setTimeout(() => {
					button.classList.remove('copied');
				}, 1500);
			} catch (err) {
				console.error('Failed to copy:', err);
			}
		}
		// Add copy buttons to all pre elements in a parent
		function addCopyButtonsToCodeBlocks(parent) {
			const preElements = parent.querySelectorAll('pre');
			preElements.forEach(pre => {
				// Skip if already has a copy button
				if (pre.querySelector('.copy-btn')) return;

				const code = pre.querySelector('code');
				const text = code ? code.textContent : pre.textContent;

				const copyBtn = document.createElement('button');
				copyBtn.className = 'copy-btn';
				copyBtn.innerHTML = '&#x1F4CB;';
				copyBtn.title = 'Copy code';
				copyBtn.onclick = function() {
					copyToClipboard(text, copyBtn);
				};
				pre.appendChild(copyBtn);
			});
		}
		const observer = new MutationObserver((mutations) => {
			for (const mutation of mutations) {
				for (const node of mutation.addedNodes) {
					if (!node.dataset || !node.dataset.b64) continue;

					const b64Chunk = node.dataset.b64.trim();
					const type = node.dataset.type || "output";

					if (b64Chunk === END_MARKER) {
						finalizeMessage();
						return;
					}

					try {
						const binaryString = atob(b64Chunk);
						const bytes = Uint8Array.from(binaryString, c => c.charCodeAt(0));
						const decoded = new TextDecoder().decode(bytes);

						let currentSec = sections[sections.length - 1];

						if (!currentSec || currentSec.type !== type) {
							const secDiv = document.createElement('div');
							secDiv.className = 'msg-' + type;
							display.appendChild(secDiv);

							currentSec = { type: type, buffer: "", element: secDiv };
							sections.push(currentSec);
						}
					currentSec.buffer += decoded;

					if (window.marked) {
						const renderBuffer = (type === "think")
							? currentSec.buffer.replace(/\s+$/, '')
							: currentSec.buffer;
						currentSec.element.innerHTML = DOMPurify.sanitize(marked.parse(renderBuffer));
						addCopyButtonsToCodeBlocks(currentSec.element);
						scrollToBottom();
					}
					} catch (e) { console.error("Stream Error:", e); }
				}
			}
		});

		function finalizeMessage() {
			if (sections.length > 0) {
				const bubble = document.createElement('div');
				bubble.className = 'message-bubble';

				sections.forEach(sec => {
					if (sec.buffer && sec.element) {
						const copyBtn = document.createElement('button');
						copyBtn.className = 'copy-btn';
						copyBtn.innerHTML = '&#x1F4CB;';
						copyBtn.title = 'Copy as Markdown';
						copyBtn.onclick = function() {
							let text = sec.buffer;
							// trim trailing whitespace/newlines/tabs
							text = text.replace(/\s+$/, '');
							copyToClipboard(text, copyBtn);
						};
						sec.element.prepend(copyBtn);
					}
				});

				while (display.firstChild) {
					bubble.appendChild(display.firstChild);
				}

				history.appendChild(bubble);
			}
			sections = [];
			display.innerHTML = "";
			sink.innerHTML = "";
			scrollToBottom();
		}

		observer.observe(sink, { childList: true });
	</script>
</body>
</html>`

	Destroy()
		{
		if .model isnt false
			UserSettings.Put(.modelSettingKey, .model.Get())
		.agent.Close()
		}
	}