// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	m: MimeMultiPart
		{
		Boundary() { "====================" }
		}
	Test_mixed()
		{
		m = (.m)()
		m.To("fred@gmail.com")
		m.From('joe@shaw.com')
		m.Subject('multipart test')
		m.Attach(MimeText("part one"))
		m.Attach(MimeBase('application', 'octet-stream').
			AddHeader('Content-Disposition', 'attachment', filename: 'test.txt').
			SetPayload("second part").Base64())
		Assert(m.ToString() is: .mixed)
		}
mixed:
'MIME-Version: 1.0
From: joe@shaw.com
To: fred@gmail.com
Subject: multipart test
Content-Type: multipart/mixed; \r\n\tboundary="===================="

--====================
Content-Type: text/plain; charset="us-ascii"
Content-Transfer-Encoding: 7bit

part one

--====================
Content-Type: application/octet-stream
Content-Transfer-Encoding: base64
Content-Disposition: attachment; filename="test.txt"

c2Vjb25kIHBhcnQ=

--====================--
'
	Test_alternative()
		{
		m = (.m)('alternative')
		m.To("fred@gmail.com")
		m.From('joe@shaw.com')
		m.Subject('multipart test')
		m.Attach(MimeText("plain"))
		m.Attach(MimeText("<h1>html</h1>", 'html'))
		Assert(m.ToString() is: .alternative)
		}
	alternative:
'MIME-Version: 1.0
From: joe@shaw.com
To: fred@gmail.com
Subject: multipart test
Content-Type: multipart/alternative; \r\n\tboundary="===================="

--====================
Content-Type: text/plain; charset="us-ascii"
Content-Transfer-Encoding: 7bit

plain

--====================
Content-Type: text/html; charset="us-ascii"
Content-Transfer-Encoding: 7bit

<h1>html</h1>

--====================--
'

	Test_type()
		{
		type = MimeMultiPart.Type
		filename = 'test1.txt'
		Assert(type(filename) is: #('text', 'plain'))
		Assert(type(filename, 'application/octet-stream')
			is: #('application', 'octet-stream'))
		}
	}
