package article

import "github/islamghany/blog/foundation/validate"

type QueryFilter struct {
	TSV_Document *string `validate:"omitempty,min=1"`
}

func (qf *QueryFilter) Validate() error {
	return validate.Check(qf)
}

func (qf *QueryFilter) WithTSV_Document(tsv_document string) {
	qf.TSV_Document = &tsv_document
}
