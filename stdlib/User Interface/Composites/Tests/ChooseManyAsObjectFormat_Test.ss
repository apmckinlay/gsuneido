// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_BuildDesc()
		{
		Assert(ChooseManyAsObjectFormat.BuildDesc(false) is: '')
		Assert(ChooseManyAsObjectFormat.BuildDesc('') is: '')

		mock = Mock()
		mock.ChooseManyAsObjectFormat_idField =
			mock.ChooseManyAsObjectFormat_displayField = 'test_field'
		mock.ChooseManyAsObjectFormat_delimiter = '|'
		desc = mock.Eval(ChooseManyAsObjectFormat.BuildDesc, #(hello, world))
		Assert(desc is: 'hello|world')

		table = .MakeTable('(test_field, test_desc) key(test_field)')
		QueryOutput(table, [test_field: 'hello', test_desc: 'HELLO'])
		QueryOutput(table, [test_field: 'world', test_desc: 'WORLD'])

		mock.ChooseManyAsObjectFormat_displayField = 'test_desc'
		mock.ChooseManyAsObjectFormat_query = table
		desc = mock.Eval(ChooseManyAsObjectFormat.BuildDesc, #(hello, world))
		Assert(desc is: 'HELLO|WORLD')
		}
	}