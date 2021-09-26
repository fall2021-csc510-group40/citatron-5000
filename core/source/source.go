package source

import "core/schema"

// Search is a generic search function for populating the database
type Search func(w *schema.Work) ([]*schema.Work, error)
