// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_BuildParam()
		{
		cl = Reporter_datadict
			{
			Reporter_datadict_getConfigRec(name)
				{
				return _configlib.GetDefault(name, false)
				}
			Reporter_datadict_outputRecord(rec, delete? = false)
				{
				_output.Add([:rec, :delete?])
				}
			}

		fn = cl.BuildParam
		_configlib = Object()
		_output = Object()

		field = .MakeDatadict()

		// field has been defined, don't overwrite it in configlib
		fn(field, '', [baseField: 'testField'])
		Assert(_output isSize: 0)

		fn(field, '_param', [baseField: 'testField'])
		Assert(_output isSize: 1)
		Assert(_output[0].rec is: [name: 'Field_' $ field $ '_param',
			text: 'Field_' $ field $ '\n\t{\n\tPromptInfo: ' $
				Display([baseField: 'testField']) $ '\n\t}'])
		Assert(not _output[0].delete?)

		// field $ '_param' record has already been created
		_configlib['Field_' $ field $ '_param'] = _output[0].rec
		fn(field, '_param', [baseField: 'testField'])
		Assert(_output isSize: 1)

		// field $ '_param' record needs to be updated
		fn(field, '_param', [baseField: 'testField', suffix: 'test'])
		Assert(_output isSize: 2)
		Assert(_output[1].rec is: [name: 'Field_' $ field $ '_param',
			text: 'Field_' $ field $ '\n\t{\n\tPromptInfo: ' $
				Display([baseField: 'testField', suffix: 'test']) $ '\n\t}'])
		Assert(_output[1].delete?)
		_configlib['Field_' $ field $ '_param'] = _output[1].rec

		for type in #('total', 'min', 'max', 'average')
			{
			_output.Delete(all:)
			fn(type $ '_' $ field, '_param', [baseField: 'testField', suffix: 'test'])
			Assert(_output isSize: 1)
			Assert(_output[0].rec is: [name: 'Field_' $ type $ '_' $ field $ '_param',
				text: 'Field_' $ field $ '\n\t{\n\tPromptInfo: ' $
					Display([baseField: 'testField', suffix: 'test']) $ '\n\t}'])
			Assert(not _output[0].delete?)
			}

		field = .TempName() $ '?'
		_output.Delete(all:)
		fn(field, '_param', [baseField: 'testField'])
		Assert(_output isSize: 1)
		Assert(_output[0].rec is: [name: 'Field_' $ field.RemoveSuffix('?') $ '_param',
			text: 'Field_string\n\t{\n\tPromptInfo: ' $
				Display([baseField: 'testField']) $ '\n\t}'])
		Assert(not _output[0].delete?)
		}
	}
