package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log/slog"
	"meet/internal/config"
	"meet/internal/pkg/app/helper"
	"meet/internal/pkg/app/model"
	"meet/internal/pkg/app/repository"
	"meet/internal/pkg/app/repository/transaction"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/guregu/null"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/schollz/progressbar/v3"
)

type countryDTO struct {
	*model.Country
	Subregion   string        `json:"subregion"`
	SubregionID null.Int      `json:"subregion_id"`
	Cities      []*model.City `json:"cities"`
}

type regionsDTO []*model.Region

func (r *regionsDTO) AppendUnique(id int, name string) {
	for _, region := range *r {
		if id == region.ID {
			return
		}
	}

	*r = append(*r, &model.Region{ID: id, Name: name})
}

var (
	regionRepository  repository.RegionRepository
	countryRepository repository.CountryRepository
	cityRepository    repository.CityRepository
)

func main() {
	_, b, _, _ := runtime.Caller(0)
	rootDir, _ := filepath.Abs(filepath.Dir(b) + "/../..")

	err := godotenv.Load(rootDir + "/.env")
	if err != nil {
		slog.Warn(err.Error())
		return
	}

	cfg := config.FromEnv(rootDir)

	logF, err := helper.OpenLogFile(rootDir)
	if err != nil {
		slog.Warn(err.Error())
		return
	}
	defer logF.Close()
	helper.ConfigureSlogger(cfg.Debug, logF)

	db, err := helper.LoadDB(cfg.DB)
	if err != nil {
		slog.Warn(err.Error())
		return
	}

	regionRepository = repository.NewRegionDBRepository(db)
	countryRepository = repository.NewCountryDBRepository(db)
	cityRepository = repository.NewCityDBRepository(db)

	f, err := os.Open(rootDir + "/cities.json")
	if err != nil {
		slog.Warn("Не удалось открыть файл", "err", err)
		return
	}
	defer f.Close()

	fillDB(f, db)
}

func fillDB(r io.Reader, db *sql.DB) {
	regions, countries, cities, err := parseJSON(r)
	if err != nil {
		slog.Warn(err.Error())
		return
	}

	ctx, tx, err := transaction.BeginTx(context.Background(), db)
	if err != nil {
		slog.Warn(err.Error())
		return
	}

	slog.Info("Добавление регионов в базу данных...")
	for _, region := range regions {
		if err := regionRepository.Add(ctx, region); err != nil {
			tx.Rollback()
			slog.Warn(err.Error())
			return
		}
	}

	slog.Info("Добавление стран в базу данных...")
	for _, country := range countries {
		if err := countryRepository.Add(ctx, country); err != nil {
			tx.Rollback()
			slog.Warn(err.Error())
			return
		}
	}

	citiesLen := len(cities)
	n := runtime.GOMAXPROCS(0)
	partSize := citiesLen / n
	remainder := citiesLen % n

	bar := progressbar.Default(int64(citiesLen), "Добавление городов в базу данных...")
	wg := &sync.WaitGroup{}
	for i := 0; i < n; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			start, end := partSize*i, partSize*(i+1)
			if i == n-1 {
				end += remainder
			}

			ch := make(chan bool)
			go insertCities(ctx, tx, cities[start:end], ch)

			for range ch {
				bar.Add(1)
			}
		}(i)

	}
	wg.Wait()

	if err := tx.Commit(); err != nil {
		slog.Warn(err.Error())
		return
	}

	slog.Info("База данных успешно заполнена!")
}

func parseJSON(r io.Reader) (regionsDTO, []*model.Country, []*model.City, error) {
	data := make([]*countryDTO, 0)

	regions := make(regionsDTO, 0)
	countries := make([]*model.Country, 0)
	cities := make([]*model.City, 0)

	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&data); err != nil {
		return nil, nil, nil, err
	}

	for _, country := range data {
		if !country.SubregionID.IsZero() {
			regions.AppendUnique(int(country.SubregionID.Int64), country.Subregion)
		}

		country.Country.RegionID = country.SubregionID
		countries = append(countries, country.Country)

		for _, city := range country.Cities {
			city.CountryID = country.ID
			cities = append(cities, city)
		}
	}

	return regions, countries, cities, nil
}

func insertCities(ctx context.Context, tx *sql.Tx, cities []*model.City, ch chan<- bool) {
	defer close(ch)

	for _, city := range cities {
		if err := cityRepository.Add(ctx, city); err != nil {
			tx.Rollback()
			slog.Warn(err.Error())
			return
		}

		ch <- true
	}
}
