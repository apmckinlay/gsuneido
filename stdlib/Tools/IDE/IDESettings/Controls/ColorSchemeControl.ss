// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: "Color Scheme Creator"
	Name: #ide_color_scheme
	Xmin: 555
	Ymin: 200
	New()
		{
		.preview = .FindControl("preview")
		.previewHtml = .FindControl("previewHtml")
		.data = .Data.Get()
		.setData(IDE_ColorScheme.GetTheme())
		.preview.SetReadOnly(true)
		.previewHtml.SetReadOnly(true)
		.data.Observer(.Record_Changed)
		}

	Controls()
		{
		.testTheme = IDE_ColorScheme.GetTheme().Copy()
		themes = Object()
		Plugins().ForeachContribution('ColorSchemes', 'theme', { themes.Add(it.name) })
		themesList = Object(#Pair, #(Static #Themes),
			Object(#ChooseList, themes.Sort!().Add('default', at: 0), width: 14,
			name: #themeList, set: IDE_ColorScheme.GetTheme().name))

		vert = Object(#Vert, themesList)
		.map.Each({ vert.Add(.buildPair(it)) })
		vert.Add(#Fill, #Skip, #(HorzEqual (Button "Output")))

		return Object(#Record, Object(#Vert,
		Object(#Horz
			vert,
			#Skip
			Object(#Vert
				Object(#WorkSpaceCode, name: "preview", scheme: .testTheme,
					Addon_show_margin:)
				Object(#ScintillaAddons, name: "previewHtml", scheme: .testTheme,
					Addon_html:)
			))))
		}

	buildPair(ob)
		{
		return Object(#Pair, Object(#Static, ob.prompt),
			Object(#ChooseColor, name: ob.ctrl))
		}

	On_Output()
		{
		themeOb = Object()
		for pair in .map
			pair.schemeFields.Each({ themeOb[it] = '0x' $ .data[it].Hex() })
		Print(themeOb)
		}

	map: (
		(ctrl: background, prompt: 'Background', schemeFields: (defaultBack, whitespace))
		(ctrl: text, prompt: 'Text',
			schemeFields: (defaultFore, braceGoodFore, braceBadFore, lineNumberFore))
		(ctrl: comment, prompt: 'Comment', schemeFields: (comment))
		(ctrl: number, prompt: 'Number', schemeFields: (number))
		(ctrl: string, prompt: 'String', schemeFields: (string))
		(ctrl: keyword, prompt: 'Keyword', schemeFields: (keyword, operator))
		(ctrl: cursorLine, prompt: 'Cursor Line', schemeFields: (cursorLine))
		(ctrl: selectedBack, prompt: 'Selected Back', schemeFields: (selectedBack))
		(ctrl: occurrence, prompt: 'Occurrence', schemeFields: (occurrence))
		(ctrl: braceGoodBack, prompt: 'Brace Good Back', schemeFields: (braceGoodBack))
		(ctrl: braceBadBack, prompt: 'Brace Bad Back', schemeFields: (braceBadBack))
		(ctrl: lineNumberBack, prompt: 'Line Number Back', schemeFields: (lineNumberBack))
		(ctrl: foldMargin, prompt: 'Fold Margin', schemeFields: (foldMargin))
		(ctrl: longLineMargin, prompt: 'Long Line Margin', schemeFields: (longLineMargin))
		(ctrl: warning, prompt: 'Warning', schemeFields: (warning))
		(ctrl: error, prompt: 'Error', schemeFields: (error))
		(ctrl: tag, prompt: 'Tag (HTML)', schemeFields: (tag))
		(ctrl: unknownTag, prompt: 'Unknown Tag (HTML)' schemeFields: (unknownTag))
		(ctrl: attr, prompt: 'Attribute (HTML)', schemeFields: (attr))
		(ctrl: unknownAttr, prompt: 'Unknown Attribute (HTML)',
			schemeFields: (unknownAttr))
		(ctrl: insideTag, prompt: 'Inside Tag (HTML)' schemeFields: (insideTag))
		(ctrl: unquotedVal, prompt: 'Unquoted Value (HTML)', schemeFields: (unquotedVal))
		(ctrl: sgml, prompt: 'SGML (HTML)', schemeFields: (sgml))
		(ctrl: entity, prompt: 'Entity (HTML)', schemeFields: (entity))
		)
	updatingExample: false
	Record_Changed(member)
		{
		if member is #preview or member is #previewHtml or .updatingExample is true
			return
		.updatingExample = true
		if member is #themeList
			.setData(IDE_ColorScheme.GetTheme(.FindControl(#themeList).Get()))
		else
			.setExample(member)
		.preview.ResetAddons()
		.previewHtml.ResetAddons()
		.updatingExample = false
		}

	setExample(field)
		{
		idx = .map.FindIf({ it.ctrl is field })
		.map[idx].schemeFields.Each({ .testTheme[it] = .FindControl(field).Get() })
		}

	setData(theme)
		{
		.map.Each()
			{
			color = IDE_ColorScheme.GetColor(it.schemeFields[0], theme.name)
			for field in it.schemeFields
				.testTheme[field] = .data[field] = color
			if false isnt colorField = .FindControl(it.ctrl)
				colorField.Set(color)
			}
		.preview.Set(.exampleText)
		.previewHtml.Set(.exampleHtml)
		}

	exampleText: `/* multi-line
	comment */
// ----------------------------------------- long line ` $
`------------------------------------------
calc = function (x, y = 0, dbl = false)
	{
	sum = x + y // single line comment
	if dbl is true
		sum *= 2
	return "result is " $ sum)
	}`
	exampleHtml: `<div class="main">
<p>Here is some text with an entity &amp;</p>
<badtag></badtag>
<p badattr="foo"></p>
<p class=unquoted></p>
</div>`

	Get()
		{ return .FindControl(#themeList).Get() }

	PostSave()
		{ IDE_ColorScheme.ResetStyles() }
	}
