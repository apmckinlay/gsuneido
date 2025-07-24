// Copyright (C) 2013 Axon Development Corporation All rights reserved worldwide.
function (name)
	{
	return Query1(LastContribution('HelpBook'), :name, path: "/res").text
	}
