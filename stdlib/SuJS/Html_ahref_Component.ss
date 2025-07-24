// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
Component
	{
	Font: ""
	Size: ""
	Weight: "normal"
	Underline: true
	Italic: false
	Xstretch: false
	styles: `
		.su-html-ahref {
			color: blue;
			cursor: pointer;
		}
		.su-html-ahref:hover {
			color: highlight;
		}`
	New(text)
		{
		LoadCssStyles('su-html-ahref.css', .styles)
		.CreateElement('a', :text, className: 'su-html-ahref')

		.SetFont(.Font, .Size, .Weight, .Underline, .Italic)
		metrics = SuRender().GetTextMetrics(.El, text)
		.Xmin = .Xmin isnt 0 ? .Xmin : metrics.width
		.Ymin = .Ymin isnt 0 ? .Ymin : metrics.height
		.SetMinSize()

		.El.AddEventListener('click', .click)
		}

	click()
		{
		.Event('LBUTTONUP')
		}
	}
