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
        if assert.Equal(t, 1, len(res.Values["name"])) {
            assert.Equal(t, "somename", res.Values["name"][0])    
            assert.Equal(t, "somename", res.Values["description"][0])   
        }        
    }
}

func TestQueryDefault2(t *testing.T) { 
    url := "http://some-api.com/api/endpoint?q=somename&p=name&p=description"   
    parser := NewParser(nil)
    res, err := parser.ParseString(url)
    if assert.Nil(t, err) {
        if assert.Equal(t, 1, len(res.Values["name"])) {
            assert.Equal(t, "somename", res.Values["name"][0])    
            assert.Equal(t, "somename", res.Values["description"][0])   
        }        
    }
}

func TestExpandItem(t *testing.T) { 
    url := "http://some-api.com/api/endpoint?expand=relation"   
    parser := NewParser(nil)
    res, err := parser.ParseString(url)
    if assert.Nil(t, err) {
        assert.NotNil(t, res.Expand["relation"])
    }
}

func TestExpandList(t *testing.T) { 
    url := "http://some-api.com/api/endpoint?expand=relation(limit:6,page:8)"   
    parser := NewParser(nil)
    res, err := parser.ParseString(url)
    if assert.Nil(t, err) {
        if assert.NotNil(t, res.Expand["relation"]) {
            assert.Equal(t, 6, res.Expand["relation"].Limit)
            assert.Equal(t, 8, res.Expand["relation"].Page)
        }
    }
}
