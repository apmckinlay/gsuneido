// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
// used by XmlParser
class
	{
	UseDeepEquals: true
	text: ""
	New(.name = "", attributes = false, .text = "", children = false)
		{
		.attributes = attributes is false ? Object() : attributes
		.children = children is false ? Object() : children
		}

	AddChild(childXmlNode)
		{
		.children.Add(childXmlNode)
		}
	Name()
		{ return .name }
	Attributes()
		{ return .attributes }
	Children()
		{ return .children }
	Text()
		{
		return .text isnt ''
			? .text
			: .children.Map(#Text).Join()
		}

	Getter_(name) // .element => list of children with that name
		{
		if Integer?(name) // handle [i]
			return .children[name]
		else if name.Prefix?('_')
			return .attributes[name[1 ..]]
		else
			return XmlNodeList(.children.Filter() { it.Name() is name })
		}
	ToString()
		{
		if .text isnt ''
			return .text
		s = '<' $ .name
		for member,value in .attributes
			s $= ' ' $ member $ '="' $ value $ '"'
		if .children.Empty?()
			return s $ ' />\n'
		s $= '>\n'
		.children.Each
			{
			for line in it.ToString().Lines()
				s $= '\t' $ line $ '\n'
			}
		s $= '</' $ .name $ '>\n'
		return s
		}

	ToObject()
		{
		ob = Object()
		if .text isnt ''
			return .text
		for member, value in .attributes
			ob[member] = value
		ob[.name] = Object()
		.children.Each
			{
			child = it.ToObject()
			if String?(child)
				ob[.name] = child
			else
				ob[.name].Add(child)
			}
		return ob
		}
	}