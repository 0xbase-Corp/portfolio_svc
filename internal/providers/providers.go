package providers

type (
	APIClient interface {
		FetchData(address string) ([]byte, error)
	}

	PriceFeedClient interface {
		FetchData(cryptoID, currency string) ([]byte, error)
	}
)
