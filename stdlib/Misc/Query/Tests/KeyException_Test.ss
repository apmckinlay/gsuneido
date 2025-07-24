// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_Translate_unhandled()
		{
		Assert(KeyException.Translate("stuff", "action") is: "stuff")

		Assert({ KeyException("stuff") } throws: "unhandled")
		}

	t(e)
		{
		// want to go through TryCatch and Translate
		KeyException.TryCatch({ throw e },
			{|e2| return KeyException.Translate(e2, 'frazzle') })
		}

	Test_Translate_duplicate_key()
		{
		Assert(.t('duplicate key: date,phone')
			is: 'Duplicate value in field Date+Phone #')

		Assert(.t('duplicate key: keyfield = [key1,key2]')
			is: 'Duplicate value in field keyfield: key1')

		Assert(.t('duplicate key: keyfield = ["A Sample Company",330689] in biz_partners')
			is: 'Duplicate value in field keyfield: "A Sample Company"')

		Assert(
			.t('duplicate key: keyfield = ["A Sample Company, 2",330689] in biz_partners')
			is: 'Duplicate value in field keyfield: "A Sample Company, 2"')

		Assert(.t('duplicate key: sulog_timestamp = [#20140311.1551,33064222] in log')
			is: 'Duplicate value in field Timestamp')

		Assert(.t('duplicate key: keyfield,keyfield2 = [key1,key2,12121212]')
			is: 'Duplicate value in field keyfield+keyfield2')

		Assert(.t('duplicate key: keyfield = ["A Samply Company, 2"]')
			is: 'Duplicate value in field keyfield: "A Samply Company, 2"')

		Assert(.t('duplicate key: keyfield,keyfield2 = ["Test, value",key2,12121212]')
			is: 'Duplicate value in field keyfield+keyfield2')

		// duplicate on empty key table
		Assert(.t('duplicate key:  = [] (from server)')
			is: 'Duplicate Entry')
		}

	Test_Translate_foreign_key_exceptions()
		{
		spy = .SpyOn(GetTableName).Return('Key Exception Test Table')
		Assert(.t('query: delete record from biz_partners blocked by foreign key ' $
			'in keyexception_test_table')
			is: 'Record cannot be updated or deleted because it is used ' $
				'(Key Exception Test Table)')

		Assert(.t('query: delete record from biz_partners blocked by foreign key ' $
			'in keyexception_test_table \r\n(from server)')
			is: 'Record cannot be updated or deleted because it is used ' $
				'(Key Exception Test Table)')

		// messages containing actual foreign key values
		Assert(.t('add record blocked by foreign key to ' $
			'keyexception_test_table [#20140128.101237271] \r\n(from server)')
			is: 'Record cannot be updated or deleted because it is used ' $
				'(Key Exception Test Table)')

		Assert(.t('add record blocked by foreign key to ' $
			'keyexception_test_table ["County of Lamb"] \r\n(from server)')
			is: 'Record cannot be updated or deleted because it is used ' $
				'(Key Exception Test Table)')

		Assert(.t('add record blocked by foreign key to keyexception_test_table ' $
			'["_paid"] (Trigger_test_table) ' $
			'(add record blocked by foreign key to keyexception_test_table ' $
			'["_paid"]) \r\n(from server)')
			is: 'Record cannot be updated or deleted because it is used ' $
				'(Key Exception Test Table)')

		// message containing "cascade" key word
		Assert(.t('output blocked by foreign key: ' $
			'keyexception_test_table index(test_field) ' $
			'in keyexception_test_table cascade')
			is: 'Record cannot be updated or deleted because it is used ' $
				'(Key Exception Test Table)')
		spy.Close()

		// table name cannot be translated
		msg = "ERROR: unable to find name for table: keyexception_test_table_xyz"
		ServerSuneido.Set('TestRunningExpectedErrors', Object(msg, msg, msg))
		Assert(.t('output blocked by foreign key: ' $
			'keyexception_test_table_xyz index(test_field) ' $
			'in keyexception_test_table_xyz cascade')
			is: 'Record cannot be updated or deleted because it is used')

		msg = "ERROR: unable to find name for table: lkassf76576$%&^"
		ServerSuneido.Set('TestRunningExpectedErrors', Object(msg, msg, msg))
		Assert(.t('blocked by foreign key unexpected text lkas(jdlkd)sf76576[*&^&^]$%&^')
			is: 'Record cannot be updated or deleted because it is used')

		Assert(.t('blocked by foreign key')
			is: 'Record cannot be updated or deleted because it is used')
		}

	Test_Translate_transaction_exceptions()
		{
		Assert(.t('transaction.Complete failed')
			is: 'Unable to frazzle\n\nAnother user has made changes')

		Assert(.t('transaction conflict write conflict with tester@127.0.0.1')
			is: 'Unable to frazzle\n\nAnother user (tester) has made changes')

		Assert(.t('record had been modified write conflict with tester@127.0.0.1')
			is: 'Unable to frazzle\n\nAnother user (tester) has made changes')

		Assert(.t('block commit failed write conflict with tester@127.0.0.1')
			is: 'Unable to frazzle\n\nAnother user (tester) has made changes')

		Assert(.t('transaction exceeded max age')
			is: 'Unable to frazzle\n\nAction timed out')
		}

	Test_Translate_timeout()
		{
		for e in #(
			"block commit failed: update transaction longer than 20 seconds",
			"can't output ended transaction (update transaction longer than 10 seconds)"
			"can't delete ended transaction (update transaction longer than 10 seconds)"
			"can't update ended transaction (update transaction longer than 10 seconds)"
			"can't query ended Transaction (update transaction longer than 10 seconds)")
			Assert(.t(e) is: 'Unable to frazzle\n\nAction timed out')
		}

	Test_Translate_too_many()
		{
		Assert(.t('aborted ut3129 - too many writes (...) in one transaction')
			is: 'Unable to frazzle\n\nToo many lines to process')
		}

	Test_Translate_server_slow()
		{
		err = `can't use ended transaction (too many overlapping update transactions) ` $
			`(from server)`
		Assert(.t(err) is: 'Unable to frazzle\n\nServer too slow responding')
		}

	Test_key_too_large()
		{
		err = `key too large, size 2513 limit 1024 (from server)`
		Assert(.t(err) is: 'Indexed field contains too much data')
		}
	}
