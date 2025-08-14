// Copyright (C) 2018 Axon Development Corporation All rights reserved worldwide.
HtmlContainer
	{
	Name: "Form"
	Left: ""
	Dir: "Horz"
	styles: `
		.su-form-container {
			display: inline-flex;
			flex-direction: row;
			align-items: baseline;
		}
		.su-form-container > div {
			padding-left: 10px;
		}
		.su-form-container > div:first-child {
			padding-left: 0;
		}
		.su-form-empty-line {
			height: 6px;
		}`
	maxGroup: 0
	New(@args)
		{
		LoadCssStyles('form-control.css', .styles)
		.CreateElement('div')
		.SetStyles(Object('display': 'inline-grid'))

		if args.GetDefault('left', '') isnt ''
			.Left = args.left
		.origXmin = .Xmin
		.origYmin = .Ymin
		.ctrls = Object()
		.containers = Object()
		.groups = Object().Set_default(Object())
		.buildControls(args)
		if .ctrls.NotEmpty?()
			.Recalc()
		}

	buildControls(args)
		{
		row = Object()
		group = -1
		for item in args.Values(list:)
			if item in ('nl', 'nL')
				{
				.ctrls.Add('nl')
				if row.Empty?()
					row = CreateElement('div', .El, className: 'su-form-empty-line')
				.containers.Add(row)
				row = Object()
				.maxGroup = Max(.maxGroup, group)
				group = -1
				}
			else // control
				{
				if Object?(item) and item.Member?("group")
					{
					Assert(item.group > group)
					group = item.group
					}
				if not row.Member?(group)
					{
					newGroup = Object(
						el: CreateElement('div', .El, className: 'su-form-container'),
						children: Object())
					row.Add(newGroup, at: group)
					.groups[group].Add(newGroup)
					}
				.TargetEl = row[group].el
				c = .Construct(item)
				.ctrls.Add(c)
				row[group].children.Add(c)
				}
		.maxGroup = Max(.maxGroup, group)
		if row.NotEmpty?()
			.containers.Add(row)
		}

	xsep: 10
	minline: 6	// min line height
	Recalc()
		{
		if .ctrls.Empty?()
			return

		.setupVert()
		.setupHorz()
		.SetStyles(Object(
			'grid-template-columns': 'auto '.Repeat(.maxGroup) $ '1fr',
			'grid-template-rows': 'repeat(' $ .containers.Size() $ ', auto)',
			'align-items': 'baseline',
			'column-gap': .xsep $ 'px'))
		.SetMinSize()
		}

	calcLineHeights()
		{
		// set x positions ignoring groups, determine line heights
		// check group numbers, track group members
		lineheights = Object()
		ymin = 0
		for i in .ctrls.Members()
			if 'nl' is c = .ctrls[i]
				{
				lineheights.Add(ymin)
				ymin = 0
				}
			else // control
				{
				ymin = Max(ymin, c.Ymin)
				if .Stretch?(i, c)
					c.formStretch? = true
				}
		lineheights.Add(ymin)
		return lineheights
		}

	// override by GridComponent
	Stretch?(i, c)
		{
		return .endofline?(i) and c.Xstretch > 0 and .Xstretch > 0
		}

	endofline?(i)
		{
		return i + 1 >= .ctrls.Size() or .ctrls[i + 1] is 'nl'
		}

	setupVert()
		{
		lineheights = .calcLineHeights()
		ymin = 0
		for line in .containers.Members()
			{
			row = .containers[line]
			if not Object?(row)
				{
				ymin += .minline // should be the gap
				row.SetStyle('grid-row', (line + 1) $ ' / ' $ (line + 2))
				continue
				}
			for div in row
				div.el.SetStyle('grid-row', (line + 1) $ ' / ' $ (line + 2))
			ymin += lineheights[line]
			}
		.Ymin = Max(.origYmin, ymin)
		}

	setupHorz()
		{
		.alignLeft()
		.justifyStretch()
		xmin = 0
		offset = .calcOffset()
		for row in .containers
			{
			if not Object?(row)
				{
				row.SetStyle('grid-column', '1 / ' $ (.maxGroup + 2 + offset))
				continue
				}
			rowXmin = 0
			gns = row.Members().Sort!()
			for i in gns.Members()
				{
				div = row[gns[i]]
				for c in div.children
					rowXmin += c.Xmin

				div.el.SetStyle('grid-column', .CalcGridColumn(i, gns, offset))
				if gns[i] is -1
					div.el.SetStyle('z-index', '1')
				}
			xmin = Max(xmin, rowXmin)
			}
		.Xmin = Max(.origXmin, xmin)
		}

	calcOffset()
		{
		for row in .containers
			{
			if Object?(row)
				{
				gns = row.Members().Sort!()
				if gns.Has?(-1) and gns.Has?(0)
					return 1
				}
			}
		return 0
		}

	// override by GridComponent
	CalcGridColumn(i, gns, offset)
		{
		next = i + 1 >= gns.Size()
			? .maxGroup + 2 + offset
			: gns[i] is -1
				? 2
				: gns[i + 1] + 1 + offset
		start = gns[i] is -1 ? 1 : gns[i] + 1 + offset
		return start $ ' / ' $ next
		}

	alignLeft()
		{
		for gn in .groups.Members().Sort!()
			{
			if gn is -1
				continue
			left = 0
			for g in .groups[gn]
				left = Max(left, g.children[0].Left)
			if .Left is ""
				.Left = left

			for g in .groups[gn]
				{
				c = g.children[0]
				c.El.SetStyle('margin-left', left - c.Left $ 'px')
				}
			}
		}

	justifyStretch()
		{
		for gn in .groups.Members().Sort!()
			{
			for g in .groups[gn]
				{
				c = g.children.Last()
				if c.GetDefault(#formStretch?, false) is false
					g.el.SetStyle('justify-self', 'start')
				else
					c.El.SetStyle('flex-grow', 1)
				}
			}
		}

	GetChildren()
		{
		return .ctrls.Filter(Instance?)
		}
	}
