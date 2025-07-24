// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
// TODO: allow ystretch, this may fix MaxHeight issues
//		set .Ystretch in Recalc ? (ystretch should be 0 if closed)
PassthruController
	{
	Name: "Expand"
	ComponentName: 'Expand'
	Ystretch: 0
	MaxHeight: 8888 // must be different from Control.MaxHeight
	CalcMaxHeight()
		{ return 8888 }

	New(label, control, open = false, leftAlign = false, .saveExpandName = false,
		.hideLine = false)
		{
		super(.make_controls(label, control, open, leftAlign))
		.ocb = .FindControl('ocbutton')
		.summary = .FindControl('summary')
		.pane = .FindControl('ExpandPane')
		.control = .pane.Vert.GetChildren()[1]
		}
	Startup()
		{
		.set_summary()
		if 0 isnt .rc = .Send('GetRecordControl')
			{
			.rc.AddObserver(.set_summary)
			.rc.AddSetObserver(.set_summary)
			}
		.loadSavedExpandSetting()
		}

	loadSavedExpandSetting()
		{
		if .saveExpandName is false
			return

		open = UserSettings.Get('ExpandControl - ' $ .saveExpandName, def: .open)
		if open is true
			.Expand(true)
		else
			.Contract()
		}

	make_controls(label, control, open, leftAlign)
		{
		.Label = label
		.open = open

		return Object('Vert',
			.hideLine ? '' : Object('EtchedLine'),
			Object('Horz'
				#(Skip 4)
				Object('Vert'
					'Fill'
					Object('EnhancedButton',
						command: 'OpenClose',
						image: open ? 'next.emf' : 'forward.emf',
						imageColor: 0x737373, mouseOverImageColor: 0x00cc00,
						imagePadding: 0.15,
						buttonHeight: ScaleWithDpiFactor(19),
						buttonWidth: ScaleWithDpiFactor(19),
						tip: .tip,
						name: 'ocbutton')
					'Fill'
					ystretch: 0)
				#(Skip 4)
				Object('Vert', #(Skip 2), .text_layout(label, leftAlign))
				#(Skip 4)
				)
			Object('ExpandPane', .open,
				Object('Vert', #(Skip 6) control, #(Skip 4)))
			)
		}
	tip: 'click to expand/contract the section'

	text_layout(label, leftAlign)
		{
		label = Object('Static' label weight: 'bold', tip: .tip, name: 'tabname')
		summary = Object('Static' '' tip: .tip, xstretch: 1, name: 'summary')
		if leftAlign is false
			summary.justify = 'RIGHT'
		return Object('Horz', label, #(Skip 30), summary)
		}

	On_OpenClose()
		{
		if not .open
			.Expand()
		else
			.Contract()
		}
	SourceFromHeading?(evtSource)
		{
		return evtSource.Name in (#summary, #tabname)
		}
	Static_LButtonUp(source)
		{
		if .SourceFromHeading?(source)
			.On_OpenClose()
		return 0
		}
	ExpandPane_Click()
		{ .On_OpenClose() }

	Expand(skipFocusFirst? = false)
		{
		if .open
			return
		.open = true
		.pane.SetOpen(true)
		if not skipFocusFirst?
			.pane.FocusFirst(.pane.Hwnd)
		.ocb.SetImage('next.emf')
		.summary.Set('')
		.Window.Refresh()
		.Send('Expand')
		}

	Contract()
		{
		if not .open
			return
		.open = false
		.pane.SetOpen(false)
		.ocb.SetImage('forward.emf')
		.set_summary()
		.Window.Refresh()
		.Send('Contract')
		}
	open: true // so Delayed set_summary won't do anything if control is destroyed
	set_summary()
		{
		if .open
			return
		summary = .Send("NeedSummary")
		if summary is 0 // no method
			if .control.Method?(#MakeSummary)
				summary = .control.MakeSummary()
		summary = String?(summary) and summary isnt "" ? summary : " "
		// NOTE: need to set summary to something (not "")
		// or else control doesn't show up correctly
		.summary.Set(summary)
		}

	control: false
	Recalc()
		{
		.Ystretch = .open and .control isnt false ? .control.Ystretch : 0
		super.Recalc()
		}

	saveExpandSetting()
		{
		if .saveExpandName is false
			return
		UserSettings.Put('ExpandControl - ' $ .saveExpandName, .open)
		}

	rc: 0
	Destroy()
		{
		.saveExpandSetting()
		if .rc isnt 0
			.rc.RemoveSetObserver(.set_summary)
		super.Destroy()
		}
	}
