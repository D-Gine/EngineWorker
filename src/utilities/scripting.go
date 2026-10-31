package utilities

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkjson"
)

func GoToStarlark(v any) (starlark.Value, error) {
	if v == nil {
		return starlark.None, nil
	}

	switch x := v.(type) {
	case string:
		return starlark.String(x), nil
	case bool:
		return starlark.Bool(x), nil
	case int:
		return starlark.MakeInt(x), nil
	case int8:
		return starlark.MakeInt64(int64(x)), nil
	case int16:
		return starlark.MakeInt64(int64(x)), nil
	case int32:
		return starlark.MakeInt64(int64(x)), nil
	case int64:
		return starlark.MakeInt64(x), nil
	case uint:
		return starlark.MakeUint(uint(x)), nil
	case uint8:
		return starlark.MakeUint64(uint64(x)), nil
	case uint16:
		return starlark.MakeUint64(uint64(x)), nil
	case uint32:
		return starlark.MakeUint64(uint64(x)), nil
	case uint64:
		return starlark.MakeUint64(x), nil
	case float32:
		return starlark.Float(x), nil
	case float64:
		return starlark.Float(x), nil
	case map[string]any:
		d := starlark.NewDict(len(x))
		for k, vv := range x {
			sv, err := GoToStarlark(vv)
			if err != nil {
				return nil, err
			}
			if err := d.SetKey(starlark.String(k), sv); err != nil {
				return nil, err
			}
		}
		return d, nil
	case []any:
		items := make([]starlark.Value, 0, len(x))
		for _, vv := range x {
			sv, err := GoToStarlark(vv)
			if err != nil {
				return nil, err
			}
			items = append(items, sv)
		}
		return starlark.NewList(items), nil
	case map[string]json.RawMessage:
		jsonBytes, err := json.Marshal(x)
		if err != nil {
			return nil, fmt.Errorf("Error marshaling map: %s", err)
		}
		jsonString := string(jsonBytes)

		thread := &starlark.Thread{Name: "main"}
		decodeFn, ok := starlarkjson.Module.Members["decode"].(starlark.Callable)
		if !ok {
			return nil, errors.New("failed to find starlark json decode function")
		}
		starlarkVal, err := starlark.Call(thread, decodeFn, starlark.Tuple{starlark.String(jsonString)}, nil)
		return starlarkVal, err
	case json.RawMessage:
		jsonString := x.String()
		thread := &starlark.Thread{Name: "main"}
		decodeFn, ok := starlarkjson.Module.Members["decode"].(starlark.Callable)
		if !ok {
			return nil, errors.New("failed to find starlark json decode function")
		}
		starlarkVal, err := starlark.Call(thread, decodeFn, starlark.Tuple{starlark.String(jsonString)}, nil)
		return starlarkVal, err

	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Map:
		if rv.Type().Key().Kind() != reflect.String {
			return nil, fmt.Errorf("map keys must be strings, got %s", rv.Type().Key().String())
		}
		d := starlark.NewDict(rv.Len())
		iter := rv.MapRange()
		for iter.Next() {
			k := iter.Key().String()
			sv, err := GoToStarlark(iter.Value().Interface())
			if err != nil {
				return nil, err
			}
			if err := d.SetKey(starlark.String(k), sv); err != nil {
				return nil, err
			}
		}
		return d, nil
	case reflect.Slice, reflect.Array:
		items := make([]starlark.Value, 0, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			sv, err := GoToStarlark(rv.Index(i).Interface())
			if err != nil {
				return nil, err
			}
			items = append(items, sv)
		}
		return starlark.NewList(items), nil
	}

	return nil, fmt.Errorf("unsupported type: %T", v)
}

func StarlarkToGo(v starlark.Value) (any, error) {
	if v == starlark.None {
		return nil, nil
	}
	if s, ok := starlark.AsString(v); ok {
		return s, nil
	}
	switch x := v.(type) {
	case starlark.Bool:
		return bool(x), nil
	case starlark.Int:
		i, ok := x.Int64()
		if !ok {
			return nil, fmt.Errorf("integer out of int64 range")
		}
		return i, nil
	case starlark.Float:
		return float64(x), nil
	case *starlark.List:
		out := make([]any, 0, x.Len())
		for i := 0; i < x.Len(); i++ {
			item, err := StarlarkToGo(x.Index(i))
			if err != nil {
				return nil, err
			}
			out = append(out, item)
		}
		return out, nil
	case starlark.Tuple:
		out := make([]any, 0, len(x))
		for _, item := range x {
			gv, err := StarlarkToGo(item)
			if err != nil {
				return nil, err
			}
			out = append(out, gv)
		}
		return out, nil
	case *starlark.Dict:
		out := map[string]any{}
		for _, item := range x.Items() {
			k, ok := starlark.AsString(item[0])
			if !ok {
				return nil, fmt.Errorf("dict key must be string, got %s", item[0].Type())
			}
			gv, err := StarlarkToGo(item[1])
			if err != nil {
				return nil, err
			}
			out[k] = gv
		}
		return out, nil
	default:
		return nil, fmt.Errorf("unsupported starlark return type: %s", v.Type())
	}
}
