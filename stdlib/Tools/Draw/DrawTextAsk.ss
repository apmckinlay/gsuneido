// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Xmin: 800
	Ymin: 400
	CallClass(title = "", content = "", font = false, hwnd = 0, justify = 'left')
		{
		return OkCancel(Object(this, content, font, justify), title, hwnd)
		}

	New(content, font, justify = 'left')
		{
		.data = .Data.Get()
		.data.alignment = justify.Capitalize()
		.data.text = content
		.data.Merge(.ConvertFontOb(font))
		}

	ConvertFontOb(font)
		{
		if not Object?(font) or font.GetDefault('name', '') is ''
			font = Object(name: 'Arial', size: 10)
		return [
			fontName: font.name
			fontStyle: .getFontStyleFromOb(font)
			fontSize: font.size]
		}

	getFontStyleFromOb(font)
		{
		italic? = font.GetDefault('italic', false) is true
		bold? = font.GetDefault('weight', FW.NORMAL) > FW.NORMAL
		if bold? and italic?
			return 'Bold Italic'
		else if bold?
			return 'Bold'
		else if italic?
			return 'Italic'
		return 'Regular'
		}
	Controls()
		{
		.infoMenu = GetContributions('DrawText_ExtraLabels').Map({ it.label }).
			Add('Page#', 'Short Date', 'Long Date')
		vert = Object('Vert'
			Object('MenuButton', 'Insert', .infoMenu, xstretch: 0)
			Object('GroupBox', 'Font'
				Object('Vert'
					Object('ChooseList',
						#('Arial', 'Times New Roman', 'Courier New'),
						mandatory:, name: 'fontName')
					Object('ChooseList',
						#('Regular', 'Italic', 'Bold', 'Bold Italic'),
						mandatory:, name: 'fontStyle')
					Object('Spinner', rangefrom: 4, rangeto: 72,
						mandatory:, name: 'fontSize')
					)
			)
			Object('RadioButtons' 'Left' 'Center' 'Right' label: 'Alignment:'
				horz:, name: 'alignment')
			)
		return Object(#Record,
			Object(#Vert
				Object(#Horz #('Editor', name: 'text') #Skip vert)
				))
		}
	On_Insert(args)
		{
		label = '<' $ args $ '>'
		.FindControl('text').ReplaceSel(label)
		.data.text = .FindControl('text').Get()
		}

	OK()
		{
		if .Data.Valid() isnt true
			return false
		return Object(text: .data.text, font: .BuildFontOb(.data),
			justify: .data.alignment.Lower())
		}

	BuildFontOb(data)
		{
		font = Object()
		font.name = data.fontName
		font.italic = data.fontStyle.Has?('Italic')
		font.weight = data.fontStyle.Has?('Bold') ? FW.SEMIBOLD : FW.NORMAL
		font.size = data.fontSize
		return font
		}
	}
