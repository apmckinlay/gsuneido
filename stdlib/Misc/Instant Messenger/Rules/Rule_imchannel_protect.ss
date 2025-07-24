// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
function()
	{
	if QueryEmpty?('im_history', imchannel_num: .imchannel_num_default)
		return false

	return #(imchannel_name:, reason: 'Channel has message history', allowDelete: false)
	}
