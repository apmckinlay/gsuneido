// Copyright (C) 2013 Axon Development Corporation All rights reserved worldwide.
// using EmailForwarder
function (relay, from, to, message, sendFile? = false) // returns true or error message
	{
	authentication = SoleContribution('SESAuthentication')()
	header = Object(
		X_Suneido_From: from,
		X_Suneido_To: to,
		Authorization: authentication.Authorization)

	//TODO use Https.Post and let it check response code
	args = Object('POST', 'https://' $ relay $ '/email', :header)
	msgMember = sendFile? ? 'fromFile' : 'content'
	args[msgMember] = message
	response = Https(@args)
	response_code = Http.ResponseCode(response.header)
	if response_code is '200' and response.content is ''
		return true

	if response_code =~ `^4\d\d$` and response_code isnt '403'
		SuneidoLog("ERROR: ForwardSendMessage - Non-Successful HTTP response",
			params: Object(:response, :header, :message, :from, :to), calls:)
	return response_code $ ' ' $ response.content
	}
