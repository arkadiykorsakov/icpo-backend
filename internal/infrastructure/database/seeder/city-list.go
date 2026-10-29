package seeder

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/jackc/pgx/v5"

	"icpo-backend/seeds"
)

const cityListFile = "cities.txt"

// CityList создаёт сидер городов: читает seeds/cities.txt и вставляет города
// в таблицу cities. Имя сидера — "cities".
//
// Формат файла — секции по странам:
//
//	[Россия]
//	Москва
//	Санкт-Петербург
//
//	[Германия]
//	Берлин
//
// Заголовок секции — название страны из справочника countries (без учёта
// регистра), строки под ним — города этой страны. Пустые строки и строки,
// начинающиеся с «#», игнорируются, пробелы по краям обрезаются.
//
// Если страны из заголовка нет в справочнике, сидер останавливается с ошибкой
// со списком ненайденных стран, поэтому в registry он должен стоять после
// сидера "countries". Дубликаты пропускаются через ON CONFLICT DO NOTHING
// (уникальный индекс cities_name_ci по LOWER(name) и country_id).
func CityList() Seeder {
	return cityList{}
}

type cityList struct{}

func (cityList) Name() string { return "cities" }

func (cityList) Run(ctx context.Context, tx pgx.Tx) (Result, error) {
	data, err := seeds.FS.ReadFile(cityListFile)
	if err != nil {
		return Result{}, fmt.Errorf("read %s: %w", cityListFile, err)
	}
	sections, err := parseCityList(cityListFile, data)
	if err != nil {
		return Result{}, err
	}

	var cityNames, countryNames []string
	for _, s := range sections {
		for _, city := range s.cities {
			cityNames = append(cityNames, city)
			countryNames = append(countryNames, s.country)
		}
	}
	if len(cityNames) == 0 {
		return Result{}, nil
	}

	if err := checkCountriesExist(ctx, tx, sections); err != nil {
		return Result{}, err
	}

	const query = `
		INSERT INTO cities (name, country_id)
		SELECT v.name, co.id
		FROM unnest($1::text[], $2::text[]) AS v(name, country)
		JOIN countries co ON LOWER(co.name) = LOWER(v.country)
		ON CONFLICT DO NOTHING`

	tag, err := tx.Exec(ctx, query, cityNames, countryNames)
	if err != nil {
		return Result{}, err
	}

	inserted := int(tag.RowsAffected())
	return Result{Inserted: inserted, Skipped: len(cityNames) - inserted}, nil
}

// checkCountriesExist проверяет, что все страны из файла есть в таблице countries.
// Без этой проверки JOIN молча отбросил бы города неизвестных стран,
// и они были бы посчитаны как «пропущенные».
func checkCountriesExist(ctx context.Context, tx pgx.Tx, sections []citySection) error {
	countries := make([]string, len(sections))
	firstLine := make(map[string]int, len(sections))
	for i, s := range sections {
		countries[i] = s.country
		if _, ok := firstLine[s.country]; !ok {
			firstLine[s.country] = s.line
		}
	}

	const query = `
		SELECT DISTINCT v.country
		FROM unnest($1::text[]) AS v(country)
		LEFT JOIN countries co ON LOWER(co.name) = LOWER(v.country)
		WHERE co.id IS NULL
		ORDER BY v.country`

	rows, err := tx.Query(ctx, query, countries)
	if err != nil {
		return err
	}
	missing, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return err
	}
	if len(missing) == 0 {
		return nil
	}

	parts := make([]string, len(missing))
	for i, c := range missing {
		parts[i] = fmt.Sprintf("%q (%s:%d)", c, cityListFile, firstLine[c])
	}
	return fmt.Errorf(
		"countries not found in table countries (seed \"countries\" first or fix the name): %s",
		strings.Join(parts, ", "),
	)
}

// citySection — страна из заголовка [Страна] и города под ним.
type citySection struct {
	country string
	line    int // номер строки заголовка, для сообщений об ошибках
	cities  []string
}

// parseCityList разбирает содержимое файла со списком городов по странам.
func parseCityList(file string, data []byte) ([]citySection, error) {
	data = bytes.TrimPrefix(data, []byte("\xef\xbb\xbf")) // BOM

	var (
		sections []citySection
		cur      = -1 // индекс текущей секции в sections
	)

	sc := bufio.NewScanner(bytes.NewReader(data))
	for line := 1; sc.Scan(); line++ {
		v := strings.TrimSpace(sc.Text())
		if v == "" || strings.HasPrefix(v, "#") {
			continue
		}

		if strings.HasPrefix(v, "[") {
			if !strings.HasSuffix(v, "]") {
				return nil, fmt.Errorf("%s:%d: section header must look like [Country]", file, line)
			}
			country := strings.TrimSpace(v[1 : len(v)-1])
			if country == "" {
				return nil, fmt.Errorf("%s:%d: empty country name in section header", file, line)
			}
			if utf8.RuneCountInString(country) > maxNameLen {
				return nil, fmt.Errorf("%s:%d: country name is longer than %d characters", file, line, maxNameLen)
			}
			sections = append(sections, citySection{country: country, line: line})
			cur = len(sections) - 1
			continue
		}

		if cur < 0 {
			return nil, fmt.Errorf("%s:%d: city %q is before the first [Country] header", file, line, v)
		}
		if utf8.RuneCountInString(v) > maxNameLen {
			return nil, fmt.Errorf("%s:%d: value is longer than %d characters", file, line, maxNameLen)
		}
		sections[cur].cities = append(sections[cur].cities, v)
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("scan %s: %w", file, err)
	}
	return sections, nil
}
