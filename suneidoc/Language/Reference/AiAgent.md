<div style="float:right"><span class="builtin">Builtin</span></div>

### AiAgent

``` suneido
(baseURL, apiKey, model, callback, prompt = '') => instance
```

AiAgent creates an AI agent that communicates with a language model API. The agent runs asynchronously, calling back to Suneido code for handling responses and tool requests.

AiAgent uses the Model Context Protocol (MCP) to provide tools to the AI.

**Parameters:**

`baseURL`
: The base URL for the AI API endpoint (e.g., "https://api.openai.com/v1")

`apiKey`
: The API key for authentication

`model`
: The model identifier to use (e.g., "gpt-4")

`callback`
: A function called with `(what, data)` for agent events. The `what` parameter indicates the event type, and `data` contains the event data.

`prompt`
: Optional system prompt to initialize the agent's context

**Methods:**

`Input(input)`
: Sends input text to the agent for processing. The agent will communicate with the AI and call the callback as needed.

`Interrupt()`
: Interrupts any ongoing agent operation. Useful for stopping long-running requests.

`SetModel(model)`
: Changes the model being used by the agent.

`ClearHistory()`
: Clears the conversation history, starting a fresh session.

`LoadConversation(file)`
: Reloads a logged conversation, ready to resume. 
It calls the callback with the contents of the conversation.

`Close()`
: Closes the agent, releasing resources. This interrupts any ongoing operations, closes the MCP client, and cleans up the agent thread. Always call Close when done with an agent.

**Example:**

```suneido
agent = AiAgent(
    "https://api.openai.com/v1",
    "your-api-key",
    "gpt-4",
    function (what, data)
        {
        Print(what, ": ", data)
        },
    "You are a helpful assistant."
    )

agent.Input("What is 2 + 2?")
// ... handle callbacks ...

agent.Close()
```

**Note:** AiAgent enables sandbox mode. The agent runs in a separate thread and communicates asynchronously via callbacks.