// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_ReportQuery()
		{
		mock = Mock(PresetsControl)
		mock.When.ReportQuery([anyArgs:]).CallThrough()
		mock.When.Report([anyArgs:]).CallThrough()
		mock.PresetsControl_report_name = 'Report Test'

		prefix = 'params where report is '
		Assert(mock.ReportQuery('test') is: prefix $ '"Report Test~presets~test"')
		Assert(mock.ReportQuery('An Odd Name')
			is: prefix $ '"Report Test~presets~An Odd Name"')
		Assert(mock.ReportQuery('') is: prefix $ '"Report Test~presets~"')
		}

	Test_MenuButton_Presets()
		{
		mock = Mock(PresetsControl)
		mock.When.MenuButton_Presets().CallThrough()
		mock.PresetsControl_initial = #(New)
		mock.PresetsControl_extraMenu = #()
		mock.PresetsControl_presets = presets = Object()
		mock.PresetsControl_listTrimmed? = false

		mb = Mock()
		mb.When.Get().Return('Preset 3')
		mock.PresetsControl_mb = mb
		Assert(mock.MenuButton_Presets() is: #(New, 'Manage...', '', 'Save As...'))

		presets.Add('Preset 1', 'Preset 2')
		Assert(mock.MenuButton_Presets()
			is: #(New, 'Manage...', '', 'Preset 1', 'Preset 2', '', 'Save As...'))

		presets.Add('Preset 3')
		Assert(mock.MenuButton_Presets()
			is: #(New, 'Manage...', '',
				'Preset 1', 'Preset 2', 'Preset 3', '',
				Save, 'Save As...'))

		mock.PresetsControl_listTrimmed? = true
		Assert(mock.MenuButton_Presets()
			is: #(New, 'Manage... (see all presets)', '',
				'Preset 1', 'Preset 2', 'Preset 3', '',
				Save, 'Save As...'))

		mock.PresetsControl_initial = #(Different, '', Initial)
		Assert(mock.MenuButton_Presets()
			is: #(Different, '', Initial, 'Manage... (see all presets)', '',
				'Preset 1', 'Preset 2', 'Preset 3', '',
				Save, 'Save As...'))

		mock.PresetsControl_extraMenu = #(Here, Is, Some, '', Extra)
		Assert(mock.MenuButton_Presets()
			is: #(Different, '', Initial, 'Manage... (see all presets)', '',
				'Preset 1', 'Preset 2', 'Preset 3', '',
				Save, 'Save As...', Here, Is, Some, '', Extra))
		}

	Test_valid()
		{
		fn = PresetsControl.PresetsControl_valid

		Assert(fn('Preset 1') is: '')
		Assert(fn('Preset 1\t') is: 'Name cannot contain newlines or tabs')
		Assert(fn('Preset\t1') is: 'Name cannot contain newlines or tabs')
		Assert(fn('Preset\r\n1') is: 'Name cannot contain newlines or tabs')
		Assert(fn('Preset 1\r\n') is: 'Name cannot contain newlines or tabs')
		Assert(fn('Longer than thirty, 11111111111')
			is: 'Name must be less than 30 characters long')
		Assert(fn('Equal to thirty, 1111111111111')
			is: 'Name must be less than 30 characters long')
		Assert(fn('Less than thirty, 11111111111') is: '')
		}

	Test_BeforeCopy()
		{
		report = 'Eta_Some_Report~presets~some preset'
		rec = [:report]

		_suffix = PresetsControl.PresetsControl_getSuffix()
		cl = PresetsControl
			{
			PresetsControl_getSuffix()
				{
				return _suffix
				}
			}
		fn = cl.BeforeCopy

		fn(rec)
		Assert(rec.report is: report $ _suffix)
		Assert(rec isSize: 1)

		_suffix = cl.PresetsControl_getSuffix()
		fn(rec)
		Assert(rec.report is: report $ _suffix)
		fn(rec)
		Assert(rec.report is: report $ _suffix)
		Assert(rec isSize: 1)
		}
	}