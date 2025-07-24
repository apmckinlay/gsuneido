// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_query()
		{
		mock = Mock(CustomTableDataSources)
		mock.When.query([anyArgs:]).CallThrough()
		mock.When.queryKeys('master_table1').Return(['a_key'])
		mock.When.queryKeys('master_table2').Return(['b_key', 'c_key'])
		mock.When.queryKeys('master_table3').Return(['d_key', '(ef_key, fe_key)'])

		mock.When.foreignKey('cust_table1').Return('a_key')
		mock.When.foreignKey('cust_table2').Return('b_key')
		mock.When.foreignKey('cust_table3').Return('d_key')

		expectedQuery = 'cust_table1 rename custtable_FK to a_key, ' $
			'custtable_num to custtable_num_new' $
			' join by(a_key) (master_table1 project a_key)'
		expectedKeys = #('a_key')
		.testTables(mock, 'cust_table1', 'master_table1', expectedQuery, expectedKeys)

		expectedQuery = 'cust_table2 rename custtable_FK to b_key, ' $
			'custtable_num to custtable_num_new' $
			' join by(b_key) (master_table2 project b_key, c_key)'
		expectedKeys = #('b_key', 'c_key')
		.testTables(mock, 'cust_table2', 'master_table2', expectedQuery, expectedKeys)

		expectedQuery = 'cust_table3 rename custtable_FK to d_key, ' $
			'custtable_num to custtable_num_new' $
			' join by(d_key) (master_table3 project d_key, ef_key, fe_key)'
		expectedKeys = #('d_key', 'ef_key', 'fe_key')
		.testTables(mock, 'cust_table3', 'master_table3', expectedQuery, expectedKeys)
		}

	testTables(mock, custTable, masterTable, expectedQuery, expectedKeys)
		{
		result = mock.query(custTable, masterTable, masterKeys = [])
		mock.Verify.queryKeys(masterTable)
		Assert(masterKeys isSize: 1)
		Assert(masterKeys[masterTable] is: expectedKeys)
		Assert(result is: expectedQuery)

		// Ensuring the masterKeys object is used rather than looking up the keys again.
		// Ensuring the result is the same when using the keys saved in masterKeys.
		result = mock.query(custTable, masterTable, masterKeys)
		mock.Verify.queryKeys(masterTable)
		Assert(masterKeys isSize: 1)
		Assert(masterKeys[masterTable] is: expectedKeys)
		Assert(result is: expectedQuery)
		}
	}