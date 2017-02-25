package qparser

import (
    "net/url"
    "strings"
    "strconv"
)

// ListOptions specifies the optional parameters for requests with pagination support
type ListOptions struct {
    // Page of results to retrieve
    Page  int
    // Max number of results to retrieve on single page
    Limit int
}

// ExpandParams
type ExpandParams map[string]ListOptions

type ParseResult struct {
    Pagination  ListOptions
    Expand      ExpandParams
    Values      url.Values
}

type ParserOptions struct {
    LimitValue   int
    PageValue    int
    LimitString  string
    PageString   string
    ExpandString string
    QueryString  string
    LeftBracket  string
    RightBracket string
    Separator    string
    KVSeparator  string
}

type Parser struct {
    options *ParserOptions
}

const (
    DefaultLimitValue   int    = 25
    DefaultPageValue    int    = 1
    DefaultLimitString  string = "limit"
    DefaultPageString   string = "page"
    DefaultExpandString string = "expand"
    DefaultQueryString  string = "q"
    DefaultLeftBracket  rune   = '('
    DefaultRightBracket rune   = ')'
    DefaultSeparator    rune   = ','
    DefaultKVSeparator  rune   = ':'
)

func ifEmptyStringAssign(s *string, val string) {
    if *s == "" {
        *s = val
    } 
}

func ifEmptyRuneAssign(s *rune, val rune) {
    if *s == '' {
        *s = val
    } 
}
func ifEmptyIntAssign(s *int, val int) {
    if *s == 0 {
        *s = val
    } 
}

func NewParser(opts *ParserOptions) *Parser {
    if opts == nil {
        opts = &ParserOptions{}
    }

    ifEmptyIntAssign(&opts.LimitValue, DefaultLimitValue)
    ifEmptyIntAssign(&opts.PageValue, DefaultPageValue)
    ifEmptyStringAssign(&opts.LimitString, DefaultLimitString)
    ifEmptyStringAssign(&opts.ExpandString, DefaultExpandString)
    ifEmptyStringAssign(&opts.QueryString, DefaultQueryString)
    ifEmptyRuneAssign(&opts.LeftBracket, DefaultLeftBracket)
    ifEmptyRuneAssign(&opts.RightBracket, DefaultRightBracket)
    ifEmptyRuneAssign(&opts.Separator, DefaultSeparator)
    ifEmptyRuneAssign(&opts.KVSeparator, DefaultKVSeparator)

    return &Parser{
        options: opts,
    }
}

func (e *ExpandParams) Get(key string) (*ListOptions, error) {
    if v, ok := e[key]; ok {
        return &v, nil
    }

    return fmt.Errorf("No such key for expanded parameters")
}

func (qp *Parser) Parse(u *url.URL) (*ParseResult, error) {
    result := &ParseResult{}
    result.Values := r.URL.Query()    
    err := result.Pagination.parse(result.Values, qp.options)
    if err != nil {
        return nil, err
    }
    err = result.Expand.parse(result.Values, qp.options)
    if err != nil {
        return nil, err
    }

    return result, nil
}

func ParseString(rawurl string) (*ParseResult, error) {
    u, err := url.Parse(rawurl)
    if err != nil {
        return nil, err
    }
    return Parse(u)
}

func (lo *ListOptions) parse(val url.Values, opts *ParserOptions) error {
    if l := val.Get(opts.LimitString); l != "" {
        lo.Limit, err = strconv.Atoi(l)
        if err != nil {
            return err
        }
    } else {
        lo.Limit = opts.LimitValue
    }
    if p := val.Get(opts.PageString); p != "" {
        lo.Page, err = strconv.Atoi(p)
        if err != nil {
            return err
        }
    } else {
        lo.Page = opts.PageValue
    }  
    return nil
}

func (ep *ExpandParams) parse(val url.Values, opts *ParserOptions) error {
    params := map[string][]string(values)
    expStr := make([]string, 0)
    for _, str := range params[opts.ExpandString] {
        open := false
        position := 0
        for i, char := range str {            
            if char == opts.LeftBracket {
                open = true
            } else if char == opts.RightBracket {
                open = false
            } else if char == opts.Separator && !open {
                expStr = append(expStr, str[position:i])
                position = i
            }
            if i == (len(str) - 1) {
                if position > 0 {
                    position++
                }
                expStr = append(expStr, str[position:])
            }
        }
    }
    for _, char := range expStr {
        splitted := strings.FieldsFunc(char, func(r rune) bool {
            return r == opts.LeftBracket || r == opts.RightBracket || r == opts.Separator
        })
        bcp := ListOptions{
            Limit: opts.LimitValue,
            Page: opts.PageValue,
        }
        if len(splitted) > 1 {
            var err error
            params := strings.Split(splitted[1], string(opts.KVSeparator))
            if params[0] == opts.LimitString {
                bcp.Limit, err = strconv.Atoi(params[1])
            } else if params[0] == opts.PageString {
                bcp.Page, err = strconv.Atoi(params[1])
            }
            if err != nil {
                return err
            }
            if len(splitted) > 2 {
                params := strings.Split(splitted[2], string(opts.KVSeparator))
                if params[0] == opts.LimitString {
                    bcp.Limit, err = strconv.Atoi(params[1])
                } else if params[0] == opts.PageString {
                    bcp.Page, err = strconv.Atoi(params[1])
                }
                if err != nil {
                    return err
                }
            }
        }
        ep[splitted[0]] = bcp
    }
    return nil
}