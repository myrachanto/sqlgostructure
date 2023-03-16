package support

type Paginator struct {
	Search   string `json:"search"`
	Page     int    `json:"page"`
	Pagesize int    `json:"pagesize"`
}
func Paginate(x []int, skip int, size int) []int {
    if skip > len(x) {
        skip = len(x)
    }

    end := skip + size
    if end > len(x) {
        end = len(x)
    }

    return x[skip:end]
}