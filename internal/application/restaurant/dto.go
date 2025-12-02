package restaurantapp

import "go-ai/internal/domain/restaurant"

type RestaurantBase struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Address     string `json:"address"`
	Category    string `json:"category"`
	City        string `json:"city"`
	District    string `json:"district"`
	LogoUrl     string `json:"logo_url"`
	BannerUrl   string `json:"banner_url"`
	PhoneNumber string `json:"phone_number"`
	WebsiteUrl  string `json:"website_url"`
	Email       string `json:"email"`
}

type RestaurantHoursBase struct {
	Day       restaurant.DayOfWeek `json:"day"`
	OpenTime  string               `json:"open_time"`
	CloseTime string               `json:"close_time"`
	// NextDay   bool                 `json:"next_day"`
}

type CreateRestaurantRequest struct {
	RestaurantBase
	Hours []RestaurantHoursBase `json:"hours"`
}

type CreateRestaurantResponse struct {
	Id int32 `json:"id"`
}

type GetRestaurantByIDRequest struct {
	Id int32 `json:"id"`
}

type GetRestaurantByIDResponse struct {
	RestaurantBase
	UserName string                `json:"user_name"`
	IsActive bool                  `json:"is_active"`
	Hours    []RestaurantHoursBase `json:"hours"`
}

type UpdateRestaurantRequest struct {
	RestaurantBase
	Hours []RestaurantHoursBase `json:"hours"`
}
