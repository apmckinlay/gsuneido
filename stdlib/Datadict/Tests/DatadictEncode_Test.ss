// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	//In general, the base Encodes return the value IF it is not coded to handle it
	Test_wrongValueTypeForField()
		{
		date = Date()
		Assert(DatadictEncode('boolean', date) is: date)
		Assert(DatadictEncode('boolean', 1) is: 1)
		Assert(DatadictEncode('boolean', 1.1) is: 1.1)
		Assert(DatadictEncode('boolean', 'Test') is: 'Test')
		Assert(DatadictEncode('boolean', Object()) is: #())
		Assert(DatadictEncode('boolean', Record()) is: [])
		}

	// Invalid / non-existent fields get encoded to strings
	Test_notAValidField()
		{
		date = Date()
		Assert(DatadictEncode('uhhh...', date) is: Display(date))
		Assert(DatadictEncode("I'm", 1) is: '1')
		Assert(DatadictEncode('not', 1.1) is: '1.1')
		Assert(DatadictEncode('a', 'Test') is: 'Test')
		Assert(DatadictEncode('real', true) is: 'true')
		Assert(DatadictEncode('thing', Object()) is: '#()')
		Assert(DatadictEncode('sadly', Record()) is: '[]')
		}

	// So far (22-12-2017) only date takes an extra param for its encode.
	// This test ensures that future changes do not lead to unexpected results.
	//		Summary: the extra params shouldn't break anything, they should be ignored
	Test_extraParamsIgnored()
		{
		Assert(DatadictEncode('string', 'test', len: 1) is: 'test')
		Assert(DatadictEncode('boolean', 'yes', inverse:))
		Assert(DatadictEncode('number', 5.0055, roundTo: 3) is: 5.0055)
		}

	Test_string()
		{
		date = Date()
		Assert(DatadictEncode('string', date) is: Display(date))
		Assert(DatadictEncode('string', 1) is: '1')
		Assert(DatadictEncode('string', 1.1) is: '1.1')
		Assert(DatadictEncode('string', 'Test') is: 'Test')
		Assert(DatadictEncode('string', true) is: 'true')
		Assert(DatadictEncode('string', Object()) is: '#()')
		Assert(DatadictEncode('string', Record()) is: '[]')
		}

	Test_boolean()
		{
		Assert(DatadictEncode('boolean', true))
		Assert(DatadictEncode('boolean', false) is: false)
		Assert(DatadictEncode('boolean', 'y'))
		Assert(DatadictEncode('boolean', 'Y'))
		Assert(DatadictEncode('boolean', 'yes'))
		Assert(DatadictEncode('boolean', 'YES'))
		Assert(DatadictEncode('boolean', 'n') is: false)
		Assert(DatadictEncode('boolean', 'N') is: false)
		Assert(DatadictEncode('boolean', 'no') is: false)
		Assert(DatadictEncode('boolean', 'NO') is: false)
		}

	Test_number()
		{
		Assert(DatadictEncode('number', 1) is: 1)
		Assert(DatadictEncode('number', 1.1) is: 1.1)
		Assert(DatadictEncode('number', '1,100') is: 1100)
		Assert(DatadictEncode('number', '1,100.01') is: 1100.01)
		Assert(DatadictEncode('number', ' 1,100.01a') is: ' 1,100.01a')
		Assert(DatadictEncode('number', '$1,100.01') is: 1100.01)
		Assert(DatadictEncode('number', ' $ 1,100.01') is: 1100.01)
		Assert(DatadictEncode('number', ' $$ 1,100.01') is: ' $$ 1,100.01')
		Assert(DatadictEncode('number', '1,100.01$') is: '1,100.01$')
		}

	Test_date()
		{
		date = Date()
		Assert(DatadictEncode('date', date) is: date)
		// Technically does nothing different as it doesn't need the fmt this time
		Assert(DatadictEncode('date', date, fmt: 'yMd') is: date)
		Assert(DatadictEncode('date', 'May 12, 2017', fmt: 'Mdy') is: #20170512)

		// Uses the fmt to decipher the string
		Assert(DatadictEncode('date', '05122017', fmt: 'Mdy') is: #20170512)
		Assert(DatadictEncode('date', 'May 12, 2017') is: #20170512)
		}
	}
