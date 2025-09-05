// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(search)
		{
		return '<!DOCTYPE html>
			<html>
			<head>
			<title>Find "' $ search $ '"</title>
			<style>
				.search_key { color: #1a0dab; font-size: 1.2em; }
				.search_key:visited { color: #609; }
				.search_snippets { color: #545454; margin-bottom: 10px;
					font-size: 1em; max-width: 1000px; }
				.container {
					display: flex;
					flex-wrap: wrap;
				}
				.block {
					width: 50%;
					overflow-wrap: break-word;
					box-sizing: border-box;
				}
				@media (max-width: 1000px) {
					.block {
						width: 100%;
					}
				}
			</style>
			</head>
			<body>' $
			WikiTemplate.Find.Replace('\$smallaction', '') $
			'<h1 style="margin-top: 0;">' $ OptContribution('WikiTemplateLogo', '') $
			'Find Results</h1><div class="container">' $
			'<div class="block">' $ .findByFtsearch(search) $ '</div>' $
			'<div class="block" id="semantic"></div>' $
			`</div><p>Return to <a href="Wiki?StartPage">Start Page</a>
			<script>
				fetch('Wiki?semantic=` $ search $ `')
					.then(response => response.text())
					.then(html => {
						document.getElementById('semantic').innerHTML = html;
					})
					.catch(error => {
						console.error('Error fetching content:', error);
						document.getElementById('semantic').innerHTML =
							'<p>Sorry, the content could not be loaded.</p>';
					});
			</script>
			</body></html>`
		}
	findByFtsearch(search)
		{
		if String(search).Blank?()
			return "<p><b>Please enter something to search for</b></p>"

		recs = .Ftsearch(search)
		results = recs.Map(.format1).Join()
		return '<p>Pages matching: <em>' $ search $ '</em></p>\n' $ results
		}

	Ftsearch(search)
		{
		search = String(search)
		idx = Suneido.GetInit(IndexWiki.Index,
			{ Ftsearch.Load(GetFile(IndexWiki.Index)) })
		resIDs = #()
		try
			{
			resIDs = idx.Search(search)
			if resIDs.Empty?()
				SuneidoLog('INFO: Ftsearch return empty. Searched - ' $ Display(search))
			}
		catch (e)
			SuneidoLog('ERROR: Ftsearch - ' $ e)
		return QueryAll('wiki remove text where num in (' $ resIDs.Join(',') $ ')')
		}

	format1(x)
		{
		name = String?(x) ? x : x.name
		return '<a class="search_key" ' $ 'href="Wiki?' $ name $ '">' $ name $ '</a> ' $
			.displayOrphanedMessage(x) $ '<br/>'
		}

	displayOrphanedMessage(x)
		{
		return (Object?(x) ? x.orphaned is true : WikiOrphans.ListOrphans().Has?(x))
			? WikiOrphans.OrphanedMessage
			: ''
		}

	Semantic(search)
		{
		results = ""
		search = Type(search) isnt 'String' ? String(search) : search
		if not search.Blank?() and KnowledgeBase.Available?()
			results = KnowledgeBase.Query(search, #('Wiki'), n: 20).
				Map({ .format1(it.link.AfterFirst('Wiki?')) }).Join()
		return `<!DOCTYPE html>
			<html>
			<head></head>
			<body>
			<p>Pages related to: <em>` $ search $ `</em></p>` $
			results $ `</p>
			</body>
			</html>`
		}
	}
