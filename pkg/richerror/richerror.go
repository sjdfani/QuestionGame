package richerror

type Kind int
type Operation string

const (
	KindInvalid Kind = iota + 1
	KindForbidden
	KindNotFound
	KindUnexpected
)

type RichError struct {
	operation    Operation
	wrappedError error
	message      string
	kind         Kind
	meta         map[string]any
}

func New(operation Operation) RichError {
	return RichError{operation: operation}
}

func (r RichError) SetMessage(message string) RichError {
	r.message = message
	return r
}

func (r RichError) SetError(err error) RichError {
	r.wrappedError = err
	return r
}

func (r RichError) SetKind(kind Kind) RichError {
	r.kind = kind
	return r
}

func (r RichError) SetMeta(meta map[string]any) RichError {
	r.meta = meta
	return r
}

func (r RichError) Error() string {
	return r.message
}

func (r RichError) GetMessage() string {
	if r.message != "" {
		return r.message
	}

	re, ok := r.wrappedError.(RichError)
	if !ok {
		return r.wrappedError.Error()
	}

	return re.GetMessage()
}

func (r RichError) GetKind() Kind {
	if r.kind != 0 {
		return r.kind
	}

	re, ok := r.wrappedError.(RichError)
	if !ok {
		return 0
	}

	return re.GetKind()
}

func (r RichError) GetMeta() map[string]any {
	return r.meta
}
