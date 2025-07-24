// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
ScintillaAddon
	{
	Init()
		{
		.SetILexer(0/*unused by Scintilla*/, CreateLexer('hypertext'))
		.SetWordChars("-_0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
		.SendMessageTextIn(SCI.SETKEYWORDS, 0, .elements $ ' ' $ .attributes)

		// base font color (fore/back)
		defaultBack = .GetSchemeColor('defaultBack')
		defaultFore = .GetSchemeColor('defaultFore')

		.DefineStyle(0, defaultFore, back: defaultBack)
		.DefineStyle(1, .GetSchemeColor('tag'), back: defaultBack)
		.DefineStyle(2, .GetSchemeColor('unknownTag'), back: defaultBack)
		.DefineStyle(3, .GetSchemeColor('attr'), back: defaultBack)
		.DefineStyle(4, .GetSchemeColor('unknownAttr'), back: defaultBack)
		.DefineStyle(5, .GetSchemeColor('number'), back: defaultBack)
		// double
		.DefineStyle(6, .GetSchemeColor('string'), back: defaultBack)
		// single
		.DefineStyle(7, .GetSchemeColor('string'), back: defaultBack)
		.DefineStyle(8, .GetSchemeColor('insideTag'), back: defaultBack)
		.DefineStyle(9, .GetSchemeColor('comment'), back: defaultBack)
		.DefineStyle(10, .GetSchemeColor('entity'), back: defaultBack)
		.DefineStyle(11, .GetSchemeColor('empty'), back: defaultBack)
		.DefineStyle(19, .GetSchemeColor('unquotedVal'),
			back: defaultBack)
		for (i = 21; i <= 31; ++i)
			.DefineStyle(i, .GetSchemeColor('sgml'), back: defaultBack)
		.DefineStyle(SC.STYLE_DEFAULT, defaultFore, back: defaultBack)

		.SetCaretFore(defaultFore)
		.SetSelBack(true, .GetSchemeColor('selectedBack'))

		.SetWrapMode(SC.WRAP_WORD)
		.SetWrapIndentMode(SC.WRAPINDENT_SAME)
		}
	elements: "a abbr acronym address applet area b base basefont
		bdo big blockquote body br button caption center
		cite code col colgroup dd del dfn dir div dl dt em
		fieldset font form frame frameset h1 h2 h3 h4 h5 h6
		head hr html i iframe img input ins isindex kbd label
		legend li link map menu meta noframes noscript
		object ol optgroup option p param pre q s samp
		script select small span strike strong style sub sup
		table tbody td textarea tfoot th thead title tr tt u ul
		var xml xmlns"
	attributes: "abbr accept-charset accept accesskey action align alink
		alt archive axis background bgcolor border
		cellpadding cellspacing char charoff charset checked cite
		class classid clear codebase codetype color cols colspan
		compact content coords
		data datafld dataformatas datapagesize datasrc datetime
		declare defer dir disabled enctype event
		face for frame frameborder
		headers height href hreflang hspace http-equiv
		id ismap label lang language leftmargin link longdesc
		marginwidth marginheight maxlength media method multiple
		name nohref noresize noshade nowrap
		object onblur onchange onclick ondblclick onfocus
		onkeydown onkeypress onkeyup onload onmousedown
		onmousemove onmouseover onmouseout onmouseup
		onreset onselect onsubmit onunload
		profile prompt readonly rel rev rows rowspan rules
		scheme scope selected shape size span src standby start style
		summary tabindex target text title topmargin type usemap
		valign value valuetype version vlink vspace width
		text password checkbox radio submit reset
		file hidden image"
	}
