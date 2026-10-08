package ecmascript

import (
	"context"
	"strings"
	"unicode/utf16"
)

// FindAll returns ordered non-overlapping matches. Empty matches advance by
// AdvanceStringIndex semantics, using code points in u or v mode.
func (p *Program) FindAll(ctx context.Context, input string, options MatchOptions) ([]Result, error) {
	view, err := makeInputViewContext(ctx, input, options.Limits)
	if err != nil {
		return nil, err
	}
	return p.findAll(ctx, view, options)
}

// FindAllUTF16 returns ordered non-overlapping matches in an exact ECMAScript
// string, including inputs containing lone surrogates.
func (p *Program) FindAllUTF16(ctx context.Context, input UTF16String, options MatchOptions) ([]Result, error) {
	view, err := makeUTF16InputViewContext(ctx, input, options.Limits)
	if err != nil {
		return nil, err
	}
	return p.findAll(ctx, view, options)
}

func (p *Program) findAll(ctx context.Context, view *inputView, options MatchOptions) ([]Result, error) {
	executor := newExecutor(ctx, p, view, options.Limits)
	results := make([]Result, 0)
	search := options.StartUTF16
	for search <= len(view.units) {
		result, matched, err := executor.find(search, p.flags.Sticky())
		if err != nil {
			return nil, err
		}
		if !matched {
			break
		}
		if uint64(len(results)+1) > options.Limits.Results {
			return nil, &LimitError{Kind: LimitResults, Limit: options.Limits.Results, Used: uint64(len(results) + 1)}
		}
		if err := executor.allocate(1); err != nil {
			return nil, err
		}
		results = append(results, result)
		start := result.Full().span.Start.UTF16
		end := result.Full().span.End.UTF16
		if start == end {
			if end == len(view.units) {
				break
			}
			search = advanceStringIndex(view, end, p.flags.unicodeMode())
		} else {
			search = end
		}
	}

	return results, nil
}

// Replace applies ECMAScript GetSubstitution tokens. The g flag selects all
// matches; without g, only the first match is replaced.
func (p *Program) Replace(ctx context.Context, input string, replacement UTF16String, options MatchOptions) (UTF16String, error) {
	view, err := makeInputViewContext(ctx, input, options.Limits)
	if err != nil {
		return UTF16String{}, err
	}
	return p.replace(ctx, view, replacement, options)
}

// ReplaceUTF16 applies ECMAScript substitution semantics to an exact
// ECMAScript string, including inputs containing lone surrogates.
func (p *Program) ReplaceUTF16(ctx context.Context, input, replacement UTF16String, options MatchOptions) (UTF16String, error) {
	view, err := makeUTF16InputViewContext(ctx, input, options.Limits)
	if err != nil {
		return UTF16String{}, err
	}
	return p.replace(ctx, view, replacement, options)
}

func (p *Program) replace(ctx context.Context, view *inputView, replacement UTF16String, options MatchOptions) (UTF16String, error) {
	if err := admitReplacement(ctx, replacement.units, options.Limits); err != nil {
		return UTF16String{}, err
	}
	executor := newExecutor(ctx, p, view, options.Limits)
	output := outputUnits{limit: options.Limits.OutputUTF16}
	search := options.StartUTF16
	lastEnd := 0
	replacements := uint64(0)
	for search <= len(view.units) {
		result, matched, err := executor.find(search, p.flags.Sticky())
		if err != nil {
			return UTF16String{}, err
		}
		if !matched {
			break
		}
		replacements++
		if replacements > options.Limits.Results {
			return UTF16String{}, &LimitError{Kind: LimitResults, Limit: options.Limits.Results, Used: replacements}
		}
		start := result.Full().span.Start.UTF16
		end := result.Full().span.End.UTF16
		if err := output.append(view.units[lastEnd:start]); err != nil {
			return UTF16String{}, err
		}
		if err := p.appendSubstitution(&output, view.units, replacement.units, result, executor); err != nil {
			return UTF16String{}, err
		}
		lastEnd = end
		if !p.flags.Global() {
			break
		}
		if start == end {
			if end == len(view.units) {
				break
			}
			search = advanceStringIndex(view, end, p.flags.unicodeMode())
		} else {
			search = end
		}
	}
	if err := output.append(view.units[lastEnd:]); err != nil {
		return UTF16String{}, err
	}

	return newUTF16String(output.units), nil
}

// Admit the template independently of the subject, before scanning tokens or
// allocating a capture name. Output limits cannot bound a consumed empty name.
func admitReplacement(ctx context.Context, units []uint16, limits MatchLimits) error {
	ctx = normalizeContext(ctx)
	if err := ctx.Err(); err != nil {
		return err
	}
	storageBytes := uint64(len(units)) * 2
	if storageBytes > limits.InputBytes {
		return &LimitError{Kind: LimitInputBytes, Limit: limits.InputBytes, Used: storageBytes}
	}
	codePoints := uint64(0)
	for index := 0; index < len(units); {
		if inputContextCheckpoint(codePoints) {
			if err := ctx.Err(); err != nil {
				return err
			}
		}
		width := 1
		if isHighSurrogate(units[index]) && index+1 < len(units) && isLowSurrogate(units[index+1]) {
			width = 2
		}
		index += width
		codePoints++
	}
	if codePoints > limits.InputRunes {
		return &LimitError{Kind: LimitInputRunes, Limit: limits.InputRunes, Used: codePoints}
	}
	return ctx.Err()
}

func (p *Program) appendSubstitution(output *outputUnits, input, replacement []uint16, result Result, executor *executor) error {
	matchStart := result.Full().span.Start.UTF16
	matchEnd := result.Full().span.End.UTF16
	for index := 0; index < len(replacement); index++ {
		if err := executor.step(); err != nil {
			return err
		}
		if replacement[index] != '$' || index+1 >= len(replacement) {
			if err := output.append(replacement[index : index+1]); err != nil {
				return err
			}
			continue
		}
		next := replacement[index+1]
		if err := executor.step(); err != nil {
			return err
		}
		switch next {
		case '$':
			if err := output.append([]uint16{'$'}); err != nil {
				return err
			}
			index++
		case '&':
			if err := output.append(result.Full().value.units); err != nil {
				return err
			}
			index++
		case '`':
			if err := output.append(input[:matchStart]); err != nil {
				return err
			}
			index++
		case '\'':
			if err := output.append(input[matchEnd:]); err != nil {
				return err
			}
			index++
		case '<':
			end := index + 2
			for end < len(replacement) {
				if err := executor.step(); err != nil {
					return err
				}
				if replacement[end] == '>' {
					break
				}
				end++
			}
			if end < len(replacement) && len(p.captureNames) > 0 {
				name, err := utf16CaptureName(replacement[index+2:end], executor)
				if err != nil {
					return err
				}
				capture, _ := result.Named(name)
				if capture.participated {
					if err := output.append(capture.value.units); err != nil {
						return err
					}
				}
				index = end
			} else {
				if err := output.append([]uint16{'$'}); err != nil {
					return err
				}
			}
		default:
			if captureIndex, ok := decimalDigit(next); ok {
				consumed := 1
				if index+2 < len(replacement) {
					if err := executor.step(); err != nil {
						return err
					}
					secondDigit, secondOK := decimalDigit(replacement[index+2])
					twoDigits := captureIndex*10 + secondDigit
					if secondOK {
						if twoDigits >= 1 && twoDigits <= p.captures {
							captureIndex = twoDigits
							consumed = 2
						}
					}
				}
				if captureIndex >= 1 && captureIndex <= p.captures {
					capture := result.captures[captureIndex]
					if capture.participated {
						if err := output.append(capture.value.units); err != nil {
							return err
						}
					}
					index += consumed
				} else {
					if err := output.append([]uint16{'$'}); err != nil {
						return err
					}
				}
			} else {
				if err := output.append([]uint16{'$'}); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// SplitValue represents either a defined string value or ECMAScript
// undefined from an unmatched separator capture.
type SplitValue struct {
	defined bool
	value   UTF16String
}

func (v SplitValue) Defined() bool      { return v.defined }
func (v SplitValue) Value() UTF16String { return newUTF16String(v.value.units) }

// Split separates input and inserts separator captures in result order.
func (p *Program) Split(ctx context.Context, input string, options MatchOptions) ([]SplitValue, error) {
	view, err := makeInputViewContext(ctx, input, options.Limits)
	if err != nil {
		return nil, err
	}
	return p.split(ctx, view, options)
}

// SplitUTF16 separates an exact ECMAScript string and preserves lone
// surrogates in both values and captures.
func (p *Program) SplitUTF16(ctx context.Context, input UTF16String, options MatchOptions) ([]SplitValue, error) {
	view, err := makeUTF16InputViewContext(ctx, input, options.Limits)
	if err != nil {
		return nil, err
	}
	return p.split(ctx, view, options)
}

func (p *Program) split(ctx context.Context, view *inputView, options MatchOptions) ([]SplitValue, error) {
	executor := newExecutor(ctx, p, view, options.Limits)
	parts := make([]SplitValue, 0)
	totalUnits := uint64(0)
	if len(view.units) == 0 {
		result, matched, err := executor.at(0)
		if err != nil {
			return nil, err
		}
		if matched && result.Full().span.Start.UTF16 == result.Full().span.End.UTF16 {
			return parts, nil
		}
		if err := appendSplitValue(&parts, UTF16String{}, true, &totalUnits, options.Limits); err != nil {
			return nil, err
		}
		return parts, nil
	}
	lastEnd := 0
	search := 0
	for search < len(view.units) {
		result, matched, err := executor.at(search)
		if err != nil {
			return nil, err
		}
		if !matched {
			search = advanceStringIndex(view, search, p.flags.unicodeMode())
			continue
		}
		start := result.Full().span.Start.UTF16
		end := result.Full().span.End.UTF16
		if start == end && start == lastEnd {
			search = advanceStringIndex(view, end, p.flags.unicodeMode())
			continue
		}
		if err := appendSplitValue(&parts, newUTF16String(view.units[lastEnd:start]), true, &totalUnits, options.Limits); err != nil {
			return nil, err
		}
		for captureIndex := 1; captureIndex < len(result.captures); captureIndex++ {
			capture := result.captures[captureIndex]
			if err := appendSplitValue(&parts, capture.value, capture.participated, &totalUnits, options.Limits); err != nil {
				return nil, err
			}
		}
		lastEnd = end
		search = nextSplitSearch(view, start, end, p.flags.unicodeMode())
	}
	if err := appendSplitValue(&parts, newUTF16String(view.units[lastEnd:]), true, &totalUnits, options.Limits); err != nil {
		return nil, err
	}

	return parts, nil
}

func nextSplitSearch(view *inputView, start, end int, unicodeMode bool) int {
	if start != end {
		return end
	}
	return advanceStringIndex(view, end, unicodeMode)
}

func decimalDigit(unit uint16) (int, bool) {
	switch unit {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return int(unit - '0'), true
	default:
		return 0, false
	}
}

func appendSplitValue(parts *[]SplitValue, value UTF16String, defined bool, total *uint64, limits MatchLimits) error {
	usedResults := uint64(len(*parts) + 1)
	if usedResults > limits.Results {
		return &LimitError{Kind: LimitResults, Limit: limits.Results, Used: usedResults}
	}
	usedOutput := *total + uint64(len(value.units))
	if usedOutput > limits.OutputUTF16 {
		return &LimitError{Kind: LimitOutputUTF16, Limit: limits.OutputUTF16, Used: usedOutput}
	}
	*parts = append(*parts, SplitValue{defined: defined, value: newUTF16String(value.units)})
	*total = usedOutput
	return nil
}

type outputUnits struct {
	units []uint16
	limit uint64
}

func (o *outputUnits) append(units []uint16) error {
	used := uint64(len(o.units) + len(units))
	if used > o.limit {
		return &LimitError{Kind: LimitOutputUTF16, Limit: o.limit, Used: used}
	}
	o.units = append(o.units, units...)
	return nil
}

func advanceStringIndex(view *inputView, index int, unicodeMode bool) int {
	if !unicodeMode || index+1 >= len(view.units) || !isHighSurrogate(view.units[index]) || !isLowSurrogate(view.units[index+1]) {
		return index + 1
	}
	return index + 2
}

func utf16CaptureName(units []uint16, executor *executor) (string, error) {
	var result strings.Builder
	for index := 0; index < len(units); index++ {
		if err := executor.step(); err != nil {
			return "", err
		}
		unit := units[index]
		if isHighSurrogate(unit) && index+1 < len(units) && isLowSurrogate(units[index+1]) {
			if err := executor.step(); err != nil {
				return "", err
			}
			result.WriteRune(utf16.DecodeRune(rune(unit), rune(units[index+1])))
			index++
		} else {
			// An unpaired surrogate cannot be a valid parsed capture name.
			// WriteRune maps it to U+FFFD rather than aliasing an ASCII name.
			result.WriteRune(rune(unit))
		}
	}
	return result.String(), nil
}
