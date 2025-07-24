// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_callclass()
		{
		// invalid list
		args = Object(listField: 'zzz_test_listfield'
			record: Record(zzz_test_listfield: false))
		Assert(ChooseList_ValidDataListFromRules(args) is: false)

		// rule provides object
		args = Object(listField: 'zzz_test_listfield'
			record: Record(zzz_test_listfield: Object('one', 'two', 'three')))
		Assert(ChooseList_ValidDataListFromRules(args) is: #(one two three))

		// rule provides string
		args = Object(listField: 'zzz_test_listfield'
			record: Record(zzz_test_listfield: 'one, two,    three'))
		Assert(ChooseList_ValidDataListFromRules(args) is: #(one two three))

		// also has allowOtherField
		args = Object(listField: 'zzz_test_listfield',
			allowOtherField: 'zzz_test_allowOtherField'
			record: Record(
				zzz_test_listfield: Object('one', 'two', 'three'),
				zzz_test_allowOtherField: Object('abc', 'def', 'ghi')))
		Assert(ChooseList_ValidDataListFromRules(args)
			is: #(one two three abc def ghi))

		// different split value
		args = Object(listField: 'zzz_test_listfield'
			allowOtherField: 'zzz_test_allowOtherField'
			splitValue: '*'
			record: Record(
				zzz_test_listfield: 'one* two*    three',
				zzz_test_allowOtherField: 'abc*    def    * 	ghi'))
		Assert(ChooseList_ValidDataListFromRules(args) is: #(one two three abc def ghi))
		}
	}