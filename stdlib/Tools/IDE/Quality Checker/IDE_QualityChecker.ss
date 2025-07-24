// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Menu_Quality_Checker()
		{
		return Object("Enable", "Disable")
		}

	On_Quality_Checker(option)
		{
		if option is "Enable"
			UserSettings.Put("QcEnabled", true)
		else if option is "Disable"
			UserSettings.Put("QcEnabled", false)
		else
			{
			Print("Something went wrong in IDE_QualityChecker. " $
				"Quality checks enabled by default.")
			UserSettings.Put("QcEnabled", true)
			}
		}
	}