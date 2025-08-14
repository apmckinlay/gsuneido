// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Field_string
	{
	Prompt: 'Zip/Postal'
	Control: (ZipPostal)
	Heading: 'Zip Code/\nPostal Code'
	Format: (Text width: 8)
	AI_Prompt: "Extract the Postal or ZIP code, Rules: " $
		"1. ZIP Code: have either 5 digits (e.g. 12345) or ZIP-4 (e.g. 12345-6789); or " $
		"2. Postal Code: will use the pattern: `[a-Z][0-9][a-Z]\s?[0-9][a-Z][0-9]`."
	}
