// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
OpenAI
	{
	// Each input must not exceed 8192 tokens in length.
	MaxEmbeddingInputSize: 32000 // = 8000 tokens * 4

	// input can be a string or an string of array
	Embeddings(input, model = 'text-embedding-ada-002', user = 'axon')
		{
		if String?(input)
			input = Object(input)
		res = .Send(.BaseUrl $ 'embeddings', Object(:input, :model, :user))
		return res.data.Map({ it.embedding })
		}
	}