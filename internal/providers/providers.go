package providers

type (
	APIClient interface {
		FetchData(address string) ([]byte, error)
	}
)
