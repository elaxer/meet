package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
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

var (
	_, b, _, _ = runtime.Caller(0)
	rootDir, _ = filepath.Abs(filepath.Dir(b) + "/../..")
)

var db *sql.DB
var cfg *config.Config

var (
	regionRepository  repository.RegionRepository
	countryRepository repository.CountryRepository
	cityRepository    repository.CityRepository
)

func init() {
	err := godotenv.Load(rootDir + "/.env")
	if err != nil {
		panic(err)
	}

	cfg = config.NewConfig(rootDir)

	db, err = helper.LoadDB(cfg.DBConfig)
	if err != nil {
		panic(err)
	}

	regionRepository = repository.NewRegionDBRepository(db)
	countryRepository = repository.NewCountryDBRepository(db)
	cityRepository = repository.NewCityDBRepository(db)
}

type countryDTO struct {
	*model.Country
	Subregion   string        `json:"subregion"`
	SubregionID null.Int      `json:"subregion_id"`
	Cities      []*model.City `json:"cities"`
}

type Regions []*model.Region

func (r *Regions) AppendUnique(id int, name string) {
	for _, region := range *r {
		if id == region.ID {
			return
		}
	}

	*r = append(*r, &model.Region{ID: id, Name: name})
}

func main() {
	f, err := os.Open(rootDir + "/cities.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var data []*countryDTO

	var regions Regions
	var countries []*model.Country
	var cities []*model.City

	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&data); err != nil {
		log.Fatal(err)
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

	ctx, tx, err := transaction.BeginTx(context.Background(), db)

	log.Println("Добавление регионов в базу данных...")
	for _, region := range regions {
		if err := regionRepository.Add(ctx, region); err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
	}

	log.Println("Добавление стран в базу данных...")
	for _, country := range countries {
		if err := countryRepository.Add(ctx, country); err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
	}

	citiesLen := len(cities)
	n := runtime.GOMAXPROCS(0)
	partSize := citiesLen / n
	remainder := citiesLen % n

	if err != nil {
		log.Fatal(err)
	}

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
		log.Fatal(err)
	}

	log.Println("База данных успешно заполнена данными!")
}

func insertCities(ctx context.Context, tx *sql.Tx, cities []*model.City, ch chan<- bool) {
	defer close(ch)

	for _, city := range cities {
		if err := cityRepository.Add(ctx, city); err != nil {
			tx.Rollback()
			log.Fatal(err)
		}

		ch <- true
	}
}
