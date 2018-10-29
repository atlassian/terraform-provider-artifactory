package hclwrite

import (
	"bytes"
	"io"
)

type File struct {
	inTree

	srcBytes []byte
	body     *node
}

// Body returns the root body of the file, which contains the top-level
// attributes and blocks.
func (f *File) Body() *Body {
	return f.body.content.(*Body)
}

// WriteTo writes the tokens underlying the receiving file to the given writer.
func (f *File) WriteTo(wr io.Writer) (int, error) {
	tokens := f.inTree.children.BuildTokens(nil)
	return tokens.WriteTo(wr)
}

// Bytes returns a buffer containing the source code resulting from the
// tokens underlying the receiving file. If any updates have been made via
// the AST API, these will be reflected in the result.
func (f *File) Bytes() []byte {
	buf := &bytes.Buffer{}
	f.WriteTo(buf)
	return buf.Bytes()
}

type comments struct {
	leafNode

	parent *node
	tokens Tokens
}

func newComments(tokens Tokens) *comments {
	return &comments{
		tokens: tokens,
	}
}

func (c *comments) BuildTokens(to Tokens) Tokens {
	return c.tokens.BuildTokens(to)
}

type identifier struct {
	leafNode

	parent *node
	token  *Token
}

func newIdentifier(token *Token) *identifier {
	return &identifier{
		token: token,
	}
}

func (i *identifier) BuildTokens(to Tokens) Tokens {
	return append(to, i.token)
}

func (i *identifier) hasName(name string) bool {
	return name == string(i.token.Bytes)
}

type number struct {
	leafNode

	parent *node
	token  *Token
}

func newNumber(token *Token) *number {
	return &number{
		token: token,
	}
}

func (n *number) BuildTokens(to Tokens) Tokens {
	return append(to, n.token)
}

type quoted struct {
	leafNode

	parent *node
	tokens Tokens
}

func newQuoted(tokens Tokens) *quoted {
	return &quoted{
		tokens: tokens,
	}
}

func (q *quoted) BuildTokens(to Tokens) Tokens {
	return q.tokens.BuildTokens(to)
}
