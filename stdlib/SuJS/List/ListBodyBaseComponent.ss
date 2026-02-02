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
			.setStyles(.generateShapeStyles(cell), shape)
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

	generateShapeStyles(cell)
		{
		styles = Object(
			'border-width': cell.thick $ 'px',
			'border-color': cell.lineColor is false ? 'black' : cell.lineColor)
		if cell.fillColor isnt false
			styles['background-color'] = cell.fillColor
		return styles
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
			el.SetAttribute('data-tip', XmlEntityDecode(cell.GetDefault(#tip, cell.data)))
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

	BuildRowContent(rec, header, cellClassName)
		{
		s = '<td data-x=0 data-type="mark-cell"></td>'
		header.ForEachHeadCol()
			{ |col, field, width|
			cell = rec[field]
			s $= '<td data-x=' $ col $ ' data-type="cell" ' $
				(width is 0 ? 'style="display: none"' : '') $'>'
			s $= '<div class="' $ cellClassName $ '" data-x=' $ col $ ' data-type="cell"'

			tip = cell.type isnt 'text' or
				cell.GetDefault(#html, false) is true or
				cell.data.Blank?()
				? ''
				: cell.GetDefault(#tip, cell.data)
			s $= ' data-tip="' $ tip $ '"'
			s $= .buildCellValue(cell)
			s $= '</div></td>'
			}
		// empty col for filling remaining space
		// invisible character to take vertical space when the row is empty
		s $= '<td data-x=' $ header.GetColsNum() $ ' data-type="empty-cell">&#8205;</td>'
		return s
		}

	buildCellValue(cell)
		{
		s = ''
		switch (cell.type)
			{
		case 'image':
			s $= '><img class="su-listbody-cell-image"
				style="text-align: center"
				src="' $ Url.Encode(cell.src) $ '" />'
		case 'text':
			s $= .buildTextNode(cell)
		case 'rect', 'circle':
			styles = .generateShapeStyles(cell)
			s $= '><div class="su-listbody-cell-' $ cell.type $ '"
				style="' $ styles.Map2({ |m, v| m $ ':' $ v }).Join(';') $ '" />'
		case 'multi':
			s $= ' style="display: flex">'
			ratios = cell.GetDefault(#ratios, Object().Set_default(1))
			total = ratios.Size() is 0 ? cell.Size(list:) : ratios.Sum()
			for i in cell.Values(list:).Members()
				{
				part = cell[i]
				part.bkColor = cell.GetDefault(#bkColor, '')
				extraStyles = Object(flex:
					ratios[i] $ ' ' $ ratios[i] $ ' ' $
						(ratios[i] / total).DecimalToPercent() $ '%')
				s $= '<div class="su-listbody-cell-part"' $
					.buildTextNode(part, :extraStyles)  $ '</div>'
				}
			}
		return s
		}

	buildTextNode(cell, extraStyles = #())
		{
		styles = .generateTextStyles(cell).Merge(extraStyles)
		s = ' style="' $
			styles.Map2({ |m, v| Opt(m, ':', v) }).Filter({ it isnt '' }).Join(';') $ '">'
		s $= cell.GetDefault(#html, false) is true
			? cell.data
			: cell.data.Replace('[\r\n\t ]', '\&nbsp;')
		return s
		}

	AddTipListener(row, cellClassName)
		{
		cells = row.QuerySelectorAll('.' $ cellClassName)
		for (i = 0; i < cells.length; i++)
			cells.Item(i).AddEventListener('mouseenter', {
				|event|
				target = event.target
				tip = target.GetAttribute('data-tip')
				if tip isnt '' and target.offsetWidth < target.scrollWidth
					target.SetAttribute('title', tip)
				else
					target.RemoveAttribute('title')
				})
		}
	}
