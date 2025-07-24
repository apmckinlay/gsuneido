// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
// NOTE: should test this as standalone window, dialog and inside book
// - Dialog(0, Object(TestStaticControls))
// - TestStaticControls()
// - and embed this in the book record
Controller
	{
	controls: #(Vert
		(Horz (Static 'Test Static', name: 'static'))
		(Horz (Static 'Test Static', size: '+3', weight: 'bold',
			color: 0xff0000, name: 'staticStyled')
			Skip)
		(Static 'Static New', font: '@mono', size: '+4', weight: 'bold', underline:)
		(Horz
			(Vert
				(Static 'Testing Static')
				(Static 'Testing Static')
				(Static 'Testing Static')
				(Static 'Testing Static')
				#(EtchedLine before: 0)
				))
		(Horz (EtchedVertLine ystretch: 0)
			(Static  'Testing Static') (EtchedVertLine ystretch: 0))

		(Horz (Static  'Testing right justify', weight: semibold
			color: 0x808080, justify: RIGHT, xmin: 200)
			(EtchedVertLine ystretch: 0))

		(Static  'Testing background color', weight: semibold
			color: 0x808080, justify: RIGHT, xmin: 200, bgndcolor: 0x00cbc0ff)

		(Static  'Testing white background', whitebgnd:)

		(Static  'Testing tab\ttab') // not handling tab properly on both controls

		(Pair (Static 'Field Selectable', name: 'new_static')
			(Field set: 'Field Selectable in Pair'))
		Skip
		(Horz
			(Vert
				(HeadingControl 'Heading')
				(HeadingControl 'Heading')
				(HeadingControl 'Heading')
				(HeadingControl 'Heading'))
		)
		(Horz
			Skip
			(Vert
				(Horz #(Skip 25) Fill
					(Title, 'Orders')
					Skip #(Skip 25) Fill
					['TitleButtons', #('Orders')])
				(Horz #(Skip 25) Fill
					(Title, 'Tickets')
					Skip #(Skip 25) Fill
					['TitleButtons', #('Tickets')])
			)
		)
		(Horz
			(Vert
				(Static 'Testing Static\r\nWith multi lines\n\nAnd empty lines\n\nmore')
				#(EtchedLine before: 0)
				)
			)
		(Record, (Horz
			(Vert
				(Highlight (Pair (Static 'Testing Highlight'),
					(StaticText textStyle: 'note', name: 'highlight')), size: '+1'),
				(EtchedLine before: 0))
			Skip
			(Vert
				(Highlight (Pair (Static 'Enter value to highlight'), (Field))),
				(EtchedLine before: 0))),
			name: 'highlightRec')
		(Horz
			(WndPane
				('Vert',
					("Static",  'test inside pane with xmin/ymin', name: 'paneStatic',
						whitebgnd: false, xmin: 180, ymin: 30)
					(EtchedLine before: 0)
					)))
		('Horz',
			("Static",  'test set', name: 'testSet')
			("Static",  'test set', name: 'testSet2')
			)
		)

	testWrap()
		{
		longStr = 'If you have time to enter a short description of what you were ' $
			'doing when this problem occurred, it will help us improve the ' $
			'software. Thank you.'
		multiLine = "The following print options have been moved\n" $
			"   Print Bill To\n" $
				"   Print Totals\n" $
				"   Print Signature\n" $
				"   Print BOL Other Charges\n" $
				"   Print ACE Commodities\n" $
				" \n"
		return Object(#Horz
			Object(#Vert
				#('Static', 'Testing long Line', textStyle: 'main'),
				Object('StaticWrap', longStr, xstretch: 1, xmin: 300)
				#('Static', 'Testing multi lines', textStyle: 'main'),
				Object('StaticWrap', multiLine, size: '+1' xmin: 300)
				#(EtchedLine before: 0)
				)
			)
		}

	testLinkButton()
		{
		return Object(#Horz
			Object(#Vert
				#('LinkButton', 'Testing Link Button'),
				#('LinkButton', 'Reset Columns', "Reset_Columns"),
				#(EtchedLine before: 0)
				)
			)
		}

	New()
		{
		super(.layout())
		.highlightRec = .FindControl('highlightRec')
		.FindControl('testSet').Set(200 /*= setting a text field to a number */)
		.FindControl('testSet2').Set(' OK Selectable')
		}

	layout()
		{
		ctrls = .controls.Copy()
		ctrls.Add(.testWrap())
		ctrls.Add(.testLinkButton())
		ctrls.Add(#Skip)
		ctrls.Add(#(Horz
			#(Button 'Get Text'), Skip
			#(Button 'Set Text') Skip
			#(Button 'Set Color') Skip
			#(Button 'Toggle Highlight') Skip
			#(Button 'Refresh Window')
			))
		return Object(#Scroll, ctrls)
		}

	getter_static()
		{
		.FindControl('new_static')
		}

	On_Toggle_Highlight()
		{
		data = .highlightRec.Get()
		data.highlight = data.highlight is '' ? 'Selected' : ''
		}

	On_Get_Text()
		{
		.AlertInfo('Get Text', .static.Get())
		}

	On_Set_Text()
		{
		.static.Set('Field New')
		.FindControl('paneStatic').Set('test inside pane, text changed')
		if .FindControl('static').Get().Has?('\n')
			{
			.FindControl('static').Set('Test Static Set')
			.FindControl('staticStyled').Set('Test Static Set')
			}
		else
			{
			txt = 'Test Static Set\nTest Static Set'
			.FindControl('static').Set(txt)
			.FindControl('staticStyled').Set(txt)
			}
		}

	On_Set_Color()
		{
		.static.SetColor(0x0000ff) /*= red */
		}

	On_Reset_Columns()
		{
		.AlertInfo("Testing New Static", 'Reset Columns Clicked')
		}

	On_Reset_Columns_ContextMenu(x, y)
		{
		.AlertInfo("Testing New Static", 'Reset Columns Context Menu ' $ x $ '-' $ y)
		}

	On_Refresh_Window()
		{
		.WindowRefresh()
		}
	}
