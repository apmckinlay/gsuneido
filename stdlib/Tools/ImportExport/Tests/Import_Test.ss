// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup()
		{
		.table = .MakeTable('(test_num, test_abbrev, test_name, test_name_lower!)
			key(test_num) key(test_name) key(test_name_lower!)')
		.tableNoLower = .MakeTable('(test2_num, test2_abbrev, test2_name)
			key(test2_num) key(test2_name)')
		}

	Test_Make_unique()
		{
		testMethod = Import.Make_unique

		rec = Record(test_name: 'abc')
		testMethod(rec, .table, "duplicate key: test_name")
		Assert(rec.test_name is: 'abc*1')

		rec = Record(test_name: 1)
		testMethod(rec, .table, "duplicate key: test_name")
		Assert(rec.test_name is: '1*1')

		rec = Record(test_name: 'abc*2')
		testMethod(rec, .table, "duplicate key: test_name")
		Assert(rec.test_name is: 'abc*3')

		rec = Record(test_abbrev: 'test', test_name: '')
		testMethod(rec, .table, "duplicate key: test_name")
		Assert(rec.test_name is: 'test')

		rec = Record(test_name_lower!: 'abc', test_name: "ABC")
		testMethod(rec, .table, "duplicate key: test_name_lower!")
		Assert(rec.test_name is: 'ABC*1')

		rec = Record(test_name_lower!: 'abc', test_name_renamed: "ABC")
		testMethod(rec, .table $ ' rename test_name to test_name_renamed',
			"duplicate key: test_name_lower!")
		Assert(rec.test_name_renamed is: 'ABC*1')
		}

	Test_EncodeConversion()
		{
		fn = Import.EncodeConversion

		field = 'number_no_decimals_custom'
		dd = Datadict(field)
		rec = [number_no_decimals_custom: 123.456]
		fn(dd, rec, field)
		Assert(rec[field] is: 123)

		field = 'number_two_decimals_custom'
		dd = Datadict(field)
		rec = [number_two_decimals_custom: 123.456]
		fn(dd, rec, field)
		Assert(rec[field] is: 123.46)

		field = 'date'
		dd = Datadict(field)
		rec = [date: date = Date()]
		fn(dd, rec, field)
		Assert(rec[field] is: date.NoTime())

		field = '_abbrev'
		dd = Datadict(field)
		rec = [_abbrev: 'UpperCase']
		fn(dd, rec, field)
		Assert(rec[field] is: 'uppercase')

		field = '_NOTabbrev'
		dd = Datadict(field)
		rec = [_NOTabbrev: 'UpperCase']
		fn(dd,  [_NOTabbrev: 'UpperCase'], field)
		Assert(rec[field] is: 'UpperCase')
		}

	cl: Import
			{
			Import_readLineLimit: 41
			New(@args)
				{
				super(false, "", false, false)
				.Tf = args.from_file
				}
			Import1(line)
				{
				return [string: line]
				}
			}
	Test_GetLine()
		{
		l1 = 'This line is under our readable size'
		l2 = 'This line is exactly readable size (39)'
		l3 = 'This line is just a bit over our readable size'
		file = FakeFile(l1 $ '\n' $ l2 $ '\r\n' $ l3)
		inst = new .cl('', false, false, false from_file: file)

		// < readLineLimit
		Assert(inst.Getline() is: l1)
		// == readLineLimit
		Assert(inst.Getline() is: l2)
		// > readLineLimit
		Assert({ inst.Getline() } throws: 'File: Readline: line too long')
		}

	Test_Import()
		{
		txt = 'passes\npasses\r\nThis line is just a bit over our readable size'
		file = FakeFile(txt)
		inst = new .cl('', false, false, false from_file: file)
		mock = .mock(inst)
		mock.Eval(inst.Import)
		mock.Verify.Times(2).Output([anyArgs:])
		mock.Verify.Import_readlineFailure('File: Readline: line too long')

		file = FakeFile(txt)
		inst = new .cl('', false, false, false from_file: file)
		mock = .mock(inst)
		mock.Verify.Never().Output([anyArgs:])
		mock.When.Import1([anyArgs:]).Throw('Unexpected Error')
		Assert({ mock.Eval(inst.Import) } throws: 'Unexpected Error')
		}

	mock(inst)
		{
		mock = Mock(inst)
		mock.When.Getline().CallThrough()
		mock.When.DoImport([anyArgs:]).CallThrough()
		mock.When.Import1([anyArgs:]).CallThrough()
		mock.When.DateFmt().CallThrough()
		mock.When.ConvertRecord([anyArgs:]).CallThrough()
		mock.When.Before().Return('')
		mock.When.Header().Return('')
		mock.When.Output([anyArgs:]).Return('')
		return mock
		}

	Test_ConvertAbbrev()
		{
		Assert(Import.ConvertAbbrev('test', '') is: '')
		Assert(Import.ConvertAbbrev('test', 'VALUE') is: 'VALUE')
		Assert(Import.ConvertAbbrev('shouldconvert_abbrev', 'VALUE') is: 'value')
		.MakeLibraryRecord([name: 'Field_shouldntconvert_abbrev',
			text: `Field_text { NonLowerAbbrev?: true }`])
		Assert(Import.ConvertAbbrev('shouldntconvert_abbrev', 'VALUE') is: 'VALUE')
		Assert(Import.ConvertAbbrev('shouldntconvert_abbrev', 'value') is: 'value')
		}

	Test_LookupNumField_with_lower!()
		{
		GetNumTable.ResetCache()

		num = Timestamp()
		name = .TempName()
		abbrev = name[.. 3].Lower()

		recToOutput = [test_num: num, test_abbrev: abbrev, test_name: name]
		QueryOutput(.table, recToOutput)

		importRec = [test_name: 'not a real record']
		Import.LookupNumField(importRec, 'test_name')
		Assert(importRec.Member?('test_num') is: false)
		Assert(importRec.test_num is: '')

		importRec.test_name = name
		Import.LookupNumField(importRec, 'test_name')
		Assert(importRec.Member?('test_num'))
		Assert(importRec.test_num is: num)

		importRec.Delete('test_num')
		importRec.test_name = name.Upper()
		Import.LookupNumField(importRec, 'test_name')
		Assert(importRec.Member?('test_num'))
		Assert(importRec.test_num is: num)

		importRec.Delete('test_num')
		importRec.test_name = name.Lower()
		Import.LookupNumField(importRec, 'test_name')
		Assert(importRec.Member?('test_num'))
		Assert(importRec.test_num is: num)
		}

	Test_LookupNumField_no_lower!()
		{
		GetNumTable.ResetCache()

		num = Timestamp()
		name = .TempName()
		abbrev = name[.. 3].Lower()

		recToOutput = [test2_num: num, test2_abbrev: abbrev, test2_name: name]
		QueryOutput(.tableNoLower, recToOutput)

		importRec = [test2_name: 'not a real record']
		Import.LookupNumField(importRec, 'test2_name')
		Assert(importRec.Member?('test2_num') is: false)
		Assert(importRec.test2_num is: '')

		importRec.test2_name = name
		Import.LookupNumField(importRec, 'test2_name')
		Assert(importRec.Member?('test2_num'))
		Assert(importRec.test2_num is: num)

		importRec.Delete('test2_num')
		importRec.test2_name = name.Upper()
		Import.LookupNumField(importRec, 'test2_name')
		Assert(importRec.Member?('test2_num') is: false)
		Assert(importRec.test2_num is: '')
		}

	Test_getRecordFields()
		{
		func = Import.Import_getRecordFields
		rec = [name: 'bob', height: '6', bday: #20020101, start_num: '1']
		sortedFields = func(rec)
		Assert(sortedFields[0] is: 'start_num')

		sortedFields = func([name: 'bob', start_num_renamed: '1'])
		Assert(sortedFields[0] is: 'start_num_renamed')

		/// no num field, so should just match Members
		rec = [name: 'bob', height: '6', bday: #20020101, start_number: '1']
		fields = rec.Members()
		sortedFields = func(rec)
		Assert(sortedFields is: fields)
		}
	}