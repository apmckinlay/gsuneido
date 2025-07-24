// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
ScintillaAddonsComponent
	{
	New(@args)
		{
		super(@args)
		.states = Object(bold: false, italic: false, underline: false, strikeout: false)
		.marks = Object()
		.AddEventListenerToCM('cursorActivity', .onCursorActivity)
		.AddEventListenerToCM('update', .onUpdate)
		}

	pendingStyles: false
	pendingUpdate: false
	styles: (
		bold: ['font-weight: normal;', 'font-weight: bold;'],
		italic: ['font-style: normal;', 'font-style: italic;'],
		line: ['text-decoration-line: none;',
			'text-decoration-line: underline;',
			'text-decoration-line: line-through;',
			'text-decoration-line: underline line-through;'])

	On_Bold()
		{
		.update('bold')
		}

	On_Italic()
		{
		.update('italic')
		}

	On_Underline()
		{
		.update('underline')
		}

	On_Strikeout()
		{
		.update('strikeout')
		}

	update(state)
		{
		from = .CM.GetCursor("from")
		to = .CM.GetCursor("to")

		if .comparePos(from, to) is 0
			{
			if .pendingStyles is false
				.pendingStyles = Object(pos: from, states: Object())
			.pendingStyles.states[state] = .states[state] = not .states[state]
			.SetFocus()
			.Event(#UpdateButtons, .states)
			return
			}

		.states[state] = not .states[state]
		if state in ('bold', 'italic')
			.updateBoldItalic(state, from, to)
		else
			.updateLine(state, from, to)
		.SetFocus()
		.Event(#UpdateButtons, .states)
		.Event(#ScintillaRichEditor_UpdateStyleObject, .getStyles())
		.Event(#EN_CHANGE)
		}

	updateBoldItalic(state, from, to)
		{
		css = .states[state] is true ? .styles[state][1] : .styles[state][0]
		.addMark(from, to, css)
		}

	addMark(from, to, css, addToHistory = true)
		{
		mark = .CM.MarkText(from, to,
			[:css, :addToHistory, inclusiveLeft: false, inclusiveRight:])
		.marks.Add(mark)
		}

	updateLine(state, from, to)
		{
		styles = .getStyles([start: from, end: to])
		for style in styles
			{
			idx = 0
			if state is 'underline'
				{
				idx |= .states[state] is true ? 1 : 0
				idx |= style.style.strikeout ? 2 : 0
				}
			else
				{
				idx |= .states[state] is true ? 2 : 0
				idx |= style.style.underline ? 1 : 0
				}
			css = .styles.line[idx]
			.addMark(style.from, style.to, css)
			}
		}

	On_ResetFont()
		{
		from = .CM.GetCursor("from")
		to = .CM.GetCursor("to")

		.states = Object(bold: false, italic: false, underline: false, strikeout: false)
		css = .styles.bold[0] $ .styles.italic[0] $ .styles.line[0]
		.addMark(from, to, css)
		.SetFocus()
		.Event(#UpdateButtons, .states)
		.Event(#ScintillaRichEditor_UpdateStyleObject, .getStyles())
		.Event(#EN_CHANGE)
		}

	onCursorActivity(@unused)
		{
		pos = .CM.GetCursor("from")

		if .pendingStyles isnt false
			{
			if .comparePos(.pendingStyles.pos, pos) is 0
				return
			.pendingStyles = false
			}

		.updateButtons(pos)
		}

	updateButtons(pos)
		{
		marks = .getMarks(at: pos)
		bold = italic = underline = strikeout = false
		for (i = marks.Size() - 1; i >= 0; i--)
			{
			markPos = marks[i].Find()
			if .comparePos(markPos.from, pos) is 0
				continue
			if bold is false
				bold = .checkBold(marks[i].css)
			if italic is false
				italic = .checkItalic(marks[i].css)
			if underline is false
				underline = .checkUnderline(marks[i].css)
			if strikeout is false
				strikeout = .checkStrikeout(marks[i].css)
			}
		.states.bold = bold is 1
		.states.italic = italic is 1
		.states.underline = underline is 1
		.states.strikeout = strikeout is 1
		.Event(#UpdateButtons, .states)
		}

	checkBold(css)
		{
		return .check(css, 'bold')
		}

	checkItalic(css)
		{
		return .check(css, 'italic')
		}

	check(css, state)
		{
		return css.Has?(.styles[state][1])
			? 1
			: css.Has?(.styles[state][0])
				? 0
				: false
		}

	checkUnderline(css)
		{
		return css.Has?('text-decoration-line')
			? css.Has?('underline') ? 1 : 0
			: false
		}

	checkStrikeout(css)
		{
		return css.Has?('text-decoration-line')
			? css.Has?('line-through') ? 1 : 0
			: false
		}

	// OnChange and onUpdate are to handle the case where styles are (un)selected without
	// any selections. CodeMirror.MarkText has no effect when from equals to
	DoOnChange(changeObj)
		{
		if not .IsSettingValue?()
			{
			if .pendingStyles isnt false and .isAddText?(changeObj) and
				.comparePos(.pendingStyles.pos, changeObj.from) is 0
				{
				css = .buildCss()
				.pendingUpdate = Object(from: .pendingStyles.pos, :css)
				.pendingStyles = false
				}
			.Event(#ScintillaRichEditor_UpdateStyleObject, .getStyles())
			}
		super.DoOnChange(changeObj)
		}

	buildCss()
		{
		css = ''
		if .hasStyle?(#bold, .pendingStyles.states)
			css $= .styles.bold[.pendingStyles.states.bold ? 1 : 0]

		if .hasStyle?(#italic, .pendingStyles.states)
			css $= .styles.italic[.pendingStyles.states.italic ? 1 : 0]

		if .hasStyle?(#line, .pendingStyles.states)
			{
			idx = .pendingStyles.states.GetDefault(#underline, false) ? 1 : 0
			idx |= .pendingStyles.states.GetDefault(#strikeout, false) ? 2 : 0
			css $= .styles.line[idx]
			}
		return css
		}

	hasStyle?(style, states)
		{
		if style in (#bold, #italic)
			return states.Member?(style)
		if style is #line
			return states.Member?(#underline) or states.Member?(#strikeout)
		return false
		}

	onUpdate(unused)
		{
		if .pendingUpdate is false
			return

		from = .pendingUpdate.from
		to = .CM.GetCursor("to")
		css = .pendingUpdate.css
		.pendingUpdate = false
		.addMark(from, to, css, addToHistory: false)
		.updateButtons(to)
		}

	isAddText?(changeObj)
		{
		return changeObj.text.Size() >= 1 and
			changeObj.removed.Size() is 1 and changeObj.removed[0] is ""
		}

	comparePos(pos1, pos2)
		{
		ob1 = Object(pos1.line, pos1.ch)
		ob2 = Object(pos2.line, pos2.ch)
		if ob1 < ob2
			return -1
		return ob1 is ob2 ? 0 : 1
		}

	SetStyleObject(styleObject)
		{
		.DoWithoutChange()
			{
			for ob in styleObject
				{
				css = ''

				if ob.style.bold is true
					css $= .styles.bold[1]
				if ob.style.italic is true
					css $= .styles.italic[1]
				idx = ob.style.underline is true ? 1 : 0
				idx |= ob.style.strikeout is true ? 2 : 0
				if idx isnt 0
					css $= .styles.line[idx]

				.addMark(ob.from, ob.to, css, addToHistory: false)
				}
			}
		}

	// group the styles between start and end
	// style[0] from       to
	// style[1]            from              to
	//          |---bold---|--italic & bold--|
	getStyles(range = false)
		{
		s = .Get()
		if s is ''
			return #()

		range = .getRange(range,  s)
		start = range.start
		end = range.end

		marks = .getMarks(:range)
		styles = Object(Object(from: start, to: end,
			style: Object(bold: false, italic: false,
				underline: false, strikeout: false)))

		for (i = marks.Size() - 1; i >= 0; i--)
			{
			if marks[i].css is ''
				continue
			style = Object(
				bold: .checkBold(marks[i].css),
				italic: .checkItalic(marks[i].css),
				underline: .checkUnderline(marks[i].css)
				strikeout: .checkStrikeout(marks[i].css))
			pos = .getPos(marks[i])
			from = .comparePos(pos.from, start) is -1 ? start : pos.from
			to = .comparePos(end, pos.to) is -1 ? end : pos.to // not inclusive
			if .comparePos(from, to) is 0
				continue

			styleFrom = .findFromStyle(from, styles)
			styleTo = .findToStyle(styleFrom, to, styles)
			.updateStyles(styleFrom, styleTo, style, styles)
			}
		return .formatAndMergeStyles(styles)
		}

	getRange(range, s)
		{
		if range isnt false
			return range

		start = [line: 0, ch: 0]
		lines = s.Lines()
		if s.Suffix?('\r\n')
			lines.Add('')
		end = [line: lines.Size() - 1, ch: lines.Last().Size()]
		return Object(:start, :end)
		}

	getMarks(at = false, range = false)
		{
		result = Object()
		for mark in .marks
			{
			res = false
			if mark.css isnt '' and
				mark.css =~ 'text-decoration-line|font-weight|font-style' and
				false isnt pos = .getPos(mark)
				res = .include?(at, range, pos)
			if res is true
				result.Add(mark)
			}
		return result
		}

	include?(at, range, pos)
		{
		if at isnt false
			return .comparePos(pos.from, at) < 0 and .comparePos(pos.to, at) >= 0
		else if range isnt false
			return not (.comparePos(pos.to, range.start) <= 0 or
				.comparePos(pos.from, range.end) >= 0)
		else
			return true
		}

	getPos(mark)
		{
		try
			{
			pos = mark.find()
			return Object(from: Object(line: pos.from.line, ch: pos.from.ch),
				to: Object(line: pos.to.line, ch: pos.to.ch))
			}
		catch // mark.find returns undefined if the mark has been removed
			return false
		}

	findFromStyle(from, styles)
		{
		styleFrom = 0
		while styleFrom < styles.Size() and
			.comparePos(styles[styleFrom].to, from) <= 0
			styleFrom++
		Assert(styleFrom isnt: styles.Size())
		if .comparePos(styles[styleFrom].from, from) isnt 0
			{
			newStyle = Object(:from, to: styles[styleFrom].to,
				style: styles[styleFrom].style.Copy())
			styles[styleFrom].to = from
			styles.Add(newStyle, at: ++styleFrom)
			}
		return styleFrom
		}

	findToStyle(startStyle, to, styles)
		{
		styleTo = startStyle
		while styleTo < styles.Size() and
			.comparePos(styles[styleTo].to, to) is -1
			styleTo++
		Assert(styleTo isnt: styles.Size())
		if .comparePos(styles[styleTo].to, to) isnt 0
			{
			newStyle = Object(from: to, to: styles[styleTo].to,
				style: styles[styleTo].style.Copy())
			styles[styleTo].to = to
			styles.Add(newStyle, at: styleTo + 1)
			}
		return styleTo
		}

	updateStyles(styleFrom, styleTo, style, styles)
		{
		for (n = styleFrom; n <= styleTo; n++)
			{
			for m in style.Members()
				{
				if style[m] is false or styles[n].style[m] isnt false
					continue
				styles[n].style[m] = style[m]
				}
			}
		}

	formatAndMergeStyles(styles)
		{
		result = Object()
		for style in styles
			{
			for m in #(bold, italic, underline, strikeout)
				style.style[m] = style.style[m] is 1
			if result.Empty?()
				result.Add(style)
			else
				{
				Assert(result.Last().to is: style.from)
				if result.Last().style is style.style
					result.Last().to = style.to
				else
					result.Add(style)
				}
			}
		return result
		}
	}
