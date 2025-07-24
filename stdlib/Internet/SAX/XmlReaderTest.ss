// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(.parse('') is: [])
		Assert(.parse('stuff') is: #(['chars', 'stuff']))
		Assert(.parse('bad &amp; chars') is: #(['chars', 'bad & chars']))
		Assert(.parse('more decode testing&lt;&gt;&amp;&quot;&apos;')
			is: #(['chars', 'more decode testing<>&"\'']))
		Assert(.parse('<tag>') is: #(['start', 'tag', #()]))
		Assert(.parse('stuff<tag>') is: #(['chars', 'stuff'], ['start', 'tag', #()]))
		Assert(.parse('<tag>stuff') is: #(['start', 'tag', #()], ['chars', 'stuff']))
		Assert(.parse('<tag></tag>') is: #(['start', 'tag', #()], ['end', 'tag']))
		Assert(.parse('<tag />') is: #(['start', 'tag', #()], ['end', 'tag']))
		Assert(.parse('<tag>stuff</tag>')
			is: #(['start', 'tag', #()], ['chars', 'stuff'], ['end', 'tag']))
		Assert(.parse('<tag>stuff<!-- comment --></tag>')
			is: #(['start', 'tag', #()], ['chars', 'stuff'], ['end', 'tag']))
		Assert(.parse('<tag><![CDATA[ a <tag> here ]]></tag>')
			is: #(['start', 'tag', #()], ['chars', ' a <tag> here '], ['end', 'tag']))
		Assert(.parse('<tag color="red">') is: #(['start', 'tag', #(color: 'red')]))
		Assert(.parse('<tag color="red" size="12">')
			is: #(['start', 'tag', #(color: 'red', size: '12')]))
		Assert(.parse('<tag\tcolor="red" size="12">')
			is: #(['start', 'tag', #(color: 'red', size: '12')]))
		Assert(.parse('<tag\tcolor="red"\r\nsize="12">')
			is: #(['start', 'tag', #(color: 'red', size: '12')]))
		Assert(.parse('<tag\r\n\tcolor="red"\r\nsize="12">')
			is: #(['start', 'tag', #(color: 'red', size: '12')]))
		Assert(.parse('<tag color ="red" size = "12">')
			is: #(['start', 'tag', #(color: 'red', size: '12')]))
		Assert(.parse('<tag nowrap color="red">')
			is: #(['start', 'tag', #(nowrap: true, color: 'red')]))
		Assert(.parse('<Tag ALIGN=Left>') is: #(['start', 'tag', #(align: Left)]))
		Assert(.parse('<tag>=</tag>')
			is: #(['start', 'tag', #()], ['chars', '='], ['end', 'tag']))
		}
	Test_gt_in_quoted_attribute_value()
		{
		Assert(Suneido.x = .parse('<p class="hello_<$= class_name $>">Hello, World!</p>')
			is: #(['start', 'p', #(class: "hello_<$= class_name $>")],
				['chars', 'Hello, World!'],
				['end', 'p']))
		}
	parse(text)
		{
		xr = new XmlReader
		xr.SetContentHandler(logger = new .logger)
		xr.Parse(text)
		return logger.Log
		}
	logger: class
		{
		New()
			{ .Log = Object() }
		StartElement(qname, atts)
			{ .Log.Add(['start', qname, atts]) }
		EndElement(qname)
			{ .Log.Add(['end', qname]) }
		Characters(string)
			{ .Log.Add(['chars', string]) }
		}
	}