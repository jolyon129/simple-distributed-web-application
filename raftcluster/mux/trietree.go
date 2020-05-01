package mux

import (
    "fmt"
    "net/http"
    "strings"
)

// Node represent a node in the trie that can be matched
type Node struct {
    parent *Node
    // For fixed pattern route, Name is same as segment value, /a -> the Name is a
    // For named parameter route, /:id -> the Name is id
    Name             string
    IsNamedParameter bool
    children         map[string]*Node
    NextChildIsNamed bool // if the next child is named parameter.
    IsEndPoint       bool // whether this node is a api endpoint
    // Each node can have multiple handlers with different method names: GET,PUT,DELETE
    handlers map[string]http.HandlerFunc
}

type Trie struct {
    root *Node
}

func NewTrie() *Trie {
    return &Trie{root: &Node{
        parent:     nil,
        children:   make(map[string]*Node),
        IsEndPoint: false,
        handlers:   nil,
    }}
}

// Parse the url and build endpoint in the Trie.
// Return the corresponding Node of the endpoint.
func (t *Trie) Parse(url string) *Node {
    url = strings.TrimPrefix(url, "/")
    segments := strings.Split(url, "/")
    node := t.root
    for idx, segment := range segments {
        if child, ok := hasChild(node, segment); !ok {
            key := segment
            name := segment
            var isNamedParameter bool
            if strings.HasPrefix(segment, ":") {
                isNamedParameter = true
                key = NAMED_PARAMETER
                name = strings.TrimPrefix(segment, ":")
                node.NextChildIsNamed = true
            }
            node.children[key] = &Node{
                parent:           node,
                children:         make(map[string]*Node),
                IsEndPoint:       idx == len(segments)-1,
                handlers:         nil,
                IsNamedParameter: isNamedParameter,
                Name:             name,
            }
            node = node.children[key]
        } else {
            node = child
        }
    }
    node.IsEndPoint = true // Have to update
    return node
}

// If has child, return the child.
// If not, return (nil,false).
func hasChild(node *Node, segment string) (*Node, bool) {
    if node.NextChildIsNamed { // If the next child is named parameter
        return node.children[NAMED_PARAMETER], true
    } else {
        ret, ok := node.children[segment]
        return ret, ok
    }
}

type EndpointWithParams struct {
    *Node
    Params map[string]string
}

// Search the endpoint according to the route
func (t *Trie) Lookup(url string) (*EndpointWithParams, error) {
    url = strings.TrimPrefix(url, "/")
    url = strings.TrimSuffix(url, "/")
    segments := strings.Split(url, "/")
    node := t.root
    params := make(map[string]string) // Store the named parameters in the route
    for _, segment := range segments {
        if child, ok := hasChild(node, segment); !ok {
            return nil, fmt.Errorf("the route does not exist for %s", url)
        } else {
            if child.IsNamedParameter { // store the parameter
                params[child.Name] = segment
            }
            node = child
        }
    }
    if node.IsEndPoint {
        return &EndpointWithParams{
            Node:   node,
            Params: params,
        }, nil
    } else {
        return nil, fmt.Errorf("the route is not an endpoint for %s", url)
    }
}

// Add handle function to the node
func (n *Node) Handle(methodName string, handler http.HandlerFunc) {
    n.handlers[methodName] = handler
}
