// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
Component
	{
	Name: 'VirtualListThumbBar'
	styles: `
	.su-vlist-thumbbar {
		position: absolute;
		right: 1.5em;
		bottom: 1.5em;
		z-index: 1;
	}
	.su-vlist-thumbbar-button {
		font-family: suneido;
		font-style: normal;
		font-weight: normal;
		font-size: 1.5em;
		padding: 0.25em;
		border: 1px solid black;
		border-radius: 20%;
		margin: 0.5em;
		user-select: none;
		background-color: white;
	}
	.su-vlist-thumbbar-button:hover {
		background-color: azure;
		border-color: deepskyblue;
	}
	.su-vlist-thumbbar-button:active,
	.su-vlist-thumbbar-button.su-vlist-thumbbar-button-pressed {
		background-color: lightblue;
		border-color: deepskyblue;
	}
	`
	select: false
	New(.disableSelectFilter)
		{
		LoadCssStyles('vlist-thumbbar.css', .styles)
		.CreateElement('div', className: 'su-vlist-thumbbar')

		home = CreateElement('div', .El, className: 'su-vlist-thumbbar-button')
		home.textContent = IconFontHelper.GetCode('arrow_home').Chr()
		home.AddEventListener('click', .onHome)

		end = CreateElement('div', .El, className: 'su-vlist-thumbbar-button')
		end.textContent = IconFontHelper.GetCode('arrow_end').Chr()
		end.AddEventListener('click', .onEnd)

		if .disableSelectFilter isnt true
			{
			.select = CreateElement('div', .El, className: 'su-vlist-thumbbar-button')
			.select.textContent = IconFontHelper.GetCode('zoom').Chr()
			.select.AddEventListener('click', .onSelect)
			}

		.showThreshold = SuRender().GetTextMetrics(home, 'M').width *
			3 /*=font + padding + margin + right*/

		.ParentEl.AddEventListener('mousemove', .mousemove)
		.show(false)
		}

	mousemove(event)
		{
		rect = SuRender.GetClientRect(.ParentEl)
		cursorX = event.clientX - rect.left
		width = rect.width
		.show(Max(0, width - cursorX) < .showThreshold)
		}

	show?: true
	show(show?)
		{
		if .show? is show?
			return
		.show? = show?
		.El.SetStyle('display', show? ? '' : 'none')
		}

	SetSelectPressed(pressed)
		{
		if .select is false
			return
		if pressed is true
			.select.classList.Add('su-vlist-thumbbar-button-pressed')
		else
			.select.classList.Remove('su-vlist-thumbbar-button-pressed')
		}

	onHome()
		{
		.EventWithOverlay('OnHome')
		}

	onEnd()
		{
		.EventWithOverlay('OnEnd')
		}

	onSelect()
		{
		.EventWithOverlay('OnSelect')
		}
	}
