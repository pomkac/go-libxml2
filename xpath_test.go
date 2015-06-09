package libxml2

import "testing"

func TestXPathContext(t *testing.T) {
	doc, err := (&Parser{}).ParseString(`<foo><bar a="b"></bar></foo>`)
	if err != nil {
		t.Errorf("Failed to parse string: %s", err)
	}
	defer doc.Free()

	root := doc.DocumentElement()
	if root == nil {
		t.Errorf("Failed to get root element")
		return
	}

	ctx, err := NewXPathContext(root)
	if err != nil {
		t.Errorf("Failed to initialize XPathContext: %s", err)
		return
	}
	defer ctx.Free()

	// Use a string
	exprString := `/*`
	nodes, err := ctx.FindNodes(exprString)
	if err != nil {
		t.Errorf("Failed to execute FindNodes: %s", err)
		return
	}

	if len(nodes) != 1 {
		t.Errorf("Expected 1 nodes, got %d", len(nodes))
		return
	}

	// Use an explicitly compiled expression
	expr, err := NewXPathExpression(exprString)
	if err != nil {
		t.Errorf("Failed to compile xpath: %s", err)
		return
	}
	defer expr.Free()

	nodes, err = ctx.FindNodesExpr(expr)
	if err != nil {
		t.Errorf("Failed to execute FindNodesExpr: %s", err)
		return
	}

	if len(nodes) != 1 {
		t.Errorf("Expected 1 nodes, got %d", len(nodes))
		return
	}
}

func TestXPathContextExpression_Number(t *testing.T) {
	ctx, err := NewXPathContext()
	if err != nil {
		t.Errorf("Failed to initialize XPathContext: %s", err)
		return
	}
	defer ctx.Free()

	res, err := ctx.FindValue("1+1")
	if err != nil {
		t.Errorf("Failed to evaluate XPath expression: %s", err)
		return
	}
	defer res.Free()

	switch res.Type() {
	case XPathNumber:
		if res.Float64() != 2 {
			t.Errorf("Expected result number to be 2, got %f", res.Float64())
		}
	default:
		t.Errorf("Expected type to be XPathObjectNumber, got %s", res.Type())
	}
}

func TestXPathContextExpression_Boolean(t *testing.T) {
	ctx, err := NewXPathContext()
	if err != nil {
		t.Errorf("Failed to initialize XPathContext: %s", err)
		return
	}
	defer ctx.Free()

	res, err := ctx.FindValue("1=1")
	if err != nil {
		t.Errorf("Failed to evaluate XPath expression: %s", err)
		return
	}
	defer res.Free()

	switch res.Type() {
	case XPathBoolean:
		if !res.Bool() {
			t.Errorf("Expected result number to be false, got %s", res.Bool())
		}
	default:
		t.Errorf("Expected type to be XPathObjectBoolean, got %s", res.Type())
	}
}

func TestXPathContextExpression_NodeList(t *testing.T) {
	doc, err := (&Parser{}).ParseString(`<foo><bar a="b">baz</bar></foo>`)
	if err != nil {
		t.Errorf("Failed to parse string: %s", err)
	}
	defer doc.Free()

	root := doc.DocumentElement()
	if root == nil {
		t.Errorf("Failed to get root element")
		return
	}

	ctx, err := NewXPathContext(root)
	if err != nil {
		t.Errorf("Failed to initialize XPathContext: %s", err)
		return
	}
	defer ctx.Free()

	res, err := ctx.FindValue("/foo/bar/text()")
	if err != nil {
		t.Errorf("Failed to evaluate XPath expression: %s", err)
		return
	}
	defer res.Free()

	switch res.Type() {
	case XPathNodeSet:
		if res.NodeList().String() != "baz" {
			t.Errorf("Expected result NodeList to be baz, got %s", res.NodeList().String())
		}
	default:
		t.Errorf("Expected type to be XPathObjectNodeSet, got %s", res.Type())
	}
}