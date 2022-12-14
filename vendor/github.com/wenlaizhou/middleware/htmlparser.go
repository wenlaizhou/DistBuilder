package middleware

/* soup package implements a simple web scraper for Go,
keeping it as similar as possible to BeautifulSoup
*/
import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// ErrorType defines types of errors that are possible from soup
type ErrorType int

const (
	// ErrUnableToParse will be returned when the HTML could not be parsed
	ErrUnableToParse ErrorType = iota
	// ErrElementNotFound will be returned when element was not found
	ErrElementNotFound
	// ErrNoNextSibling will be returned when no next sibling can be found
	ErrNoNextSibling
	// ErrNoPreviousSibling will be returned when no previous sibling can be found
	ErrNoPreviousSibling
	// ErrNoNextElementSibling will be returned when no next element sibling can be found
	ErrNoNextElementSibling
	// ErrNoPreviousElementSibling will be returned when no previous element sibling can be found
	ErrNoPreviousElementSibling
	// ErrCreatingGetRequest will be returned when the get request couldn't be created
	ErrCreatingGetRequest
	// ErrInGetRequest will be returned when there was an error during the get request
	ErrInGetRequest
	// ErrCreatingPostRequest will be returned when the post request couldn't be created
	ErrCreatingPostRequest
	// ErrMarshallingPostRequest will be returned when the body of a post request couldn't be serialized
	ErrMarshallingPostRequest
	// ErrReadingResponse will be returned if there was an error reading the response to our get request
	ErrReadingResponse
)

// Error allows easier introspection on the type of error returned.
// If you know you have a Error, you can compare the Type to one of the exported types
// from this package to see what kind of error it is, then further inspect the Error() method
// to see if it has more specific details for you, like in the case of a ErrElementNotFound
// type of error.
type Error struct {
	Type ErrorType
	msg  string
}

func (se Error) Error() string {
	return se.msg
}

func newError(t ErrorType, msg string) Error {
	return Error{Type: t, msg: msg}
}

// Root is a structure containing a pointer to an html node, the node value, and an error variable to return an error if one occurred
type Root struct {
	Pointer   *html.Node
	NodeValue string
	Error     error
}

// HTMLParse parses the HTML returning a start pointer to the DOM
func HTMLParse(s string) Root {
	r, err := html.Parse(strings.NewReader(s))
	if err != nil {
		return Root{Error: newError(ErrUnableToParse, "unable to parse the HTML")}
	}
	for r.Type != html.ElementNode {
		switch r.Type {
		case html.DocumentNode:
			r = r.FirstChild
		case html.DoctypeNode:
			r = r.NextSibling
		case html.CommentNode:
			r = r.NextSibling
		}
	}
	return Root{Pointer: r, NodeValue: r.Data}
}

// Find finds the first occurrence of the given tag name,
// with or without attribute key and value specified,
// and returns a struct with a pointer to it
func (r Root) Find(args ...string) Root {
	temp, ok := findOnce(r.Pointer, args, false, false)
	if !ok {
		return Root{Error: newError(ErrElementNotFound, fmt.Sprintf("element `%s` with attributes `%s` not found", args[0], strings.Join(args[1:], " ")))}
	}
	return Root{Pointer: temp, NodeValue: temp.Data}
}

// FindAll finds all occurrences of the given tag name,
// with or without key and value specified,
// and returns an array of structs, each having
// the respective pointers
func (r Root) FindAll(args ...string) []Root {
	temp := findAllofem(r.Pointer, args, false)
	if len(temp) == 0 {
		return []Root{}
	}
	pointers := make([]Root, 0, len(temp))
	for i := 0; i < len(temp); i++ {
		pointers = append(pointers, Root{Pointer: temp[i], NodeValue: temp[i].Data})
	}
	return pointers
}

// FindStrict finds the first occurrence of the given tag name
// only if all the values of the provided attribute are an exact match
func (r Root) FindStrict(args ...string) Root {
	temp, ok := findOnce(r.Pointer, args, false, true)
	if !ok {
		return Root{nil, "", newError(ErrElementNotFound, fmt.Sprintf("element `%s` with attributes `%s` not found", args[0], strings.Join(args[1:], " ")))}
	}
	return Root{Pointer: temp, NodeValue: temp.Data}
}

// FindAllStrict finds all occurrences of the given tag name
// only if all the values of the provided attribute are an exact match
func (r Root) FindAllStrict(args ...string) []Root {
	temp := findAllofem(r.Pointer, args, true)
	if len(temp) == 0 {
		return []Root{}
	}
	pointers := make([]Root, 0, len(temp))
	for i := 0; i < len(temp); i++ {
		pointers = append(pointers, Root{Pointer: temp[i], NodeValue: temp[i].Data})
	}
	return pointers
}

// FindNextSibling finds the next sibling of the pointer in the DOM
// returning a struct with a pointer to it
func (r Root) FindNextSibling() Root {
	nextSibling := r.Pointer.NextSibling
	if nextSibling == nil {
		return Root{Error: newError(ErrNoNextSibling, "no next sibling found")}
	}
	return Root{Pointer: nextSibling, NodeValue: nextSibling.Data}
}

// FindPrevSibling finds the previous sibling of the pointer in the DOM
// returning a struct with a pointer to it
func (r Root) FindPrevSibling() Root {
	prevSibling := r.Pointer.PrevSibling
	if prevSibling == nil {
		return Root{Error: newError(ErrNoPreviousSibling, "no previous sibling found")}
	}
	return Root{Pointer: prevSibling, NodeValue: prevSibling.Data}
}

// FindNextElementSibling finds the next element sibling of the pointer in the DOM
// returning a struct with a pointer to it
func (r Root) FindNextElementSibling() Root {
	nextSibling := r.Pointer.NextSibling
	if nextSibling == nil {
		return Root{Error: newError(ErrNoNextElementSibling, "no next element sibling found")}
	}
	if nextSibling.Type == html.ElementNode {
		return Root{Pointer: nextSibling, NodeValue: nextSibling.Data}
	}
	p := Root{Pointer: nextSibling, NodeValue: nextSibling.Data}
	return p.FindNextElementSibling()
}

// FindPrevElementSibling finds the previous element sibling of the pointer in the DOM
// returning a struct with a pointer to it
func (r Root) FindPrevElementSibling() Root {
	prevSibling := r.Pointer.PrevSibling
	if prevSibling == nil {
		return Root{Error: newError(ErrNoPreviousElementSibling, "no previous element sibling found")}
	}
	if prevSibling.Type == html.ElementNode {
		return Root{Pointer: prevSibling, NodeValue: prevSibling.Data}
	}
	p := Root{Pointer: prevSibling, NodeValue: prevSibling.Data}
	return p.FindPrevElementSibling()
}

// Children returns all direct children of this DOME element.
func (r Root) Children() []Root {
	child := r.Pointer.FirstChild
	var children []Root
	for child != nil {
		children = append(children, Root{Pointer: child, NodeValue: child.Data})
		child = child.NextSibling
	}
	return children
}

// Attrs returns a map containing all attributes
func (r Root) Attrs() map[string]string {
	if r.Pointer.Type != html.ElementNode {
		return nil
	}
	if len(r.Pointer.Attr) == 0 {
		return nil
	}
	return getKeyValue(r.Pointer.Attr)
}

// Text returns the string inside a non-nested element
func (r Root) Text() string {
	k := r.Pointer.FirstChild
checkNode:
	if k != nil && k.Type != html.TextNode {
		k = k.NextSibling
		if k == nil {
			return ""
		}
		goto checkNode
	}
	if k != nil {
		r, _ := regexp.Compile(`^\s+$`)
		if ok := r.MatchString(k.Data); ok {
			k = k.NextSibling
			if k == nil {
				return ""
			}
			goto checkNode
		}
		return k.Data
	}
	return ""
}

// HTML returns the HTML code for the specific element
func (r Root) HTML() string {
	var buf bytes.Buffer
	if err := html.Render(&buf, r.Pointer); err != nil {
		return ""
	}
	return buf.String()
}

// FullText returns the string inside even a nested element
func (r Root) FullText() string {
	var buf bytes.Buffer

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n == nil {
			return
		}
		if n.Type == html.TextNode {
			buf.WriteString(n.Data)
		}
		if n.Type == html.ElementNode {
			f(n.FirstChild)
		}
		if n.NextSibling != nil {
			f(n.NextSibling)
		}
	}

	f(r.Pointer.FirstChild)

	return buf.String()
}

func matchElementName(n *html.Node, name string) bool {
	return name == "" || name == n.Data
}

// Using depth first search to find the first occurrence and return
func findOnce(n *html.Node, args []string, uni bool, strict bool) (*html.Node, bool) {
	if uni == true {
		if n.Type == html.ElementNode && matchElementName(n, args[0]) {
			if len(args) > 1 && len(args) < 4 {
				for i := 0; i < len(n.Attr); i++ {
					attr := n.Attr[i]
					searchAttrName := args[1]
					searchAttrVal := args[2]
					if (strict && attributeAndValueEquals(attr, searchAttrName, searchAttrVal)) ||
						(!strict && attributeContainsValue(attr, searchAttrName, searchAttrVal)) {
						return n, true
					}
				}
			} else if len(args) == 1 {
				return n, true
			}
		}
	}
	uni = true
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		p, q := findOnce(c, args, true, strict)
		if q != false {
			return p, q
		}
	}
	return nil, false
}

// Using depth first search to find all occurrences and return
func findAllofem(n *html.Node, args []string, strict bool) []*html.Node {
	var nodeLinks = make([]*html.Node, 0, 10)
	var f func(*html.Node, []string, bool)
	f = func(n *html.Node, args []string, uni bool) {
		if uni == true {
			if n.Type == html.ElementNode && matchElementName(n, args[0]) {
				if len(args) > 1 && len(args) < 4 {
					for i := 0; i < len(n.Attr); i++ {
						attr := n.Attr[i]
						searchAttrName := args[1]
						searchAttrVal := args[2]
						if (strict && attributeAndValueEquals(attr, searchAttrName, searchAttrVal)) ||
							(!strict && attributeContainsValue(attr, searchAttrName, searchAttrVal)) {
							nodeLinks = append(nodeLinks, n)
						}
					}
				} else if len(args) == 1 {
					nodeLinks = append(nodeLinks, n)
				}
			}
		}
		uni = true
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c, args, true)
		}
	}
	f(n, args, false)
	return nodeLinks
}

// attributeAndValueEquals reports when the html.Attribute attr has the same attribute name and value as from
// provided arguments
func attributeAndValueEquals(attr html.Attribute, attribute, value string) bool {
	return attr.Key == attribute && attr.Val == value
}

// attributeContainsValue reports when the html.Attribute attr has the same attribute name as from provided
// attribute argument and compares if it has the same value in its values parameter
func attributeContainsValue(attr html.Attribute, attribute, value string) bool {
	if attr.Key == attribute {
		for _, attrVal := range strings.Fields(attr.Val) {
			if attrVal == value {
				return true
			}
		}
	}
	return false
}

// Returns a key pair value (like a dictionary) for each attribute
func getKeyValue(attributes []html.Attribute) map[string]string {
	var keyvalues = make(map[string]string)
	for i := 0; i < len(attributes); i++ {
		_, exists := keyvalues[attributes[i].Key]
		if exists == false {
			keyvalues[attributes[i].Key] = attributes[i].Val
		}
	}
	return keyvalues
}
