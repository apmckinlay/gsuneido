// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_AddHeader()
		{
		ob = Object(Fields: Object(),
			MimeBase_hdr: Object(),
			MimeBase_extra: Object())
		ob.Eval(MimeBase.AddHeader, 'Content-Transfer-Encoding', 'base64')
		Assert(ob.Fields is: #('Content-Transfer-Encoding'))
		Assert(ob.MimeBase_hdr is: #('Content-Transfer-Encoding': 'base64'))
		Assert(ob.MimeBase_extra is: #())

		ob.Eval(MimeBase.AddHeader, 'Content-Transfer-Encoding', 'json')
		Assert(ob.Fields is: #('Content-Transfer-Encoding'))
		Assert(ob.MimeBase_hdr is: #('Content-Transfer-Encoding': 'json'))
		Assert(ob.MimeBase_extra is: #())

		ob.Eval(MimeBase.AddHeader,
			'Content-Disposition', 'attachment', 'filename': 'test.txt')
		Assert(ob.Fields is: #('Content-Transfer-Encoding', 'Content-Disposition'))
		Assert(ob.MimeBase_hdr is: #('Content-Transfer-Encoding': 'json',
			'Content-Disposition': 'attachment'))
		Assert(ob.MimeBase_extra is: #('Content-Disposition': '; filename="test.txt"'))
		}
	}