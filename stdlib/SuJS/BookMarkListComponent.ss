// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
Component
	{
	Name:		"BookMarkList"
	Xstretch:		1
	Ystretch:		1

	styles: `
		.su-bookmark-list {
			flex-basis: 0;
			overflow: auto;
		}
		.su-bookmark-item {
			position: 		relative;
			padding-left: 	3px;
			padding-right: 	3px;
			margin-left: 	3px;
			cursor: 		pointer;
			user-select: 	none;
			white-space: 	nowrap;
			overflow: 		hidden;
			text-overflow: 	ellipsis;
			text-align: 	end;
		}
		.su-bookmark-item.su-bookmark-item-active {
			font-weight: 	bold;
			margin-left: 	0px;
		}
		.su-bookmark-color-input {
			position: 		absolute;
			left: 			0px;
			visibility: 	hidden;
			height:			0px;
		}`
	New()
		{
		LoadCssStyles('su-bookmark-list.css', .styles)
		.CreateElement('div', className: 'su-bookmark-list')
		.marks = Object()
		.El.AddEventListener('dblclick', .doubleClick)
		.El.AddEventListener('dragover', .dragover)
		}

	doubleClick(event)
		{
		if event.target isnt .El
			return
		.Event('DoubleClick')
		}

	SetState(marks)
		{
		.El.innerText = ''
		.marks.Delete(all:)
		for mark in marks
			.AddMark(mark)
		}

	handlerFactory(mark, handler)
		{
		return { |event| handler(event, mark) }
		}

	rightclick(event, mark)
		{
		event.PreventDefault()
		event.StopPropagation()
		mark.colorPicker.Click()
		}

	active: false
	click(event, mark)
		{
		if event.button isnt 0 or event.target is mark.colorPicker
			return
		.EventWithOverlay(#Click, mark.path)
		}

	colorInput(event, mark)
		{
		mark.el.SetStyle('background-color', event.target.value)
		}

	colorChange(event, mark)
		{
		mark.color = ToCssColor.Reverse(event.target.value)
		.Event(#UpdateColor, mark.Project(#color, #path))
		}

	RemoveMark(path)
		{
		if false is i = .findMark(path)
			return
		if .active is .marks[i]
			.active = false
		.marks[i].el.Remove()
		.marks.Delete(i)
		}

	AddMark(mark)
		{
		color = ToCssColor(mark.color)
		el = CreateElement('div', .El, className: 'su-bookmark-item')
		el.title = el.innerHTML = mark.path.RemovePrefix('/') $ '/'
		el.dir = 'rtl'
		el.SetStyle('background-color', color)
		el.AddEventListener(#contextmenu, .handlerFactory(mark, .rightclick))
		el.AddEventListener(#click, .handlerFactory(mark, .click))
		el.draggable = true
		el.AddEventListener(#dragstart, .handlerFactory(mark, .dragstart))
		el.AddEventListener(#dragend, .handlerFactory(mark, .dragend))
		mark.el = el

		picker = CreateElement('input', mark.el, className: 'su-bookmark-color-input')
		picker.type = 'color'
		picker.value = color
		picker.AddEventListener(#input, .handlerFactory(mark, .colorInput))
		picker.AddEventListener(#change, .handlerFactory(mark, .colorChange))
		mark.colorPicker = picker

		.marks.Add(mark)
		}

	dragging: false
	dragstart(event/*unused*/, mark)
		{
		.dragging = mark
		}

	dragend(event/*unused*/, mark/*unused*/)
		{
		.dragging = false
		.Event(#UpdateMarks, .marks.Map({ it.Project(#path, #color) }))
		}

	dragover(event)
		{
		if .dragging is false
			{
			event.dataTransfer.dropEffect = "none";
			return
			}
		event.PreventDefault()
		event.StopPropagation();
		event.dataTransfer.dropEffect = "move";

		dy = event.clientY - .El.GetBoundingClientRect().top
		markHeight = .dragging.el.GetBoundingClientRect().height
		i = Min((dy / markHeight).Round(0), .marks.Size())
		cur = .findMark(.dragging.path)

		if i is cur or i is cur + 1
			return
		if i < .marks.Size()
			.El.InsertBefore(.dragging.el, .marks[i].el)
		else
			.El.AppendChild(.dragging.el)
		if i < cur
			{
			for (j = cur - 1; j >= i; j--)
				.marks[j + 1] = .marks[j]
			.marks[i] = .dragging
			}
		else
			{
			for (j = cur + 1; j < i; j++)
				.marks[j - 1] = .marks[j]
			.marks[i - 1] = .dragging
			}
		}

	GotoPath(path)
		{
		i = .findMark(path)
		.updateActive(.marks.GetDefault(i, false))
		}

	findMark(path)
		{
		return .marks.FindIf({ it.path is path })
		}

	updateActive(mark)
		{
		if .active is mark
			return
		if .active isnt false
			.active.el.classList.Remove('su-bookmark-item-active')
		if false isnt .active = mark
			.active.el.classList.Add('su-bookmark-item-active')
		}
	}
