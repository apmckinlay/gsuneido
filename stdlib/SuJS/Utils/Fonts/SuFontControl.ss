// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
Controller
	{
	CallClass(font)
		{
		ToolDialog(0, Object(this, font), 'Font')
		}

	weights: (
		Regular: 400,
		Medium: 500,
		Bold: 700,
		Extrabold: 800)
	New(.font)
		{
		if not .font.Member?(#lfWeight) or false is style = .weights.Find(.font.lfWeight)
			style = 'Regular'
		if .font.GetDefault(#lfItalic, 0) isnt 0
			style = style is 'Regular' ? 'Italic' : style $ ' Italic'
		rec = [
			sufont_font: .font.lfFaceName,
			sufont_style: style,
			sufont_size: StdFonts.PtSize(.font.lfHeight)]
		.Data.Set(rec)
		.sample = .FindControl('sample')
		.updateSample()
		}

	Controls()
		{
		return Object('Record', Object('Vert',
			'sufont_font',
			'sufont_style',
			'sufont_size',
			#(GroupBox, 'Sample', #('Static', name: 'sample'))
			#Skip,
			#OkCancel))
		}

	Record_NewValue(@unused)
		{
		.convert()
		.updateSample()
		}

	updateSample()
		{
		.sample.Set('AaBbYyZz', .font, refreshRequired?:)
		}

	On_OK()
		{
		if not .Data.Dirty?()
			{
			.On_Cancel()
			return
			}

		if .Data.Valid() isnt true
			return

		.convert()
		.Window.Result(.font)
		}

	convert()
		{
		rec = .Data.Get()
		.font.lfFaceName = rec.sufont_font
		.font.fontPtSize = rec.sufont_size
		.font.lfHeight = StdFontsSize.LfSize(rec.sufont_size, WinDefaultDpi)
		.font.lfItalic = rec.sufont_style.Has?('Italic') ? 1 : 0
		if false is weight = .weights.Members().FindOne({ rec.sufont_style.Has?(it) })
			weight = 'Regular'
		.font.lfWeight = .weights[weight]
		}

	On_Cancel()
		{
		.Window.Result(false)
		}
	}