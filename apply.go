package gojsondiff

import (
	"errors"
	"fmt"
)

// ApplyPatch applies a Diff to an JSON object.
// This method destruct the original JSON object.
func (differ *Differ) ApplyPatch(json interface{}, patch Diff) (patched interface{}, err error) {
	return differ.applyDelta(patch.Delta(), json)
}

func (differ *Differ) applyDelta(delta Delta, value interface{}) (interface{}, error) {
	switch d := delta.(type) {
	case *Object:
		o, ok := value.(map[string]interface{})
		if !ok {
			return nil, errors.New("Type mismatch")
		}

		return differ.applyObject(d, o)

	case *Array:
		a, ok := value.([]interface{})
		if !ok {
			return errors.New("Type mismatch"), nil
		}

		return differ.applyArray(d, a)

	case *Added:
		return d.Value, nil
	case *Modified:
		return d.NewValue, nil
	case *TextDiff:
		return d.NewValue, nil
	case *Deleted:
		return nil, nil
	default:
		return nil, fmt.Errorf("unknown Delta type detected: %#v", delta)
	}
}

func (differ *Differ) applyObject(delta *Object, value interface{}) (interface{}, error) {
	v, ok := value.(map[string]interface{})
	if !ok {
		return nil, errors.New("Type mismatch")
	}

	var err error
	for key, d := range delta.Deltas {
		switch d := d.(type) {
		case *Object:
			v[key], err = differ.applyObject(d, v[key])
			if err != nil {
				return nil, err
			}

		case *Array:
			v[key], err = differ.applyArray(d, v[key])
			if err != nil {
				return nil, err
			}

		case *Added:
			v[key] = d.Value
		case *Modified:
			v[key] = d.NewValue
		case *TextDiff:
			v[key] = d.NewValue
		case *Deleted:
			delete(v, key)
		default:
			return nil, fmt.Errorf("unsupported Delta type detected for an object: %#v", delta)
		}
	}

	return v, nil
}

func (differ *Differ) applyArray(delta *Array, value interface{}) (interface{}, error) {
	v, ok := value.([]interface{})
	if !ok {
		return errors.New("Type mismatch"), nil
	}

	pre, post := delta.WithoutMoved()

	var err error
	keys := sortedKeysInt(pre)
	for i := len(keys) - 1; i >= 0; i-- {
		key := keys[i]
		switch pre[key].(type) {
		case *Deleted:
			v = append(v[:key], v[key+1:]...)
		default:
			return nil, fmt.Errorf("unsupported pre-indexed Delta type detected for a array: %#v", delta)
		}
	}

	keys = sortedKeysInt(post)
	for _, key := range keys {
		switch d := post[key].(type) {
		case *Object:
			v[key], err = differ.applyObject(d, v[key])
			if err != nil {
				return nil, err
			}

		case *Array:
			v[key], err = differ.applyArray(d, v[key])
			if err != nil {
				return nil, err
			}

		case *Added:
			v = append(
				v[:key],
				append(
					[]interface{}{d.Value}, v[key:]...,
				)...,
			)
			v[key] = d.Value
		case *Modified:
			v[key] = d.NewValue
		case *TextDiff:
			v[key] = d.NewValue

		default:
			return nil, fmt.Errorf("unsupported post-indexed Delta type detected for a array: %#v", delta)
		}
	}

	return v, nil
}
