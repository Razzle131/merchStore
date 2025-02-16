package merchRepo

type MerchRepoInterface interface {
	GetMerchPrice(item string) (int, error)
}
