package scarlet

import (
	"sort"
	"strconv"
	"strings"

	"github.com/OneOfOne/xxhash"
)

// Force interface implementation
var _ Renderable = (*CSSRule)(nil)

// CSSRule ...
type CSSRule struct {
	Selector
	Statements []*CSSStatement
	Duplicates []*CSSRule
	Parent     *CSSRule
}

// Render renders the CSS rule into the output stream.
func (rule *CSSRule) Render(output *strings.Builder, pretty bool) {
	if len(rule.Statements) == 0 {
		return
	}

	output.WriteString(strings.TrimSpace(rule.SelectorPath(pretty)))

	if len(rule.Duplicates) > 0 {
		for _, duplicate := range rule.Duplicates {
			output.WriteString(",")

			if pretty {
				output.WriteString(" ")
			}

			output.WriteString(strings.TrimSpace(duplicate.SelectorPath(pretty)))
		}
	}

	if pretty {
		output.WriteString(" ")
	}

	output.WriteString("{")

	if pretty {
		output.WriteString("\n")
	}

	for index, statement := range rule.Statements {
		if pretty {
			output.WriteString("\t")
		}

		output.WriteString(statement.Property)
		output.WriteString(":")

		if pretty {
			output.WriteString(" ")
		}

		output.WriteString(statement.Value)

		// Remove last semicolon
		if pretty || index != len(rule.Statements)-1 {
			output.WriteString(";")
		}

		if pretty {
			output.WriteString("\n")
		}
	}

	output.WriteString("}")

	if pretty {
		output.WriteString("\n\n")
	}
}

// Root ...
func (rule *CSSRule) Root() *CSSRule {
	parent := rule

	for {
		nextParent := parent.Parent

		if nextParent == nil {
			return parent
		}

		parent = nextParent
	}
}

// Copy ...
func (rule *CSSRule) Copy() *CSSRule {
	return &CSSRule{
		Selector:   rule.Selector,
		Statements: rule.Statements,
		Parent:     rule.Parent,
	}
}

// SelectorPath returns the selector string for the rule (recursive, returns absolute path).
func (rule *CSSRule) SelectorPath(pretty bool) string {
	return rule.Selector.Render()
}

// StatementsHash returns a hash of all the statements which is used to find duplicate CSS rules.
func (rule *CSSRule) StatementsHash() string {
	sort.Slice(rule.Statements, func(i, j int) bool {
		return rule.Statements[i].Property < rule.Statements[j].Property
	})

	hash := xxhash.NewS64(0)

	for _, statement := range rule.Statements {
		_, _ = hash.WriteString(statement.Property)
		_, _ = hash.WriteString(statement.Value)
	}

	return strconv.FormatUint(hash.Sum64(), 16)
}
