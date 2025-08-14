// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
Component
	{
	styles: `
		.su-listhead th {
			position: sticky;
			top: 0px;
			border: none;
			padding-left: 4px;
			padding-right: 4px;
			background: white;
			z-index: 1;
			text-align: left;
			overflow: hidden;
			white-space: nowrap;
			text-overflow: ellipsis;
			font-weight: normal;
		}
		.su-listhead tr {
			background-color: white;
		}
		.su-splitter {
			top: 0;
			bottom: 0;
			right: 0;
			width: 3px;
			position: absolute;
			cursor: col-resize;
			user-select: none;
			padding-left: 2px;
			border-right: 1px solid lightgrey;
		}
		th[data-type="mark-cell"] {
			background-color: var(--su-color-buttonface);
		}
		.su-header-splitline {
			position: absolute;
			border-left: 2px dashed black;
			top: 0;
			bottom: 0;
		}
		.su-sort-indicator {
			color: #aaa;
			font-size: 0.75em;
			font-style: normal;
			font-weight: normal;
			pointer-events: none;
		}
		.su-sortable:hover {
			background-color: #d9ebf9;
		}
		.su-sortable:active {
			background-color: #bcdcf4;
		}
		.su-drag-to {
			position: absolute;
			border-left: 2px solid darkblue;
			top: 0px;
			bottom: 0px;
			left: 0px;
			display: none;
		}
		.su-drag-shadow {
			position: fixed;
			opacity: .75;
			z-index: 1;
			display: none;
		}
		.su-listhead.su-listhead-button th:hover {
			background-color: azure;
		}
		.su-listhead.su-listhead-button th:active {
			background-color: lightblue;
		}
		`
	New(.parentEl, .buttonStyle = false, .noDragDrop = false, .noHeaderButtons = false,
		.stretch = false, .noHeader = false)
		{
		LoadCssStyles('list-header-control.css', .styles)
		.thead = CreateElement('thead', parentEl,
			className: 'su-listhead' $
				(.buttonStyle is true ? ' su-listhead-button' : ''))

		.thead.AddEventListener('mousedown', .mousedown)
		.thead.AddEventListener('contextmenu', .contextMenu)
		}

	Reset()
		{
		.thead.innerText = ''
		}

	markCol: false
	offset: 0
	headCols: #()
	showSortIndicator: false
	emptyHdr: false
	Update(.headCols, .markCol = false, virtualList? = false, .showSortIndicator = false)
		{
		.Reset()
		headRow = CreateElement('tr', .thead)
		// plus .offset when convert control col to .headCols index
		// virtual list column index doesn't include mark col
		// list control column index include mark col if it has one
		.offset = virtualList?
			? 1
			: .markCol is true
				? 0 : 1
		i = 0
		markHeader = CreateElement('th', headRow)
		.SetStyles(.markCol is true
			? #(width: '1em', padding: '')
			: #(width: '0px', padding: '0px'), markHeader)
		markHeader.SetAttribute('data-x', 0)
		markHeader.SetAttribute('data-type', 'mark-cell')
		.headCols.Add(Object(el: markHeader, field: false), at: 0)
		i++

		stretchCol = .StretchCol()
		for j in i .. .headCols.Size()
			.updateHeaderValues(headRow, j, stretchCol)
		if false is stretchCol
			.addEmptyHdr(headRow)
		}

	updateHeaderValues(headRow, i, stretchCol)
		{
		col = .headCols[i]
		col.el = header = CreateElement('th', headRow)
		.setWidth(header, col.width, i is stretchCol)
		if .noHeader
			{
			header.SetStyle('visibility', 'hidden')
			header.SetStyle('line-height', 0)
			return
			}
		sortIndicator = Opt(
			'<span data-type="header" data-x="' $ i $ '" class="su-sort-indicator">',
			col.sort is 1 ? '&#x25B2;' : col.sort is -1 ? '&#x25BC;' : '',
			'</span>')
		header.innerHTML = sortIndicator $ col.text
		if col.Member?('format')
			header.SetStyle('text-align', col.format)
		header.SetAttribute('data-x', i)
		header.SetAttribute('data-type', 'header')
		if .showSortIndicator and not .noHeaderButtons
			header.className = 'su-sortable'
		.AddToolTip(col.tip isnt false ? col.tip : col.text, header)
		.addSplitter(header, i, stretchCol)
		}

	setWidth(header, width, stretch?)
		{
		header.SetStyle('width', stretch?
			? '100%'
			: String?(width) and width.Prefix?('calc')
				? width
				: width $ 'px')
		header.SetStyle('display', width is 0 ? 'none' : '')
		}

	StretchCol()
		{
		if .stretch is false
			return false
		if .stretch is true
			return .headCols.Size() - 1
		return .headCols.FindIf({ it.field is .stretch })
		}

	SetStetchCol(.stretch) { }

	addSplitter(header, i, stretchCol)
		{
		splitter = CreateElement('div', header)
		splitter.className = 'su-splitter'
		splitter.SetStyle('display',
			stretchCol isnt false and i is .headCols.Size() - 1 ? 'none' : '')
		splitter.SetAttribute('data-x', i)
		splitter.SetAttribute('data-type', 'splitter')
		splitter.AddEventListener('dblclick', .factory(i))
		}

	addEmptyHdr(headRow)
		{
		// empty col for filling remaining space
		.emptyHdr = el = CreateElement('th', headRow)
		el.SetStyle('width', '100%')
		el.SetAttribute('data-x', .GetColsNum())
		el.SetAttribute('data-type', 'empty-header')
		}

	factory(i)
		{
		return { .dblclick(i) }
		}

	dblclick(i)
		{
		if not .inList?()
			return
		width = .Parent.MeasureWidth(i)
		.Parent.EventWithOverlay('HeaderDividerDoubleClick', .ToControlColIndex(i), width)
		}

	inList?()
		{
		return .Parent.Base?(ListComponent) or .Parent.Base?(VirtualListGridComponent)
		}

	SetMaxWidth(i)
		{
		.SetColWidth(i, .Parent.MeasureWidth(.ToWebColIndex(i)))
		}

	GetField(col)
		{
		return .headCols[col].field
		}

	ForEachHeadCol(block)
		{
		for (i = 1; i < .headCols.Size(); i++)
			block(col: i, field: .headCols[i].field, width: .headCols[i].width)
		}

	HasMarkCol?()
		{
		return .markCol
		}

	curCol: false
	curX: false
	splitLine: false
	splitPos: false
	dragCol: false
	preDragCol: false
	mousedown(event)
		{
		if event.button isnt 0
			return
		target = event.target

try
	{
		type = target.GetAttribute('data-type')
		if type is 'empty-header'
			return
	}
catch (e)
	{
	SuRender().Event(false, 'SuneidoLog', Object(
		'ERROR: (CAUGHT) ' $ e,
		params: [target: Display(target), class: target.className,
			innerHTML: target.innerHTML,
			outerHTML: target.outerHTML], caughtMsg: 'for debugging 33767'))
	return
	}

		if target.GetAttribute('data-type') is 'splitter'
			{
			.curCol = Number(target.GetAttribute('data-x'))
			.curX = event.x
			.splitLine = CreateElement('div', .parentEl, className: 'su-header-splitline')
			containerRect = .parentEl.GetBoundingClientRect()
			.splitLine.SetStyle('left', (.splitPos = event.x - containerRect.x) $ 'px')
			.StartMouseTracking(.splitterMouseup, .splitterMousemove)
			}
		else if .noDragDrop isnt true
			{
			dragCol = Number(target.GetAttribute('data-x'))
			if dragCol is 0
				return

			.preDragCol = dragCol

			.dragToLine = CreateElement(
				'div', .headCols[dragCol].el, className: 'su-drag-to')

			dragColEl = .headCols[.preDragCol].el
			rect = dragColEl.GetBoundingClientRect()
			.moving = CreateElement('div', .thead, className: 'su-drag-shadow')
			.moving.AppendChild(dragColEl.CloneNode(true))
			.dragStart = [x: event.x, y: event.y]
			.draggingOffset = event.x - rect.left
			.draggingTop = rect.top
			.StartMouseTracking(.dragColMouseup, .dragColMousemove)
			}
		else
			.StartMouseTracking(.sortMouseUp)
		}

	dragToLine: false
	moving: false
	draggingOffset: false
	dragStart: false
	dragColMousemove(event)
		{
		if .dragCol is false
			{
			// mouse not moved
			if event.x is .dragStart.x and event.y is .dragStart.y
				return
			.dragCol = .preDragCol
			}
		oldCol = .dragCol

		.moving.SetStyle('display', 'block')
		.moving.SetStyle('left', (event.x - .draggingOffset) $ 'px')
		.moving.SetStyle('top', .draggingTop $ 'px')

		if false is i = .getNewIndex(oldCol, event)
			return
		if i is 0
			return
		.dragToLine.SetStyle('display', 'block')
		cellWithPipe = .headCols.Member?(i) ? .headCols[i].el : .emptyHdr
		if cellWithPipe isnt false
			cellWithPipe.AppendChild(.dragToLine)
		}

	dragColMouseup(event)
		{
		if .dragCol is false
			{
			.sortMouseUp(event)
			.clearDrag()
			return
			}
		oldCol = .dragCol
		if false isnt i = .getNewIndex(oldCol, event)
			{
			from = .ToControlColIndex(oldCol)
			to = Max(.ToControlColIndex(i), 0)
			if to > from
				--to
			.Parent.EventWithOverlay('HeaderReorder', from, to)
			}
		.clearDrag()
		}

	clearDrag()
		{
		if .dragToLine isnt false
			{
			.dragToLine.Remove()
			.dragToLine = false
			}
		if .moving isnt false
			{
			.moving.Remove()
			.moving = false
			}
		.dragCol = false
		.preDragCol = false
		.dragStart = false
		.StopMouseTracking()
		}

	getNewIndex(oldCol, event)
		{
		i = 0
		for (col in .headCols)
			{
			rect = col.el.GetBoundingClientRect()
			if oldCol is i and event.x > rect.left and event.x < rect.right
				return false
			if event.x < rect.right
				{
				if event.x > rect.left + rect.width / 2
					{
					++i
					}
				break
				}
			i++
			}
		return i
		}

	splitterMousemove(event)
		{
		containerRect = .parentEl.GetBoundingClientRect()
		if event.x - containerRect.x >= 0
			{
			.splitLine.SetStyle('left', (.splitPos = event.x - containerRect.x) $ 'px')
			}
		}

	splitterMouseup(event)
		{
		col = .curCol
		stretchCol = .StretchCol()
		// resize the next col if the resized col is or after the stretch col
		if stretchCol isnt false and col >= stretchCol
			.Parent.Event('HeaderResize', .ToControlColIndex(col + 1),
				Max(0, .headCols[col + 1].width - event.x + .curX))
		else
			.Parent.Event('HeaderResize', .ToControlColIndex(col),
				Max(0, .headCols[col].width + event.x - .curX))
		.curCol = false
		.curX = false
		.splitLine.Remove()
		.splitLine = .splitPos = false
		.StopMouseTracking()
		}

	sortMouseUp(event)
		{
		.StopMouseTracking()
		if .skipMouseUp?(event)
			return

		target = event.target
		/*
		https://www.w3schools.com/jsref/prop_node_nodetype.asp
		 * If the node is an element node, the nodeType property will return 1.
		 * If the node is an attribute node, the nodeType property will return 2.
		 * If the node is a text node, the nodeType property will return 3.
		 * If the node is a comment node, the nodeType property will return 8.
		*/
		if target.nodeType isnt 1 or not .thead.Contains(target) or
			target.GetAttribute('data-type') in ('splitter', 'empty-header')
			return

		col = Number(target.GetAttribute('data-x'))
		if col is 0
			return
		.Parent.EventWithOverlay('HeaderClick', col: .ToControlColIndex(col))
		}

	skipMouseUp?(event)
		{
		if event.button isnt 0 or .Destroyed?() or
			.dragCol isnt false or .curCol isnt false
			return true

		return .noHeaderButtons is true or
			.buttonStyle is false and .showSortIndicator is false
		}

	SetColWidth(col, width)
		{
		headCol = .ToWebColIndex(col)
		if headCol is 0 // mark column
			return
		.headCols[headCol].width = width
		.setWidth(.headCols[headCol].el, width, .StretchCol() is headCol)
		if width is 0 and .inList?()
			.Parent.HideCol(headCol)
		}

	GetColsNum()
		{
		return .headCols.Size()
		}

	GetOffsetHeight()
		{
		return .thead.offsetHeight
		}

	GetColRect(col)
		{
		col = .ToWebColIndex(col)
		return Object(
			left: .headCols[col].el.offsetLeft,
			width: .headCols[col].el.offsetWidth)
		}

	contextMenu(event)
		{
		col = .ToControlColIndex(Number(event.target.GetAttribute('data-x')))
		.Parent.RunWhenNotFrozen()
			{
			.Parent.EventWithOverlay('CONTEXTMENU_HEADER', event.clientX, event.clientY,
				:col)
			}
		event.StopPropagation()
		event.PreventDefault()
		}

	ToWebColIndex(index)
		{
		return index + .offset
		}

	ToControlColIndex(index)
		{
		return index - .offset
		}
	}
