// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
Component
	{
	Name: 'ChooseListBox'
	NumLines: 10
	styles: `
		.su-chooselist {
			background-color: white;
			border: 1px solid black;
			border-radius: 0.3em;
			user-select: none;
			box-shadow: 3px 3px 3px grey;
		}
		.su-chooselist:focus {
			outline: none;
		}
		.su-chooselist div {
			text-decoration: none;
		}
		.su-chooselist-selected {
			background-color: lightblue;
		}`

	New(.list, select, .listSeparator, fieldHwnd)
		{
		LoadCssStyles('su-choose-list.css', .styles)
		.selected = select is false ? 0 : select
		.buildList()

		field = SuRender().GetRegisteredComponent(fieldHwnd)
		.SetStyles(Object(
			'max-height': .NumLines $ 'em',
			'min-width': field.Xmin $ 'px',
			overflow: 'auto'))
		}

	Startup()
		{
		.items[.selected].ScrollIntoView()
		}

	buildList()
		{
		.CreateElement('div', className: 'su-chooselist')
		.El.SetAttribute('translate', 'no')
		.items = Object()
		for (i = 0; i < .list.Size(); i++)
			{
			text = String(.list[i])
			el = CreateElement('div', .El,
				className: .selected is i ? 'su-chooselist-selected' : false)
			el.innerText = text
			el.AddEventListener('click', .selectFactory(i, .selectItem))
			el.AddEventListener('mouseover', .selectFactory(i, .hover))
			.items.Add(el)
			}
		.El.tabIndex = "0"
		.El.AddEventListener('keydown', .KEYDOWN)
		}

	selectFactory(item, fn)
		{
		return { fn(item) }
		}

	selectItem(i)
		{
		.Result(i)
		}

	hover(i)
		{
		if .items is false
			return
		.items[.selected].className = ''
		.selected = i
		.items[.selected].className = 'su-chooselist-selected'
		}

	KEYDOWN(event)
		{
		switch (event.key)
			{
		case 'ArrowUp':
			.move(-1)
		case 'ArrowDown':
			.move(1)
		case 'Escape':
			.Result(false)
		case 'Enter':
			.Result(.selected)
		default:
			}
		event.PreventDefault()
		event.StopPropagation()
		}

	move(offset)
		{
		if .items is false
			return
		pos = Max(0, Min(.items.Size() - 1, .selected + offset))
		.items[.selected].className = ''
		.selected = pos
		.items[.selected].className = 'su-chooselist-selected'
		.items[.selected].ScrollIntoView(offset < 0)
		}

	Result(i)
		{
		item = false
		if i isnt false
			{
			item = String(.list[i])
			if (.listSeparator isnt '')
				item = item.BeforeFirst(.listSeparator)
			}
		.Event(#CHOOSE, item)
		}
	}
