// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		examples = Md_Examples
		for (i = 41; i < 632; i++)
			{
			if i in (
				//inline
				55, 64, 65, 75, 79, 80, 81, 87, 101, 105, 154, 166, 167, 175,
				176, 187, 120, 137, 144, 147, 151
				// hard line break
				225
				// loose
				324, 325
				// autolink
				345

				// excape
				631, 605
				// url encode
				597, 598, 604
				)
				continue
			// link ref def
			if i >= 191 and i <= 217 or i is 316
				continue

			if i >= 349 and i < 593
				continue
			result = MarkdownToHtml(examples[i].markdown, noIndent?:)
			Assert(result is: examples[i].html, msg: 'example ' $ examples[i].id)
			}
		}
	}