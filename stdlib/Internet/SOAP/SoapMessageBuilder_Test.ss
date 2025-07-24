// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		soapBuilder = new SoapMessageBuilder

		// test XML Header
		Assert(soapBuilder.XMLHeader() is: '<?xml version="1.0" encoding="utf-8"?>')
		soapBuilder.SetXMLVersion("1.1")
		Assert(soapBuilder.XMLHeader() is: '<?xml version="1.1" encoding="utf-8"?>')
		soapBuilder.SetXMLEncoding("other")
		Assert(soapBuilder.XMLHeader() is: '<?xml version="1.1" encoding="other"?>')

		// test envelope
		Assert(soapBuilder.Envelope() is: '<soap:Envelope' $
			' xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"' $
			' xmlns:xsd="http://www.w3.org/2001/XMLSchema"' $
			' xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">')

		soapBuilder.SetSOAP_xsd("xsd_url")
		Assert(soapBuilder.Envelope() is: '<soap:Envelope' $
			' xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"' $
			' xmlns:xsd="xsd_url"' $
			' xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">')

		soapBuilder.SetSOAP_xsi("xsi_url")
		Assert(soapBuilder.Envelope() is: '<soap:Envelope' $
			' xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"' $
			' xmlns:xsd="xsd_url"' $
			' xmlns:xsi="xsi_url">')

		soapBuilder.AddEnvelopeNameSpaces(#(test_ns1: "Test Name Space 1",
			test_ns2: "Test Name Space 2",
			test_ns3: "Test Name Space 3"))

		Assert(soapBuilder.Envelope() is: '<soap:Envelope' $
			' xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"' $
			' xmlns:xsd="xsd_url"' $
			' xmlns:xsi="xsi_url"' $
			' xmlns:test_ns1="Test Name Space 1"' $
			' xmlns:test_ns2="Test Name Space 2"' $
			' xmlns:test_ns3="Test Name Space 3">')

		soapBuilder.SetSOAPEnvelopeNS("another SOAP Env NS")
		Assert(soapBuilder.Envelope() is: '<soap:Envelope' $
			' xmlns:soap="another SOAP Env NS"' $
			' xmlns:xsd="xsd_url"' $
			' xmlns:xsi="xsi_url"' $
			' xmlns:test_ns1="Test Name Space 1"' $
			' xmlns:test_ns2="Test Name Space 2"' $
			' xmlns:test_ns3="Test Name Space 3">')

		// test header
		Assert(soapBuilder.Header() is: "")
		Assert(soapBuilder.Header("") is: "")
		Assert(soapBuilder.Header("Test Header XML")
			is: '<soap:Header>Test Header XML</soap:Header>')

		// test body, cody also closes the envelope
		Assert(soapBuilder.Body("") is: "<soap:Body></soap:Body></soap:Envelope>")
		Assert(soapBuilder.Body("<xml_testtag> Test Info</xml_testtag>")
			is: "<soap:Body><xml_testtag> Test Info</xml_testtag></soap:Body>" $
				"</soap:Envelope>")
		}
	}