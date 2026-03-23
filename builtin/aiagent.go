// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"log"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/llm"
)

// @immutable
type suAgent struct {
	ValueBase[suAgent]
	agent    *llm.Agent
	th       *Thread
	callback Value
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

	a := &suAgent{
		th:       t2,
		callback: callback,
		agent: llm.NewAgent(baseURL, apiKey, model, prompt,
			outputCallback(t2, callback)),
	}
	return a
}

// @immutable
type suToolApproval struct {
	ValueBase[suToolApproval]
	approval *llm.ToolApproval
}

var _ Value = (*suToolApproval)(nil)

func (ta *suToolApproval) Equal(other any) bool {
	return ta == other
}

func (ta *suToolApproval) Lookup(_ *Thread, method string) Value {
	return toolApprovalMethods[method]
}

func (*suToolApproval) SetConcurrent() {
	// safe: wraps thread-safe llm.ToolApproval
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

func outputCallback(th *Thread, callback Value) func(what, data string, approval *llm.ToolApproval) {
	return func(what, data string, approval *llm.ToolApproval) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("ERROR agent callback: ", err)
				panic(err)
			}
		}()
		if approval == nil {
			th.Call(callback, SuStr(what), SuStr(data))
			return
		}
		ta := &suToolApproval{approval: approval}
		th.Call(callback, SuStr(what), SuStr(data), ta)
	}
}

var agentMethods = methods("agent")
var toolApprovalMethods = methods("toolapproval")

var _ = method(toolapproval_Allow, "(text = '')")

func toolapproval_Allow(this Value, text Value) Value {
	this.(*suToolApproval).approval.Allow(ToStr(text))
	return nil
}

var _ = method(toolapproval_Deny, "(text = '')")

func toolapproval_Deny(this Value, text Value) Value {
	this.(*suToolApproval).approval.Deny(ToStr(text))
	return nil
}

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

var _ = method(agent_LoadConversation, "(filename)")

func agent_LoadConversation(th *Thread, this Value, args []Value) Value {
	a := this.(*suAgent)
	err := a.agent.LoadConversation(ToStr(args[0]), 
		outputCallback(a.th, a.callback))
	if err != nil {
		th.ReturnThrow = true
		return SuStr("LoadConversation: " + err.Error())
	}
	return True
}

var _ = method(agent_Usage, "()")

func agent_Usage(this Value) Value {
	usage := this.(*suAgent).agent.LastUsage()
	if usage == nil {
		return False
	}
	return IntVal(usage.TotalTokens)
}

var _ = method(agent_Close, "()")

func agent_Close(this Value) Value {
	a := this.(*suAgent)
	a.agent.Interrupt()
	a.th.Close()
	DisableSandbox()
	return nil
}
