// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_empty_string()
		{
		Assert(Razor("") is: "")
		}
	Test_just_html()
		{
		s = "<h1>hello</h1>"
		Assert(Razor(s) is: s)
		}
	Test_email_address()
		{
		s = '<a href="mailto:fred@@ibm.com">fred@@ibm.com</a>'
		Assert(Razor(s) is: s.Replace('@@', '@'))
		}
	Test_var()
		{
		s = "<b>@.name</b>"
		x = #(name: Fred)
		Assert(Razor(s, x) is: "<b>Fred</b>")
		}
	Test_expr()
		{
		s = "<b>@(.x + .y)</b>"
		x = #(x: 12, y: 34)
		Assert(Razor(s, x) is: "<b>46</b>")
		}
	Test_nested_code()
		{
		s = "<a>@{ x=123 <b> @(x) </b> }</a>"
		Assert(Razor(s) is: "<a><b> 123 </b></a>")
		}
	Test_var_member()
		{
		s = "<b>@.x.name</b>"
		context = #(x: (name: fred))
		Assert(Razor(s, context) is: '<b>fred</b>')
		}
	Test_function_call()
		{
		s = "<b>@.name.Size()</b>"
		context = #(name: fred)
		Assert(Razor(s, context) is: '<b>4</b>')
		}
	Test_subscript()
		{
		s = "<b>@.name[0]</b>"
		context = #(name: fred)
		Assert(Razor(s, context) is: '<b>f</b>')
		}
	Test_chained()
		{
		s = "<b>@.name[0].Upper()</b>"
		context = #(name: fred)
		Assert(Razor(s, context) is: '<b>F</b>')
		}
	Test_var_in_attribute()
		{
		s = '<a href="@.url">link</a>'
		context = #(url: "http://ibm.com")
		Assert(Razor(s, context) is: s.Replace("@.url", context.url))
		}
	Test_nested_outer_tag()
		{
		s = "<span>one<span>two</span>three</span>"
		Assert(Razor(s) is: s)
		}
	Test_if()
		{
		s = "<p>@if (.flag) { <b>hello</b> }</p>"
		Assert(Razor(s, #(flag: true)) is: "<p><b>hello</b></p>")
		Assert(Razor(s, #(flag: false)) is: "<p></p>")
		}
	Test_encode()
		{
		s = "<p>@.name</p>"
		x = #(name: "<joe & sue>")
		Assert(Razor(s, x) is: "<p>" $ XmlEntityEncode(x.name) $ "</p>")
		}
	Test_nested_expr()
		{
		s = `<tr>@for(field in .fields){@field}</tr>`
		Assert(Razor(s, [fields: #(a, b)]) is: '<tr>ab</tr>')

		s = `<tr>@for(field in .fields){@HtmlString(field)}</tr>`
		Assert(Razor(s, [fields: #(a, b)]) is: '<tr>ab</tr>')

		s = `<tr>@for(field in .fields){@field\n<td>@field</td>}</tr>`
		Assert(Razor(s, [fields: #(a, b)]) is: '<tr>a<td>a</td>b<td>b</td></tr>')

		s = `<tr>@for(field in .fields){@field@field}</tr>`
		Assert(Razor(s, [fields: #(a, b)]) is: '<tr>aabb</tr>')

		s = `<tr>@for(field in .fields){@field\n<td>@field</td>@field}</tr>`
		Assert(Razor(s, [fields: #(a, b)]) is: '<tr>a<td>a</td>ab<td>b</td>b</tr>')

		s = `<tr>@if(true){@.field}</tr>`
		Assert(Razor(s, [field: #a]) is: '<tr>a</tr>')

		s = `<tr>@{i=0}@while(i<3){@.field;i++}</tr>`
		Assert(Razor(s, [field: #a]) is: '<tr>aaa</tr>')
		}
	Test_text()
		{
		s = "<text>Hello world</text>"
		Assert(Razor(s, #()) is: "Hello world")
		}
	}