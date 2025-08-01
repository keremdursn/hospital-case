package usecase

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
	"github.com/keremdursn/hospital-case/internal/dto"
	"github.com/keremdursn/hospital-case/internal/repository"
)

type LocationUsecase interface {
	ListAllCities() ([]dto.CityLookup, error)
	ListDistrictsByCity(cityID uint) ([]dto.DistrictLookup, error)
}

type locationUsecase struct {
	repo  repository.LocationRepository
	redis *redis.Client
}

func NewLocationUsecase(repo repository.LocationRepository, redis *redis.Client) LocationUsecase {
	return &locationUsecase{
		repo:  repo,
		redis: redis,
	}
}

func (u *locationUsecase) ListAllCities() ([]dto.CityLookup, error) {
	ctx := context.Background()
	cacheKey := "cities"

	cached, err := u.redis.Get(ctx, cacheKey).Result()
	if err == nil && cached != "" {
		var resp []dto.CityLookup
		if err := json.Unmarshal([]byte(cached), &resp); err == nil {
			return resp, nil
		}
	}

	cities, err := u.repo.GetAllCities()
	if err != nil {
		return nil, err
	}

	resp := make([]dto.CityLookup, 0, len(cities))
	for _, c := range cities {
		resp = append(resp, dto.CityLookup{
			ID:   c.ID,
			Name: c.Name,
		})
	}

	if data, err := json.Marshal(resp); err == nil {
		_ = u.redis.Set(ctx, cacheKey, data, 0).Err()
	}

	return resp, nil
}

func (u *locationUsecase) ListDistrictsByCity(cityID uint) ([]dto.DistrictLookup, error) {
	ctx := context.Background()
	cacheKey := "districts_by_city_" + string(rune(cityID))

	cached, err := u.redis.Get(ctx, cacheKey).Result()
	if err == nil && cached != "" {
		var resp []dto.DistrictLookup
		if err := json.Unmarshal([]byte(cached), &resp); err == nil {
			return resp, nil
		}
	}

	districts, err := u.repo.GetDistrictsByCity(cityID)
	if err != nil {
		return nil, err
	}

	resp := make([]dto.DistrictLookup, 0, len(districts))
	for _, d := range districts {
		resp = append(resp, dto.DistrictLookup{
			ID:     d.ID,
			Name:   d.Name,
			CityID: d.CityID,
		})
	}

	if data, err := json.Marshal(resp); err == nil {
		_ = u.redis.Set(ctx, cacheKey, data, 0).Err()
	}

	return resp, nil
}
