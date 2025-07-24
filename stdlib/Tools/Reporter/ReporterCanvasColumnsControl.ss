// Copyright (C) 2023 Axon Development Corporation All rights reserved worldwide.
// This is a dummy class to work with Reporter, as ReporterColumns class cannot be used
// with ReporterCanvas
Controller
	{
	Name: 'ReporterCanvasColumns'
	Controls()
		{
		return #('Vert')
		}

	SetCanvas(.canvas) {}

	Get()
		{
		columns = Object()
		for col in .canvas.GetDAFs()
			if false isnt prompt = .canvas.FieldToPrompt(col)
				columns.Add(Object(
					text: prompt,
					width: 11 /* = canvas doesn't need width,
						but in case this is loaded into the design tab*/))
		return columns
		}

	Set(unused) {}
	ContextMenu(x /*unused*/, y /*unused*/) {}
	}
