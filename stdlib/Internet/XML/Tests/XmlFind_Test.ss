// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	xmlAll: '<stuff attrib = "1">
			<things1>
					<things3>a</things3>
					<things4>c</things4>
					<things8>e</things8>
			</things1>
			<things2>
					<things3>b</things3>
					<things5>d</things5>
					<stuff attrib = "2" id = "5"></stuff>
			</things2>
			<things3>
					<things7>z</things7>
					<things7 attrib = "2" id = "5">x</things7>
			</things3>
			<things1>
					<things8>t</things8>
					<things8>u</things8>
					<things8>v</things8>
			</things1>
		</stuff>'
	Setup()
		{
		.parsed = XmlParser(.xmlAll)
		}

	Test_All()
		{
		nodes = XmlFind.All(.parsed, #(stuff))
		Assert(nodes isSize: 1)
		result = nodes[0]
		Assert(result.Name() is: 'stuff')
		Assert(result.Attributes().attrib is: '1')
		Assert(result.Children() isSize: 4)

		.testProcess(#(stuff, things1, things3), #(a))
		.testProcess(#(stuff, things2, things3), #(b))
		.testProcess(#(stuff, things2, things5), #(d))
		.testProcess(#(stuff, things1), #(ace, tuv))
		.testProcess(#(stuff, things1, things8), #(e,t,u,v))
		}

	testProcess(nodePath, texts)
		{
		nodes = XmlFind.All(.parsed, nodePath)
		Assert(nodes isSize: texts.Size())
		n = 0
		for node in nodes
			Assert(node.Text() is: texts[n++])
		return nodes
		}

	Test_Attributes()
		{
		nodes = .testProcess(#(stuff, things2, stuff), #(''))
		Assert(nodes[0].Attributes().id is: '5')
		Assert(nodes[0].Attributes().attrib is: '2')

		nodes = .testProcess(#(stuff, things3, things7), #(z, x))
		Assert(nodes[0].Attributes() hasnt: 'id')
		Assert(nodes[0].Attributes() hasnt: 'attrib')
		Assert(nodes[1].Attributes().id is: '5')
		Assert(nodes[1].Attributes().attrib is: '2')
		}

	Test_Invalids()
		{
		Assert(XmlFind.All(.parsed, #()) is: #())
		Assert(XmlFind.All(.parsed, #(things1)) is: #())
		Assert(XmlFind.All(.parsed, #(nonexistent)) is: #())
		Assert(XmlFind.All(.parsed, #(xml)) is: #())
		Assert(XmlFind.All(.parsed, #(nonexistent1, nonexistent2)) is: #())
		Assert(XmlFind.All(.parsed, #(stuff, things1, things6)) is: #())
		Assert(XmlFind.All(.parsed, #(stuff, things1, things5)) is: #())
		Assert(XmlFind.All(.parsed, #(stuff, things1, things3, things11)) is: #())
		}



	xmlFirst: '<stuff attrib = "1">
			<things1>
					<things3>a</things3>
<things3>aa</things3>
					<things4>c</things4>
			</things1>
			<things2>
					<things3>b</things3>
					<things5>d</things5>
					<stuff attrib = "2" id = "5"></stuff>
			</things2>
		</stuff>'
	Test_First()
		{
		parsed = XmlParser(.xmlFirst)

		// invalid
		Assert(XmlFind.First(parsed, #(things1)) is: false)
		Assert(XmlFind.First(parsed, #(nonexistent)) is: false)
		Assert(XmlFind.First(parsed, #(nonexistent1, nonexistent2)) is: false)
		Assert(XmlFind.First(parsed, #(stuff, things1, things6)) is: false)
		Assert(XmlFind.First(parsed, #(stuff, things1, things5)) is: false)
		Assert(XmlFind.First(parsed, #(stuff, things1, things3, things11)) is: false)

		result = XmlFind.First(parsed, #(stuff))
		Assert(result.Name() is: 'stuff')
		Assert(result.Attributes().attrib is: '1')
		Assert(result.Children() isSize: 2)

		Assert(XmlFind.First(parsed, #(stuff, things1, things3)).Text() is: 'a')
		Assert(XmlFind.First(parsed, #(stuff, things2, things3)).Text() is: 'b')
		Assert(XmlFind.First(parsed, #(stuff, things2, things5)).Text() is: 'd')
		Assert(XmlFind.First(parsed, #(stuff, things2, stuff)).Text() is: '')
		Assert(XmlFind.First(parsed, #(stuff, things2, stuff)).Attributes().id is: '5')
		Assert(XmlFind.First(parsed, #(stuff, things2, stuff)).Attributes().attrib
			is: '2')
		}
	}