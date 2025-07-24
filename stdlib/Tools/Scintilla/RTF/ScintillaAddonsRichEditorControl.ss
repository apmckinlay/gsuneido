// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
ScintillaAddonsControl
	{
	undoStyleBuffer: false
	redoStyleBuffer: false
	stylingMask: 0x1f // default of 5 style bits, and 3 indicator bits = 31
	wordChars: "_0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	New(@args)
		{
		super(@.processArgs(args))
		.SetWordChars(.wordChars)
		// setting the lexer to use CONTAINER (instead of the default NULL)
		// ensures that the SCN_STYLENEEDED message gets sent.
		// This prevents scintilla from losing formats when the user deletes text.
		.SetLexer(SCLEX.CONTAINER)
		.defineStyles()
		.SetUndoCollection(true)
		.SetModEventMask(.GetModEventMask() | SC.MOD_CONTAINER | SC.MOD_BEFOREDELETE)
		.undoStyleBuffer = Stack()
		.redoStyleBuffer = Stack()
		if .readonly = args.GetDefault(#readonly, false)
			.setBackground(GetSysColor(COLOR.BTNFACE))
		}

	processArgs(args)
		{
		.zoom = args.GetDefault('zoom', false)
		zoomArgs = [zoom: .zoom, zoom_ctrl: ScintillaRichWordZoomControl]
		return args.MergeNew([wrap:, margin: 8, exStyle: WS_EX.STATICEDGE,
			Addon_url:, Addon_multiple_selection:, Addon_zoom: zoomArgs,
			font: 		Suneido.logfont.lfFaceName,
			fontSize: 	'+0',
			weight: 	FW.NORMAL,
			italic: 	false])
		}

	STRIKETHROUGHINDIC: 1
	defineStyles()
		{
		.DefineStyle(SCIRT.NORMAL)
		.DefineStyle(SCIRT.BOLD, bold:)
		.DefineStyle(SCIRT.ITALIC, italic:)
		.DefineStyle(SCIRT.BOLD_ITALIC, bold:, italic:)
		.DefineStyle(SCIRT.UNDERLINE, underline:)
		.DefineStyle(SCIRT.BOLD_UNDERLINE, bold:, underline:)
		.DefineStyle(SCIRT.ITALIC_UNDERLINE, italic:, underline:)
		.DefineStyle(SCIRT.BOLD_ITALIC_UNDERLINE, bold:, italic:, underline:)
		.DefineStyle(SCIRT.STRIKETHROUGH)
		.DefineStyle(SCIRT.STRIKETHROUGH_BOLD, bold:)
		.DefineStyle(SCIRT.STRIKETHROUGH_ITALIC, italic:)
		.DefineStyle(SCIRT.STRIKETHROUGH_BOLD_ITALIC, bold:, italic:)
		.DefineStyle(SCIRT.STRIKETHROUGH_UNDERLINE, underline:)
		.DefineStyle(SCIRT.STRIKETHROUGH_BOLD_UNDERLINE, bold:, underline:)
		.DefineStyle(SCIRT.STRIKETHROUGH_ITALIC_UNDERLINE, italic:, underline:)
		.DefineStyle(SCIRT.STRIKETHROUGH_BOLD_ITALIC_UNDERLINE, bold:, italic:,
			underline:)
		.DefineIndicator(.STRIKETHROUGHINDIC, SC.INDIC_STRIKE, fore: 0)
		}

	EN_CHANGE()
		{
		.getCached = false
		return super.EN_CHANGE()
		}

	Dirty?(dirty = "")
		{
		if dirty is true
			.getCached = false
		return super.Dirty?(dirty)
		}

	getCached: false
	Get()
		{
		if .getCached isnt false
			return .getCached

		txt = .getStyledText(0, .GetLength())
		// don't save HTML tags if no text to save
		if txt.Blank?()
			return ""
		str = '<span style="font-weight:normal;font-style:normal">'
		.parseStyledText(txt,
			{|curStyle, pos /*unused*/|
			str $= '</span>' $ '<' $ .setFontStyle(curStyle) $ '>'
			},
			{|chr, pos /*unused*/, lastStyle /*unused*/|
			str $= ScintillaRichEditorEncodeMap.GetDefault(chr, chr)
			})
		str $= '</span>'
		.getCached = str
		return str
		}

	getStyledText(start, end)
		{
		s = ""
		chunk = 10000
		for (i = start; i < end; i += chunk)
			s $= SendMessageTextRange(.Hwnd, SCI.GETSTYLEDTEXT, i, Min(end, i + chunk), 2)
		return s
		}

	parseStyledText(text, handleStyleBlock, handleCharBlock)
		{
		iter = text.Iter()
		// need to initialize to -1 so that text that starts with
		// no formatting still gets handled correctly
		lastStyle = -1
		pos = 0
		while iter isnt chr = iter.Next()
			{
			if iter isnt style = iter.Next()
				{
				curStyle = style.Asc()
				if curStyle isnt lastStyle
					{
					handleStyleBlock(curStyle, pos)
					lastStyle = curStyle
					}
				}
			handleCharBlock(chr, pos, lastStyle)
			pos++
			}
		}

	GetText(includeStrikeThrough = false)
		{
		return .getText(0, .GetLength(), includeStrikeThrough)
		}

	getText(start, end, includeStrikeThrough = false)
		{
		if includeStrikeThrough is true
			return .GetRange(start, end)

		txt = .getStyledText(start, end)
		str = ""
		.parseStyledText(txt,
			{|curStyle /*unused*/, pos /*unused*/| /* do nothing */ },
			{|chr, pos /*unused*/, lastStyle|
				if (includeStrikeThrough is true or
					(lastStyle & SCIRT.STRIKETHROUGH) is 0)
					str $= chr
			})
		return str
		}

	SearchText()
		{
		return .GetText(includeStrikeThrough:)
		}

	Set(s)
		{
		.getCached = s
		if false isnt parser = ParseHTMLRichText.GetParsedText(s)
			.setStyled(parser)
		else
			.superSet(s)
		}
	setStyled(parser)
		{
		str = ""
		styleObject = Object()
		for child in parser.Children()
			str $= .buildStyled(child, styleObject, start: str.Size())
		.superSet(str)
		for style in styleObject
			{
			.StartStyling(style.start, .stylingMask)
			.SetStyling(style.length, style.styleBit)
			.SetIndicatorCurrent(.STRIKETHROUGHINDIC)
			if style.indicator is true
				.IndicatorFillRange(style.start, style.length)
			else
				.IndicatorClearRange(style.start, style.length)
			}
		}
	buildStyled(child, styleObject, start)
		{
		styleBit = 0
		indicator = false
		// need to handle text that has been appended to the end
		// might not have rtf formatting
		format = child.Attributes().Member?('style')
			? child.Attributes().style
			: ""
		if format =~ (':bold(;|$)')
			styleBit += SCIRT.BOLD
		if format =~ (':italic(;|$)')
			styleBit += SCIRT.ITALIC
		if format.Has?(':underline')
			styleBit += SCIRT.UNDERLINE
		if format.Has?(' line-through')
			{
			styleBit += SCIRT.STRIKETHROUGH
			indicator = true
			}
		txt = child.Text().Replace('&crlf;', '\r\n')
		styleObject.Add(Object(:start, length: txt.Size(), :styleBit, :indicator))
		return txt
		}

	// factored out so we can test
	superSet(s)
		{
		super.Set(s)
		}

	setFontStyle(curStyle)
		{
		base = 'span style="'
		bold = "font-weight:" $ ((curStyle & SCIRT.BOLD) isnt 0 ? "bold" : "normal")
		italic = "font-style:" $ ((curStyle & SCIRT.ITALIC) isnt 0 ? "italic" : "normal")
		underline = (curStyle & SCIRT.UNDERLINE) isnt 0 ? "underline" : ""
		strikethrough = (curStyle & SCIRT.STRIKETHROUGH) isnt 0 ? " line-through" : ""
		return base $ bold $ ';' $ italic $
			Opt(';text-decoration:', (underline $ strikethrough)) $ '"'
		}

	CHARADDED(lParam)
		{
		if .GetReadOnly() isnt true
			{
			selection = .getSelection()
			val = .getStyleVal(.bold_state, .italic_state, .underline_state,
				.strikeout_state)
			startPoint = selection.start-1
			.StartStyling(startPoint, .stylingMask)
			.SetStyling(1, val)

			if .strikeout_state is true
				{
				.SetIndicatorCurrent(.STRIKETHROUGHINDIC)
				.IndicatorFillRange(startPoint, selection.end - startPoint)
				.toggleSelectedTextFormat(strikethrough:)
				}
			}

		return super.CHARADDED(lParam)
		}

	getStyleVal(bold = false, italic = false, underline = false, strikethrough = false)
		{
		return (bold is true ? SCIRT.BOLD : 0) +
			(italic is true ? SCIRT.ITALIC : 0) +
			(underline is true ? SCIRT.UNDERLINE : 0) +
			(strikethrough is true ? SCIRT.STRIKETHROUGH : 0)
		}

	UPDATEUI()
		{
		.update_buttons()
		return super.UPDATEUI()
		}

	// TODO handle selection range not just start
	currentSelect: false
	update_buttons()
		{
		if .currentSelect is select = .getSelection()
			return
		.currentSelect = select

		x = .GetStyleAt(select.start - 1)
		.bold_state = (x & SCIRT.BOLD) isnt 0
		.italic_state = (x & SCIRT.ITALIC) isnt 0
		.underline_state = (x & SCIRT.UNDERLINE) isnt 0
		.Send('ToggleFontButton', 'bold', .bold_state)
		.Send('ToggleFontButton', 'italic', .italic_state)
		.Send('ToggleFontButton', 'underline', .underline_state)

		y = .IndicatorValueAt(.STRIKETHROUGHINDIC, select.start - 1)
		// 0 or 1 is 1 - strikethrough is off
		// 1 or 1 is 0 - strikethrough is on
		.strikeout_state = y is 1
		.Send('ToggleFontButton', 'strikeout', .strikeout_state)
		}

	Hasfocus?: false
	HasFocus?()
		{
		return .Hasfocus? or super.HasFocus?()
		}

	SetReadOnly(readOnly)
		{
		if .readonly
			return
		super.SetReadOnly(readOnly)
		.setBackground(readOnly is true ? GetSysColor(COLOR.BTNFACE) : 0xffffff)
		}

	setBackground(bgnd)
		{
		.StyleSetBack(0, bgnd)
		for style in SCIRT.Members()
			.StyleSetBack(SCIRT[style], bgnd)
		.StyleSetBack(SC.STYLE_DEFAULT, bgnd)
		}

	bold_state: false
	italic_state: false
	underline_state: false
	strikeout_state: false

	On_Bold()
		{
		if .GetReadOnly() is true
			return
		.bold_state = not .bold_state
		.Send('ToggleFontButton', 'bold', .bold_state)
		.toggleSelectedTextFormat(bold:)
		.SetFocus()
		}

	On_Italic()
		{
		if .GetReadOnly() is true
			return
		.italic_state = not .italic_state
		.Send('ToggleFontButton', 'italic', .italic_state)
		.toggleSelectedTextFormat(italic:)
		.SetFocus()
		}

	On_Underline()
		{
		if .GetReadOnly() is true
			return
		.underline_state = not .underline_state
		.Send('ToggleFontButton', 'underline', .underline_state)
		.toggleSelectedTextFormat(underline:)
		.SetFocus()
		}

	On_Strikeout()
		{
		if .GetReadOnly() is true
			return
		.strikeout_state = not .strikeout_state
		.Send('ToggleFontButton', 'strikeout', .strikeout_state)
		.toggleSelectedTextIndicator(strikethrough:)
		.SetFocus()
		}

	On_ResetFont()
		{
		// Might be able to use SCI_CLEARDOCUMENTSTYLE here
		if .GetReadOnly() is true
			return
		.bold_state = .italic_state = .underline_state = .strikeout_state =  false
		.Send('ToggleFontButton', 'bold', .bold_state)
		.Send('ToggleFontButton', 'italic', .italic_state)
		.Send('ToggleFontButton', 'underline', .underline_state)
		.Send('ToggleFontButton', 'strikeout', .strikeout_state)

		selection = .getSelection()
		.StartStyling(selection.start, .stylingMask)
		.SetStyling(selection.end - selection.start, 0)
		.SetIndicatorCurrent(.STRIKETHROUGHINDIC)
		.IndicatorClearRange(selection.start, selection.end - selection.start)
		.Dirty?(true)
		.SetFocus()
		}

	toggleSelectedTextFormat(bold = false, italic = false,
		underline = false, strikethrough = false)
		{
		selection = .getSelection()
		styleVal = .getStyleVal(bold, italic, underline, strikethrough)
		.toggleTextFormat(selection, styleVal, undoable:)
		}

	// undoable flag is set false by actions performed by undo itself
	// prevents undo action from adding itself to the undo buffer
	toggleTextFormat(selection, styleVal, undoable = false)
		{
		textFormat = .getCurrentTextFormat(selection, styleVal)
		styles = textFormat.styles
		forceStyleVal? = textFormat.forceStyleVal?

		if undoable is true
			.singleUndoableStyleAction(styles, forceStyleVal?, 0)
				{ |style, forceStyleVal?|
				.setStyle(style, forceStyleVal?)
				}
		else
			for style in styles.Members()
				.setStyle(styles[style], forceStyleVal?)

		.Dirty?(true)
		}

	getCurrentTextFormat(selection, styleVal)
		{
		text = .getStyledText(selection.start, selection.end)
		forceStyleVal? = false
		styles = Object()

		.parseStyledText(text,
			{|curStyle, pos|
			// if we find any text in the selection that does not have styleVal
			// bit enabled (not formatted with that format) - then we want to
			// force turning this on (i.e. do not turn off sections already
			// formatted with styleVal bit)
			if ((forceStyleVal? is false) and (curStyle & styleVal) is 0)
				forceStyleVal? = true
			newStart = selection.start+pos
			if not styles.Empty?()
				styles.Last().end = newStart
			styles.Add(Object(start: newStart, end: NULL, style: curStyle, :styleVal))
			}, {|chr /*unused*/, pos /*unused*/, lastStyle /*unused*/| /* do nothing */ })

		if not styles.Empty?()
			styles.Last().end = selection.end

		return Object(:styles, :forceStyleVal?)
		}

	singleUndoableStyleAction(styles, forceStyleval?, token, block)
		{
		.BeginUndoAction()
		.forUndoableStyles(styles, forceStyleval?, token, block)
		.EndUndoAction()
		}

	forUndoableStyles(styles, forceStyleVal?, token, block)
		{
		for style in styles.Members()
			{
			block(styles[style], forceStyleVal?)
			if .beginDelete is false
				.AddUndoAction(token, 1)
			.undoStyleBuffer.Push(styles[style])
			}
		}

	setStyle(style, forceStyleVal?)
		{
		.StartStyling(style.start, .stylingMask)
		.SetStyling(style.end - style.start, forceStyleVal?
			? style.style | style.styleVal : style.style ^ style.styleVal )
		}

	addDeletedFormatToBuffer(selection)
		{
		textFormat = .getCurrentTextFormat(selection, styleVal: 0)
		styles = textFormat.styles
		forceStyleVal? = textFormat.forceStyleVal?
		.forUndoableStyles(styles, forceStyleVal?, 1, {|@unused| /* do nothing */ })
		}

	toggleSelectedTextIndicator(strikethrough = false)
		{
		selection = .getSelection()

		if strikethrough is true
			{
			.SetIndicatorCurrent(.STRIKETHROUGHINDIC)
			y = .IndicatorValueAt(.STRIKETHROUGHINDIC, (selection.end-1))
			if y is 0
				{
				.IndicatorFillRange(selection.start, selection.end - selection.start)
				}
			else
				{
				if selection.end - selection.start > 0
					{
					.IndicatorClearRange(selection.start, selection.end - selection.start)
					}
				}
			.toggleSelectedTextFormat(strikethrough:)
			}
		}

	getSelection()
		{
		start = .GetSelectionStart()
		end = .GetSelectionEnd()
		return Object(:start, :end)
		}

	KEYDOWN(wParam)
		{
		if super.KEYDOWN(wParam) is 0
			return 0
		return .Eval(EditorKeyDownHandler, wParam,
			zoomArgs: [this, .zoom, ScintillaRichWordZoomControl])
		}

	// token
	// - first bit indicates remove or add -
	// 0 tells undo to remove the format
	// 1 tells undo to add the format
	// second bit indicates whether an undo action was prompted
	// from the backspace key or not
	// 0 tells undo it is a normal undo
	// 1 tells undo it came from backspace - and needs to be treated special

	token: #(
		add: 0x1
		backspace: 0x2
		)

	beginDelete: false
	preBuffer: false
	styles: false

	SCN_MODIFIED(lParam)
		{
		scn = SCNotification(lParam)
		.handleSCNotification(scn)
		return super.SCN_MODIFIED(lParam)
		}

	handleSCNotification(scn)
		{
		modType = scn.modificationType
		// if user deletes (need to distinquish between a system delete,
		// which happens on undo)
		if .userPerformedAction(modType)
			.handleUserDelete(scn, modType)
		if .nonDeleteAction(modType) and .beginDelete is true
			.beginDelete = false
		if .userPerformedUndoRedo(modType)
			.handleUserUndoRedo(scn, modType)
		else
			{
			// action comes from undo of multi-single char deletes
			if .systemPerformedUndoDelete(modType) and .styles isnt false
				.undoFormatAction(scn.token is 2
					? .styles.styleVal: .styles.style)
			}
		}

	handleUserDelete(scn, modType)
		{
		// the scintilla undo buffer behaves differently if single char is deleted
		// vs if a block of text is deleted. Need to handle each case separatly
		if ((modType & SC.MOD_BEFOREDELETE) is SC.MOD_BEFOREDELETE)
			.beforeDelete(scn)
		if ((modType & SC.MOD_DELETETEXT) is SC.MOD_DELETETEXT)
			{
			if .beginDelete is true and scn.length is 1
				.AddUndoAction(.token.add | .token.backspace, 1)
			}
		}

	beforeDelete(scn)
		{
		if scn.length is 1
			.beginDelete = true
		else
			if .beginDelete isnt true
				.BeginUndoAction()

		// put the styles from the deleted text on the stack
		// need to ignore formatting on newlines
		if not (.getText(scn.position, scn.position+scn.length,
			includeStrikeThrough:) is '\r\n')
			.addDeletedFormatToBuffer(Object(start: scn.position,
				end: scn.position + scn.length))

		if scn.length > 1 and .beginDelete isnt true
			.EndUndoAction()
		}

	handleUserUndoRedo(scn, modType)
		{
		if ((modType & SC.PERFORMED_UNDO) is SC.PERFORMED_UNDO)
			{
			.styles = .undoStyleBuffer.Pop()
			if (((scn.token & .token.add) isnt .token.add) or
				((scn.token & .token.backspace) isnt .token.backspace))
				.undoFormatAction(scn.token is 0
					? .styles.styleVal : .styles.style)
			}
		if ((modType & SC.PERFORMED_REDO) is SC.PERFORMED_REDO)
			{
			.styles = .redoStyleBuffer.Pop()
			.redoFormatAction(scn.token is 0
				? .styles.styleVal : .styles.style)
			}
		}

	undoFormatAction(styleVal)
		{
		.redoStyleBuffer.Push(.styles)
		.toggleTextFormat(Object(start: .styles.start, end: .styles.end),
			styleVal)
		.scn_HandleIndicators(.styles, .styles.style)
		.styles = false
		}

	redoFormatAction(styleVal)
		{
		.undoStyleBuffer.Push(.styles)
		.toggleTextFormat(Object(start: .styles.start, end: .styles.end),
			styleVal)
		.scn_HandleIndicators(.styles, .styles.styleVal)
		.styles = false
		}

	scn_HandleIndicators(styles, styleVal)
		{
		.SetIndicatorCurrent(.STRIKETHROUGHINDIC)
		if ((styleVal & SCIRT.STRIKETHROUGH) is SCIRT.STRIKETHROUGH)
			.IndicatorFillRange(styles.start, styles.end - styles.start)
		else
			.IndicatorClearRange(styles.start, styles.end - styles.start)
		}

	// SCN_MODIFIED ACTION TESTS
	// used to keep testing for action types in SCN_MODIFIED human readable
	userPerformedAction(modType)
		{
		return ((modType & SC.PERFORMED_USER) is SC.PERFORMED_USER)
		}

	nonDeleteAction(modType)
		{
		return ((modType & SC.MOD_DELETETEXT) isnt SC.MOD_DELETETEXT) and
			((modType & SC.MOD_BEFOREDELETE) isnt SC.MOD_BEFOREDELETE) and
			((modType & SC.MOD_CHANGEINDICATOR) isnt SC.MOD_CHANGEINDICATOR)
		}

	userPerformedUndoRedo(modType)
		{
		return ((modType & SC.MOD_CONTAINER) is SC.MOD_CONTAINER)
		}

	systemPerformedUndoDelete(modType)
		{
		return (((modType & SC.PERFORMED_UNDO) is SC.PERFORMED_UNDO) and
			((modType & SC.MOD_INSERTTEXT) is SC.MOD_INSERTTEXT))
		}
	}
