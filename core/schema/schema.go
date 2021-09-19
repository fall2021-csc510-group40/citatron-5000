package schema

type Work struct {
	ID   string
	Hash string

	Type string

	DOI   string
	Arxiv string
	ISBN  string

	Title   string
	Authors []string

	Version string
	Venue   string
	Page    string

	Year  int
	Month int
	Day   int

	Keywords []string
}

type SearchRequest struct {
	Query *Work
	Page  int
}

type SearchResponse struct {
	Results []*Work
	Error   string
}

type FormatRequest struct {
	ID     string
	Work   *Work
	Format string
}

type FormatResponse struct {
	Result string
	Error  string
}
