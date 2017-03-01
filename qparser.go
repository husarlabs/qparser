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

type QueryValues map[string][]string

type ParseResult struct {
    Pagination  *ListOptions
    Expand      *ExpandParams
    Values      *QueryValues
}

type ParserOptions struct {
    LimitValue   int
    PageValue    int
    LimitString  string
    PageString   string
    ExpandString string
    QueryString  string
    ParamString  string
    LeftBracket  rune
    RightBracket rune
    Separator    rune
    KVSeparator  rune
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
    DefaultParamString  string = "p"
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
    if *s == 0 {
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
    ifEmptyStringAssign(&opts.PageString, DefaultPageString)
    ifEmptyStringAssign(&opts.ExpandString, DefaultExpandString)
    ifEmptyStringAssign(&opts.QueryString, DefaultQueryString)
    ifEmptyStringAssign(&opts.ParamString, DefaultParamString)
    ifEmptyRuneAssign(&opts.LeftBracket, DefaultLeftBracket)
    ifEmptyRuneAssign(&opts.RightBracket, DefaultRightBracket)
    ifEmptyRuneAssign(&opts.Separator, DefaultSeparator)
    ifEmptyRuneAssign(&opts.KVSeparator, DefaultKVSeparator)

    return &Parser{
        options: opts,
    }
}

func (e *ExpandParams) Get(key string) (*ListOptions) {
    if v, ok := (*e)[key]; ok {
        return &v
    }

    return nil
}

func (e *QueryValues) Get(key string) ([]string) {
    if v, ok := (*e)[key]; ok {
        return v
    }

    return nil
}

func (qp *Parser) Parse(u *url.URL) (*ParseResult, error) {
    result := &ParseResult{
        Pagination:  &ListOptions{},
        Expand:      &ExpandParams{},
        Values:      &QueryValues{},
    }
    values := u.Query()    
    err := result.Pagination.parse(values, qp.options)
    if err != nil {
        return nil, err
    }
    err = result.Expand.parse(values, qp.options)
    if err != nil {
        return nil, err
    }
    err = result.Values.parse(values, qp.options)
    if err != nil {
        return nil, err
    }
    return result, nil
}

func (qp *Parser) ParseString(rawurl string) (*ParseResult, error) {
    u, err := url.Parse(rawurl)
    if err != nil {
        return nil, err
    }
    return qp.Parse(u)
}

func (lo *ListOptions) parse(val url.Values, opts *ParserOptions) error {
    var err error
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
    params := map[string][]string(val)
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
        if *ep == nil {
            (*ep) = make(map[string]ListOptions)
        }
        (*ep)[splitted[0]] = bcp
    }
    return nil
}

func (qv *QueryValues) parse(val url.Values, opts *ParserOptions) error {
    params := map[string][]string(val)
    paramsStr := make([]string, 0)
    for _, str := range params[opts.ParamString] {
        paramsStr = append(paramsStr, strings.Split(str, string(opts.Separator))...)
    }

    for _, str := range paramsStr {
        if *qv == nil {
            (*qv) = make(map[string][]string)
        } else if (*qv)[str] == nil {
            (*qv)[str] = make([]string, 0)
        }
        (*qv)[str] = append((*qv)[str], params[opts.QueryString][0])
    }

    return nil
}