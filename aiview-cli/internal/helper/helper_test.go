package helper

import "testing"

func TestGetString(t *testing.T) {
	m := map[string]interface{}{
		"name": "test",
		"nested": map[string]interface{}{
			"key": "value",
		},
	}
	if v := GetString(m, "name"); v != "test" {
		t.Errorf("expected 'test', got '%s'", v)
	}
	if v := GetString(m, "missing"); v != "" {
		t.Errorf("expected '', got '%s'", v)
	}
	if v := GetString(m, "nested", "key"); v != "value" {
		t.Errorf("expected 'value', got '%s'", v)
	}
	if v := GetString(nil, "key"); v != "" {
		t.Errorf("expected '', got '%s'", v)
	}
}

func TestGetInt(t *testing.T) {
	m := map[string]interface{}{
		"count": 42.0,
		"nested": map[string]interface{}{
			"num": 7.0,
		},
	}
	if v := GetInt(m, "count"); v != 42 {
		t.Errorf("expected 42, got %d", v)
	}
	if v := GetInt(m, "missing"); v != 0 {
		t.Errorf("expected 0, got %d", v)
	}
	if v := GetInt(m, "nested", "num"); v != 7 {
		t.Errorf("expected 7, got %d", v)
	}
	if v := GetInt(nil, "key"); v != 0 {
		t.Errorf("expected 0, got %d", v)
	}
}

func TestGetMap(t *testing.T) {
	m := map[string]interface{}{
		"data": map[string]interface{}{"name": "test"},
	}
	if v := GetMap(m, "data"); v == nil || v["name"] != "test" {
		t.Errorf("expected map, got %v", v)
	}
	if v := GetMap(m, "missing"); v != nil {
		t.Errorf("expected nil, got %v", v)
	}
	if v := GetMap(nil, "key"); v != nil {
		t.Errorf("expected nil, got %v", v)
	}
}

func TestGetSlice(t *testing.T) {
	m := map[string]interface{}{
		"items": []interface{}{"a", "b"},
	}
	if v := GetSlice(m, "items"); len(v) != 2 {
		t.Errorf("expected len 2, got %d", len(v))
	}
	if v := GetSlice(m, "missing"); v != nil {
		t.Errorf("expected nil, got %v", v)
	}
	if v := GetSlice(nil, "key"); v != nil {
		t.Errorf("expected nil, got %v", v)
	}
}