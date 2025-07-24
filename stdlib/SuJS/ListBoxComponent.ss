// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
Component
	{
	Name:		"ListBox"
	Xmin: 		100
	Ymin: 		100
	Xstretch: 	1
	Ystretch: 	1
	styles: `
		.su-listbox-container {
			position: relative;
			overflow: auto;
			border: solid 1px black;
		}
		.su-listbox-container:focus {
			outline: none;
		}
		.su-listbox {
			position: absolute;
			top: 0px;
			left: 0px;
			width: 100%;
			height: 100%;
			display: inline-flex;
			flex-direction: column;
			user-select: none;
		}`

	New(.multicolumn = false, font = '', size = '', weight = '')
		{
		LoadCssStyles('su-listbox.css', .styles)
		.items = Object()
		.CreateElement('div', className: 'su-listbox-container')
		.El.SetAttribute('translate', 'no')
		.SetFont(font, size, weight)

		.buildList()
		.SetMinSize()
		}

	selected: false
	buildList()
		{
		.list = CreateElement('div', .El, className: 'su-listbox')
		if .multicolumn is true
			.list.SetStyle('flex-wrap', 'wrap')

		.El.tabIndex = "0"
		.El.AddEventListener('keydown', .KEYDOWN)
		}

	InsertItem(s, i)
		{
		el = CreateElement('div', .list, at: i)
		el.innerText = s
		el.SetStyle('text-decoration', 'none')
		el.SetStyle('white-space', 'nowrap')
		el.AddEventListener('click', .selectFactory(el, .SELCHANGE))
		el.AddEventListener('dblclick', .selectFactory(el, .LBN_DBLCLK))
		el.AddEventListener('contextmenu', .selectFactory(el, .contextMenu))
		.items.Add(el at: i)
		return el
		}

	selectFactory(item, fn)
		{
		return { |event| fn(item, :event) }
		}

	SELCHANGE(el)
		{
		.Select(el)
		.Event('SELCHANGE', .items.Find(el))
		}

	Select(el = false)
		{
		.clearSelect()
		if el isnt false
			.select(el)
		}

	LBN_DBLCLK(el)
		{
		.Event('LBN_DBLCLK', .items.Find(el))
		}

	clearSelect()
		{
		if .selected is false
			return

		.selected.SetStyle('background-color', '')
		.selected.SetStyle('color', '')
		.selected = false
		}

	select(el)
		{
		.selected = el
		.selected.SetStyle('background-color', '#0076d7')
		.selected.SetStyle('color', 'white')
		.scrollRowToView(el)
		}

	contextMenu(el, event)
		{
		if false is i = .items.Find(el)
			return
		.EventWithOverlay('CONTEXTMENU', event.clientX, event.clientY, :i)
		}

	DeleteItem(i)
		{
		el = .items[i]
		if el is .selected
			.clearSelect()
		.items.Delete(i)
		el.Remove()
		}

	SetCurSel(i)
		{
		.clearSelect()
		if i isnt -1
			.select(.items[i])
		}

	scrollRowToView(rowEl)
		{
		rowHeight = rowEl.offsetHeight
		rowOffsetTop = rowEl.offsetTop
		scrollTop = .El.scrollTop
		scrollHeight = .El.clientHeight
		if scrollTop > rowOffsetTop
			.El.scrollTop = rowOffsetTop
		else if scrollTop + scrollHeight < rowOffsetTop + rowHeight
			.El.scrollTop = rowOffsetTop + rowHeight - scrollHeight
		}

	KEYDOWN(event)
		{
		if .items.Empty?()
			return
		selectChanged = true

		if event.key in ('ArrowUp','ArrowDown')
			.processArrowKeys(event.key)
		else if .alphaNumeric?(event.key)
			.findItemPos(event.key.Lower())
		else
			selectChanged = false

		if selectChanged is true
			{
			event.PreventDefault()
			.Event('SELCHANGE', .items.Find(.selected))
			}
		}

	processArrowKeys(key)
		{
		if key is 'ArrowUp'
			{
			if .selected is false
				.select(.items.Last())
			else
				{
				i = .items.Find(.selected)
				.clearSelect()
				.select(.items[Max(0, i - 1)])
				}
			}
		else if key is 'ArrowDown'
			{
			if .selected is false
				.select(.items.First())
			else
				{
				i = .items.Find(.selected)
				.clearSelect()
				.select(.items[Min(.items.Size() - 1, i + 1)])
				}
			}
		}

	alphaNumeric?(key)
		{
		return String?(key) and key.Size() is 1 and key =~ `[[:alnum:]]`
		}

	findItemPos(key)
		{
		startText = false
		foundStart = false
		if .selected isnt false
			startText = .selected.innerText
		else
			foundStart = true

		found = false
		firstFound = false
		for item in .items
			{
			text = item.innerText.Lower()
			if firstFound is false and text.Prefix?(key[0])
				firstFound = item
			if foundStart
				if text.Prefix?(key[0])
					{
					.clearSelect()
					.select(item)
					found = true
					return
					}
			if item.innerText is startText
				foundStart = true
			}
		.selectFirstFound(found, firstFound)
		}

	selectFirstFound(found, firstFound)
		{
		if found is false and firstFound isnt false
			{
			.clearSelect()
			.select(firstFound)
			}
		}
	}