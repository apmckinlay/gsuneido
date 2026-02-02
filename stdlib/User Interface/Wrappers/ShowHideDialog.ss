// Copyright (C) 2006 Suneido Software Corp. All rights reserved worldwide.
// IDE > Show/Hide Override
Controller
	{
	CallClass(@unused)
		{
		Window(ShowHideDialog, title: "Show/Hide", exStyle: WS_EX.TOPMOST)
		}
	New()
		{
		switch Suneido.GetDefault('ShowAll', '')
			{
		case true:	value = 'Show'
		case false:	value = 'Hide'
		default:	value = 'Normal'
			}
		.Data.Set(Record(RadioButtons: value))
		}
	Controls:
		(Record (Border 10
			(Vert
				(RadioButtons 'Normal' 'Show' 'Hide', horz:)
				Skip
				(HorzEqual Fill (Skip 40)
					(Button 'OK') Skip (Button 'Apply') Skip (Button 'Cancel')
					xstretch: 0)
				)
			))
	On_Apply()
		{
		switch .Data.Get().RadioButtons
			{
		case 'Show' :	Suneido.ShowAll = true
		case 'Hide' :	Suneido.ShowAll = false
		default :		Suneido.Delete('ShowAll')
			}
		}
	On_OK()
		{
		.On_Apply()
		.Window.Destroy()
		}
	On_Cancel()
		{ .Window.Destroy() }
	}
