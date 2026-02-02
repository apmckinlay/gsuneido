// Copyright (C) 2025 Axon Development Corporation All rights reserved worldwide.
Component
	{
	Name: OverviewBar
	Xmin: 7
	Ystretch: 1

	Styles: `
		.su-overview-bar {
			position: relative;
			overflow: hidden;
			background-color: white;
		}
		.su-overview-bar-item {
			position: absolute;
			width: 100%;
		}`
	New(.priorityColor = false, .partnerCtrlHwnd = false)
		{
		LoadCssStyles('su-overview-bar.css', .Styles)
		.CreateElement('div', className: 'su-overview-bar')
		.SetMinSize()
		.marks = Object()
		}

	Startup()
		{
		if .partnerCtrlHwnd isnt false
			.partner = SuRender().GetRegisteredComponent(.partnerCtrlHwnd)
		else
			{
			siblings = .Parent.GetChildren()
			i = siblings.Find(this)
			.partner = siblings[i - 1]
			}
		.partnerEl = .partner.Member?(#CMEl) ? .partner.CMEl : .partner.El
		if .partner.Member?(#PaddingTop)
			.El.SetStyle('margin-top', .partner.PaddingTop $ 'px')

		.resizeObserver = SuUI.MakeWebObject('ResizeObserver', .onResize)
		.resizeObserver.Observe(.partnerEl)
		}

	partner: false
	hasVerticalScroll: ''
	hasHorizontalScroll: ''
	onResize(@unused)
		{
		if .partner is false
			return
		dimension = .partner.GetDimension()
		hasVerticalScroll = dimension.scrollHeight > dimension.clientHeight
		hasHorizontalScroll = dimension.scrollWidth > dimension.clientWidth
		if ((hasVerticalScroll isnt .hasVerticalScroll) or
			(hasHorizontalScroll isnt .hasHorizontalScroll))
			{
			scrollbarWidth = SuRender().GetScrollbarWidth()
			bottom = (hasVerticalScroll ? scrollbarWidth : 0) +
				(hasHorizontalScroll ? scrollbarWidth : 0)
			.El.SetStyle('margin-bottom', bottom $ 'px')

			top = (hasVerticalScroll ? scrollbarWidth : 0) +
				(.partner.Member?(#PaddingTop) ? .partner.PaddingTop : 0)
			.El.SetStyle('margin-top', top $ 'px')

			.hasVerticalScroll = hasVerticalScroll
			.hasHorizontalScroll = hasHorizontalScroll
			}
		}

	numRows: 0
	SetNumRows(.numRows) {}

	SetMaxRowHeight(hwnd, method)
		{
		that = SuRender().GetRegisteredComponent(hwnd)
		.maxRowHeight = (that[method])()
		}

	AddMark(row, color)
		{
		row++
		.marks[row] = el = CreateElement('div', .El, className: 'su-overview-bar-item')
		.SetStyles([
			'background-color': ToCssColor(color),
			'height': 'max(3px, ' $
				'min(' $ .maxRowHeight $ 'px, calc(100% / ' $ .numRows $ ')))',
			'top': 'min(calc(' $ .maxRowHeight $ 'px * ' $ row $ '), ' $
				'calc(100% * ' $ row $ ' / ' $ .numRows $ '))'], :el)
		}

	RemoveMark(row)
		{
		if not .marks.Member?(row)
			return

		.marks[row].Remove()
		.marks.Delete(row)
		}

	ClearMarks()
		{
		.El.InnerHTML = ''
		.marks = Object()
		}

	Destroy()
		{
		.resizeObserver.Unobserve(.El)
		.resizeObserver = false
		super.Destroy()
		}
	}
