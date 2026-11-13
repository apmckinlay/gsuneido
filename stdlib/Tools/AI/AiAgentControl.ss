// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Xmin:  800
	Ymin:  800
	Title: "AI Chat"
	url:   "https://openrouter.ai/api/v1"
	keyUrl: "https://appserver.internal.axonsoft.com:8088/Wiki?" $
		"OpenRouterApiKeyForAiAgentControl"
	modelSettingKey: #AiAgentControl_model
	models: (
		"~deepseek/deepseek-flash-latest": [in: .22, out: .40],
		"google/gemini-3.8-flash": [in: .75, out: 3.75],
		"google/gemini-4-hydrogen": [in: 1.5, out: 7.5],
		"meituan/longcat-2.0": [in: .30, out: 1.20],
		"minimax/minimax-m3": [in: .23, out: .96],
		"nvidia/nemotron-3-ultra-550b-a55b": [in: .50, out: 2.20],
		"~openai/gpt-luna-latest": [in: .10, out: .50],
		"~openai/gpt-terra-latest": [in: 1, out: 5],
		"~openai/gpt-sol-latest": [in: 2, out: 10],
		"openrouter/auto": [in: .5, out: 2],
		"qwen/qwen3.8-max-0902": [in: 2, out: 6],
		"qwen/qwen3.8-flash": [in: .15, out: .47],
		"typesafe/jev-router": [in: .5, out: 2],
		"xiaomi/mimo-v2.6-pro": [in: .32, out: 0.65],
		"z-ai/glm-5.3": [in: 1.40, out: 4.40],
		"z-ai/glm-5.3-flash": [in: .08, out: .25]
		)
	defaultModel: "~deepseek/deepseek-flash-latest"

	CallClass(setText = "")
		{
		DeleteOldFiles(".ai/", -7) /*= one week */
		aiPrompt = TableExists?(#suneidoc)
			? Query1(#suneidoc, path: "/res", name: #AiPrompt)
			: false
		if aiPrompt is false
			{
			Alert("suneidoc table not found or suneidoc /res AiPrompt not found",
				title: .Title, flags: MB.ICONERROR)
			return
			}
		prompt = aiPrompt.text
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

	New(.key, model, prompt, setText = "")
		{
		.agent = AiAgent(.url, key, model, .output, prompt)
		.model = .FindControl(#model)
		.model.Set(model)
		.vert = .Vert.VertSplit.Vert
		.FindControl(#enterToSend).
			Set(.enterToSend = UserSettings.Get(#AiAgentControl_EnterToSend, true))
		.editor = .FindControl(#Editor)
		if setText isnt ""
			.editor.Set(setText)
		.status = .FindControl(#statusbar)
		.startTime = Date()
		.Defer(
			{
			.updateOpenRouterUsage()
			.updateStatus()
			})
		Defer({ .editor.SetFocus() })
		}

	Commands: ((Users_Manual, F1), (Go_To_Definition, F12))
	Menu: ()
	Controls()
		{
		[#Vert,
			[#VertSplit,
				[#Mshtml, .page, name: #webView, xstretch: 1, ystretch: 4],
				[#Vert,
					#Skip,
					.normalButtons(),
					#Skip,
					#(ScintillaAddons, name: Editor, wrap:, xstretch: 1)]],
			#(Statusbar, name: statusbar)]
		}

	normalButtons(model = false)
		{
		if model is false
			model = .defaultModel
		return [#Horz,
			#Fill,
			#(Button, Send, tip: "send a message to the AI"),
			#Skip,
			#(CheckBox, tip: "Enter key sends", tabover:, name: enterToSend),
			#(Fill),
			#(Button, Stop, tip: "interrupt the AI"),
			#Fill, #Skip,
			#(Button, New, tip: "start a new conversation"),
			#Fill, #Skip,
			#(Button, Load, tip: "load a previous conversation"),
			#Fill, #Skip,
			[#ChooseButton, model, name: #model,
				list: .models.Members().Sort!()],
			#Skip,
			#(LinkButton, '?', modelHelp),
			#Fill]
		}

	ApproveButtons: (Horz,
		Fill,
		(EnhancedButton, Allow, tip: "let the action go ahead",
			mouseEffect:, buttonStyle:, pad: 20, weight: bold, textColor: 0x00_7700),
		Fill,
		(EnhancedButton, Deny, tip: "block the action",
			mouseEffect:, buttonStyle:, pad: 20, weight: bold, textColor: 0x00_00ff),
		Fill)
	On_modelHelp()
		{
		AiAgentModelHelp(.Window.Hwnd, .models)
		}

	On_Go_To_Definition()
		{
		if false isnt webview = .FindControl(#webView)
			webview.TriggerKeyDown(VK.F12)
		}

	GoTo(name) // from SuneidoAPP
		{
		pub = GlobalRegExForGoTo
		if name !~ ("^(" $ pub $ "\.)?" $ pub $ '$')
			name = false
		if name is false or Gotofind(name.BeforeFirst('.')).Empty?()
			{
			Beep()
			return ""
			}
		GotoLibView(name)
		return ""
		}

	sending: false
	On_Send()
		{
		text = .userText()
		if text is ""
			return
		.startTime = Date()
		.FindControl(#Send).SetEnabled(false)
		.sending = true
		.agent.Input(text)
		}

	userText()
		{
		text = .editor.Get().Trim()
		if text isnt ""
			.AppendMd(text $ "\n\n", #user)
		.editor.Set("")
		.editor.SetFocus()
		return text
		}

	Enter_Pressed(pressed = false)
		{
		if not .enterToSend
			return 0 // allow default (newline)
		if KeyPressed?(VK.SHIFT, :pressed)
			return 0 // allow default (newline)
		.On_Send()
		return false
		}

	output(what, data, approve = false)
		{
		switch what
			{
		case #user:
			.AppendMd("**You:** " $ data, what)
		case #responseModel:
			.AppendMd("** Response Model** " $ data $ "<br/>", #think)
		case #think, #tool, #output:
			.AppendMd(data, what)
			.updateStatus()
		case #complete:
			.AppendMd(.endMarker)
			if false isnt sendBtn = .FindControl(#Send)
				sendBtn.SetEnabled(true)
			.sending = false
			.Defer(// need defer to get outside sandbox for openrouter usage
				{
				.updateOpenRouterUsage()
				.updateStatus()
				})
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
		try // in case the exe doesn't have Usage or Cost yet
			// using ending spaces to avoid overlapping with the resizing handler
			.status.Set("\t\tContext: " $ .percentUsed(.agent.Usage()) $ "  |  Cost: " $
				.agent.Cost().Format("##.##") $ "  |  Budget: " $ .creditsStr() $
				"  |  Time: " $ .formatTime(Date().MinusSeconds(.startTime)) $ "      ")
		}

	percentUsed(usage)
		{
		contextLimit = 1_000_000
		percent = 100
		usage = (usage / contextLimit * percent).Round(0)
		return usage $ '%'
		}

	formatTime(secs)
		{
		oneMinute = 60
		secs = secs.Round(0)
		return (.fmt((secs / oneMinute).Floor(), #m) $
			Opt(' ', .fmt((secs % oneMinute).Round(0), #s))).Trim()
		}

	fmt(number, unit)
		{
		return number is 0 ? "" : number $ unit
		}

	credits: false

	creditsStr()
		{
		if false is c = .credits
			return ""
		if c.limit is false or c.limit is 0 or c.limit is ""
			return ""
		percent = 100
		usage = (c.usage / c.limit * percent).Round(0)
		return usage $ "% / $" $ c.limit
		}

	updateOpenRouterUsage()
		{
		.credits = false
		try
			{
			url = "https://openrouter.ai/api/v1/auth/key"
			header = Object(Authorization: "Bearer " $ .key)
			result = HttpClient2(#GET, url, :header)
			json = Json.Decode(result.content, handleNull: #skip)
			.credits = [usage: json.data.usage_daily, limit: json.data.limit]
			}
		catch (e)
			Print(e)
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
				case #allow:
					.On_Allow()
				case #deny:
					.On_Deny()
				default:
					}
				}
			.editLib = .editName = ""
			}
		}

	restoreNormalButtons()
		{
		model = .selectedModel
		.replaceBottomRow(.normalButtons(model))
		.model = .FindControl(#model)
		.model.Set(model)
		.FindControl(#Send).SetEnabled(false)
		.FindControl(#enterToSend).Set(.enterToSend)
		}

	On_Allow()
		{
		if .pendingUpdate is false
			return
		update = .pendingUpdate
		.pendingUpdate = false
		update.Allow(.userText())
		.restoreNormalButtons()
		PubSub.PublishConsolidate(#LibraryTreeChange, force:,
			table: .editLib, name: .editName)
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
		.FindControl(#Send).SetEnabled(true)
		.sending = false
		}

	On_New()
		{
		.agent.ClearHistory()
		.FindControl(#webView).Set(.page)
		.status.Set("")
		}

	On_Load()
		{
		if false isnt filename = AiAgentLoadControl()
			{
			.FindControl(#webView).Set(.page)
			.agent.LoadConversation(filename)
			}
		}

	agent:         false
	model:         false
	selectedModel: false
	NewValue(value, source)
		{
		if source is .model and .agent isnt false
			{
			.agent.SetModel(value)
			.selectedModel = value
			}
		if source is ets = .FindControl(#enterToSend)
			.enterToSend = ets.Get()
		}

	appendDeferred: false
	editLib:        ""
	editName:       ""
	AppendMd(chunk, type = #output)
		{
		.queue.Add([chunk, type])
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
		if false is webview = .FindControl(#webView)
			return
		chunk = item[0]
		type = item[1]
		// base64 encode to avoid unicode issues
		b64 = Base64.Encode(chunk)
		html = `<i data-b64="` $ b64 $ `" data-type="` $ type $ `"></i>`
		if type is #tool and
			(chunk.Prefix?("**Edit Code** ") or chunk.Prefix?("**Create Code** "))
			{
			.editLib = chunk.AfterFirst('`').BeforeFirst('`')
			.editName = chunk.AfterFirst("` `").BeforeFirst('`')
			}
		webview.InsertAdjacentHTML("base64-sink", #beforeend, html)
		}

	page: `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<style>
		body {
			font-family: sans-serif; line-height: 1.5;
			margin: 0; display: flex; }
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
		#prompt-summary {
			width: 180px; min-width: 180px; overflow-y: auto;
			border-left: 1px solid #ddd; padding: 10px 8px;
			font-size: 0.85em; position: sticky; top: 0;
			height: 100vh; box-sizing: border-box;
			transition: width 0.2s, min-width 0.2s, padding 0.2s, border 0.2s;
		}
		#prompt-summary.collapsed { width: 0; min-width: 0; padding: 0;
			border-left: none; overflow: hidden; }
		#prompt-summary:empty:not(.collapsed) { display: none; }
		.prompt-item {
			padding: 6px 8px; margin-bottom: 4px; cursor: pointer;
			border-radius: 4px; color: #0056b3;
			white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
		}
		.prompt-item:hover { background: #e8f0fe; }
		#chat-wrapper { flex: 1; min-width: 0; padding: 0 15px 20px; overflow-y: auto;
			height: 100vh; box-sizing: border-box; position: relative; }
		.fab-wrap { position: fixed; bottom: 10px; right: 10px; z-index: 100;
			display: none; }
		#toggle-sidebar {
			width: 36px; height: 36px; border-radius: 50%;
			border: 1px solid #ccc; background: #fff; cursor: pointer;
			font-size: 16px; box-shadow: 0 1px 3px rgba(0,0,0,0.2);
		}
		#toggle-sidebar:hover { background: #f0f0f0; }
	</style>
</head>
<body>
	<div id="chat-wrapper">
		<div id="chat-history"></div>
		<div id="active-display"></div>
		<div id="base64-sink" style="display:none;"></div>
		<div id="bottom-anchor"></div>
		<div class="fab-wrap"><button id="toggle-sidebar" title="Toggle prompt list">
			&#9776;</button></div>
	</div>
	<div id="prompt-summary" class="collapsed"></div>

	<script src="https://cdn.jsdelivr.net/npm/marked/marked.min.js"></script>
	<script src="https://cdn.jsdelivr.net/npm/dompurify/dist/purify.min.js"></script>

	<script>
		const sink = document.getElementById('base64-sink');
		const display = document.getElementById('active-display');
		const history = document.getElementById('chat-history');
		const anchor = document.getElementById('bottom-anchor');
		const chatWrapper = document.getElementById('chat-wrapper');

		const END_MARKER = "W0VORF9PRl9NRVNTQUdFXQ==";
		let isAtBottom = true;
		let sections = [];

		const sidebar = document.getElementById('prompt-summary');
		let promptCount = 0;
		document.getElementById('toggle-sidebar').onclick = function() {
			sidebar.classList.toggle('collapsed');
		};
		const visibilityObserver = new IntersectionObserver(([entry]) => {
			isAtBottom = entry.isIntersecting;
		}, { threshold: 0, root: chatWrapper, rootMargin: "0px 0px 50px 0px" });
		visibilityObserver.observe(anchor);

		document.addEventListener('keydown', function (event) {
			if ((event.key !== 'F12' && event.keyCode !== 123) ||
				event.ctrlKey || event.altKey || event.shiftKey)
				return;
			event.preventDefault();
			if (!document.hasFocus())
				return; // e.g. focus is in the editor, not the output
			const sel = window.getSelection();
			const text = (sel ? sel.toString() : '').trim();
			if (text === '')
				return;
			if (text.length > 100 /*=max name length*/)
				return;
			// image beacon so the browser does not navigate (see SuneidoAPP)
			new Image().src = 'suneido:/eval?AiAgentControl.GoTo("' +
				encodeURIComponent(text) + '")';
		});
		function scrollToBottom() {
			if (isAtBottom) chatWrapper.scrollTop = chatWrapper.scrollHeight;
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

						if (type === "user") {
							promptCount++;
							currentSec.promptId = 'prompt-' + promptCount;
							secDiv.id = currentSec.promptId;
						}
					}
					currentSec.buffer += decoded;

					if (type === "user" && !currentSec.sidebarAdded) {
						const label = decoded.replace(/\*\*/g, '').
							replace(/^You:\s*/, '').replace(/\s+/g, ' ').trim();
						if (label) {
							const item = document.createElement('div');
							item.className = 'prompt-item';
							item.textContent = promptCount + '. ' + (label.length > 40 ?
								label.substring(0, 40) + '...' : label);
							item.title = label;
							const pid = currentSec.promptId;
							item.onclick = function() {
								const el = document.getElementById(pid);
								if (el) el.scrollIntoView(
									{ behavior: 'smooth', block: 'start' });
							};
							sidebar.appendChild(item);
							document.querySelector('.fab-wrap').style.display = 'block';
							currentSec.sidebarAdded = true;
						}
					}

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
		UserSettings.Put(#AiAgentControl_EnterToSend, .enterToSend)
		.agent.Close()
		}
	}
