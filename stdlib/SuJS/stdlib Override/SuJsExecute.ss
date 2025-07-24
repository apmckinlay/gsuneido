// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
function(lpVerb, lpFile, lpParameters = '')
	{
	switch (lpVerb)
		{
	case 'open':
		if lpFile.Match('^\w+?:') is false // not has scheme
			lpFile = `https://` $ lpFile
		SuUI.Open(lpFile)
	case 'download':
		SuDownloadFile(lpFile, lpParameters)
	default:
		Print("unsupported excute verb: " $ lpVerb)
		}
	}
