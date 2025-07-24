// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_makeTitle()
		{
		mock = Mock(SvcDisplayControl)
		mock.When.makeTitle([anyArgs:]).CallThrough()
		mock.When.buildPath([anyArgs:]).Return('path')

		oldRec = [lib_committed: #20230801.1100, lib_modified: #20230802.1200]
		newRec = [lib_committed: #20230803.1300, lib_modified: #20230804.1400]
		base = false
		Assert(mock.makeTitle('', oldRec, newRec, base)	is: '')

		titleTemplate = 'LOCAL <Npath>'
		Assert(mock.makeTitle(titleTemplate, oldRec, newRec, base)
			is: 'LOCAL path')

		titleTemplate = 'MASTER <Npath>'
		Assert(mock.makeTitle(titleTemplate, oldRec, newRec, base)
			is: 'MASTER path')

		titleTemplate = 'AS OF <Ndate> <Npath>'
		Assert(mock.makeTitle(titleTemplate, oldRec, newRec, base)
			is: 'AS OF Com: ' $ newRec.lib_committed.ShortDateTime() $ ' path')

		titleTemplate = 'LOCAL tableName DEFINITION <Npath>'
		Assert(mock.makeTitle(titleTemplate, oldRec, newRec, base)
			is: 'LOCAL tableName DEFINITION path')

		base = []
		titleTemplate = '<Odate>'
		Assert(mock.makeTitle(titleTemplate, oldRec, newRec, base)
			is: 'Com: ' $ oldRec.lib_committed.ShortDateTime())

		titleTemplate = 'LOCAL <Opath>'
		Assert(mock.makeTitle(titleTemplate, oldRec, newRec, base)
			is: 'LOCAL path')

		titleTemplate = 'MASTER <Opath>'
		Assert(mock.makeTitle(titleTemplate, oldRec, newRec, base)
			is: 'MASTER path')

		titleTemplate = 'TABLENAME <Opath>'
		Assert(mock.makeTitle(titleTemplate, oldRec, newRec, base)
			is: 'TABLENAME path')

		titleTemplate = 'PREVIOUS <Odate> <Opath>'
		Assert(mock.makeTitle(titleTemplate, oldRec, newRec, base)
			is: 'PREVIOUS Com: ' $ oldRec.lib_committed.ShortDateTime() $ ' path')
		}

	Test_makeRightTitle()
		{
		fn = SvcDisplayControl.SvcDisplayControl_makeRightTitle
		oldRec = [lib_committed: #20230801.1100, lib_modified: #20230802.1200]
		newRec = [lib_committed: #20230803.1300, lib_modified: #20230804.1400]
		Assert(fn('', oldRec, newRec) is: '')

		titleRightTemplate = '<Ndate>'
		Assert(fn(titleRightTemplate, oldRec, newRec) is:
			'Com: ' $ newRec.lib_committed.ShortDateTime())

		titleRightTemplate = '<Nmdate>  <Ndate>'
		Assert(fn(titleRightTemplate, oldRec, newRec) is:
			'Mod: ' $ newRec.lib_modified.ShortDateTime() $
				'  Com: ' $ newRec.lib_committed.ShortDateTime())

		titleRightTemplate = '<Odate>'
		Assert(fn(titleRightTemplate, oldRec, newRec) is:
			'Com: ' $ oldRec.lib_committed.ShortDateTime())

		titleRightTemplate = '<Omdate>  <Odate>'
		Assert(fn(titleRightTemplate, oldRec, newRec) is:
			'Mod: ' $ oldRec.lib_modified.ShortDateTime() $
				'  Com: ' $ oldRec.lib_committed.ShortDateTime())
		}
	}