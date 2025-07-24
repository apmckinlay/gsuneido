// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
// In the future we may need to handle URL's i.e. with scheme and host
// See also Url.Split
class
	{
	New(.request)
		{
		r = .SplitRequestLine(request)
		.method = r.method
		.path = r.path
		.query = r.query
		.headers = Object().Set_default(false)
		}
	BadRequest: 'HttpRequest: Bad Request'
	SplitRequestLine(line) // static - doesn't require instance
		{
		line = line.Tr(' \t', ' ').Trim() // standardize whitespace
		uri = line.AfterFirst(' ').BeforeLast(' ')
		versionStr = line.AfterLast('HTTP/').Trim()
		if not versionStr.Number?() and versionStr isnt ''
			throw .BadRequest
		return Object(
			method: line.BeforeFirst(' ').Upper(),
			path: uri.BeforeFirst('?'),
			query: uri.AfterFirst('?'),
			version: Number(versionStr)
			)
		}
	fields: ('Content-Length', 'Authorization', 'User-Agent')
	Default(field, value = false)
		{
		field = field.Tr('_', '-')
		if not .fields.Has?(field)
			throw "HttpRequest: method not found: " $ field
		if value isnt false
			.headers[field] = value
		return .headers.GetDefault(field, .headers[field.Lower()])
		}
	remote_user: ''
	RemoteUser(value = false)
		{
		if value isnt false
			.remote_user = value
		return .remote_user
		}
	Request()
		{ return .request }
	Path()
		{ return .path }
	Method()
		{ return .method }
	Query()
		{ return .query }
	Getter_Headers()
		{ return .headers }

	QueryValues()
		{
		return Url.SplitQuery(.query)
		}

	body: false
	Body(value = false)
		{
		if value isnt false
			.body = value
		return .body
		}

	ToString() // for debugging
		{
		s = 'HttpRequest\n' $
			'\tRequest: ' $ Display(.request) $ '\n' $
			'\tMethod: ' $ Display(.method) $ '\n' $
			'\tPath: ' $ Display(.path) $ '\n' $
			'\tQuery: ' $ Display(.query) $ '\n' $
			'\tRemoteUser: ' $ Display(.remote_user)
		return s
		}
	}
