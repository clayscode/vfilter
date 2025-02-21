package types

import (
	"context"
	"log"
	"runtime"

	"github.com/Velocidex/ordereddict"
)

// A scope is passed inside the evaluation context.  Although this is
// an interface, there is currently only a single implementation
// (scope.Scope). The interface exposes the public methods.
type Scope interface {

	// Duplicate the scope to a completely new scope - this is a
	// deep copy not a subscope!  Very rarely used.
	NewScope() Scope

	// Copy the scope and create a subscope child.
	Copy() Scope

	// The scope context is a global k/v store
	GetContext(name string) (Any, bool)
	SetContext(name string, value Any)

	// Replace the entire context dict.
	SetContextDict(context *ordereddict.Dict)
	ClearContext()

	// Extract debug string about the current scope state.
	PrintVars() string

	// Scope manages the protocols
	Bool(a Any) bool
	Eq(a Any, b Any) bool
	Lt(a Any, b Any) bool
	Gt(a Any, b Any) bool
	Add(a Any, b Any) Any
	Sub(a Any, b Any) Any
	Mul(a Any, b Any) Any
	Div(a Any, b Any) Any
	Membership(a Any, b Any) bool
	Associative(a Any, b Any) (Any, bool)
	GetMembers(a Any) []string
	GetVars() []Row
	Match(a Any, b Any) bool
	Iterate(ctx context.Context, a Any) <-chan Row

	// The scope's top level variable. Scopes search backward
	// through their parents to resolve names from these vars.
	AppendVars(row Row) Scope
	Resolve(field string) (interface{}, bool)

	// Program a custom sorter
	SetSorter(sorter Sorter)
	SetGrouper(grouper Grouper)

	// We can program the scope's protocols
	AddProtocolImpl(implementations ...Any) Scope
	AppendFunctions(functions ...FunctionInterface) Scope
	AppendPlugins(plugins ...PluginGeneratorInterface) Scope

	// Logging and performance monitoring.
	SetLogger(logger *log.Logger)
	SetTracer(logger *log.Logger)
	GetLogger() *log.Logger
	GetStats() *Stats

	Log(format string, a ...interface{})
	Trace(format string, a ...interface{})

	// Introspection
	GetFunction(name string) (FunctionInterface, bool)
	GetPlugin(name string) (PluginGeneratorInterface, bool)
	GetSimilarPlugins(name string) []string
	Describe(type_map *TypeMap) *ScopeInformation
	CheckForOverflow() bool

	// Destructors are called when the scope is Close(). If the
	// scope is already closed adding the destructor may fail.
	AddDestructor(fn func()) error
	Close()
}

// Utilities to do with scope.
func RecoverVQL(scope Scope) {
	r := recover()
	if r != nil {
		scope.Log("PANIC: %v\n", r)
		buffer := make([]byte, 4096)
		n := runtime.Stack(buffer, false /* all */)
		scope.Log("%s", buffer[:n])
	}
}
