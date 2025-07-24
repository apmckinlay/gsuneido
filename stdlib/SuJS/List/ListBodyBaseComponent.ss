// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
Component
	{
	New()
		{
		LoadCssStyles('list-body-control.css', ListBodyStyles)

		.rowHeight = SuRender().GetTextMetrics(.Parent.El, 'M').height + 4/*=padding*/
		.Parent.El.style.setProperty('--su-row-height', .rowHeight $ 'px')
		}

	GetRowHeight()
		{
		return .rowHeight
		}

	SetCellValue(el, cell)
		{
		switch (cell.type)
			{
		case 'image':
			el.innerHTML = ''
			imageEl = CreateElement('img', el, className: 'su-listbody-cell-image')
			imageEl.SetAttribute(#src, cell.src)
			.setStyles(#('text-align': 'center'), el)
		case 'text':
			.setTextNode(el, cell)
		case 'rect', 'circle':
			el.innerHTML = ''
			shape = CreateElement('div', el, className: 'su-listbody-cell-' $ cell.type)
			.generateShapeStyles(cell, shape)
		case 'multi':
			el.innerHTML = ''
			ratios = cell.GetDefault(#ratios, Object().Set_default(1))
			total = ratios.Size() is 0 ? cell.Size(list:) : ratios.Sum()
			for i in cell.Values(list:).Members()
				{
				part = cell[i]
				part.bkColor = cell.GetDefault(#bkColor, '')
				partEl = CreateElement('div', el, className: 'su-listbody-cell-part')
				.setTextNode(partEl, part)
				partEl.SetStyle('flex', ratios[i] $ ' ' $ ratios[i] $ ' ' $
					(ratios[i] / total).DecimalToPercent() $ '%')
				}
			.setStyles(#(display: 'flex'), el)
			}
		.setTooltip(el, cell)
		}

	setTextNode(el, cell)
		{
		if cell.GetDefault(#html, false) is true
			el.innerHTML = cell.data
		else
			el.innerHTML = cell.data.Replace('[\r\n\t ]', '\&nbsp;')
		.setStyles(.generateTextStyles(cell), el)
		}

	generateTextStyles(cell)
		{
		styles = Object()
		styles['text-align'] = cell.GetDefault(#justify, 'left')
		styles['text-overflow'] = cell.GetDefault(#ellipsis?, false) ? 'ellipsis' : 'clip'
		styles['color'] = cell.GetDefault(#color, '')
		SetFontStyles(cell.GetDefault(#font, false), styles)
		styles['background-color'] = cell.GetDefault(#bkColor, '')
		if cell.GetDefault(#extra, false) isnt false
			styles.Merge(cell.extra)
		return styles
		}

	generateShapeStyles(cell, el)
		{
		styles = Object(
			'border-width': cell.thick $ 'px',
			'border-color': cell.lineColor is false ? 'black' : cell.lineColor)
		if cell.fillColor isnt false
			styles['background-color'] = cell.fillColor
		.setStyles(styles, el)
		}

	setStyles(styles, el)
		{
		el.style = styles.Map2({ |m, v| m $ ':' $ v }).Join(';')
		}

	setTooltip(el, cell)
		{
		if cell.type isnt 'text' or
			cell.GetDefault(#html, false) is true or
			cell.data.Blank?()
			el.SetAttribute('data-tip', '')
		else
			el.SetAttribute('data-tip', cell.GetDefault(#tip, cell.data))
		}

	CreateCellElement(parent, cell, col, className)
		{
		el = CreateElement('div', parent, className)
		.SetCellValue(el, cell)
		.SetCellAttributes(parent, x: col, type: 'cell')
		.SetCellAttributes(el, x: col, type: 'cell')

		el.AddEventListener('mouseenter', {
			|event|
			target = event.target
			tip = target.GetAttribute('data-tip')
			if tip isnt '' and target.offsetWidth < target.scrollWidth
				target.SetAttribute('title', tip)
			else
				target.RemoveAttribute('title')
			})
		return el
		}

	SetCellAttributes(el, x = false, y = false, type = false)
		{
		if x isnt false
			el.SetAttribute('data-x', x)
		if y isnt false
			el.SetAttribute('data-y', y)
		if type isnt false
			el.SetAttribute('data-type', type)
		}
	}
