// Copyright (C) 2004 Suneido Software Corp. All rights reserved worldwide.
function (source, newtable)
	{
	Database("create " $ Schema(source).Replace(source, newtable, 1))
	QueryDo("insert " $ source $ " into " $ newtable)
	}