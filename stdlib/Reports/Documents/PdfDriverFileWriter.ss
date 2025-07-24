// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New(.filename)
		{
		.f = .file(filename)
		.curPos = 0
		.locations = Object()
		.s = ''
		}

	file(filename)
		{
		return File(filename, mode: 'w')
		}

	Reserve()
		{
		.locations.Add(false)
		}

	AddObject(content, id = false)
		{
		if id is false
			{
			nextObj = .NextObjId
			.locations.Add(.curPos)
			}
		else
			{
			Assert(.locations[id] is: false)
			nextObj = id
			.locations[id] = .curPos
			}
		s = (nextObj is 0 ? '' : nextObj $ ' 0 obj ') $ content
		.s $= s
		.curPos += s.Size()
		return nextObj
		}

	Write(content)
		{
		.s $= content
		.curPos += content.Size()
		}

	Getter_NextObjId()
		{
		return .locations.Size()
		}

	Getter_Locations()
		{
		return .locations
		}

	Getter_TotalLength()
		{
		return .curPos
		}

	Flush()
		{
		.f.Write(.s)
		.s = ''
		}

	Finish()
		{
		.Flush()
		.f.Close()
		}

	Abort()
		{
		.f.Close()
		DeleteFile(.filename)
		}
	}