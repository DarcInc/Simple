package service

import (
	"context"
	"github.com/darcinc/Simple/model"
	"github.com/darcinc/Simple/reflex"
	"html/template"
	"net/http"
)

type ImageSearchRequest struct {
}

// ImageSearchResponse is a holding of data that our output representation can understand.
// We pass this along to the template to format as a page.
type ImageSearchResponse struct {
}

type ImageSearcher struct {
	Repository model.ImageRepository
}

// TODO: Need a permanent reference model for internal objects.
// TODO: Permanent reference model needs a service to back the data.

// Search is the boundary layer between internal and external.  The returned image model
// is still an internal concept/abstraction.  It exists within the scope of how we reason
// about our system - not how the world necessarily reasons about our software.
//
// Our notion of image is the subject-matter captured, which may be represented with
// multiple encodings, which could each have multiple stored data.  E.g. a picture of a
// man on a boat could have a jpeg and a tiff encoding, and the jpeg could be stored in
// two places, meaning there are two copies of the same JPEG.
//
// Here's where we also need to generate permalinks, if they don't already exist.
func (is ImageSearcher) Search(ctx context.Context, _ ImageSearchRequest) (ImageSearchResponse, error) {
	qp := model.QueryParameters{}

	images, err := is.Repository.Find(ctx, qp)
	if err != nil {
		// TODO: Wrap error appropriately
		return ImageSearchResponse{}, err
	}

	response := ImageSearchResponse{}
	for range images {
		// TODO: Handle adding an image to the response.
	}

	return response, nil
}

type ImageSearchHandler struct {
	SearchPage *template.Template
}

func (ish ImageSearchHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	theReflex := reflex.GlobalReflex()
	searcher, ok := theReflex.MustGet("ImageSearcher").(ImageSearcher)
	if !ok {
		w.WriteHeader(500)
		// TODO: Fill in the error page.
		return
	}

	// TODO: Extract search request from parameters
	isr := ImageSearchRequest{}

	// TODO: Set appropriate context (e.g. timeout)
	results, err := searcher.Search(context.Background(), isr)
	if err != nil {
		w.WriteHeader(500)
		// TODO: Fill in error handling
	}

	if err := ish.SearchPage.Execute(w, results); err != nil {
		// TODO: Handle page execution error
	}
}
