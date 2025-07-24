// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
class
	{
	stickyFields: false
	New(stickyFields = false)
		{
		Assert(Object?(stickyFields) or (stickyFields is false))
		if stickyFields is false
			{
			.stickyFields = false
			return
			}
		.stickyFields = Object()
		for field in stickyFields
			.stickyFields[field] = ''
		}

	GetStickyFields()
		{
		return .stickyFields
		}

	UpdateStickyField(record, member, new? = false)
		{
		if new? is false
			new? = record.New?()
		if (new? and Object?(.stickyFields) and .stickyFields.Member?(member))
			.stickyFields[member] = record.Copy()[member]
		}

	SetRecordStickyFields(record)
		{
		if .stickyFields is false
			return
		for (stickyField in .stickyFields.Members())
			if ((not record.Member?(stickyField) or record[stickyField] is "") and
				.stickyFields[stickyField] isnt "")
				record[stickyField] = .stickyFields[stickyField]
		}

	ClearStickyFieldValues()
		{
		if not Object?(.stickyFields)
			return
		for field in .stickyFields.Members()
			.stickyFields[field] = ''
		}

	RemoveInvalidStickyValue(field, invalidVal)
		{
		if Object?(.stickyFields) and .stickyFields.Member?(field) and
			.stickyFields[field] is invalidVal
			.stickyFields[field] = ''
		}
	}
