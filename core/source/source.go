/*
Copyright (c) 2021 contributors of the Citatron-5000 Project. All Rights Reserved
*/
package source

import "core/schema"

// Search is a generic search function for populating the database
type Search func(w *schema.Work) ([]*schema.Work, error)
