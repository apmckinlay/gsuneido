// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
// returns an object of ALL nodes that match the specified nodePath
class
	{
	First(xml, nodePath)
		{
		nodes = .run(xml, nodePath, { |name, path, empty?| name is path and empty?})
		return nodes.Empty?() ? false : nodes[0]
		}

	All(xml, nodePath)
		{
		nodes = .run(xml, nodePath, { |name, path, empty?/*unused*/| name is path})
		return nodes
		}

	run(xml, nodePath, continue?)
		{
		nodes = Object()
		if not nodePath.Empty?()
			.recursive(xml, nodePath, nodes, continue?)
		return nodes
		}

	recursive(xml, nodePath, nodes, continue?)
		{
		if continue?(xml.Name(), nodePath[0], nodes.Empty?())
			nodePath.Size() is 1
				? nodes.Add(xml)
				: xml.Children().Each({ .recursive(it, nodePath[1..], nodes, continue?) })
		}
	}