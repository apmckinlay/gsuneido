// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
Field_num
	{
	Prompt: 'Channel'
	Control: (Id "im_channels"
		columns: (imchannel_name, imchannel_abbrev, imchannel_desc),
		field: imchannel_num,
		width: 30)
	Format: (Id query: 'im_channels' numField: imchannel_num, width: 16)
	}
