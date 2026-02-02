// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	cases: ((md: '
This is a paragraph
Term1
: Definition 1 - a
: Definition 1 - b
Definition 1 - b cont.

Term2
:	Definition 2

Another paragraph
TermA
:\t
', sbe: '<p>This is a paragraph</p>
<dl>
<dt>Term1</dt>
<dd>Definition 1 - a</dd>
<dd>Definition 1 - b
Definition 1 - b cont.</dd>
<dt>Term2</dt>
<dd>Definition 2</dd>
</dl>
<p>Another paragraph</p>
<dl>
<dt>TermA</dt>
</dl>'))
	Test_one()
		{
		for ob in .cases
			Assert(MarkdownToHtml(ob.md, addons: Object(Md_Addon_Definition))
				like: ob.sbe)
		}
	}