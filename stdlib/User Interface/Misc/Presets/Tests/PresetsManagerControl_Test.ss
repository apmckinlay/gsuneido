// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_compatiblePresets?()
		{
		ob = PresetsManagerControl
			{
			PresetsManagerControl_baseName: false
			PresetsManagerControl_reportName: "test presets"
			}
		Assert(ob.PresetsManagerControl_compatiblePresets?('test') is: false)
		Assert(ob.PresetsManagerControl_compatiblePresets?('test presets'))

		ob = PresetsManagerControl
			{
			PresetsManagerControl_baseName: "test"
			PresetsManagerControl_reportName: "test presets"
			}
		Assert(ob.PresetsManagerControl_compatiblePresets?('test'))
		Assert(ob.PresetsManagerControl_compatiblePresets?('test presets'))
		Assert(ob.PresetsManagerControl_compatiblePresets?('text') is: false)
		}

	Test_changeAccessForPreset()
		{
		fn = PresetsManagerControl.PresetsManagerControl_changeAccessForPreset

		user = ''
		admin? = false
		private? = true
		Assert(fn(rec = [], user, admin?, private?))
		Assert(rec.report_options.private?)

		private? = false
		Assert(fn(rec = [], user, admin?, private?))
		Assert(rec.report_options.private? is: private?)

		Assert(fn(rec = [user: 'admin'], user, admin?, private?))
		Assert(rec hasntMember: 'report_options')

		user = 'admin'
		Assert(fn(rec = [user: 'admin'], user, admin?, private?))
		Assert(rec.report_options.private? is: private?)

		user = 'clerk'
		Assert(fn(rec = [user: 'admin'], user, admin?, private?))
		Assert(rec hasntMember: 'report_options')

		private? = true
		Assert(fn(rec = [user: 'admin'], user, admin?, private?) is: false)
		Assert(rec hasntMember: 'report_options')

		admin? = true
		Assert(fn(rec = [user: 'admin'], user, admin?, private?))
		Assert(rec.report_options.private? is: private?)
		}
	}