// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CopyFieldName(fieldName)
		{
		if fieldName isnt false
			WndProc.Copy_Field_Name(fieldName)
		}

	GoToFieldDefinition(fieldName)
		{
		if fieldName isnt false
			GotoLibView('Field_' $ fieldName)
		}
	}