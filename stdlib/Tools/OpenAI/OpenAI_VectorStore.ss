// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
OpenAI
	{
	Create(name = "", expires_after = false, file_ids = #())
		{
		data = Object()

		if name isnt ""
			data.name = name

		if expires_after
			data.expires_after = expires_after

		if file_ids.Size() > 0
			data.file_ids = file_ids

		response = .Send(.BaseUrl $ "vector_stores", data)
		id = response.id
		return Object(:id, :response)
		}

	List(after = false, before = false)
		{
		params = Object()

		if after isnt false
			params.after = after
		if before isnt false
			params.before = before

		url = Url.Encode(.BaseUrl $ "vector_stores", params)
		return .Send(url, body: false, method: "GET")
		}

	DeleteStore(vectorStoreId)
		{
		return .Send(.BaseUrl $ "vector_stores/" $ vectorStoreId, body: false,
			method: "DELETE")
		}

	AddFiles(vectorStoreId, file_ids, attributes = false)
		{
		data = Object("file_ids": file_ids)
		if attributes isnt false
			data.attributes = attributes
		return .Send(.BaseUrl $ "vector_stores/" $ vectorStoreId $ "/files", data)
		}

	RemoveFile(vectorStoreId, fileId)
		{
		return .Send(.BaseUrl $ "vector_stores/" $ vectorStoreId $ "/files/" $ fileId,
			body: false, method: "DELETE")
		}

	ListFiles(vectorStoreId, after = false, before = false, filter = false)
		{
		params = Object()

		if after isnt false
			params.after = after
		if before isnt false
			params.before = before
		if filter isnt false
			params.filter = filter

		url = Url.Encode(.BaseUrl $ vectorStoreId $ "/files", params)
		return .Send(url, body: false, method: "GET")
		}
	}