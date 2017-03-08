package qparser

import (
    "testing"
    "github.com/stretchr/testify/assert"
)


func TestLimitPage(t *testing.T) { 
    url := "http://some-api.com/api/endpoint?limit=1&page=2"   
    parser := NewParser(nil)
    res, err := parser.ParseString(url)
    if assert.Nil(t, err) {
        assert.Equal(t, 1, res.Pagination.Limit)
        assert.Equal(t, 2, res.Pagination.Page)
    }
}

func TestLimitPageDefault(t *testing.T) { 
    url := "http://some-api.com/api/endpoint"   
    parser := NewParser(nil)
    res, err := parser.ParseString(url)
    if assert.Nil(t, err) {
        assert.Equal(t, 25, res.Pagination.Limit)
        assert.Equal(t, 1, res.Pagination.Page)
    }
}

func TestLimitPageCustom(t *testing.T) { 
    opts := &ParserOptions {
        LimitString: "l",
        PageString: "pg",
    }
    url := "http://some-api.com/api/endpoint?l=1&pg=2"   
    parser := NewParser(opts)
    res, err := parser.ParseString(url)
    if assert.Nil(t, err) {
        assert.Equal(t, 1, res.Pagination.Limit)
        assert.Equal(t, 2, res.Pagination.Page)
    }
}

func TestQueryDefault(t *testing.T) { 
    url := "http://some-api.com/api/endpoint?q=somename&p=name,description"   
    parser := NewParser(nil)
    res, err := parser.ParseString(url)
    if assert.Nil(t, err) {
        if assert.Equal(t, 2, len(res.Values.Search.Keys)) {
            assert.Equal(t, "name", res.Values.Search.Keys[0])    
            assert.Equal(t, "description", res.Values.Search.Keys[1])   
        }        
        assert.Equal(t, "somename", res.Values.Search.Value)    
    }
}

func TestQueryDefault2(t *testing.T) { 
    url := "http://some-api.com/api/endpoint?q=somename&p=name&p=description"   
    parser := NewParser(nil)
    res, err := parser.ParseString(url)
    if assert.Nil(t, err) {
        if assert.Equal(t, 2, len(res.Values.Search.Keys)) {
            assert.Equal(t, "name", res.Values.Search.Keys[0])    
            assert.Equal(t, "description", res.Values.Search.Keys[1])   
        }        
        assert.Equal(t, "somename", res.Values.Search.Value) 
    } 
}

func TestQueryAllParams(t *testing.T) { 
    url := "http://some-api.com/api/endpoint?q=somename"   
    parser := NewParser(nil)
    res, err := parser.ParseString(url)
    if assert.Nil(t, err) {
        if assert.Equal(t, 0, len(res.Values.Search.Keys)) {
            assert.Equal(t, "somename", res.Values.Search.Value) 
        }        
    }
}

func TestQueryParticular(t *testing.T) { 
    url := "http://some-api.com/api/endpoint?fruit=apple&color=red"   
    parser := NewParser(nil)
    res, err := parser.ParseString(url)
    if assert.Nil(t, err) {
        if assert.Equal(t, 2, len(res.Values.Filter)) {
            assert.Equal(t, "apple", res.Values.Filter["fruit"][0]) 
            assert.Equal(t, "red", res.Values.Filter["color"][0]) 
        }        
    }
}


func TestExpandItem(t *testing.T) { 
    url := "http://some-api.com/api/endpoint?expand=relation"   
    parser := NewParser(nil)
    res, err := parser.ParseString(url)
    if assert.Nil(t, err) {
        assert.NotNil(t, res.Expand.Get("relation"))
    }
}

func TestExpandList(t *testing.T) { 
    url := "http://some-api.com/api/endpoint?expand=relation(limit:6,page:8)"   
    parser := NewParser(nil)
    res, err := parser.ParseString(url)
    if assert.Nil(t, err) {
        if assert.NotNil(t, res.Expand.Get("relation")) {
            assert.Equal(t, 6, res.Expand.Get("relation").Limit)
            assert.Equal(t, 8, res.Expand.Get("relation").Page)
        }
    }
}

func TestOrder(t *testing.T) { 
    url := "http://some-api.com/api/endpoint?order=field1(asc),field2(desc)"   
    parser := NewParser(nil)
    res, err := parser.ParseString(url)
    if assert.Nil(t, err) {
        if assert.NotNil(t, res.Values.Order) {
            v1 := (*res.Values.Order)[0]
            v2 := (*res.Values.Order)[1]
            assert.Equal(t, "field1", v1.Field)
            assert.Equal(t, ASC, v1.Order)
            assert.Equal(t, "field2", v2.Field)
            assert.Equal(t, DESC, v2.Order)
        }
    }
}
