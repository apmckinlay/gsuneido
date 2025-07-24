// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: 'Color keywords in CLR'
	New()
		{
		super(.controls())
		.Vert.Status.Set('Double click on a color rectangle to see color')
		.horz = .Vert.GroupBox.Horz
		}
	numCols: 4
	controls()
		{
		clrnames = CLR.Members().Sort!()
		maxclrnames = clrnames.Size() - 1
		form = Object('Form')
		for (i = 0; i < clrnames.Size(); i += .numCols)
			{
			group = 0
			row1 = Object()
			row2 = Object()
			clr = i
			for (x = 0; x < .numCols; x += 1)
				{
				if clr > maxclrnames
					break
				color = CLR[clrnames[clr]]
				s = .descolor(color)
				row2.Add(Object('Field', set: s, readonly:, width: 26, :group))

				row1.Add(Object('ColorRect', :color, rounded:,
					xmin: 49, ymin: 24, :group))
				row1.Add(Object('Field', set: clrnames[clr], readonly:,
					width: 16, group: ++group))
				++group
				++clr
				}
			.addColorRow(form, row1, row2)
			if clr > maxclrnames
				break
			}
		return Object('Vert'
			Object('Pane'
				Object('Scroll'
					Object('Vert' 'Skip', form, 'Skip'))),
			'Skip', .chooseColorCtrl, 'Skip', 'Statusbar')
		}

	addColorRow(form, row1, row2)
		{
		for r in row1
			form.Add(r)
		form.Add('nl')
		for r in row2
			form.Add(r)
		form.Add('nl')
		form.Add(Object('Static' '' group: 0))
		form.Add('nl')
		}

	chooseColorCtrl: (GroupBox 'choose a color' (Horz
		(ChooseColor tip: 'choose a color')
		(Skip 1) (Static 'press Choose' name: Stxt1)
		(Static '... button' name: Stxt2)
		(Skip 1)
		(Static ' to change ' ymin: 19 name: Stxt3)
		(Skip 1)
		(Field set: 'text color', name: Stxt4)
		(Skip 1)
		(EnhancedButton 'Button', buttonStyle:, mouseEffect:)
		(Field set: '', readonly:, name: Fcol width: 41))
		)

	descolor(color)
		{
		rgb = color.ToRGB()
		dec = RGB(rgb[0], rgb[1], rgb[2])
		return 'RGB(' $ rgb[0] $ ', ' $ rgb[1] $ ', ' $ rgb[2] $ ') - ' $
			dec $ ' - ' $ '0x' $ dec.Hex().LeftFill(6, '0') /*= number of hex digit */
		}
	NewValue(value, source){
		if source.Name is 'ChooseColor'
			{
			.settext(value)
			}
		}
	settext(color)
		{
		.horz.Stxt1.SetBgndColor(color)
		.horz.Stxt2.SetBgndColor(color)
		.horz.Stxt2.SetColor(CLR.WHITE)
		.horz.Stxt3.SetColor(color)
		.horz.Stxt4.SetTextColor(color)
		.horz.ChooseColor.Set(color)
		.horz.Button.SetTextColor(color)
		s = (CLR.Find(color) is false) ? 'not in CLR - ' : CLR.Find(color) $ ' - '
		.horz.Fcol.Set(s $ .descolor(color))
		}
	ColorRect_DoubleClick(color)
		{
		.settext(color)
		}
	}