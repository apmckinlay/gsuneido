// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		.validCases.Each({ Assert(ValidEmailAddress?(it)) })
		.invalidCases.Each({ Assert(ValidEmailAddress?(it) is: false) })
		}

	validCases:
		(
		'a@b.com',
		'a.b@c.com',
		'Andrew74+junk@b.c.dd',
		'a.b-c@d.ee',
		'x <a@b.cc>',
		'a@b.cc ',
		' a@b.cc',
		' a@b.cc	',
		"A.B'Cd@e.ff",
		"#!$%&'*+-/=?^_`{}|~@a.bb",
		"This. is (a) test<#!$%&'*+-/=?^_`{}|~@a.bb>",
		'Tester, Test. - Company LLC <tester@testing.com>',
		'example-test@dash-example.com',
		'More @ Tests <example-test@dash-example.com>',
		'test@domain.this-is-totest-limit',
		'test@domain.c123om',
		'test@domain.c1o2m3'
		)

	invalidCases:
		(
		'"a.b c"@d.e',
		'abc',
		'a@.c',
		'a@b',
		'a@b.',
		'.a@b.c',
		'a.@b.c',
		'abc@',
		'@abc',
		'A@b@c@d.e',
		'a@b@c@morelikelydomain.com',
		'x <a@b.c',
		'x a@b.c>',
		'a..b@c.d',
		'a"b(c)d,e:f;gi[j\k]l@example.com',
		'test"invalid"quotes@example.com',
		'this is"also\invalid@example.com',
		'1234567890123456789012345678901234567890123456789012345678901234' $
			'andonandon@justtesting.com',
		'tester@abcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd' $
			'efabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdef' $
			'abcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefab' $
			'cdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdef.com',
		'john..doe@testing.com',
		'john.doe@testing..com',
		// the following case is just to ensure the regex match size limit is not exceeded
		'<shiporder orderid="889923" xmlns:xsi="http://www.w3.org/2001/XMLSchema-' $
			'instance" xsi:noNamespaceSchemaLocation="shiporder.xsd"><orderperson>' $
			'John Smith</orderperson><shipto><name>Ola Nordmann</name><address>Langgt' $
			' 23</address><city>4000 Stavanger</city><country>Norway</country></shipto>' $
			'<item><title>Empire Burlesque</title><note>Special Edition</note>' $
			'<quantity>1</quantity><price>10.90</price></item><item><title>Hide ' $
			'your heart</title><quantity>1</quantity><price>9.90</price></item>' $
			'</shiporder>',
		// another case to test regex match size limit (user actually entered similar)
		"`````````````````````````````````````````````````````````````````````````````" $
			"```````````````````````````````````````asdsad`````````````````````````````" $
			"`````````````````````````````````````````````````````````````````",
		'abc is"not\valid@domain.com',
		'abc\ is\"not\valid@domain.com',
		'test@domain..com',
		'Abc.example.com',
		'Test @ Description <Abc.example.com>',
		true,
		#20200101,
		113,
		'test@domain.123',
		'test@domain.123com',
		'test@domain.123com...',
		'test@domain.',
		'test@domain.a',
		'test@domain.1',
		'test@domain...',
		'test@domain.this-is-to-test-limit',
		'test@domain.com-'
		)
	}