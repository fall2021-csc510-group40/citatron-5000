package source

import "core/schema"

type Source interface {
	Search(w *schema.Work) ([]*schema.Work, error)
}
