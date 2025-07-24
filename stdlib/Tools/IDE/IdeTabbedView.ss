// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
// Used by FindReferencesControl and VersionHistoryControl
Controller
	{
	CallClass(ctrlname, name, extra = false)
		{
		sname = 'IdeTabbedView_' $ ctrlname
		if false is me = Suneido.GetDefault(sname, false)
			{
			me = Window([this, ctrlname, name, extra], keep_placement:).Ctrl
			Suneido[sname] = me
			}
		else
			me.Show(name, extra)
		return me
		}

	images: false
	New(.ctrlname, name, extra)
		{
		super(['Border',
			['Tabs', .control(ctrlname, name, extra), close_button: 1,
				selectedTabColor: IDESettings.Get('ide_selected_tab_color', false),
				selectedTabBold: IDESettings.Get('ide_selected_tab_bold', true)],
			border: 3])
		.tabs = .FindControl('Tabs')
		.Title = .tabs.GetControl().Title
		.tabs.SetTabData(0, [tooltip: .tooltip(name, extra)])
		}

	control(ctrlname, name, extra)
		{
		x = [ctrlname, name, Tab: name]
		if extra isnt false
			x.Add(extra)
		return x
		}

	tooltip(name, extra)
		{ return Opt(extra is false ? '' : extra, ':') $ name }

	Show(name, extra)
		{
		tooltip = .tooltip(name, extra)
		if false isnt idx = .tabs.FindTabBy('tooltip', tooltip)
			.tabs.Select(idx)
		else
			.tabs.Insert(name, .control(.ctrlname, name, extra), image: 0,
				data: [:tooltip])
		WindowActivate(.Window.Hwnd) // in case it's minimized
		}

	Tab_Close(tab)
		{
		if .tabs.GetTabCount() > 1
			{
			if tab is .tabs.GetSelected()
				.tabs.Select(tab is 0 ? 1 : tab - 1)
			.tabs.Remove(tab)
			}
		else // last tab
			DestroyWindow(.Window.Hwnd)
		}

	Destroy()
		{
		if .images isnt false
			ImageList_Destroy(.images)
		name = 'IdeTabbedView_' $ .ctrlname
		Suneido.Delete(name)
		super.Destroy()
		}
	}
