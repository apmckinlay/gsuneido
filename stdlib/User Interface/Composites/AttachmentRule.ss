// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
function (value, position)
	{
	perRow = AttachmentsRepeatControl.PerRow
	--position // convert from 1 based to 0 based
	row = (position / perRow).Int()
	if not Object?(value) or row >= value.Size()
		return ''
	col = position % perRow
	return value[row]['attachment' $ col]
	}
