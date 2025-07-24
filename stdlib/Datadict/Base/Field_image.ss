// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Prompt: ''
	Heading: 'Image'
	Control: (OpenImage xmin: 300 ymin: 300)
	Format: (Image width: 1440)
	AllowCustomizableOptions: #(formula: false)
	Encode(val)
		{
		return String(val)
		}
	}