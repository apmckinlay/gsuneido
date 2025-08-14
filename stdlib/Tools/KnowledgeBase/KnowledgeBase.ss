// Copyright (C) 2023 Axon Development Corporation All rights reserved worldwide.
class
	{
	table: 'knowledge_base_info'
	Ensure()
		{
		Database('ensure ' $ .table $ ' (openAI_api, index_name, vectorDB_api) key ()')
		}

	Setup(openAI_api, index_name, vectorDB_api)
		{
		.Ensure()
		QueryOutput(.table, [:openAI_api, :index_name, :vectorDB_api])
		Query1CacheReset()
		}

	Available?()
		{
		return not TestRunner.RunningTests?() and .getInfo() isnt false
		}

	getInfo()
		{
		if not TableExists?(.table)
			return false
		return Query1Cached(.table)
		}

	batch: 10
	Update(sourceName, forEach)
		{
		if false is info = .getInfo()
			return 0

		vdb = new VectorDB(info.index_name, info.vectorDB_api)
		openAI = new OpenAI_Embeddings(info.openAI_api)

		c = 0
		outputs = Object()
		forEach({ |title, content, key, extra|
			for i, s in .wrapText(content, openAI.MaxEmbeddingInputSize)
				outputs.Add([
					text: Opt('TITLE: ', title, '\r\n') $ s,
					id: .formatId(sourceName, key, i)
					metadata: Object(source: sourceName, key: String(key)).Merge(extra)])
			if outputs.Size() >= .batch
				{
				c += .update(outputs, openAI, vdb)
				outputs = Object()
				}
			})
		c += .update(outputs, openAI, vdb)
		return c
		}

	wrapText(s, width)
		{
		res = Object()
		while s.Size() > width
			{
			if ((false is pos = s.FindLast1of("\r\n", width)) and
				(false is pos = s.FindLast1of(" \t", width)))
				pos = width - 1
			sub = s[..pos+1].Trim()
			if sub isnt ''
				res.Add(sub)
			s = s[pos+1..]
			}
		res.Add(s.Trim())
		return res
		}

	sep: ' #@# '
	formatId(sourceName, name, part)
		{
		return sourceName $ .sep $ String(name) $ .sep $ part
		}

	update(outputs, openAI, vdb)
		{
		if outputs.Empty?()
			return 0
		vectors = #()
		extra = 0
		do
			{
			try
				{
				retry = false
				// Need to instantiate the map because outputs can be modified in
				// .handleLongContent. If the map sequence is iterated/instantiated
				// after the modification, we will get the "object modified during iteration" error
				embeddings = openAI.Embeddings(outputs.Map({ it.text }).Instantiate())
				vectors = embeddings.Map2({ |i, values|
					Object(id: outputs[i].id, :values, metadata: outputs[i].metadata) })
				vdb.Upsert(vectors)
				}
			catch (e)
				{
				if e.Prefix?('OpenAI') and e.Has?('Please reduce your prompt')
					{
					extra += .handleLongContent(outputs, openAI, vdb)
					retry = outputs.NotEmpty?()
					}
				else if e.Lower().Has?("bad gateway")
					continue
				else
					SuneidoLog('ERRATIC: KnowledgeBase.Update - ' $ e,
						params: outputs.Map({ it.id }), calls:)
				}
			}
			while retry is true
		return vectors.Size() + extra
		}

	handleLongContent(outputs, openAI, vdb)
		{
		longest = 0
		for (i = 1; i < outputs.Size(); i++)
			if outputs[i].text.Size() > outputs[longest].text.Size()
				longest = i
		tempOutputs = Object()
		lines = .wrapText(outputs[longest].text,
			(outputs[longest].text.Size() * 2 / 3/*=reduce length*/).Ceiling())
		for i, s in lines
			tempOutputs.Add([
				text: s,
				id: outputs[longest].id $ '.' $ i
				metadata: outputs[longest].metadata])
		c = .update(tempOutputs, openAI, vdb)
		outputs.Delete(longest)
		return c
		}

	Delete(source, key)
		{
		if false is info = .getInfo()
			return 0

		vdb = new VectorDB(info.index_name, info.vectorDB_api)
		prefix = .formatId(source, key, '')

		return vdb.DeleteVector(prefix)
		}

	QueryById(source, key, sources = false, n = 10)
		{
		if false is info = .getInfo()
			return 'not available'

		vdb = new VectorDB(info.index_name, info.vectorDB_api)
		prefix = .formatId(source, key, '')

		vectors = vdb.List(prefix)
		if vectors.Empty?()
			return 'not found'

		results = Object()
		filter = Object(key: Object('$ne': key))
		if sources isnt false
			filter['source'] = Object('$in': sources)
		for vector in vectors
			for result in vdb.QueryById(vector.id, filter, n)
				if not results.Member?(result.metadata.key) or
					results[result.metadata.key].score < result.score
					results[result.metadata.key] = result
		return .formatQueryResults(results.Values().Sort!(By(#score)).Reverse!().Take(n))
		}

	Query(q, sources = false, n = 10)
		{
		if false is info = .getInfo()
			return #()

		vdb = new VectorDB(info.index_name, info.vectorDB_api)
		openAI = new OpenAI_Embeddings(info.openAI_api)
		try
			vector = openAI.Embeddings(q)[0]
		catch (e)
			{
			SuneidoLog('ERRATIC: KnowledgeBase.Query - ' $ e, params: Object(:q), call:)
			return #()
			}
		filter = Object()
		if sources isnt false
			filter['source'] = Object('$in': sources)

		return .formatQueryResults(vdb.Query(vector, filter, n))
		}

	formatQueryResults(results)
		{
		res = Object()
		for result in results
			{
			source = Global(result.metadata.source $ '_KnowledgeBaseSource')
			if false isnt formatted = source.FormatQueryResult(result.metadata.key,
				metadata: result.metadata)
				{
				formatted.score = result.score
				formatted.source = result.metadata.source
				res.Add(formatted)
				}
			}
		return res
		}
	}
