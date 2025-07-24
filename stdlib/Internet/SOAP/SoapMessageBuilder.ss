// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
class
	{
	xmlVersion: "1.0"
	SetXMLVersion(.xmlVersion)
		{
		}

	xmlEncoding: "utf-8"
	SetXMLEncoding(.xmlEncoding)
		{
		}

	XMLHeader()
		{
		return "<?xml version=" $ Display(.xmlVersion) $
			" encoding=" $ Display(.xmlEncoding) $ "?>"
		}

	soapEnvelope: "http://schemas.xmlsoap.org/soap/envelope/"
	SetSOAPEnvelopeNS(.soapEnvelope)
		{
		}

	xsd: "http://www.w3.org/2001/XMLSchema"
	SetSOAP_xsd(.xsd)
		{
		}

	xsi: "http://www.w3.org/2001/XMLSchema-instance"
	SetSOAP_xsi(.xsi)
		{
		}

	extraEnvelopeNameSpaces: #()
	// object where members are namespace names, vals are the urls
	AddEnvelopeNameSpaces(.extraEnvelopeNameSpaces)
		{
		}

	Envelope()
		{
		env = '<soap:Envelope' $
			' xmlns:soap=' $ Display(.soapEnvelope) $
			' xmlns:xsd=' $ Display(.xsd) $
			' xmlns:xsi=' $ Display(.xsi)

		for ns in .extraEnvelopeNameSpaces.Members().Sort!()
			env $= ' xmlns:' $ ns $ '=' $ Display(.extraEnvelopeNameSpaces[ns])

		return env $ '>'
		}

	// header XML attribute tags that must be processed by server should contain: mustUnderstand="true"
	Header(hdrXML = '')
		{
		return Opt('<soap:Header>', hdrXML, '</soap:Header>')
		}

	Body(bodyXML)
		{
		body = '<soap:Body>' $ bodyXML $ '</soap:Body>'
		return body $ '</soap:Envelope>' // close the envelope
		}
	}
