package flow

// MediaKind identifier transport-agnostic media type supported by flow view model.
type MediaKind string

const (
	MediaKindImage MediaKind = "image"
)

// MediaView describes one optional media attachment for a screen.
//
// Flow stays transport-agnostic: the concrete transport decides
// how to render the provided source.
type MediaView struct {
	Kind   MediaKind
	Source string
	Alt    string
}

// NewImageMedia builds one image media view.
//
// Empty source returns nil so caller can attach media conditionally.
func NewImageMedia(source, alt string) *MediaView {
	if source == "" {
		return nil
	}

	return &MediaView{
		Kind:   MediaKindImage,
		Source: source,
		Alt:    alt,
	}
}
