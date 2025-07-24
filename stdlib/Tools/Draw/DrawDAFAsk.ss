// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Xmin: 100
	Ymin: 200
	Title: 'Data Field'
	CallClass(canvas, value = '', font = false, justify = 'left', showPrompt? = false)
		{
		return OkCancel(Object(this, canvas, value, font, justify, showPrompt?), .Title)
		}

	New(.canvas, value, font, justify, showPrompt?)
		{
		super(.layout(showPrompt?))
		.data = .Data.Get()
		if value isnt '' and false isnt prompt = .canvas.Send('FieldToPrompt', value)
			.data.value = prompt
		.data.alignment = justify.Capitalize()
		.data.Merge(DrawTextAsk.ConvertFontOb(font))
		}

	layout(showPrompt?)
		{
		cols = .canvas.Send('GetDAFAvailableCols')
		layout = Object('Record', Object('Vert'
				Object('ChooseList', cols, mandatory:, name: 'value')
				Object('Skip' amount: 20)
				#('GroupBox', 'Font', ('Vert'
					('ChooseList',
						#('Arial', 'Times New Roman', 'Courier New'),
						mandatory:, name: 'fontName', set: 'Arial')
					('ChooseList',
						#('Regular', 'Italic', 'Bold', 'Bold Italic'),
						mandatory:, name: 'fontStyle', set: 'Regular')
					('Spinner', rangefrom: 4, rangeto: 72,
						mandatory:, name: 'fontSize', set: 10)
					))
				Object('RadioButtons' 'Left' 'Center' 'Right' label: 'Alignment:'
					horz:, name: 'alignment')
			))
		if showPrompt?
			layout[1].Add(Object('CheckBox' 'Show Prompt' name: 'showPrompt'), at: 2)
		return layout
		}

	OK()
		{
		if .Data.Valid() isnt true
			return false
		field = .canvas.Send('PromptToField', .data.value)
		return Object(:field, prompt: .data.value, font: DrawTextAsk.BuildFontOb(.data),
			showPrompt: .data.showPrompt, justify: .data.alignment.Lower())
		}
	}
