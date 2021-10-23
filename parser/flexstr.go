package parser

import (
	"bytes"
	"strings"
	"unicode"
	"unicode/utf8"
)

// A FlexStr describes an identifier of an Entity (Message, Field, Enum, Service,
// Field). It can be converted to multiple forms using the provided helper
// methods, or a custom transform can be used to modify its behavior.
type FlexStr string

// String satisfies the strings.Stringer interface.
func (n FlexStr) String() string { return string(n) }

// UpperCamelCase converts FlexStr n to upper camelcase, where each part is
// title-cased and concatenated with no separator.
func (n FlexStr) UpperCamelCase() FlexStr { return n.Transform(strings.Title, strings.Title, "") }

// LowerCamelCase converts FlexStr n to lower camelcase, where each part is
// title-cased and concatenated with no separator except the first which is
// lower-cased.
func (n FlexStr) LowerCamelCase() FlexStr { return n.Transform(strings.Title, strings.ToLower, "") }

// ScreamingSnakeCase converts FlexStr n to screaming-snake-case, where each part
// is all-caps and concatenated with underscores.
func (n FlexStr) ScreamingSnakeCase() FlexStr {
	return n.Transform(strings.ToUpper, strings.ToUpper, "_")
}

// LowerSnakeCase converts FlexStr n to lower-snake-case, where each part is
// lower-cased and concatenated with underscores.
func (n FlexStr) LowerSnakeCase() FlexStr { return n.Transform(strings.ToLower, strings.ToLower, "_") }

// LowerDashNotation converts FlexStr n to lower-dash-notation, where each part is
// lower-cased and concatenated with dash.
func (n FlexStr) LowerDashNotation() FlexStr {
	return n.Transform(strings.ToLower, strings.ToLower, "-")
}

// UpperSnakeCase converts FlexStr n to upper-snake-case, where each part is
// title-cased and concatenated with underscores.
func (n FlexStr) UpperSnakeCase() FlexStr { return n.Transform(strings.Title, strings.Title, "_") }

// SnakeCase converts FlexStr n to snake-case, where each part preserves its
// capitalization and concatenated with underscores.
func (n FlexStr) SnakeCase() FlexStr { return n.Transform(ID, ID, "_") }

// LowerDotNotation converts FlexStr n to lower dot notation, where each part is
// lower-cased and concatenated with periods.
func (n FlexStr) LowerDotNotation() FlexStr {
	return n.Transform(strings.ToLower, strings.ToLower, ".")
}

// UpperDotNotation converts FlexStr n to upper dot notation, where each part is
// title-cased and concatenated with periods.
func (n FlexStr) UpperDotNotation() FlexStr { return n.Transform(strings.Title, strings.Title, ".") }

// Split breaks apart FlexStr n into its constituent components. Precedence
// follows dot notation, then underscores (excluding underscore prefixes), then
// camelcase. Numbers are treated as standalone components.
func (n FlexStr) Split() (parts []string) {
	ns := string(n)

	switch {
	case ns == "":
		return []string{""}
	case strings.LastIndex(ns, ".") >= 0:
		return strings.Split(ns, ".")
	case strings.LastIndex(ns, "-") >= 0:
		return strings.Split(ns, "-")
	case strings.LastIndex(ns, "_") > 0: // leading underscore does not count
		parts = strings.Split(ns, "_")
		if parts[0] == "" {
			parts[1] = "_" + parts[1]
			return parts[1:]
		}
		return
	default: // camelCase
		buf := &bytes.Buffer{}
		var capt, lodash, num bool
		for _, r := range ns {
			uc := unicode.IsUpper(r) || unicode.IsTitle(r)
			dg := unicode.IsDigit(r)

			if r == '_' && buf.Len() == 0 && len(parts) == 0 {
				lodash = true
			}

			if uc && !capt && buf.Len() > 0 && !lodash { // new upper letter
				parts = append(parts, buf.String())
				buf.Reset()
			} else if dg && !num && buf.Len() > 0 && !lodash { // new digit
				parts = append(parts, buf.String())
				buf.Reset()
			} else if !uc && capt && buf.Len() > 1 { // upper to lower
				if ss := buf.String(); len(ss) > 1 &&
					(len(ss) != 2 || ss[0] != '_') {
					pr, _ := utf8.DecodeLastRuneInString(ss)
					parts = append(parts, strings.TrimSuffix(ss, string(pr)))
					buf.Reset()
					buf.WriteRune(pr)
				}
			} else if !dg && num && buf.Len() >= 1 {
				parts = append(parts, buf.String())
				buf.Reset()
			}

			num = dg
			capt = uc
			buf.WriteRune(r)
		}
		parts = append(parts, buf.String())
		return
	}
}

// FlexStrTransformer is a function that mutates a string. Many of the methods in
// the standard strings package satisfy this signature.
type FlexStrTransformer func(string) string

// ID is a FlexStrTransformer that does not mutate the string.
func ID(s string) string { return s }

// Chain combines the behavior of two Transformers into one. If multiple
// transformations need to be performed on a FlexStr, this method should be used
// to reduce it to a single transformation before applying.
func (n FlexStrTransformer) Chain(t FlexStrTransformer) FlexStrTransformer {
	return func(s string) string { return t(n(s)) }
}

// Transform applies a transformation to the parts of FlexStr n, returning a new
// FlexStr. Transformer first is applied to the first part, with mod applied to
// all subsequent ones. The parts are then concatenated with the separator sep.
// For optimal efficiency, multiple FlexStrTransformers should be Chained together
// before calling Transform.
func (n FlexStr) Transform(mod, first FlexStrTransformer, sep string) FlexStr {
	parts := n.Split()

	for i, p := range parts {
		if i == 0 {
			parts[i] = first(p)
		} else {
			parts[i] = mod(p)
		}
	}

	return FlexStr(strings.Join(parts, sep))
}
