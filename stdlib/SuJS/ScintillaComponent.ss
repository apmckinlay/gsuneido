// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
Component
	{
	Name:		"Editor"
	Xmin: 		200
	Ymin: 		100
	Xstretch:	1
	Ystretch:	1
	DefaultFontSize: 11
	Unsortable: true
	PaddingTop: 4

	styles: `
		.su-code-wrapper {
			position: relative;
		}
		.su-code-mirror {
			position: absolute;
			box-sizing: border-box;
			top: 0;
			bottom: 0;
			left: 0;
			right: 0;
			height: 100%;
			padding-top: 4px;
			padding-bottom: 4px;
		}
		.su-code-gutter {
			width: 1em;
		}
		.su-code-gutter-container {
			position: relative;
		}
		.su-code-mirror .CodeMirror-lines {
			padding: 0;
		}
		.CodeMirror-overwrite .CodeMirror-cursor {
			border-color: red;
		}
		.CodeMirror-selected  { background-color: #308dfc !important; }
		.CodeMirror-selectedtext { color: white; }
		.CodeMirror-activeline-background {
			background-color: var(--active-line-bg-color)
		}`

	New(.readonly = false, height = false, tabthrough = false)
		{
		if height is 1
			tabthrough = true

		LoadCssStyles('code-mirror.css', .styles)
		.CreateElement('div', className: 'su-code-wrapper')
		.textarea = CreateElement('textarea', .El)

		.CodeMirror = SuUI.GetCodeMirror()
		.CM = .CodeMirror.FromTextArea(.textarea)
		.CMEl = .CM.GetWrapperElement()
		.CMEl.classList.Add('su-code-mirror')
		.changes = Object()

		.CM.SetOption("mode", "null")
		.CM.SetOption("readOnly", .readonly)
		.CM.SetOption("lineSeparator", "\r\n")
		.CM.SetOption("styleSelectedText", true)
		if readonly
			.CM.SetOption("tabindex", "-1")

		extraKeys = ['Shift-Tab': 'indentLess']
		if tabthrough is true
			extraKeys['Tab'] = false
		.CM.SetOption("extraKeys", extraKeys)

		.SetStyles(Object('border': '1px solid black'), .CMEl)

		.AddEventListenerToCM('change', .OnChange)
		.CodeMirror.OnBeforeChange(.CM, .eventFactory(.ensureCRLF))
		.AddEventListenerToCM('beforeChange', .onBeforeChange)
		.AddEventListenerToCM('focus', .OnFocus)
		.AddEventListenerToCM('blur', .blur)
		.AddEventListenerToCM('beforeSelectionChange', .onSelectionChange)
		.AddEventListenerToCM('contextmenu', .onContextMenu)
		.AddEventListenerToCM('scroll', .onScroll)
		.El.AddEventListener('mouseup', .onMouseUp)
		.AddEventListenerToCM('keydown', .OnKeyDown)
		.El.AddEventListener('dragend', .dragEnd)

		.indicators = Object()
		.markers = Object()
		.gutters = Object()
		.SetMinSize()
		}

	AddEventListenerToCM(event, fn)
		{
		.CM.On(event, .eventFactory(fn))
		}

	eventFactory(fn)
		{
		return { |@args|
			if not .Destroyed?()
				fn(@args)
			}
		}

	TabOver?()
		{
		// use switch to handle undefined
		switch (.CM.GetOption("tabindex"))
			{
		case "-1", -1:
			return true
		default:
			return false
			}
		}

	onContextMenu(unused, event)
		{
		.SetFocus()
		selections = .CM.ListSelections()
		pos = .CM.CoordsChar(Object(left: event.pageX, top: event.pageY))
		if not selections.Any?({ .inSelect?(pos, it.anchor, it.head) })
			.CM.SetCursor(pos)
		.UpdateUI()
		.OnContextMenu(event)
		}

	inSelect?(pos, anchor, head)
		{
		less? = .less?(anchor, head)
		from = less? ? anchor : head
		to = less? ? head : anchor
		return not .less?(pos, from) and .less?(pos, to)
		}

	less?(pos1, pos2)
		{
		ob1 = Object(pos1.line, pos1.ch)
		ob2 = Object(pos2.line, pos2.ch)
		return ob1 < ob2
		}

	UpdateUI()
		{
		.Event(#UPDATEUI)
		}

	refreshing?: false
	refresh()
		{
		if .refreshing?
			return
		.refreshing? = true
		// needed to combine multiple screen refreshes
		// into one to reduce layout reflow
		.refreshTimer = SuDelayed(0, .doRefresh)
		}
	refreshTimer: false
	doRefresh()
		{
		.refreshing? = false
		if .Member?(#CM)
			.CM.Refresh()
		}
	WindowResize()
		{
		// to fix the issue where CM is not draw when WindowComponent is resized
		.refresh()
		}

	Recalc()
		{
		// to fix the issue where CM is not draw when changing from hide to show
		.refresh()
		}

	TextHeight()
		{
		return .CM.DefaultTextHeight()
		}

	SetReadOnly(readOnly)
		{
		if (.readonly)
			return
		.CM.SetOption("readOnly", readOnly is true)
		super.SetReadOnly(readOnly)
		}

	GetReadOnly()
		{
		return .CM.GetOption("readOnly")
		}

	StyleSetBack(unused, back)
		{
		.CMEl.SetStyle('background-color', ToCssColor(back))
		}

	SetlineWrapping(wrap)
		{
		.CM.SetOption("lineWrapping", wrap)
		}

	SetTABWIDTH(width)
		{
		.CM.SetOption('tabSize', width)
		.CM.SetOption('indentUnit', width)
		.CM.SetOption('indentWithTabs', true)
		}

	wordchars: 'zyxwvutsrqponmlkjihgfedcba_ZYXWVUTSRQPONMLKJIHGFEDCBA?9876543210!'
	SetWordChars(.wordChars)
		{
		if not Suneido.Member?(#RegisteredWC)
			Suneido.RegisteredWC = Object()
		if false is name = Suneido.RegisteredWC.GetDefault(wordChars, false)
			{
			name = 'suneido_' $ Suneido.RegisteredWC.Size()
			regex = SuUI.MakeWebObject('RegExp', '[' $ wordChars $ ']')
			.CodeMirror.RegisterHelper('wordChars', name, regex)
			Suneido.RegisteredWC[wordChars] = name
			}
		mode = .CodeMirror.GetMode(.CM)
		mode.wordChars = name
		}

	GetWordChars()
		{
		return .wordchars
		}

	setValue?: false
	Set(s)
		{
		.setValue? = true
		.CM.SetValue(s)
		.CM.ClearHistory()
		.setValue? = false
		}

	AppendText(s)
		{
		.setValue? = true
		.CM.ReplaceRange(s, .getEndPos())
		.Event(#EN_CHANGE)
		.setValue? = false
		.scrollToEnd()
		}

	getEndPos()
		{
		line = .CM.LastLine()
		ch = .CM.GetLine(line).Size()
		return Object(:line, :ch)
		}

	scrollToEnd()
		{
		.CM.ScrollIntoView(.getEndPos())
		}

	IsSettingValue?()
		{
		return .setValue?
		}

	markerCss: false
	DefineMarker(n, style, fore, back)
		{
		.markers[n] = [:style, :fore, :back]
		.markerCss = false
		if .gutterMarker?(style)
			.ensureGutter('su-code-gutter')
		}

	gutterMarker?(style)
		{
		return style is SC.MARK_SHORTARROW or
			style is SC.MARK_ROUNDRECT or
			style is #diffMarker
		}

	DefineXPMMarker(n, style, fore, back)
		{
		.DefineMarker(n, style, fore, back)
		}

	MarkerName: 'default'
	ensureMarkerCss()
		{
		if .markerCss isnt false
			return

		css = ''
		for n in .markers.Members()
			{
			className = 'su-code-' $ .MarkerName $ '-' $ n
			marker = .markers[n]
			if marker.style is SC.MARK_BACKGROUND
				{
				css $= '.' $ className $ '{
						background-color: ' $ ToCssColor(marker.back) $ ';
					}\r\n'
				}
			else if marker.style is SC.MARK_SHORTARROW
				css $= .buildShortArrow(marker, className, n)
			else if marker.style is #diffMarker
				css $= .buildDiffMarker(marker, className, n)
			else if marker.style is SC.MARK_ROUNDRECT
				css $= .buildRoundRect(marker, className, n)
			else
				Print('marker not handled: ', marker)
			}
		LoadCssStyles('su-code-' $ .MarkerName $ '.css', .markerCss = css, override?:)
		}

	ensureGutter(className)
		{
		if not .gutters.Has?(className)
			{
			.gutters.Add(className)
			.CM.SetOption(#gutters, .gutters)
			}
		}

	buildShortArrow(marker, className, n)
		{
		shape = '<svg viewBox="0 0 100 100" style="width: 100%; height: 100%;">
				<polygon
					points="5,35 40,35 40,10 95,50 40,90 40,65 5,65"
					fill="' $ ToCssColor(marker.back) $ '"
					stroke="black"
					stroke-width="6"
					stroke-linejoin="round"
				/>
			</svg>'
		return .buildMarker(marker, className, n, :shape,
			extraCss: 'justify-content: center; align-items: center;')
		}

	buildRoundRect(marker, className, n)
		{
		shape = '<svg viewBox="0 0 100 100" style="width: 100%; height: 100%;">
				<rect
					width="60" height="80"
					x="20" y="10" rx="20" ry="20"
					fill="' $ ToCssColor(marker.back) $ '"
					stroke="black"
					stroke-width="6"
				/>
			</svg>'
		return .buildMarker(marker, className, n, :shape,
			extraCss: 'justify-content: center; align-items: center;')
		}

	buildDiffMarker(marker, className, n)
		{
		return .buildMarker(marker, className, n,
			extraCss: 'width: 33%; right: 0px; background-color: ' $
				ToCssColor(marker.back) $ ';')
		}

	buildMarker(marker, className, n, shape = '', extraCss = '')
		{
		if not marker.Member?(#div)
			{
			div = CreateElement('div', className: className $ '-div')
			div.innerHTML = shape
			marker.div = div
			}
		return '.' $ className $ '-div {
				display: flex;
				position: absolute;
				height: ' $ .CM.DefaultTextHeight() $ 'px;
				z-index: ' $ n $ ';
				' $ extraCss $ '
			}\r\n'
		}

	MarkerAdd(row, n)
		{
		if false is marker = .markers.GetDefault(n, false)
			return
		.ensureMarkerCss()
		className = 'su-code-' $ .MarkerName $ '-' $ n
		if marker.style is SC.MARK_BACKGROUND
			.CM.AddLineClass(row, "background", className)
		else if .gutterMarker?(marker.style)
			{
			info = .CM.LineInfo(row)
			try
				container = info['gutterMarkers']['su-code-gutter']
			catch
				{
				container = CreateElement('div', className: 'su-code-gutter-container')
				.CM.SetGutterMarker(row, 'su-code-gutter', container)
				}
			container.AppendChild(marker.div.CloneNode(true))
			}
		}

	MarkerDelete(row, n)
		{
		if false is marker = .markers.GetDefault(n, false)
			return

		className = 'su-code-' $ .MarkerName $ '-' $ n
		if marker.style is SC.MARK_BACKGROUND
			.CM.RemoveLineClass(row, 'background', className)
		else if .gutterMarker?(marker.style)
			{
			info = .CM.LineInfo(row)
			container = false
			try container = info['gutterMarkers']['su-code-gutter']
			if container isnt false
				{
				list = container.QuerySelectorAll('.' $ className $ '-div')
				for i in .. list.length
					list.Item(i).Remove()
				}
			}
		}

	MarkerDeleteAll(n)
		{
		if false is marker = .markers.GetDefault(n, false)
			return

		className = 'su-code-' $ .MarkerName $ '-' $ n
		if marker.style is SC.MARK_BACKGROUND
			{
			c = .CM.LineCount()
			for row in c
				.CM.RemoveLineClass(row, 'background', className)
			}

		className = 'su-code-' $ .MarkerName $ '-' $ n
		list = .CMEl.QuerySelectorAll('.' $ className $ '-div')
		for i in .. list.length
			list.Item(i).Remove()
		}

	UpdateOverview(ovbarHwnds, markersInfo)
		{
		ovbars = Object()
		for type in ovbarHwnds.Members()
			ovbars[type] = SuRender().GetRegisteredComponent(ovbarHwnds[type])
		ovbars.Each(#ClearMarks)

		curLine = .CM.GetCursor().line
		for row in ...CM.LineCount()
			{
			if row is curLine
				ovbars.Each({ it.AddMark(row, CLR.GRAY) })
			else
				.updateMarkersAtLine(row, markersInfo, ovbars)
			}
		}

	updateMarkersAtLine(row, markersInfo, ovbars)
		{
		info = .CM.LineInfo(row)
		container = false
		try container = info['gutterMarkers']['su-code-gutter']
		if container isnt false and container.childElementCount isnt 0
			{
			for type in markersInfo.Members()
				{
				for idx in markersInfo[type]
					{
					list = container.QuerySelectorAll(
						'.su-code-' $ .MarkerName $ '-' $ idx $ '-div')
					if list.length isnt 0
						{
						ovbars[type].AddMark(row, .markers[idx].back)
						break
						}
					}
				}
			}
		}

	braceStyle: ''
	braceMark1: false
	braceMark2: false
	braceBadStyle: ''
	braceBadMark: false
	DefineStyle(n, fore = false, back = false, bold = false,
		underline/*unused*/ = false, italic/*unused*/ = false)
		{
		if n is SC.STYLE_BRACELIGHT
			.braceStyle = .buildStyle(fore, back, bold)
		else if n is SC.STYLE_BRACEBAD
			.braceBadStyle = .buildStyle(fore, back, bold)
		}

	buildStyle(fore = false, back = false, bold = false)
		{
		css = ''
		if fore isnt false
			css $= ' color: ' $ ToCssColor(fore) $ ';'
		if back isnt false
			css $= ' background-color: ' $ ToCssColor(back) $ ';'
		if bold isnt false
			css $= ' font-weight: bold;'
		return css
		}

	BraceHighlight(pos1, pos2)
		{
		.clearBraceHighlights()
		.braceMark1 = .addBrace(pos1, .braceStyle)
		.braceMark2 = .addBrace(pos2, .braceStyle)
		}

	BraceBadLight(pos)
		{
		.clearBraceHighlights()
		.braceBadMark = .addBrace(pos, .braceBadStyle)
		}

	addBrace(pos, css)
		{
		if pos is -1
			return false
		return .CM.MarkText(.CM.PosFromIndex(pos), .CM.PosFromIndex(pos + 1),
			[:css, inclusiveLeft: false, inclusiveRight: false])
		}

	clearBraceHighlights()
		{
		.clearBrace(.braceMark1)
		.clearBrace(.braceMark2)
		.clearBrace(.braceBadMark)
		.braceMark1 = .braceMark2 = .braceBadMark = false
		}

	clearBrace(mark)
		{
		if mark is false
			return
		mark.Clear()
		}

	DefineIndicator(n, style, fore = false)
		{
		.indicators[n] = [styles: Object(:style, :fore),
			css: .buildCss(style, fore), marks: Object()]
		}

	IndicSetAlpha(n, alpha)
		{
		.indicators[n].styles.alpha = alpha.Hex()
		.indicators[n].css = .buildCss(@.indicators[n].styles)
		}

	IndicSetOutlineAlpha(n, alpha)
		{
		.indicators[n].styles.outlineAlpha = alpha.Hex()
		.indicators[n].css = .buildCss(@.indicators[n].styles)
		}

	IndicSetHoverStyle(@unused) {}

	buildCss(style, fore, alpha = '1A', outlineAlpha = '80')
		{
		css = ''
		color = fore is false ? 'black' : ToCssColor(fore)
		if style is INDIC.SQUIGGLE
			css = 'text-decoration: underline wavy ' $ color $ ' 1px;'
		else if style is INDIC.TEXTFORE
			css = 'color: ' $ color $ '; cursor: pointer;'
		else if style is INDIC.ROUNDBOX
			css = 'border-radius: 2%; background-color: ' $ color $ alpha $ ';'
		else if style is INDIC.STRAIGHTBOX
			css = 'background-color: ' $ color $ alpha $ '; ' $
				'outline: 1px solid ' $ color $ outlineAlpha $ ';'
		else if style is INDIC.HIDDEN
			css = ''
		return css
		}

	SetIndicator(indic, pos, len)
		{
		from = .CM.PosFromIndex(pos)
		to = .CM.PosFromIndex(pos + len)
		mark = .CM.MarkText(from, to,
			[css: .indicators[indic].css, inclusiveLeft: false, inclusiveRight:])
		.indicators[indic].marks.Add(mark)
		}

	curIndic: false
	SetIndicatorCurrent(.curIndic) { }
	IndicatorFillRange(pos, len)
		{
		if .curIndic is false
			return
		.SetIndicator(.curIndic, pos, len)
		}

	ClearIndicator(indic, pos = 0, len = false)
		{
		if len is false
			{
			.indicators[indic].marks.Each({ it.Clear() })
			.indicators[indic].marks = Object()
			return
			}
		toClear = Object()
		for i in .indicators[indic].marks.Members()
			{
			mark = .indicators[indic].marks[i]
			if false isnt markPos = .getMarkPos(mark)
				{
				markFrom = .CM.IndexFromPos(markPos.from)
				markTo = .CM.IndexFromPos(markPos.to)
				if markFrom >= pos + len or markTo < pos
					continue
				}
			toClear.Add(i)
			mark.Clear()
			}
		.indicators[indic].marks.Delete(@toClear)
		}

	HasIndicator?(pos, indic)
		{
		for mark in .indicators[indic].marks
			{
			if false isnt markPos = .getMarkPos(mark)
				{
				markFrom = .CM.IndexFromPos(markPos.from)
				markTo = .CM.IndexFromPos(markPos.to)
				if markFrom < pos + 1 and markTo >= pos
					return true
				}
			}
		return false
		}

	getMarkPos(mark)
		{
		try
			{
			// mark.Find() return undefined if the mark is no longer in the document
			markPos = mark.Find()
			markPos.from
			markPos.to
			return markPos
			}
		catch
			return false
		}

	Get()
		{
		.CM.GetValue()
		}

	ignoreChange: false
	DoWithoutChange(block)
		{
		.ignoreChange = true
		block()
		.ignoreChange = false
		}

	SetOption(option, value)
		{
		.CM.SetOption(option, value)
		}

	SetFont(font = "", size = "", weight = "", italic = false)
		{
		super.SetFont(font, size, weight, :italic, el: .CMEl)
		}

	SetFocus()
		{
		.CM.Focus()
		}

	OnChange(unused, changeObj)
		{
		if .ignoreChange
			return

		.DoOnChange(changeObj)
		}

	DoOnChange(changeObj)
		{
		from = false
		if false isnt i = .changes.FindIf({
			it.line is changeObj.from.line and it.ch is changeObj.from.ch })
			{
			from = .changes[i].fromIdx
			.changes.Delete(i)
			}
		Assert(from isnt: false)
		added = changeObj.text.Join('\r\n')
		removed = changeObj.removed.Join('\r\n')
		to = from + removed.Size()

		if .IsSettingValue?() isnt true
			{
			.Event(#Scintilla_UpdateValue,  change: [:from, :to, text: added])
			.Event(#EN_CHANGE)
			}
		else
			// this is to trigger the IdleTimer without setting dirty
			.Event(#Scintilla_SetValue)

		toAfterChange = to + added.Size() - removed.Size()
		.Event(#SCN_MODIFIED, [
			modificationType: added is ''
				? SC.MOD_DELETETEXT
				: SC.MOD_INSERTTEXT,
			position: from,
			length: toAfterChange - from + 1
			])
		}

	ensureCRLF(unused, changeObj)
		{
		try
			{
			if not changeObj.text.Any?({ it.Has1of?('\r\n') })
				return
			newText = Object()
			for line in changeObj.text
				newText.Append(line.Tr('\r').Split('\n'))
			(changeObj.update)(changeObj.from, changeObj.to, newText)
			}
		catch (e)
			{
			SuRender().Event(false, 'SuneidoLog', Object(
				'ERROR (CAUGHT): ScintillaComponent.ensureCRLF - ' $ e,
				params: changeObj.text, caughtMsg: 'continue'))
			}
		}

	onBeforeChange(unused, changeObj)
		{
		if .ignoreChange
			return

		fromIdx = .CM.IndexFromPos(changeObj.from)
		.changes.Add([:fromIdx, line: changeObj.from.line, ch: changeObj.from.ch])
		}

	OnFocus(@unused)
		{
		super.OnFocus()
		.Event(#SCEN_SETFOCUS)
		}

	// event is undefined when right click on a codemirror
	blur(unused, event = false)
		{
		.OnBlur(event)
		.EventWithFreeze(#SCEN_KILLFOCUS)
		}

	onSelectionChange(unused, obj)
		{
		anchor = .CM.IndexFromPos(obj.ranges.Last().anchor)
		head = .CM.IndexFromPos(obj.ranges.Last().head)
		.Event(#SU_UPDATESELECT, Object(:anchor, :head))
		.UpdateUI()
		}

	onScroll(unused)
		{
		scrollInfo = .CM.GetScrollInfo()
		line = .CM.LineAtHeight(scrollInfo.top, 'local')
		.Event(#SU_SYNCFIRSTVISIBLELINE, line)
		}

	SetFirstVisibleLine(line, centerInScreen? = false)
		{
		offset = 0
		if centerInScreen? is true
			offset = .CM.GetScrollInfo().clientHeight / 2
		top = .CM.HeightAtLine(line, 'local')
		.CM.ScrollTo(0, top - offset)
		}

	scrollInfo: false
	SavePos()
		{
		.scrollInfo = .CM.GetScrollInfo()
		}

	RestorePos()
		{
		if .scrollInfo is false
			return
		.CM.ScrollTo(.scrollInfo.left, .scrollInfo.top)
		.scrollInfo = false
		}

	EnsureRangeVisible(from, to)
		{
		.CM.ScrollIntoView(.CM.PosFromIndex(from), .CM.PosFromIndex(to))
		}

	SetSel(from, to)
		{
		if to is -1
			to = .Get().Size()
		from = .CM.PosFromIndex(from)
		to = .CM.PosFromIndex(to)
		.CM.SetSelection(from, to)
		}

	CharLeftExtend()
		{
		selections = .CM.ListSelections()
		if selections.Size() is 1
			{
			anchor = selections[0].anchor
			head = selections[0].head
			if head.ch > 0
				head.ch--
			.CM.SetSelection(anchor, head)
			}
		}

	AddSel(head, anchor)
		{
		head = .CM.PosFromIndex(head)
		anchor = .CM.PosFromIndex(anchor)
		.CM.AddSelection(anchor, head)
		}

	CUT()
		{
		if '' isnt str = .CM.GetSelection()
			{
			SuClipboardWriteString(str, 'Cut').Then(
				{|res|
				if false isnt res and not .GetReadOnly()
					.CM.ReplaceSelection('')
				})
			}
		}
	COPY()
		{
		if .CM.GetSelection() is ''
			.CM.ExecCommand('selectAll')
		if '' isnt str = .CM.GetSelection()
			SuClipboardWriteString(str, 'Copy')
		}
	PASTE()
		{
		SuClipboardPasteString(this, .Paste)
		}
	UNDO()
		{
		if .GetReadOnly()
			return
		.CM.ExecCommand('undo')
		}
	REDO()
		{
		if .GetReadOnly()
			return
		.CM.ExecCommand('redo')
		}
	DELETE()
		{
		if .GetReadOnly()
			return
		if .CM.GetSelection() isnt ''
			{
			.CM.ReplaceSelection('')
			return
			}
		from = .CM.GetCursor('from')
		to = Object(ch: from.ch + 1, line: from.line)
		.CM.SetSelection(from, to)
		if .CM.GetSelection() is ''
			{
			to = Object(line: from.line + 1, ch: 1)
			.CM.SetSelection(from, to)
			}
		.CM.ReplaceSelection('')
		}
	SELECTALL()
		{
		.CM.ExecCommand('selectAll')
		}

	ReplaceSel(s)
		{
		.CM.ReplaceSelection(s)
		}

	ReplaceSelFromServer(s, start, end)
		{
		.setValue? = true
		.CM.SetSelection(.CM.PosFromIndex(start), .CM.PosFromIndex(end))
		.CM.ReplaceSelection(s)
		.Event(#EN_CHANGE)
		.setValue? = false
		}

	Paste(s)
		{
		.ReplaceSel(s)
		}

	HasFocus?()
		{
		return .CM.HasFocus()
		}

	onMouseUp(event)
		{
		if event.button isnt 0
			return

		.Event('LBUTTONUP', .CM.IndexFromPos(.CM.GetCursor()))
		if event.detail isnt 0 and event.detail % 2 is 0
			.Event(#SCN_DOUBLECLICK)
		}

	SetCaretLineVisible(enable)
		{
		.CM.SetOption("styleActiveLine", enable)
		}

	SetCaretLineBack(color)
		{
		.SetStyleProperty('--active-line-bg-color', ToCssColor(color))
		}

	SetStyleProperty(name, value)
		{
		.CMEl.style.setProperty(name, value)
		}

	GetListPos() // called by AutoChooseListComponent
		{
		coord = .CM.CursorCoords()
		return coord
		}

	GetCharPos(i) // called by GotoLibView.pt
		{
		pos = .CM.PosFromIndex(i)
		coord = .CM.CharCoords(pos)
		return coord
		}

	autocOpen?: false
	SyncAutocStatus(.autocOpen?) {}
	OnKeyDown(unused, event)
		{
		if .autocOpen? and event.key in (#ArrowUp, #ArrowDown, #Tab, #Enter, #Escape)
			{
			.Event(#AUTOC_KEYDOWN, event.key)
			event.PreventDefault()
			event.StopPropagation()
			}
		}

	dragEnd(event /*unused*/)
		{
		.EventWithFreeze(#SCEN_KILLFOCUS)
		}

	MoveSelectedLinesUp()
		{
		cursor = .CM.GetCursor()
		line = cursor.line

		if line <= 0
			return

		.swap(line, line - 1, cursor.ch)
		}

	MoveSelectedLinesDown()
		{
		cursor = .CM.GetCursor()
		line = cursor.line

		if line >= .CM.LineCount() - 1
			return

		.swap(line, line + 1, cursor.ch)
		}

	swap(from, to, ch)
		{
		fromLine = .CM.GetLine(from)
		toLine = .CM.GetLine(to)

		.CM.StartOperation()
		.CM.replaceRange(fromLine, [line: to, ch: 0], [line: to, ch: toLine.Size()]);
		.CM.replaceRange(toLine, [line: from, ch: 0], [line: from, ch: fromLine.Size()]);
		.CM.SetCursor([line: to, :ch])
		.CM.EndOperation()
		}

	GetDimension()
		{
		if not .Member?(#CM)
			return #(scrollWidth: 0, scrollHeight: 0, clientWidth: 0, clientHeight: 0)
		info = .CM.GetScrollInfo()
		return Object(scrollWidth: info.width, scrollHeight: info.height,
			clientWidth: info.clientWidth, clientHeight: info.clientHeight)
		}

	Getter_(member)
		{
		return { |@args| Print('Scintilla', member, args) }
		}

	Destroy()
		{
		.CM.ToTextArea()
		super.Destroy()
		}
	}
