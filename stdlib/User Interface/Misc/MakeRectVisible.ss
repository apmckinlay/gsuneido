// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
function (rect)
	{
	rcWork = GetWorkArea(rect)
	if rect.right > rcWork.right
		{
		rect.left -= rect.right - rcWork.right
		rect.right -= rect.right - rcWork.right
		}
	if rect.left < rcWork.left
		{
		rect.right += rcWork.left - rect.left
		rect.left += rcWork.left - rect.left
		}
	if rect.bottom > rcWork.bottom
		{
		rect.top -= rect.bottom - rcWork.bottom
		rect.bottom -= rect.bottom - rcWork.bottom
		}
	if rect.top < rcWork.top
		{
		rect.bottom += rcWork.top - rect.top
		rect.top += rcWork.top - rect.top
		}
	return rect
	}
