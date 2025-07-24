// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(attachments)
		{
		if not Object?(attachments)
			return ''
		attachmentList = Object()
		for row in attachments
			for m in row.Members().Sort!()
				attachmentList.AddUnique(row[m])
		return attachmentList.Join(", ")
		}
	}