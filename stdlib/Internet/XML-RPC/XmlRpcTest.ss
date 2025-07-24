// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	data:
		(
		(true, "<value><boolean>1</boolean></value>")
		(false, "<value><boolean>0</boolean></value>")
		(123, "<value><int>123</int></value>")
		(123.456, "<value><double>123.456</double></value>")
		("hello", "<value><string>hello</string></value>")
		('&<>"', "<value><string>&amp;&lt;&gt;&quot;</string></value>")
		(#20030825,
			"<value><dateTime.iso8601>20030825T00:00:00</dateTime.iso8601></value>")
		(#20030825.1613,
			"<value><dateTime.iso8601>20030825T16:13:00</dateTime.iso8601></value>")
		(#(1, hello),
			"<value><array><data><value><int>1</int></value>" $
			"<value><string>hello</string></value></data></array></value>")
		(#(a: 1),
			"<value><struct><member><name>a</name>" $
			"<value><int>1</int></value></member></struct></value>")
		)
	Test_EncodeDecode()
		{
		for d in .data
			{
			Assert(XmlRpc.EncodeValue(d[0]) is: d[1])
			Assert((new XmlRpc).Decode2(d[1])[0] is: d[0])
			}
		}
	Test_EncodeCall()
		{
		Assert(XmlRpc.EncodeCall(#(addr F hello))
			like: '<?xml version="1.0"?>\r\n' $
				'<methodCall><methodName>F</methodName>\r\n' $
				'<params><param><value><string>hello</string></value></param>\r\n' $
				'</params>\r\n</methodCall>\r\n')
		}
	Test_EncodeResponse()
		{
		Assert(XmlRpc.EncodeResponse(123)
			like: '<?xml version="1.0"?>\r\n<methodResponse>' $
				'<params><param><value><int>123</int></value></param></params>' $
				'</methodResponse>\r\n')
		}
	Test_Decode2()
		{
		Assert(
			(new XmlRpc).Decode2('<?xml version="1.0"?>\r\n' $
				'<methodCall><methodName>F</methodName>\r\n' $
				'<params><param><value><string>hello</string></value></param>\r\n' $
				'</params>\r\n' $
				'</methodCall>')
			is: #("F", "hello"))

		msg = '<?xml version="1.0"?>
<methodResponse><params><param><value><string /></value></param></params></methodResponse>
'
		Assert((new XmlRpc).Decode2(msg) is: #(''))
		}
	}
