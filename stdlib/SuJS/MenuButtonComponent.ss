// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
ButtonComponent
	{
	styles: `
		.su-menubutton-caret {
			position: absolute;
			right: 2px;
			bottom: 1px;
			width: 0;
			height: 0;
			display: inline-block;
			vertical-align: middle;
			margin-left: 2px;
			border-top: 4px dashed;
			border-left: 4px solid transparent;
			border-right: 4px solid transparent;
		}
		`
	disabled: false
	New(text, tip = false, tabover = false, left = false, width = false)
		{
		super(text, pad: 36, :tip, :tabover, :width)
		LoadCssStyles('menu-button.css', .styles)
		CreateElement('span', .El, className: 'su-menubutton-caret')
		if left is true
			.SetStyles(#('text-align': 'left', 'padding': '0px 18px'))
		.El.AddEventListener('keydown', .keydown)
		}

	keydown(event)
		{
		if .disabled is true
			return
		.Keydown(event)
		}

	Keydown(event)
		{
		if event.key in (' ', 'ArrowDown')
			{
			r = SuRender.GetClientRect(.El)
			rcExclude = Object(left: -9999, right: 9999, top: r.top, bottom: r.bottom)
			.pullDown(r, rcExclude)
			}
		}

	pullDown(r, rcExclude)
		{
		.RunWhenNotFrozen()
			{
			.EventWithFreeze('MenuButton_PullDown', r.left, r.bottom, rcExclude, r)
			}
		}

	CLICKED()
		{
		if .disabled is true
			return

		r = SuRender.GetClientRect(.El)
		rcExclude = Object(left: -9999, right: 9999, top: r.top, bottom: r.bottom)
		.pullDown(r, rcExclude)
		}

	Set(text)
		{
		super.Set(text)
		CreateElement('span', .El, className: 'su-menubutton-caret')
		}

	Disable(.disabled)
		{
		.El.SetAttribute('data-disable', disabled is true ? 'true' : '')
		}
	// TODO: handle display
	grayed: false
	Grayed(.grayed)
		{
		}
	}
