// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
TableModel
	{
	Table: 'im_history'

	Columns: (im_num, im_from, im_to, im_message, imchannel_num)
	Keys: (im_num)
	ForeignKeys: (im_channels: ((from: imchannel_num)))
	}
