// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
TableModel
	{
	Table: 'im_channels'
	Name: 'Instant Messenger Channels'

	Columns: (imchannel_name, imchannel_num, imchannel_abbrev, imchannel_status,
		imchannel_desc)
	Keys: (imchannel_name, imchannel_num)
	UniqueIndexes: (imchannel_abbrev)
	}
