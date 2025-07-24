// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
class
	{
	styles: `
		.su-overlay-anchor {
			padding: 0;
			display: inline-block;
			border: none;
			border-spacing: 0px;
		}
		.su-overlay-anchor::backdrop {
			background-color: transparent;
			cursor: wait;
		}
		.su-overlay-msg-container {
			display: none;
			padding: 20px;
			background-color: white;
			border: 1px solid black;
			font-size: 1.5em;
		}
		.su-overlay-button-container {
			display: none;
		}
		.su-overlay-button-container .su-button {
			padding: 0.5em;
			flex: 1;
			border-radius: 0px;
		}
		.su-overlay-container {
			position: fixed;
			background-color: transparent;
			display: inline-block;
			border: none;
			left: 0px;
			bottom: 0px;
		}
		.su-overlay-loading {
			border: 5px solid #f3f3f3; /* Light grey */
			border-top: 5px solid #3498db; /* Blue */
			border-radius: 50%;
			width: 20px;
			height: 20px;
			animation: spin 2s linear infinite;
		}
		@keyframes spin {
			0% { transform: rotate(0deg); }
			100% { transform: rotate(360deg); }
		}`
	overlay: false
	curOverlay: false
	New()
		{
		LoadCssStyles('su-overlay.css', .styles)
		.overlay = CreateElement('dialog', SuUI.GetCurrentDocument().body,
			className: 'su-overlay-anchor')
		.overlay.style.display = 'none'
		.container = CreateElement('div', .overlay, className: 'su-overlay-container')
		CreateElement('div', .container, className: 'su-overlay-loading')

		.msgContainer = CreateElement('div', .overlay,
			className: 'su-overlay-msg-container')
		.initButtons()

		.overlayList = Object()
		}

	initButtons()
		{
		LoadCssStyles('su-button.css', ButtonComponent.Styles)
		.buttonContainer = CreateElement('div', .overlay,
			className: 'su-overlay-button-container')
		.cancelEl = el = CreateElement('button', .buttonContainer, className: 'su-button')
		el.innerText = 'Cancel'
		el.AddEventListener('click', .onCancel)
		.okEl = el = CreateElement('button', .buttonContainer, className: 'su-button')
		el.innerText = 'OK'
		el.AddEventListener('click', .onOK)
		}

	Status: #Closed
	Show(id = 'default', msg = false, level = 0,
		okHandler = false, cancelHandler = false)
		{
		if .overlay isnt false and .overlay.open isnt true
			{
			.Status = #Opening
			.overlay.ShowModal()
			.overlay.style.display = ''
			.Status = #Open
			}
		item = Object(:id, msg: '', :level)
		if false isnt i = .overlayList.FindIf({ it.id is id })
			item = .overlayList.Extract(i)

		.overlayList.Add(item)
		.overlayList.Sort!(By(#level))

		if okHandler isnt false
			item.okHandler = okHandler
		if cancelHandler isnt false
			item.cancelHandler = cancelHandler
		if msg isnt false
			item.msg = msg
		.updateOverlay()
		}

	onOK()
		{
		if .overlayList.Empty?()
			return

		item = .overlayList.Last()
		if false is okHandler = item.GetDefault(#okHandler, false)
			return

		okHandler()
		.Close(item.id)
		}

	onCancel()
		{
		if .overlayList.Empty?()
			return

		item = .overlayList.Last()
		if false is cancelHandler = item.GetDefault(#cancelHandler, false)
			return

		cancelHandler()
		.Close(item.id)
		}

	SetMsg(msg, id = 'default')
		{
		if false is i = .overlayList.FindIf({ it.id is id })
			return
		.overlayList[i].msg = msg
		if i is .overlayList.Size() - 1
			.updateOverlay()
		}

	updateOverlay()
		{
		if .overlayList.Empty?()
			{
			.hide()
			return
			}

		item = .overlayList.Last()
		if item.msg is ''
			.hide()
		else
			{
			.msgContainer.style.display = 'block'
			.msgContainer.textContent = item.msg

			el = false
			if item.GetDefault(#cancelHandler, false) isnt false
				{
				el = .cancelEl
				.cancelEl.style.display = ''
				}
			else
				.cancelEl.style.display = 'none'
			if item.GetDefault(#okHandler, false) isnt false
				{
				el = .okEl
				.okEl.style.display = ''
				}
			else
				.okEl.style.display = 'none'

			.buttonContainer.style.display = el isnt false ? 'flex' : ''
			if el isnt false
				{
				el.dataset.highlight = true
				el.Focus()
				}
			}
		}

	hide()
		{
		.msgContainer.style.display = ''
		.buttonContainer.style.display = ''
		}

	Close(id = 'default')
		{
		if false is i = .overlayList.FindIf({ it.id is id })
			return
		.overlayList.Delete(i)
		.updateOverlay()
		if .overlayList.Empty?()
			.close()
		}

	close()
		{
		if .overlay isnt false and .overlay.open is true
			{
			.Status = #Closing
			.overlay.Close()
			.overlay.style.display = 'none'
			.Status = #Closed
			}
		}
	}
