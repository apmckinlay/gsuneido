// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Name: "BookPath"
	New()
		{
		.horz = .Vert.Horz
		}
	Controls:
		(Vert
			(Skip 6)
			(Horz Fill (Html_ahref_ text: Contents href: "/Contents") (Static " ")))
	SetPath(path)
		{
		while (.horz.GetChildren().Size() > 1)
			.horz.Remove(1)
		path = path.Split('/')
		path.Delete(path.Size() - 1)
		while ((pathsize = path.Size()) > 1)
			{
			.horz.Insert(1, Object("Html_ahref_Control",
				text: path[pathsize - 1], href: path.Join('/')))
			if (pathsize > 2)
				.horz.Insert(1, #(Static ' > '))
			path.Delete(pathsize - 1)
			}
		}
	Goto(address)
		{ .Send("Goto", address) }
	}
