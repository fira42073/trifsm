package generator

type templateData struct {
	Command string
	Version string
	Package string
	FSMs    []fsmData
}

type sortedFSMs []fsmData

func (s sortedFSMs) Len() int           { return len(s) }
func (s sortedFSMs) Less(i, j int) bool { return s[i].CustomTypeName < s[j].CustomTypeName }
func (s sortedFSMs) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

type fsmData struct {
	CustomTypeName      string
	ConstructorName     string
	InitialState        string
	TransitionEvents    string
	TransitionEventDefs []transitionEvent
	States              []stateDecl
}

type stateDecl struct {
	Name    string
	Type    string
	Value   string
	Comment string
}

// uniqueStateDecls is a helper struct to ensure that state declarations are unique
type uniqueStateDecls struct {
	// cache for faster lookup
	namesIndex map[string]int
	decls      []stateDecl
}

func newUniqueStateDecls(capacity int) *uniqueStateDecls {
	return &uniqueStateDecls{
		decls:      make([]stateDecl, 0, capacity),
		namesIndex: make(map[string]int),
	}
}

func (s *uniqueStateDecls) add(decl stateDecl) {
	if s == nil || s.decls == nil {
		panic("nil stateDecls is unusable, please use the constructor")
	}

	if _, ok := s.namesIndex[decl.Name]; ok {
		// already exists
		return
	}

	s.decls = append(s.decls, decl)
	s.namesIndex[decl.Name] = len(s.decls) - 1
}

func (s *uniqueStateDecls) patchOne(name string, fn func(decl *stateDecl)) { //nolint:unused // TODO: will be used for comments
	if s == nil || s.decls == nil {
		panic("nil stateDecls is unusable, please use a constructor")
	}

	if idx, ok := s.namesIndex[name]; ok {
		fn(&s.decls[idx])
	}
}

func (s *uniqueStateDecls) collect() []stateDecl {
	return s.decls
}

type transitionEvent struct {
	Name    string
	Type    string
	Value   string
	Comment string
}
