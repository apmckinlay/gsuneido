// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
// reference: amazon_sqss_helper.rb
function (type, response)
	{
	tagContent = function (s, tag)
		{
		return s.AfterFirst('<' $ tag $ '>').BeforeFirst('</' $ tag $ '>')
		}
	switch type
		{
	case 'send' :
		return Object(
			md5: tagContent(response, 'MD5OfMessageBody')
			id: tagContent(response, 'MessageId'))
	case 'receive' :
		return response.Split('</Message>').
			Filter({ it.Has?('<Message>') }).
			Map({ Object(
				md5: tagContent(it, 'MD5OfBody'),
				id: tagContent(it, 'MessageId')
				receipt: tagContent(it, 'ReceiptHandle')
				body: XmlEntityDecode(tagContent(it, 'Body'))) })
	case 'delete':
		return true
		}
	}
