// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		cl = RecursiveDir
			{
			RecursiveDir_dir(path)
				{
				return _paths.GetDefault(path, Object()).Map({ [name: it] })
				}
			}
		_paths = #(
			'test1/': ('test1_1', 'test1_2/', 'test1_3/'),
			'test1/test1_2/': ('test1_2_1')
			'test1/test1_3/': ('test1_3_1', 'test1_3_2')
			)
		result = Object()
		block = { result.Add(it.name) }

		cl('test', :block)
		Assert(result is: #())

		cl('test1', :block)
		Assert(result is: #('test1/test1_1',
			'test1/test1_2/test1_2_1',
			'test1/test1_3/test1_3_1', 'test1/test1_3/test1_3_2'))

		result.Delete(all:)
		skipFn? = { |cur| cur.Has?('test1_2') or cur.Has?('test1_3_2') }
		cl('test1', :block, :skipFn?)
		Assert(result is: #('test1/test1_1', 'test1/test1_3/test1_3_1'))
		}
	}