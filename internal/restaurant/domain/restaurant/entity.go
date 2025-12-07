package restaurant

import (
	"go-ai/pkg/utils"
	"strings"

	"github.com/google/uuid"
)

type Entity struct {
	ID          uuid.UUID
	Name        string
	Description string
	Address     string
	Category    string
	City        string
	District    string
	LogoUrl     Url
	BannerUrl   Url
	PhoneNumber Phone
	WebsiteUrl  Url
	Email       utils.Email
	CreatedBy   uuid.UUID
	UpdateBy    uuid.UUID
	Hours       []Hours
}

type Hours struct {
	Day       DayOfWeek
	OpenTime  string
	CloseTime string
}

func NewEntity(
	name string,
	description string,
	address string,
	category string,
	city string,
	district string,
	logoUrl Url,
	bannerUrl Url,
	phone Phone,
	website Url,
	email utils.Email,
	userID uuid.UUID,
	hours []Hours,
) (*Entity, error) {

	// Required primitive validations
	if strings.TrimSpace(name) == "" {
		return nil, ErrNameRequired
	}
	if strings.TrimSpace(description) == "" {
		return nil, ErrDescriptionRequired
	}
	if strings.TrimSpace(address) == "" {
		return nil, ErrAddressRequired
	}
	if strings.TrimSpace(category) == "" {
		return nil, ErrCategoryRequired
	}
	if strings.TrimSpace(city) == "" {
		return nil, ErrCityRequired
	}
	if strings.TrimSpace(district) == "" {
		return nil, ErrDistrictRequired
	}
	if userID == uuid.Nil {
		return nil, ErrUserIDRequired
	}

	// Validate hours
	if err := ValidateHours(hours); err != nil {
		return nil, err
	}

	// Build entity
	return &Entity{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		Address:     address,
		Category:    category,
		City:        city,
		District:    district,
		LogoUrl:     logoUrl,
		BannerUrl:   bannerUrl,
		PhoneNumber: phone,
		WebsiteUrl:  website,
		Email:       email,
		CreatedBy:   userID,
		UpdateBy:    userID,
		Hours:       hours,
	}, nil
}

func (e *Entity) Validate() error {
	if strings.TrimSpace(e.Name) == "" {
		return ErrNameRequired
	}
	if strings.TrimSpace(e.Address) == "" {
		return ErrAddressRequired
	}
	if err := ValidateHours(e.Hours); err != nil {
		return err
	}
	return nil
}
