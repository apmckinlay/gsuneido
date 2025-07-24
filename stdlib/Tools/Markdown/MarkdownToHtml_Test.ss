// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		test = {|md, html|
			Assert(MarkdownToHtml(.S(md)).RightTrim() is: .S(html))
			}
		test("", "")
		test("hello world",
			"<p>hello world</p>")
		test("hello
			big
			world",
			"<p>hello
			big
			world</p>")
		test(
			"hello

			world",
			"<p>hello</p>

			<p>world</p>")
		test("```
			foo
			bar",
			"<pre>foo
			bar
			</pre>")
//		test("```
//			foo
//			bar
//			```",
//			"<pre>foo
//			bar
//			</pre>")
		}
	S(s)
		{
		return s.Replace('^\t\t\t').Tr('\r')
		}
	}