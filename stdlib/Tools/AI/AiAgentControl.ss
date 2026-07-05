// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
// BuiltDate > 20260219
Controller
	{
	Xmin: 800
	Ymin: 800
	Title: "AI Chat"
	url: "https://openrouter.ai/api/v1"
	keyUrl: "https://appserver.internal.axonsoft.com:8088/Wiki?" $
		"OpenRouterApiKeyForAiAgentControl"
	modelSettingKey: "AiAgentControl_model"
	models: #(
		"arcee-ai/trinity-large-thinking":
			{ context: "256K", in: .22, out: .85 },
		"deepseek/deepseek-v4-flash":
			{ context: "1M", in: 0.10, out: 0.20 },
		"deepseek/deepseek-v4-pro":
			{ context: "1M", in: 0.44, out: 0.87 },
		"google/gemini-3.1-flash-lite":
			{ context: "1M", in: 0.25, out: 1.50 },
		"minimax/minimax-m3":
			{ context: "1M", in: 0.60, out: 2.40 },
		"moonshotai/kimi-k2.7-code":
			{ context: "256K", in: 0.75, out: 3.50 },
		"nvidia/nemotron-3-super-120b-a12b":
			{ context: "1M", in: 0.09, out: 0.45 },
		"nvidia/nemotron-3-ultra-550b-a55b":
			{ context: "1M", in: 0.50, out: 2.50 },
		"qwen/qwen3.7-plus":
			{ context: "1M", in: 0.32, out: 1.28 },
		"qwen/qwen3.7-max":
			{ context: "1M", in: 1.25, out: 3.75 },
		"xiaomi/mimo-v2.5-pro":
			{ context: "1M", in: 0.44, out: 0.87 },
		)
	defaultModel: "deepseek/deepseek-v4-pro"

	CallClass(setText = '')
		{
		DeleteOldFiles('.ai/', -7) /*= one week */
		prompt = Query1("suneidoc", path: "/res", name: "AiPrompt").text
		model = UserSettings.Get(.modelSettingKey, .defaultModel)
		if not .models.Member?(model)
			model = .defaultModel
		// cache the key since the ai sandbox will prevent fetching it again
		key = Suneido.GetInit(#AIAGENT_API_KEY, .getApiKey)
		super.CallClass(key, model, prompt, setText)
		}
	getApiKey()
		{
		try
			{
			x = HttpClient2(#GET, .keyUrl)
			return x.content.Extract(`sk-or-v1-[0-9a-f]+`)
			}
		catch (e)
			throw "error getting api key from wiki: " $ e
		}
	New(key, model, prompt, setText = '')
		{
		.agent = AiAgent(.url, key, model, .output, prompt)
		.model = .FindControl("model")
		.model.Set(model)
		.vert = .Vert.VertSplit.Vert
		.editor = .FindControl("Editor")
		if setText isnt ''
			.editor.Set(setText)
		.status = .FindControl("statusbar")
		Defer({ .editor.SetFocus() })
		}

	Commands: ((Users_Manual,	"F1"))
	Menu: ()
	Controls()
		{
		["Vert",
			["VertSplit",
				[#Mshtml, .page, name: "webView", xstretch: 1, ystretch: 4],
				[#Vert,
					#Skip,
					.normalButtons(),
					#Skip,
					#(ScintillaAddons, name: "Editor", wrap:, xstretch: 1),
					]
				]
			#(Statusbar, name: "statusbar")
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
			[#ChooseButton model, name: "model",
				list: .models.Members().Sort!()],
			#Skip
			#(LinkButton, "?", modelHelp)
			#Fill
			]
		}
	ApproveButtons: #(Horz,
		Fill,
		(EnhancedButton, "Allow", tip: "let the action go ahead"
			mouseEffect:, buttonStyle:, pad: 20, weight: bold, textColor: 0x007700)
		Fill,
		(EnhancedButton, "Deny", tip: "block the action"
			mouseEffect:, buttonStyle:, pad: 20, weight: bold, textColor: 0x0000ff)
		Fill
		)
	On_modelHelp()
		{
		w0 = 40
		w1 = 10
		w2 = 6
		w3 = 6
		s = "Id".RightFill(w0) $ "Context".LeftFill(w1) $
			"In".LeftFill(w2) $ "Out".LeftFill(w3) $ "\n"
		s $= "-".Repeat(w0 + w1 + w2 + w3) $ "\n"
		for m in .models.Members().Sort!()
			{
			x = .models[m]
			s $= m.RightFill(w0) $
				x.context.LeftFill(w1) $
				x.in.Format('##.##').LeftFill(w2) $
				x.out.Format('##.##').LeftFill(w3) $ "\n"
			}
		Alert(s, font: "@mono", title: "Models")
		}
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

	output(what, data, approve = false)
		{
		switch what
			{
		case "user":
			.AppendMd("**You:** " $ data, what)
		case "think", "tool", "output":
			.AppendMd(data, what)
			.updateStatus()
		case "complete":
			.AppendMd(.endMarker)
			if false isnt sendBtn = .FindControl("Send")
				sendBtn.SetEnabled(true)
			.sending = false
			.updateStatus()
		default:
			}
		if approve isnt false
			.Defer({ .approval(approve) })
		}
	endMarker: "[END_OF_MESSAGE]"

	updateStatus()
		{
		if .model.Destroyed?()
			return
		if '' is model = .model.Get()
			return
		contextLimit = .models[model].context
		try // in case the exe doesn't have Usage or Cost yet
			{
			// using ending spaces to avoid overlapping with the resizing handler
			.status.Set("\t\tContext: " $ .k(.agent.Usage()) $ " / " $ contextLimit $
				"  |  Cost: " $ .agent.Cost().Format("##.##") $ '      ' )
			}
		}
	k(n)
		{
		k = 1000
		return (n / k).RoundToPrecision(2) $ "k"
		}

	buttonsRowIndex: 1
	replaceBottomRow(row)
		{
		.vert.Remove(.buttonsRowIndex)
		.vert.Insert(.buttonsRowIndex, row)
		}

	approval(approve)
		{
		.pendingUpdate = approve
		.selectedModel = .model.Get()
		.replaceBottomRow(.ApproveButtons)
		before = approve.Before()
		after = approve.After()
		if after isnt ""
			{
			if false isnt response = before is ""
				? AiAgentView(.Window.Hwnd, after, .editLib, .editName)
				: AiAgentDiff(.Window.Hwnd, before, after, .editLib, .editName)
				{
				.editor.Set(response.feedback)
				switch response[0]
					{
				case "allow":
					.On_Allow()
				case "deny":
					.On_Deny()
				default:
					}
				}
			.editLib = .editName = ''
			}
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
		PubSub.PublishConsolidate('LibraryTreeChange', force:)
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
		.FindControl("webView").Set(.page)
		.status.Set("")
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
	selectedModel: false
	NewValue(value, source)
		{
		if source is .model and .agent isnt false
			{
			.agent.SetModel(value)
			.selectedModel = value
			}
		}

	appendDeferred: false
	editLib: ''
	editName: ''
	AppendMd(chunk, type = "output")
		{
		.queue.Add(Object(chunk, type))
		if .appendDeferred is true
			return
		.appendDeferred = true
		.Defer()
			{
			while not Same?(.queue, first = .queue.PopFirst())
				.appendMd(first)
			.appendDeferred = false
			}
		return
		}

	getter_queue()
		{
		.queue = Object()
		}

	appendMd(item)
		{
		if false is webview = .FindControl("webView")
			return
		chunk = item[0]
		type = item[1]
		// base64 encode to avoid unicode issues
		b64 = Base64.Encode(chunk)
		html = `<i data-b64="` $ b64 $ `" data-type="` $ type $ `"></i>`
		if type is 'tool' and
			(chunk.Prefix?("**Edit Code** ") or chunk.Prefix?("**Create Code** "))
			{
			.editLib = chunk.AfterFirst('`').BeforeFirst('`')
			.editName = chunk.AfterFirst('` `').BeforeFirst('`')
			}
		webview.InsertAdjacentHTML("base64-sink", "beforeend", html)
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
			tab-size: 4;
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
						currentSec.element.innerHTML =
							DOMPurify.sanitize(marked.parse(renderBuffer));
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
		if .selectedModel isnt false
			UserSettings.Put(.modelSettingKey, .selectedModel)
		.agent.Close()
		}
	}