package formatter

import (
	"encoding/json"
	"fmt"

	gdiff "github.com/yudai/gojsondiff"
)

const (
	DeltaDelete   = 0
	DeltaTextDiff = 2
	DeltaMove     = 3
)

func NewDeltaFormatter() *DeltaFormatter {
	return &DeltaFormatter{
		PrintIndent: true,
	}
}

type DeltaFormatter struct {
	PrintIndent bool
}

func (f *DeltaFormatter) Format(diff gdiff.Diff) (result string, err error) {
	jsonObject, err := f.formatValue(diff.Delta())
	if err != nil {
		return "", err
	}
	var resultBytes []byte
	if f.PrintIndent {
		resultBytes, err = json.MarshalIndent(jsonObject, "", "  ")
	} else {
		resultBytes, err = json.Marshal(jsonObject)
	}
	if err != nil {
		return "", err
	}

	return string(resultBytes) + "\n", nil
}

func (f *DeltaFormatter) FormatAsJson(diff gdiff.Diff) (deltaJson interface{}, err error) {
	return f.formatValue(diff.Delta())
}

func (f *DeltaFormatter) formatValue(delta gdiff.Delta) (deltaJson interface{}, err error) {
	switch d := delta.(type) {
	case *gdiff.Object:
		return f.formatObject(d)
	case *gdiff.Array:
		return f.formatArray(d)
	case *gdiff.Added:
		return []interface{}{d.Value}, nil
	case *gdiff.Modified:
		return []interface{}{d.OldValue, d.NewValue}, nil
	case *gdiff.TextDiff:
		return []interface{}{d.DiffString(), 0, DeltaTextDiff}, nil
	case *gdiff.Deleted:
		return []interface{}{d.Value, 0, DeltaDelete}, nil
	case *gdiff.Moved:
		return []interface{}{"", d.NewPosition, DeltaMove}, nil
	default:
		return nil, fmt.Errorf("unknown Delta type detected: %#v", delta)
	}
}

func (f *DeltaFormatter) formatObject(delta *gdiff.Object) (deltaJson interface{}, err error) {
	result := map[string]interface{}{}
	for key, d := range delta.Deltas {
		j, err := f.formatValue(d)
		if err != nil {
			return nil, fmt.Errorf("failed to process property `%s`: %s", key, err)
		}
		result[key] = j
	}
	return result, nil
}

func (f *DeltaFormatter) formatArray(delta *gdiff.Array) (deltaJson map[string]interface{}, err error) {
	result := map[string]interface{}{
		"_t": "a",
	}

	for index, d := range delta.PostDeltas {
		key := fmt.Sprintf("%d", index)
		j, err := f.formatValue(d)
		if err != nil {
			return nil, fmt.Errorf("failed to process at `%d`: %s", index, err)
		}
		result[key] = j
	}

	for index, d := range delta.PreDeltas {
		key := fmt.Sprintf("_%d", index)
		j, err := f.formatValue(d)
		if err != nil {
			return nil, fmt.Errorf("failed to process at `%d`: %s", index, err)
		}
		result[key] = j
	}

	return result, nil
}
