package backends

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/ymgyt/cloudops/core"
)

// NewPromptConfirmer -
func NewPromptConfirmer(w io.Writer, r io.Reader, yesNo string, yesAnswers []string) (core.Confirmer, error) {
	return &promptConfirmer{w: w, r: r, yesNo: yesNo, yesAnswers: yesAnswers}, nil
}

type promptConfirmer struct {
	w          io.Writer
	r          io.Reader
	yesNo      string
	yesAnswers []string
}

// Confirm -
func (p *promptConfirmer) Confirm(operation string, resources core.Resources) (bool, error) {
	msg := p.msg(operation, resources)
	io.WriteString(p.w, msg)
	b := bufio.NewReader(p.r)
	answer, err := b.ReadString('\n')
	if err != nil {
		return false, err
	}

	answer = strings.TrimSpace(answer)
	for _, want := range p.yesAnswers {
		if answer == want {
			return true, nil
		}
	}
	return false, nil
}

func (p *promptConfirmer) msg(operation string, resources core.Resources) string {
	b := core.Buffer()
	defer func() { core.PutBuffer(b) }()

	for _, r := range resources {
		b.WriteString(r.URI() + "\n")
	}
	b.WriteString(fmt.Sprintf("%s %s ", operation, p.yesNo))
	return b.String()
}
