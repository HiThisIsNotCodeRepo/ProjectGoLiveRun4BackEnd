package router

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"paotui.sg/app/handler/my_info"
)

const (
	SpendingURL             = "/spending"
	YesterdaySpendingURL    = "/yesterday"
	TwoDaysAgoSpendingURL   = "/two-days-ago"
	ThreeDaysAgoSpendingURL = "/three-days-ago"
	ThisWeekURL             = "/this-week"
	LastWeekURL             = "/last-week"
	BuyNecessityURL         = "/buy-necessity"
	FoodDeliveryURL         = "/food-delivery"
	SendDocumentURL         = "/send-document"
	OtherURL                = "/other"
	SummaryURL              = "/summary"
	userIDURL               = "/{userID}"
)

func SpendCardRouter(s *mux.Router) *mux.Router {
	dateArray := []string{YesterdaySpendingURL, TwoDaysAgoSpendingURL, ThreeDaysAgoSpendingURL}
	categoryArray := []string{BuyNecessityURL, FoodDeliveryURL, SendDocumentURL, OtherURL}
	for _, tempDateUrl := range dateArray {
		for _, tempCategoryUrl := range categoryArray {
			url := fmt.Sprintf("%s%s%s%s", SpendingURL, tempDateUrl, tempCategoryUrl, userIDURL)
			fmt.Println(url)
			s.HandleFunc(url, my_info.GetSpendingCard).Methods(http.MethodGet, http.MethodOptions)
		}
	}
	return s
}
func SpendSummaryRouter(s *mux.Router) *mux.Router {
	dateArray := []string{ThisWeekURL, LastWeekURL}
	for _, tempDateUrl := range dateArray {
		url := fmt.Sprintf("%s%s%s%s", SpendingURL, tempDateUrl, SummaryURL, userIDURL)
		fmt.Println(url)
		s.HandleFunc(url, my_info.GetSummary).Methods(http.MethodGet, http.MethodOptions)
	}
	return s
}
func SpendTaskRouter(s *mux.Router) *mux.Router {
	url := fmt.Sprintf("%s%s%s", SpendingURL, "/tasks", userIDURL)
	fmt.Println(url)
	s.HandleFunc(url, my_info.GetSpendingDataSource).Methods(http.MethodGet, http.MethodOptions)
	return s
}
