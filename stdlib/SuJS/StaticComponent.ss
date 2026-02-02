// Copyright (C) 2018 Axon Development Corporation All rights reserved worldwide.
Component
	{
	Name: 'Static'
	SkipSetFocus: true
	ContextMenu: true
	New(text = "", .font = "", .size = "", .weight = "", justify = "LEFT",
		.underline = false, color = "", tip = false, tabstop = false,
		bgndcolor = "", hidden = false, textStyle = false)
		{
		.untranslated = text
		.CreateElement('div')
		.SetHidden(hidden)
		.updateAlignment()

		.SetStyles(.buildStyles(justify, color, bgndcolor, textStyle))

		.Orig_xmin = .Xmin
		.Orig_ymin = .Ymin

		.SetFont(.font, .size, .weight, .underline)
		.El.innerHTML = .ConvertText(text)
		.Recalc()

		if tabstop or .inListEdit?()
			.El.tabIndex = "0"
		.AddToolTip(tip)
		.El.AddEventListener('mousedown', .mousedown)
		}

	buildStyles(justify, color, bgndcolor, textStyle)
		{
		styles = Object('text-align': justify, 'overflow': 'hidden', 'line-height': 1)
		if textStyle isnt false and StaticTextStyles.Member?(textStyle)
			{
			.size = .size is '' ? StaticTextStyles[textStyle].size : .size
			.weight = .weight is '' ? StaticTextStyles[textStyle].weight : .weight
			color = color isnt '' ? color : StaticTextStyles[textStyle].color
			}
		if color isnt ""
			styles.color = ToCssColor(TranslateColor(color))
		if bgndcolor not in ("", "none")
			styles['background-color'] = ToCssColor(bgndcolor)
		return styles
		}

	ConvertText(text)
		{
		return text.Lines().Map(XmlEntityEncode).Join('<br>').Replace(' ', '\&nbsp')
		}

	Recalc()
		{
		metrics = SuRender().GetTextMetrics(.El, .untranslated)
		.Xmin = .Orig_xmin isnt 0 ? .Orig_xmin : metrics.width
		.Ymin = .Orig_ymin isnt 0 ? .Orig_ymin : metrics.height
		.SetMinSize()
		if .Orig_xmin isnt 0
			.El.SetStyle('width', .Orig_xmin $ 'px')
		if .Orig_ymin isnt 0
			.El.SetStyle('height', .Orig_ymin $ 'px')
		}

	CalcXminByControls(@args)
		{
		.Orig_xmin = .DoCalcXminByControls(@args)
		.Recalc()
		}

	SetColor(color)
		{
		.El.SetStyle('color', ToCssColor(TranslateColor(color)))
		}

	Get()
		{
		return .untranslated
		}

	Set(text, logfont = false, refreshRequired? = false)
		{
		text = String(text)
		if logfont isnt false
			{
			.font = logfont.lfFaceName
			.size = logfont.fontPtSize
			.weight = logfont.GetDefault(#lfWeight, FW.NORMAL)
			.SetFont(.font, .size, .weight, .underline)
			refreshRequired? = true
			}
		if .untranslated is text and refreshRequired? is false
			return
		.untranslated = text
		.updateAlignment()
		.El.innerHTML = .ConvertText(text)
		.calcX()
		.WindowRefresh()
		}

	calcX()
		{
		if .Orig_xmin isnt 0
			return
		metrics = SuRender().GetTextMetrics(.El, .untranslated)
		.Xmin = metrics.width
		}

	updateAlignment()
		{
		// empty div's baseline is its bottom
		if .untranslated.Blank?()
			.El.SetStyle('align-self', 'flex-start')
		else
			.El.SetStyle('align-self', '')
		}

	BestFit(text, xmin)
		{
		lines = Object()
		str = text
		while false isnt fitText = TextBestFit(xmin, str, .measure)
			{
			lines.Add(fitText)
			if "" is str = str[fitText.Size() ..]
				break
			}
		return lines
		}
	measure(str)
		{
		return SuRender().GetTextMetrics(.El, str).width
		}

	inListEdit?()
		{
		return .Controller.Base?(ListEditWindowComponent) or
			(.Parent.Member?("Parent") and .Parent.Parent.Base?(ListEditWindowComponent))
		}

	mousedown()
		{
		.RunWhenNotFrozen({ .EventWithFreeze('LBUTTONDOWN') })
		}

	SelectAll()
		{
		SuClearFocus()
		range = SuUI.GetCurrentDocument().createRange()
		range.SelectNodeContents(.El)
		selection = SuUI.GetCurrentWindow().GetSelection()
		selection.RemoveAllRanges()
		selection.AddRange(range)
		}

	On_Copy()
		{
		.selectAllIfNoSelection()
		SuUI.GetCurrentDocument().ExecCommand('copy')
		}

	selectAllIfNoSelection()
		{
		try
			{
			selection = SuUI.GetCurrentWindow().GetSelection()
			// selection's start and end points are at the same position
			if selection.isCollapsed is true or selection.rangeCount is 0
				{
				.SelectAll()
				return
				}
			curRange = selection.GetRangeAt(0)
			nodeRange = SuUI.GetCurrentDocument().createRange()
			nodeRange.SelectNodeContents(.El)

			// no overlapping
			if curRange.CompareBoundaryPoints(3/*=Range.END_TO_START*/, nodeRange) >= 0 or
				curRange.CompareBoundaryPoints(1/*=Range.START_TO_END*/, nodeRange) <= 0
				.SelectAll()
			}
		}

	SetReadOnly(unused)		// stub to override Component
		{ }
	GetReadOnly()			// stub to override Component
		{ return true }
	}
