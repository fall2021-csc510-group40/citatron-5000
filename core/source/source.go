package source

import "core/schema"

type Search func(w *schema.Work) ([]*schema.Work, error)
