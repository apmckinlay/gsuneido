// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
function (companyNum = false)
	{
	portFunc = OptContribution('HttpPort', function(){ return 80 /*= default httpport*/})
	return (portFunc)(:companyNum)
	}
