// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
PassthruController
	{
	Name: Color
	Controls: (Horz
		#(EnhancedButton command: 'N', image: 'null.bmp', mouseDownImage: 'nullin.bmp',
			name: 'N')
		#(EnhancedButton command: 'BL', image: 'black.bmp', mouseDownImage: 'blackin.bmp',
			name: 'BL')
		#(EnhancedButton command: 'B', image: 'blue.bmp', mouseDownImage: 'bluein.bmp',
			name: 'B')
		#(EnhancedButton command: 'G', image: 'green.bmp', mouseDownImage: 'greenin.bmp',
			name: 'G')
		#(EnhancedButton command: 'W', image: 'white.bmp', mouseDownImage: 'whitein.bmp',
			name: 'W')
		#(EnhancedButton command: 'R', image: 'red.bmp', mouseDownImage: 'redin.bmp',
			name: 'R')
		#(EnhancedButton command: 'custom', image: 'custom.bmp',
			mouseDownImage: 'customin.bmp', name: 'custom')
		Skip
		#(Static 'LINE' weight: bold, name: line)
		Skip
		#(Static 'FILL' weight: bold, name: fill)
		)
	SetButtons(buttonPushed)
		{
		buttons = .Horz.GetChildren()
		for (button in buttons)
			{
			if not button.Base?(EnhancedButtonControl)
				continue
			button.Name is buttonPushed
				? button.Pushed?(true)
				: button.Pushed?(false)
			button.Repaint()
			}
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
	// Color functions
	On_N()
		{
		.SetButtons('N')
		.Horz.fill.SetColor(CLR.WHITE)
		.Send('SetColor', false)
		}
	On_BL()
		{
		.SetButtons('BL')
		.Horz.fill.SetColor(CLR.BLACK)
		.Send('SetColor', CLR.BLACK)
		}
	On_B()
		{
		.SetButtons('B')
		.Horz.fill.SetColor(CLR.BLUE)
		.Send('SetColor', CLR.BLUE)
		}
	On_G()
		{
		.SetButtons('G')
		.Horz.fill.SetColor(CLR.GREEN)
		.Send('SetColor', CLR.GREEN)
		}
	On_W()
		{
		.SetButtons('W')
		.Horz.fill.SetColor(CLR.WHITE)
		.Send('SetColor', CLR.WHITE)
		}
	On_R()
		{
		.SetButtons('R')
		.Horz.fill.SetColor(CLR.RED)
		.Send('SetColor', CLR.RED)
		}
	custColors: false
	On_custom()
		{
		if (false isnt (color = .custom_colors(.Horz.fill.GetColor())))
			{
			.Horz.fill.SetColor(color)
			.Send('SetColor', color)
			}
		}
	custom_colors(prevcolor = 0)
		{
		if .custColors is false
			.custColors = Object()
		.SetButtons('custom')
		color = false
		if false isnt result = ChooseColorWrapper(prevcolor, .Window.Hwnd,
			custColors: .custColors)
			color = result is "" ? 0 : result

		return color
		}
	EnhancedRButtonDown(command)
		{
		name = command[3..]
		color = CLR.BLACK
		switch (name)
			{
		case 'N' :
			color = false
		case 'BL' :
			color = CLR.BLACK
		case 'B' :
			color = CLR.BLUE
		case 'G' :
			color = CLR.GREEN
		case 'W' :
			color = CLR.WHITE
		case 'R' :
			color = CLR.RED
		case 'custom' :
			if false is color = .custom_colors(.Horz.line.GetColor())
				return
			}
		.SetButtons(name)
		.Horz.line.SetColor(color is false ? CLR.WHITE : color)
		.Send('SetLineColor', color)
		}
	}
