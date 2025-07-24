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
		if tabthrough is true
			.CM.SetOption("extraKeys", [Tab: false])

		.SetStyles(Object('border': '1px solid black'), .CMEl)

		.AddEventListenerToCM('change', .OnChange)
		.AddEventListenerToCM('beforeChange', .onBeforeChange)
		.AddEventListenerToCM('focus', .OnFocus)
		.AddEventListenerToCM('blur', .blur)
		.AddEventListenerToCM('beforeSelectionChange', .onSelectionChange)
		.AddEventListenerToCM('contextmenu', .onContextMenu)
		.AddEventListenerToCM('scroll', .onScroll)
		.El.AddEventListener('mouseup', .onMouseUp)
		.AddEventListenerToCM('keydown', .OnKeyDown)

		.indicators = Object()
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
		.updateUI()
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

	updateUI()
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

	DefineIndicator(n, style, fore = false)
		{
		css = ''
		color = fore is false ? 'black' : ToCssColor(fore)
		if style is INDIC.SQUIGGLE
			css = 'text-decoration: underline wavy ' $ color $ ' 1px;'
		else if style is INDIC.TEXTFORE
			css = 'color: ' $ color $ '; cursor: pointer;'
		.indicators[n] = [:css, marks: Object()]
		}

	SetIndicator(indic, pos, len)
		{
		from = .CM.PosFromIndex(pos)
		to = .CM.PosFromIndex(pos + len)
		mark = .CM.MarkText(from, to,
			[css: .indicators[indic].css, inclusiveLeft: false, inclusiveRight:])
		.indicators[indic].marks.Add(mark)
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

	getMarkPos(mark)
		{
		try
			{
			// mark.Find() return undefined if the mark is no longer in the document
			markPos = mark.Find()
			markPos.from
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
			modificationType: changeObj.text.Empty?()
				? SC.MOD_DELETETEXT
				: SC.MOD_INSERTTEXT,
			position: from,
			length: toAfterChange - from + 1
			])
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
		.updateUI()
		}

	onScroll(unused)
		{
		scrollInfo = .CM.GetScrollInfo()
		line = .CM.LineAtHeight(scrollInfo.top, 'local')
		.Event(#SU_SYNCFIRSTVISIBLELINE, line)
		}

	SetFirstVisibleLine(line)
		{
		top = .CM.HeightAtLine(line, 'local')
		.CM.ScrollTo(0, top)
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
			if not .GetReadOnly()
				.CM.ReplaceSelection('')
			SuClipboardWriteString(str)
			}
		}
	COPY()
		{
		if .CM.GetSelection() is ''
			.CM.ExecCommand('selectAll')
		if '' isnt str = .CM.GetSelection()
			SuClipboardWriteString(str)
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
