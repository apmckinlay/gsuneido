// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		Assert(AllowInsertRecord?([], false))
		Assert(AllowInsertRecord?([], .TempName()))
		Assert(AllowInsertRecord?([protect: true], 'protect'))
		Assert(AllowInsertRecord?([protect: 'protected'], 'protect'))
		Assert(AllowInsertRecord?([protect: #()], 'protect'))
		Assert(AllowInsertRecord?([protect: #(noInsert: false)], 'protect'))
		Assert(AllowInsertRecord?([protect: #(noInsert:)], 'protect') is: false)
		}
	}