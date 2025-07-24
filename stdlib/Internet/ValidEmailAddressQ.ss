// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
function (addr)
	{
	if not String?(addr)
		return false

	if false isnt s = addr.Extract('^[^<>]*?<([^<>]*?)>$')
		addr = s
	addr = addr.Trim()
	if addr.BeforeFirst('@').Size() > 64 /* = user max length */ or
		addr.AfterFirst('@').Size() > 255 /* = domain max length */
		return false

	// the actual size limit for top level domain is 63
	// see https://tools.ietf.org/id/draft-liman-tld-names-00.html
	// 18 is the longest used at the moment
	if addr.AfterLast(`.`).Trim().Size() > 20 /*= top level domain max length*/
		return false

	// was made using two sources:
	//	 https://en.wikipedia.org/wiki/Email_address#Valid_email_addresses
	//	 http://www.regular-expressions.info/email.html
	// this regex does not support the use of unicode characters, double quotes in the
	// local part, or an IP as the domain part, though these are all technically valid
	userElement = "[a-z0-9!#$%&'*+/=?^_`{|}~-]+?"
	userPart = userElement $ "(\." $ userElement $ ")*?"
	domainElement = "[a-z0-9]([a-z0-9-]*?[a-z0-9])?"
	topDomainElement = "[a-z][a-z0-9-]*?[a-z0-9]"
	domainPart = "(" $ domainElement $ "\.)+?" $ topDomainElement
	addressRegex = "(?i)^" $ userPart $ "@" $ domainPart $ "$"

	return addr =~ addressRegex
	}