// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		Md_Examples.Each()
			{ |example|
			result = MarkdownToHtml(example.markdown, noIndent?:)
			Assert(result is: example.html, msg: 'example ' $ example.id)

			result = MarkdownToHtml(example.markdown, noIndent?:,
				addons: [Md_Addon_Table, Md_Addon_Definition])
			Assert(result is: example.html, msg: 'example ' $ example.id)
			}
		}

	tableCases: (
	(md: '
| title 1 | title 2 | title 3 | title 4 |
|:--------|--------:|:-------:| ------- |
| left    |   right | center  | default |', sbe: '
<table>
	<tr>
		<th style="text-align: left;">title 1</th>
		<th style="text-align: right;">title 2</th>
		<th style="text-align: center;">title 3</th>
		<th style="text-align: left;">title 4</th>
	</tr>
	<tr>
		<td style="text-align: left;">left</td>
		<td style="text-align: right;">right</td>
		<td style="text-align: center;">center</td>
		<td style="text-align: left;">default</td>
	</tr>
</table>'),
	// omit leading/trailing pipes
	(md: '
 title 1 | title 2
:--------|--------:', sbe: '
<table>
	<tr>
		<th style="text-align: left;">title 1</th>
		<th style="text-align: right;">title 2</th>
	</tr>
</table>'),
	// escaped pipes
	(md: '
| name \\| id | description |
|- | - |', sbe: '
<table>
	<tr>
		<th style="text-align: left;">name | id</th>
		<th style="text-align: left;">description</th>
	</tr>
</table>'),
	// inline formatting
	(md: '
| ***name*** | note |
|------|------|
| *foo* | `bar` |', sbe: '
<table>
	<tr>
		<th style="text-align: left;"><em><strong>name</strong></em></th>
		<th style="text-align: left;">note</th>
	</tr>
	<tr>
		<td style="text-align: left;"><em>foo</em></td>
		<td style="text-align: left;"><code>bar</code></td>
	</tr>
</table>'),
	// Uneven column counts (extra cell ignored)
	(md: '
| a | b |
|---|---|
| 1 |   | 3 |
| 4 |', sbe: '
<table>
	<tr>
		<th style="text-align: left;">a</th>
		<th style="text-align: left;">b</th>
	</tr>
	<tr>
		<td style="text-align: left;">1</td>
		<td style="text-align: left;"></td>
	</tr>
	<tr>
		<td style="text-align: left;">4</td>
	</tr>
</table>'),
	// table between other items
	(md: '
This is a paragraph
| \\x\\| |
| - |
- list', sbe: '
<p>This is a paragraph</p>
<table>
	<tr>
		<th style="text-align: left;">\\x|</th>
	</tr>
</table>
<ul>
	<li>list </li>
</ul>'))
	Test_table()
		{
		for ob in .tableCases
			Assert(MarkdownToHtml(ob.md, addons: Object(Md_Addon_Table)) like: ob.sbe)
		}

	Test_misc()
		{
		Assert(MarkdownToHtml('
* list item line

*line*') like: '
<ul>
	<li>list item line </li>
</ul>
<p><em>line</em></p>')
		}
	}