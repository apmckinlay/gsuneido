// Copyright (C) 2023 Axon Development Corporation All rights reserved worldwide.
class
	{
	New(indexName, .token)
		{
		.url = 'https://' $ indexName
		}
	/* vectors = [
		{ id: string, values: float[], metadata: {} }
		]
	*/
	Upsert(vectors, namespace = '')
		{
		header = Object('Api-Key': .token,
			'Content-Type': 'application/json')
		content = Json.Encode([:vectors, :namespace])
		res = Json.Decode(Https.Post(.url $ '/vectors/upsert', content, :header))
		return res.upsertedCount
		}

	DeleteVector(prefix, namespace = '')
		{
		header = Object('Api-Key': .token,
			'Content-Type': 'application/json',
			'Accept': 'application/json')
		listResult = Json.Decode(Https.Get(Url.Encode(.url $ '/vectors/list',
			[:namespace, :prefix]), :header))
		ids = listResult.vectors.Map({ it.id })
		if ids.Empty?()
			return 0

		content = Json.Encode([deleteAll: false, :namespace, :ids])
		Https.Post(.url $ '/vectors/delete', content, :header)
		return ids.Size()
		}

	DeleteAll(namespace = '')
		{
		header = Object('Api-Key': .token,
			'Content-Type': 'application/json',
			'Accept': 'application/json')
		content = Json.Encode([deleteAll: true, :namespace])
		Https.Post(.url $ '/vectors/delete', content, :header)
		SuneidoLog('VectorDB DeleteAll', params: [.url], calls:)
		}

	Query(vector, filter = #(), topK = 10, namespace = '')
		{
		return .query(Object(:vector), filter, topK, namespace)
		}

	QueryById(id, filter = #(), topK = 10, namespace = '')
		{
		return .query(Object(:id), filter, topK, namespace)
		}

	query(query, filter = #(), topK = 10, namespace = '')
		{
		header = Object('Api-Key': .token,
			'Content-Type': 'application/json',
			'Accept': 'application/json')
		body = [:namespace, :topK,
			includeValues: false,
			includeMetadata:].Merge(query)
		if filter.NotEmpty?()
			body.filter = filter
		res = Json.Decode(Https.Post(.url $ '/query', Json.Encode(body), :header))
		return res.matches
		}

	List(prefix, namespace = '')
		{
		header = Object('Api-Key': .token,
			'Content-Type': 'application/json',
			'Accept': 'application/json')
		queries = Object(
			:namespace,
			:prefix)
		url = Url.Encode(.url $ '/vectors/list', queries)
		res = Json.Decode(Https.Get(url, :header))
		return res.vectors
		}
	}
