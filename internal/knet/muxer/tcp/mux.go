package tcp

import (
	"fmt"

	"sort"

	"github.com/cloudwego/kitex/pkg/discovery"
	"github.com/llsw/ikunet/internal/kitex_gen/transport"
	"github.com/rs/zerolog/log"
	"github.com/traefik/traefik/v3/pkg/rules"
	"github.com/vulcand/predicate"
)

type Handler interface {
	Serve()
}

// Data contains TCP connection metadata.
type Data struct {
	Req      *transport.Transport
	Instance discovery.Instance
}

// Muxer defines a muxer that handles TCP routing with rules.
type Muxer struct {
	routes routes
	parser predicate.Parser
}

// NewMuxer returns a TCP muxer.
func NewMuxer() (*Muxer, error) {
	var matcherNames []string
	for matcherName := range tcpFuncs {
		matcherNames = append(matcherNames, matcherName)
	}

	parser, err := rules.NewParser(matcherNames)
	if err != nil {
		return nil, fmt.Errorf("error while creating rules parser: %w", err)
	}

	return &Muxer{
		parser: parser,
	}, nil
}

// Match returns the handler of the first route matching the connection metadata,
// and whether the match is exactly from the rule HostSNI(*).
func (m Muxer) Match(meta *Data) bool {
	for _, route := range m.routes {
		if route.matchers.match(meta) {
			return true
		}
	}
	return false
}

// AddRoute adds a new route, associated to the given handler, at the given
// priority, to the muxer.
func (m *Muxer) AddRoute(rule string) error {
	var parse interface{}
	var err error
	var matcherFuncs map[string]func(*matchersTree, ...string) error

	parse, err = m.parser.Parse(rule)
	if err != nil {
		return fmt.Errorf("error while parsing rule %s: %w", rule, err)
	}

	matcherFuncs = tcpFuncs

	buildTree, ok := parse.(rules.TreeBuilder)
	if !ok {
		return fmt.Errorf("error while parsing rule %s", rule)
	}

	ruleTree := buildTree()

	var matchers matchersTree
	err = matchers.addRule(ruleTree, matcherFuncs)
	if err != nil {
		return fmt.Errorf("error while adding rule %s: %w", rule, err)
	}

	// if ruleTree.RuleLeft == nil && ruleTree.RuleRight == nil && len(ruleTree.Value) == 1 {
	// 	catchAll = ruleTree.Value[0] == "*" && strings.EqualFold(ruleTree.Matcher, "HostSNI")
	// }

	newRoute := &route{
		matchers: matchers,
	}
	m.routes = append(m.routes, newRoute)
	sort.Sort(m.routes)
	return nil
}

// HasRoutes returns whether the muxer has routes.
func (m *Muxer) HasRoutes() bool {
	return len(m.routes) > 0
}

// routes implements sort.Interface.
type routes []*route

// Len implements sort.Interface.
func (r routes) Len() int { return len(r) }

// Swap implements sort.Interface.
func (r routes) Swap(i, j int) { r[i], r[j] = r[j], r[i] }

// Less implements sort.Interface.
func (r routes) Less(i, j int) bool { return false }

// route holds the matchers to match TCP route,
// and the handler that will serve the connection.
type route struct {
	// matchers tree structure reflecting the rule.
	matchers matchersTree
	// handler responsible for handling the route.
	handler Handler
}

// matchersTree represents the matchers tree structure.
type matchersTree struct {
	// matcher is a matcher func used to match connection properties.
	// If matcher is not nil, it means that this matcherTree is a leaf of the tree.
	// It is therefore mutually exclusive with left and right.
	matcher func(*Data) bool
	// operator to combine the evaluation of left and right leaves.
	operator string
	// Mutually exclusive with matcher.
	left  *matchersTree
	right *matchersTree
}

func (m *matchersTree) match(meta *Data) bool {
	if m == nil {
		// This should never happen as it should have been detected during parsing.
		log.Warn().Msg("Rule matcher is nil")
		return false
	}

	if m.matcher != nil {
		return m.matcher(meta)
	}

	switch m.operator {
	case "or":
		return m.left.match(meta) || m.right.match(meta)
	case "and":
		return m.left.match(meta) && m.right.match(meta)
	default:
		// This should never happen as it should have been detected during parsing.
		log.Warn().Str("operator", m.operator).Msg("Invalid rule operator")
		return false
	}
}

type matcherFuncs map[string]func(*matchersTree, ...string) error

func (m *matchersTree) addRule(rule *rules.Tree, funcs matcherFuncs) error {
	switch rule.Matcher {
	case "and", "or":
		m.operator = rule.Matcher
		m.left = &matchersTree{}
		err := m.left.addRule(rule.RuleLeft, funcs)
		if err != nil {
			return err
		}

		m.right = &matchersTree{}
		return m.right.addRule(rule.RuleRight, funcs)
	default:
		err := rules.CheckRule(rule)
		if err != nil {
			return err
		}

		err = funcs[rule.Matcher](m, rule.Value...)
		if err != nil {
			return err
		}

		if rule.Not {
			matcherFunc := m.matcher
			m.matcher = func(meta *Data) bool {
				return !matcherFunc(meta)
			}
		}
	}

	return nil
}
