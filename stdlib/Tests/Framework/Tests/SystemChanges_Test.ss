// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_TableStates()
		{
		tables_before = SystemChanges.TableStates(false)
		Assert(tables_before hasMember: 'stdlib')
		Assert(tables_before.stdlib.nrows isnt: 0)
		Assert(tables_before.stdlib.totalsize isnt: 0)
		Assert(tables_before hasntMember: 'Test_lib')

		extraTable = .MakeTable('(col1, col2, col3, colKey) key(colKey) index(col1)',
			[[col1: 'Column 1', colKey: 'Key Value']])
		tables_after = SystemChanges.TableStates(false)
		Assert(tables_after hasMember: extraTable)
		Assert(tables_after hasntMember: 'Test_lib')

		tables_after = SystemChanges.TableStates(Object(extraTable))
		Assert(tables_after hasntMember: extraTable)
		Assert(tables_after hasntMember: 'Test_lib')
		Assert(tables_after.Members() equalsSet: tables_before.Members())
		Assert(SystemChanges.
			SystemChanges_checkTableDifferences(tables_before, tables_after) is: '')
		}

	Test_checkTableDifferences()
		{
		mock = Mock(SystemChanges)
		mock.When.checkTableDifferences([anyArgs:]).CallThrough()
		mock.When.totalsizeThreshold([anyArgs:]).Return(Object(lower: 0, upper: 0))
		tables_before = Object(
			table0: [nrows: 0,  totalsize: 0],
			table1: [nrows: 10, totalsize: 100],
			table2: [nrows: 20, totalsize: 200],
			table3: [nrows: 30, totalsize: 300],
			table4: [nrows: 40, totalsize: 400],
			table5: [nrows: 50, totalsize: 500])
		tables_after = tables_before.DeepCopy()

		// Test table size checking with no discrepancies
		Assert(mock.checkTableDifferences(tables_before, tables_after) is: '')

		// Create discrepancies between tables_before and tables_after
		tables_after.Delete('table0')
		tables_after.table7 = Object(nrows: 70, totalsize: 700)
		tables_after.table3.nrows = 35
		tables_after.table3.totalsize = 350

		// Test table size checking
		// table3: nrows change and size change detected
		result = mock.checkTableDifferences(tables_before, tables_after).Split('\r\n\t- ')
		Assert(result isSize: 4)
		Assert(result[0] is: 'Table Discrepancies:')
		Assert(result has: 'table0: deleted')
		Assert(result has: 'table3: nrows: 5, totalsize: 50')
		Assert(result has: 'table7: created')

		// table3: size change detected
		tables_after.table3.nrows = 30
		result = mock.checkTableDifferences(tables_before, tables_after).Split('\r\n\t- ')
		Assert(result isSize: 4)
		Assert(result[0] is: 'Table Discrepancies:')
		Assert(result has: 'table0: deleted')
		Assert(result has: 'table3: totalsize: 50')
		Assert(result has: 'table7: created')
		}

	Test_changedFiles()
		{
		c = SystemChanges.SystemChanges_changedFiles
		files_before = Object()
		files_after = Object()
		errorStr = c(files_before, files_after)
		Assert(errorStr like: "")

		//no change to files
		name = 'xxxETAfile.txt'
		date = Timestamp()
		size = 123
		file1 = Object(:name, :date, :size)
		files_before = Object(file1)
		files_after = Object(file1)
		errorStr = c(files_before, files_after)
		Assert(errorStr like: "")

		//changed date and size
		date2 = Timestamp()
		file1After = Object(:name, date: date2, size: 234)
		files_after = Object(file1After)
		errorStr = c(files_before, files_after)
		Assert(errorStr startsWith: "CHANGED xxxETAfile.txt: ")
		Assert(errorStr
			has: "date: before " $ Display(date) $ ", after " $ Display(date2) $ ";")
		Assert(errorStr
			has: "size: " $ "before 123, after 234;")

		//multiple files, one file changing
		file2Date = Timestamp()
		file2 = Object(name: "xxxETAfile2.txt", date: file2Date, size: 456)
		files_before = Object(file1, file2)
		file2DateAfter = Timestamp()
		file2After = Object(name: "xxxETAfile2.txt", date: file2DateAfter, size: 456)
		files_after = Object(file1, file2After)
		errorStr = c(files_before, files_after)
		Assert(errorStr
			like: "CHANGED xxxETAfile2.txt: date: before " $ Display(file2Date) $
				", after " $ Display(file2DateAfter) $ ";")

		//deleted file
		files_before = Object(file1)
		files_after = Object()
		errorStr = c(files_before, files_after)
		Assert(errorStr like: "DELETED xxxETAfile.txt")
		}
	Test_findFile()
		{
		f = SystemChanges.SystemChanges_findFile
		fileList = #((name: one, date: #20151229, size: 123),
						(name: two, date: #20151225, size: 456))
		Assert(f(fileList, 'one') is: 0)
		Assert(f(fileList, 'two') is: 1)
		Assert(f(fileList, 'abc') is: false)
		}

	Test_removeExcludedFiles()
		{
		m = SystemChanges.SystemChanges_removeExcludedFiles
		excludedFiles = #('test.txt', 'another_file', 'this_prefix.abc', 'gs')
		files = Object(#(name: 'abc.txt'), #(name: 'another_file'), #(name: 'good.exe'),
			#(name: 'test.txt'), #(name: 'this_prefix.abc_p11'),
			#(name: 'this_prefix.abc_p360'), #(name: 'gs1732057725.tmp'))
		m(files, excludedFiles)
		Assert(files is: Object(#(name: 'abc.txt'), #(name: 'good.exe')))
		}
	}
