// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
PassthruController
	{
	Name: Palette
	New(.extraItems = #(), .horizontal? = false)
		{
		}
	Controls()
		{
		size = '+2'
		mouseOverImageColor = CLR.Highlight

		layout = Object(.horizontal? ? 'Horz' : 'Vert'
			Object('EnhancedButton' command: #Select image: 'select.emf', :size,
				name: 'select', tip: 'Select', mouseEffect:, :mouseOverImageColor)
			Object('EnhancedButton' command: #Line image: 'line.emf', :size,
				name: 'line', tip: 'Draw a Line',mouseEffect:, :mouseOverImageColor)
			Object('EnhancedButton' command: #Rectangle image: 'rectangle.emf',:size,
				name: 'rectangle', tip: 'Draw a Rectangle',mouseEffect:,
				:mouseOverImageColor)
			Object('EnhancedButton' command: #RoundRectangle
				image: 'rounded-rectangle.emf', :size, name: 'roundrect',
				tip: 'Draw a Round Rectangle', mouseEffect:, :mouseOverImageColor)
			Object('EnhancedButton' command: #Ellipse image: 'ellipse.emf',:size,
				name: 'ellipse', tip: 'Draw a Circle',mouseEffect:, :mouseOverImageColor)
			Object('EnhancedButton' command: #Text image: 'type.emf', :size,
				name: 'text', tip: 'Add Text',mouseEffect:, :mouseOverImageColor)
			Object('EnhancedButton' command: #Arc image: 'curve.emf', :size,
				name: 'arc', tip: 'Draw an Arc',mouseEffect:, :mouseOverImageColor)
			Object('EnhancedButton' command: #Image image: 'image.emf', :size,
				name: 'image', tip: 'Add an Image',mouseEffect:, :mouseOverImageColor))


		for item in .extraItems
			{
			layout.Add(Object('EnhancedButton', :size, mouseEffect:,
				:mouseOverImageColor).Merge(item.button))
			}

		return layout
		}
	readOnly: false
	SetReadOnly(readOnly)
		{
		.readOnly = readOnly
		super.SetReadOnly(readOnly)
		}
	Msg(args)
		{
		if args.GetDefault(0, '') is 'EnhancedButtonAllowPush'
			return not .readOnly
		if .readOnly
			return 0
		super.Msg(args)
		}
	SetButtons(buttonPushed)
		{
		parent = .horizontal? ? .Horz : .Vert
		buttons = parent.GetChildren()
		for (button in buttons)
			{
			button.Name is buttonPushed
				? button.Pushed?(true)
				: button.Pushed?(false)
			button.Repaint()
			}
		}
	}
