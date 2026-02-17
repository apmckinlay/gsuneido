// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"log"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/llm"
	"github.com/apmckinlay/gsuneido/mcp"
)

// @immutable
type suAgent struct {
	ValueBase[suAgent]
	agent     *llm.Agent
	mcpClient *llm.MCPClient
	th        *Thread
}

var _ = builtin(AiAgent, "(baseURL, apiKey, model, callback, prompt = '')")

func AiAgent(th *Thread, args []Value) Value {
	EnableSandbox()
	baseURL := ToStr(args[0])
	apiKey := ToStr(args[1])
	model := ToStr(args[2])
	callback := args[3]
	prompt := ToStr(args[4])
	callback.SetConcurrent()
	t2 := NewThread(th)
	th.Call(callback, SuStr(model), EmptyStr)

	mcpClient, err := llm.NewMCPClient(mcp.Server())
	if err != nil {
		log.Println("ERROR creating MCP client: ", err)
		panic(err)
	}

	return &suAgent{
		th:        t2,
		mcpClient: mcpClient,
		agent: llm.NewAgent(baseURL, apiKey, model, prompt, mcpClient,
			func(what, data string) {
				defer func() {
					if err := recover(); err != nil {
						log.Println("ERROR agent callback: ", err)
						panic(err)
					}
				}()
				t2.Call(callback, SuStr(what), SuStr(data))
			}),
	}
}

var _ Value = (*suAgent)(nil)

func (a *suAgent) Equal(other any) bool {
	return a == other
}

func (a *suAgent) Lookup(_ *Thread, method string) Value {
	return agentMethods[method]
}

func (a *suAgent) SetConcurrent() {
	// ok since immutable
}

var agentMethods = methods("agent")

var _ = method(agent_Input, "(input)")

func agent_Input(this, input Value) Value {
	this.(*suAgent).agent.Input(ToStr(input))
	return nil
}

var _ = method(agent_Interrupt, "()")

func agent_Interrupt(this Value) Value {
	this.(*suAgent).agent.Interrupt()
	return nil
}

var _ = method(agent_SetModel, "(model)")

func agent_SetModel(this, model Value) Value {
	this.(*suAgent).agent.SetModel(ToStr(model))
	return nil
}

var _ = method(agent_ClearHistory, "()")

func agent_ClearHistory(this Value) Value {
	this.(*suAgent).agent.ClearHistory()
	return nil
}

var _ = method(agent_Close, "()")

func agent_Close(this Value) Value {
	a := this.(*suAgent)
	a.agent.Interrupt()
	if a.mcpClient != nil {
		a.mcpClient.Close()
	}
	a.th.Close()
	DisableSandbox()
	return nil
}
