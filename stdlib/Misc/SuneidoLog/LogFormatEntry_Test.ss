// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_formatObjectMaxSize()
		{
		outerObMaxSize = LogFormatEntry.LogFormatEntry_outerObjectSizeLimit
		nestedObMaxSize = LogFormatEntry.LogFormatEntry_nestedObjectSizeLimit
		ob = Object()
		for(i = 0; i < outerObMaxSize + 10; i++)
			ob[i] = Display(i) $ "'s position"

		newOb = LogFormatEntry(ob)
		Assert(newOb isSize: outerObMaxSize + 1) // + 1 for ellipsis
		Assert(newOb['...'] is: '...')

		ob = ob[.. 20]
		newOb = LogFormatEntry(ob)
		Assert(newOb isSize: 20)
		Assert(newOb.Last() is: "19's position")

		ob2 = Object()
		for(i = 0; i < nestedObMaxSize + 10; i++)
			ob2[i] = Display(i) $ "'s position"
		ob.innerOb = ob2

		newOb = LogFormatEntry(ob)
		Assert(newOb.innerOb isSize: nestedObMaxSize + 1) // + 1 for ellipsis
		Assert(newOb.innerOb['...'] is: '...')
		}

	Test_formatString()
		{
		func = LogFormatEntry.LogFormatEntry_formatString
		maxStrSize = 100
		str = ''
		Assert(func(str, maxStrSize) is: '')

		str = 'this should not do anything'
		Assert(func(str, maxStrSize) is: 'this should not do anything')

		str = 'this will get ellipsed ' $ 'a'.Repeat(100)
		Assert(func(str, maxStrSize) has: '...')
		Assert(func(str, maxStrSize) isSize: 103)

		str = '<html><stuff><body>'.Repeat(10) $
			'some other part of the message which isnt html so there is no closing tag'
		Assert(func(str, maxStrSize)
			is: '<html><stuff><body><html><stuff><body><html>...' $
				'message which isnt html so there is no closing tag')
		str = '<html><stuff><body><!- here is a really long commment that should not' $
			'get included in the expected message because the closing tag will ' $
			'be ellipsis>'$
			'we only parse><the first portion as it might cut off a closing tag>'
		Assert(func(str, maxStrSize) is: '<html><stuff><body>...' $
				'e first portion as it might cut off a closing tag>')

		str = 'this is a generic message that has a random set of braces near the end
			but this isnt html so we should not parse anything >>><<<<<<><< okay?'
		Assert(func(str, maxStrSize)
			is: 'this is a generic message that has a random set of...' $
				'so we should not parse anything >>><<<<<<><< okay?')
		}

	Test_format_object()
		{
		Assert(LogFormatEntry(#()) is: #())
		Assert(LogFormatEntry(#(123, a: 'hello')) is: #(123, a: 'hello'))
		// test object inside object but not empty
		Assert(LogFormatEntry(#(ob: (a: 1, ob2: (b: 2))))
			is: #(ob: (a: 1, ob2: '<object>')))
		// test object inside object but empty
		Assert(LogFormatEntry(#(ob: (a: 1, ob2: ())))
			is: #(ob: (a: 1, ob2: '<emptyobject>')))
		Transaction(read:)
			{ |t|
			Assert(LogFormatEntry(Object(:t)) is: #(t: "<Transaction>"))
			}
		}

	Test_remove_private_data()
		{
		ob = #(one: 1, pass: 'abc123', two: false, three: 'text')
		Assert(LogFormatEntry(ob) is: #(one: 1, pass: '***', two: false, three: 'text'))

		ob = #(one: 1, password: 'abc123', two: false, three: 'aa')
		Assert(LogFormatEntry(ob) is: #(one: 1, password: '***', two: false, three: 'aa'))

		ob = #(one: 1, passWord: 'abc123', two: false, three: 'aa')
		Assert(LogFormatEntry(ob) is: #(one: 1, passWord: '***', two: false, three: 'aa'))

		ob = #(one: 1, passWD: 'abc123', two: false, three: 'aa')
		Assert(LogFormatEntry(ob) is: #(one: 1, passWD: '***', two: false, three: 'aa'))

		ob = #(one: 1, pw: 'abc123', two: false, three: 'aa')
		Assert(LogFormatEntry(ob) is: #(one: 1, pw: '***', two: false, three: 'aa'))

		ob = #(one: 1, passphrase: 'abc123', two: false, three: 'aa')
		Assert(LogFormatEntry(ob)
			is: #(one: 1, passphrase: '***', two: false, three: 'aa'))

		ob = #(one: 1, opw: 'abc123', two: false, three: 'aa')
		Assert(LogFormatEntry(ob) is: #(one: 1, opw: 'abc123', two: false, three: 'aa'))

		ob = #(one: 1, sin_ssn: 'abc123', two: false, three: 'aa')
		Assert(LogFormatEntry(ob) is: #(one: 1, sin_ssn: '***', two: false, three: 'aa'))

		ob = #(one: 1, employee_sin_ssn: 'abc123', two: false, three: 'aa')
		Assert(LogFormatEntry(ob)
			is: #(one: 1, employee_sin_ssn: '***', two: false, three: 'aa'))

		ob = #(one: 1, employee_sin_ssn_display: 'abc123', two: false, three: 'aa')
		Assert(LogFormatEntry(ob)
			is: #(one: 1, employee_sin_ssn_display: '***', two: false, three: 'aa'))

		ob = #(one: 1, sin_ssn: 'abc123', two: false, three: 'aa')
		Assert(LogFormatEntry(ob) is: #(one: 1, sin_ssn: '***', two: false, three: 'aa'))

		ob = #(one: 1, ssn: 'abc123', two: false, three: 'aa')
		Assert(LogFormatEntry(ob) is: #(one: 1, ssn: '***', two: false, three: 'aa'))

		ob = #(one: 1, sin: 'abc123', two: false, three: 'aa')
		Assert(LogFormatEntry(ob) is: #(one: 1, sin: '***', two: false, three: 'aa'))

		ob = #(one: 1, single: 'abc123', two: false, three: 'aa')
		Assert(LogFormatEntry(ob)
			is: #(one: 1, single: 'abc123', two: false, three: 'aa'))

		ob = #(one: 1, credit: 'abc123', two: false, three: 'aa')
		Assert(LogFormatEntry(ob)
			is: #(one: 1, credit: 'abc123', two: false, three: 'aa'))

		ob = #(one: 1, card: 'abc123', two: false, three: 'aa')
		Assert(LogFormatEntry(ob) is: #(one: 1, card: '***', two: false, three: 'aa'))

		ob = #(one: 1, creditcard: 'abc123', two: false, three: 'aa')
		Assert(LogFormatEntry(ob)
			is: #(one: 1, creditcard: '***', two: false, three: 'aa'))

		ob = #(one: 1, credit_card: 'abc123', two: false, three: 'aa')
		Assert(LogFormatEntry(ob)
			is: #(one: 1, credit_card: '***', two: false, three: 'aa'))

		ob = #(one: 1, creditCard: 'abc123', two: false, three: 'aa')
		Assert(LogFormatEntry(ob)
			is: #(one: 1, creditCard: '***', two: false, three: 'aa'))

		ob = #(one: 1, creditCard: 'abc123', two: #(a: 'a', b: 'b'))
		Assert(LogFormatEntry(ob)
			is: #(one: 1, creditCard: '***', two: #(a: 'a', b: 'b')))

		ob = #(one: 1, two: '2', three: #(password: 'test', b: 'b'))
		Assert(LogFormatEntry(ob)
			is: #(one: 1, two: '2', three: #(password: '***', b: 'b')))

		ob = #(one: 1, passhash: '123xyz', two: #(a: 'a', b: 'b', authorization: 'abc'))
		Assert(LogFormatEntry(ob)
			is: #(one: 1, passhash: '***', two: #(a: 'a', b: 'b', authorization: '***')))

		ob = #(one: 1, test_passhash: 'test123', two: #(a: 'a', b: 'b'))
		Assert(LogFormatEntry(ob)
			is: #(one: 1, test_passhash: '***', two: #(a: 'a', b: 'b')))
		}

	Test_privateData?()
		{
		m = LogFormatEntry.LogFormatEntry_privateData?
		shouldBePrivate = #(newPassword test_account, accountTst, sinssn, bank_aa_account,
			employee_sin_ssn, sample_fuelcard_num6, test_t4a_rcpnt_sin,
			test_t4_proprietor_sin, fuel_card, test_ftp_pass, ftp_pass, card, cardnum,
			sample_fuel_card_locate, customer_password, co_fuelcard, towing_card_number,
			emp_cardnum, bank_account_decrypt, sample_account_decrypt, account,
			bank_account, cus_cardnum_display, token, extraToken, currentPassword,
			creditCard, credit_card, cardNum, cvv, cardCVV, card_cvv, pwd, test_pwd
			passhash, someprefix_passhash, authorization)
		for member in shouldBePrivate
			Assert(m(member))

		shouldntBePrivate = #(x, y, w, h, sample, rec, data, coord, record, num, test_num)
		for member in shouldntBePrivate
			Assert(m(member) is: false)
		}

	Test_nonObjectArgument()
		{
		Assert(LogFormatEntry('test') is: 'test')
		}

	Test_bufferLimit()
		{
		cl = LogFormatEntry
			{
			LogFormatEntry_maxSize: 10
			}

		result = cl(#())
		Assert(result is: #())

		result = cl("")
		Assert(result is: "")

		// .maxSize does not affect a single string, they are restricted by the
		// maxStrSize which is passed as an argument
		result = cl("string is too long")
		Assert(result is: "string is too long")

		result = cl("string is too long", 6)
		Assert(result is: "str...ong")

		result = cl(#(1,2,3,4,5))
		Assert(result is: #(1,2,3,4,5))

		result = cl(#(1,2,3,4,5,6,7,8,9,10))
		Assert(result
			is: #(1,2,3,4,5,"Stopped logging to prevent buffer overflow"))

		result = cl(#('longer string'))
		Assert(result is: #("Stopped logging to prevent buffer overflow"))

		result = cl(#(#20210101))
		Assert(result is: #(#20210101))

		result = cl(#(20210101.123401001))
		Assert(result is: #("Stopped logging to prevent buffer overflow"))

		result = cl(#(1, (1,2,3)))
		Assert(result is: #(1, (1,2,3)))

		result = cl(#(1, (1,2,3,4,5,6,7,8,9,10)))
		Assert(result is: #(1, (1,2,3, "Stopped logging to prevent buffer overflow")))

		result = cl(#((1), (2), (3)))
		Assert(result is: #((1), (2), (3)))

		result = cl(#((1), (2), (3), (4)))
		Assert(result
			is: #((1), (2), (3), ("Stopped logging to prevent buffer overflow")))

		result = cl(#((("nested", "more", "stuff"))))
		Assert(result is: #(("<object>")))

		result = cl(#(mem: "val", ab: "a"))
		Assert(result equalsSet: #(mem: "val", ab: "a"))

		result = cl(#(mem: "val", ab: "abcde"))
		Assert(result.Has?("Stopped logging to prevent buffer overflow"),
			msg: 'missing stopped logging from long values')

		result = cl(#(mem: "val", abcde: "ab"))
		Assert(result.Has?("Stopped logging to prevent buffer overflow"),
			msg: 'missing stopped logging from long members')
		}
	}
