package daysteps

import (
	"bytes"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type DayStepsTestSuite struct {
	suite.Suite
}

func TestDayStepsSuite(t *testing.T) {
	suite.Run(t, new(DayStepsTestSuite))
}

func (suite *DayStepsTestSuite) TestParsePackage() {
	tests := []struct {
		name         string
		input        string
		wantSteps    int
		wantDuration time.Duration
		wantErr      bool
	}{
		// Корректные значения
		{
			name:         "корректный ввод",
			input:        "678,0h50m",
			wantSteps:    678,
			wantDuration: 50 * time.Minute,
			wantErr:      false,
		},
		{
			name:         "корректный ввод с часами и минутами",
			input:        "1000,1h30m",
			wantSteps:    1000,
			wantDuration: 90 * time.Minute,
			wantErr:      false,
		},
		{
			name:         "положительное число с плюсом",
			input:        "+12345,1h30m",
			wantSteps:    12345,
			wantDuration: 90 * time.Minute,
			wantErr:      false,
		},
		// Корректные значения продолжительности
		{
			name:         "продолжительность - только минуты",
			input:        "1000,30m",
			wantSteps:    1000,
			wantDuration: 30 * time.Minute,
			wantErr:      false,
		},
		{
			name:         "продолжительность - только часы",
			input:        "1000,2h",
			wantSteps:    1000,
			wantDuration: 2 * time.Hour,
			wantErr:      false,
		},
		{
			name:         "продолжительность - дробные часы",
			input:        "1000,1.5h",
			wantSteps:    1000,
			wantDuration: 90 * time.Minute,
			wantErr:      false,
		},
		{
			name:         "продолжительность - дробные минуты",
			input:        "1000,30.5m",
			wantSteps:    1000,
			wantDuration: 30*time.Minute + 30*time.Second,
			wantErr:      false,
		},
		// Ошибки формата
		{
			name:         "неверный формат - неправильное количество параметров",
			input:        "678",
			wantSteps:    0,
			wantDuration: 0,
			wantErr:      true,
		},
		{
			name:         "неверный формат - три параметра",
			input:        "678,1h30m,extra",
			wantSteps:    0,
			wantDuration: 0,
			wantErr:      true,
		},
		{
			name:         "пустой ввод",
			input:        "",
			wantSteps:    0,
			wantDuration: 0,
			wantErr:      true,
		},
		// Ошибки в количестве шагов
		{
			name:         "неверные шаги - не числовое значение",
			input:        "abc,1h30m",
			wantSteps:    0,
			wantDuration: 0,
			wantErr:      true,
		},
		{
			name:         "неверные шаги - ноль",
			input:        "0,1h30m",
			wantSteps:    0,
			wantDuration: 0,
			wantErr:      true,
		},
		{
			name:         "неверные шаги - отрицательное значение",
			input:        "-100,1h30m",
			wantSteps:    0,
			wantDuration: 0,
			wantErr:      true,
		},
		{
			name:         "неверные шаги - только знак минус",
			input:        "-,1h30m",
			wantSteps:    0,
			wantDuration: 0,
			wantErr:      true,
		},
		{
			name:         "неверные шаги - только знак плюс",
			input:        "+,1h30m",
			wantSteps:    0,
			wantDuration: 0,
			wantErr:      true,
		},
		{
			name:         "неверные шаги - пробелы в начале",
			input:        " 12345,1h30m",
			wantSteps:    0,
			wantDuration: 0,
			wantErr:      true,
		},
		{
			name:         "неверные шаги - пробелы в конце",
			input:        "12345 ,1h30m",
			wantSteps:    0,
			wantDuration: 0,
			wantErr:      true,
		},
		{
			name:         "неверные шаги - некорректные символы",
			input:        "123abc,1h30m",
			wantSteps:    0,
			wantDuration: 0,
			wantErr:      true,
		},
		// Ошибки в продолжительности
		{
			name:         "неверный формат продолжительности",
			input:        "678,invalid",
			wantSteps:    0,
			wantDuration: 0,
			wantErr:      true,
		},
		{
			name:         "неверная продолжительность - ноль",
			input:        "678,0h0m",
			wantSteps:    0,
			wantDuration: 0,
			wantErr:      true,
		},
		{
			name:         "неверная продолжительность - отрицательное значение",
			input:        "678,-1h30m",
			wantSteps:    0,
			wantDuration: 0,
			wantErr:      true,
		},
		{
			name:         "неверная продолжительность - отрицательные минуты",
			input:        "678,1h-30m",
			wantSteps:    0,
			wantDuration: 0,
			wantErr:      true,
		},
		{
			name:         "неверная продолжительность - неверная единица измерения",
			input:        "678,1.5d",
			wantSteps:    0,
			wantDuration: 0,
			wantErr:      true,
		},
		{
			name:         "неверная продолжительность - пробел между числом и единицей",
			input:        "678,1 h30m",
			wantSteps:    0,
			wantDuration: 0,
			wantErr:      true,
		},
		{
			name:         "неверная продолжительность - пропущена единица измерения",
			input:        "678,30",
			wantSteps:    0,
			wantDuration: 0,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			gotSteps, gotDuration, err := parsePackage(tt.input)

			if tt.wantErr {
				assert.Error(suite.T(), err, "parsePackage() для строки данных %q ожидалась ошибка, но её нет", tt.input)
			} else {
				assert.NoError(suite.T(), err, "parsePackage() неожиданная ошибка для строки данных %q: %v", tt.input, err)
			}

			assert.Equal(suite.T(), tt.wantSteps, gotSteps, "parsePackage() полученное количество шагов: %v, ожидается %v", gotSteps, tt.wantSteps)
			assert.Equal(suite.T(), tt.wantDuration, gotDuration, "parsePackage() полученная продолжительность прогулки: %v, ожидается %v", gotDuration, tt.wantDuration)
		})
	}
}

func (suite *DayStepsTestSuite) TestDayActionInfo() {
	var buf bytes.Buffer
	log.SetOutput(&buf)

	defer log.SetOutput(os.Stderr)

	tests := []struct {
		name          string
		input         string
		weight        float64
		height        float64
		want          string
		wantLogOutput bool
	}{
		{
			name:   "нормальная нагрузка - один час",
			input:  "6000,1h",
			weight: 75.0,
			height: 1.75,
			// Расчет по формуле из spentcalories:
			// Дистанция = 6000 * (1.75 * 0.45) / 1000 = 4.725 км ≈ 4.72 км
			// Скорость = 4.725 / 1 = 4.725 км/ч
			// Калории = (75 * 4.725 * 60) / 60 * 0.5 = 177.1875 ≈ 177.19
			want:          "Количество шагов: 6000.\nДистанция составила 4.72 км.\nВы сожгли 177.19 ккал.",
			wantLogOutput: false,
		},
		{
			name:   "нормальная нагрузка - полчаса",
			input:  "3000,30m",
			weight: 75.0,
			height: 1.75,
			// Дистанция = 3000 * (1.75 * 0.45) / 1000 = 2.3625 км ≈ 2.36 км
			// Скорость = 2.3625 / 0.5 = 4.725 км/ч
			// Калории = (75 * 4.725 * 30) / 60 * 0.5 = 88.59375 ≈ 88.59
			want:          "Количество шагов: 3000.\nДистанция составила 2.36 км.\nВы сожгли 88.59 ккал.",
			wantLogOutput: false,
		},
		{
			name:   "высокая нагрузка",
			input:  "20000,1h",
			weight: 75.0,
			height: 1.75,
			// Дистанция = 20000 * (1.75 * 0.45) / 1000 = 15.75 км
			// Скорость = 15.75 / 1 = 15.75 км/ч
			// Калории = (75 * 15.75 * 60) / 60 * 0.5 = 590.625 ≈ 590.62
			want:          "Количество шагов: 20000.\nДистанция составила 15.75 км.\nВы сожгли 590.62 ккал.",
			wantLogOutput: false,
		},
		{
			name:   "низкая нагрузка",
			input:  "1000,2h",
			weight: 75.0,
			height: 1.75,
			// Дистанция = 1000 * (1.75 * 0.45) / 1000 = 0.7875 км ≈ 0.79 км
			// Скорость = 0.7875 / 2 = 0.39375 км/ч
			// Калории = (75 * 0.39375 * 120) / 60 * 0.5 = 29.53125 ≈ 29.53
			want:          "Количество шагов: 1000.\nДистанция составила 0.79 км.\nВы сожгли 29.53 ккал.",
			wantLogOutput: false,
		},
		{
			name:   "другой вес и рост",
			input:  "6000,1h",
			weight: 60.0,
			height: 1.85,
			// Дистанция = 6000 * (1.85 * 0.45) / 1000 = 4.995 км ≈ 5.00 км
			// Скорость = 4.995 / 1 = 4.995 км/ч
			// Калории = (60 * 4.995 * 60) / 60 * 0.5 = 149.85 ≈ 149.85
			want:          "Количество шагов: 6000.\nДистанция составила 5.00 км.\nВы сожгли 149.85 ккал.",
			wantLogOutput: false,
		},
		{
			name:          "некорректный формат",
			input:         "not valid",
			weight:        75.0,
			height:        1.75,
			want:          "",
			wantLogOutput: true,
		},
		{
			name:          "пустая строка",
			input:         "",
			weight:        75.0,
			height:        1.75,
			want:          "",
			wantLogOutput: true,
		},
		{
			name:          "отрицательные шаги",
			input:         "-1000,1h",
			weight:        75.0,
			height:        1.75,
			want:          "",
			wantLogOutput: true,
		},
		{
			name:          "ноль шагов",
			input:         "0,1h",
			weight:        75.0,
			height:        1.75,
			want:          "",
			wantLogOutput: true,
		},
		{
			name:          "отрицательная продолжительность",
			input:         "1000,-1h",
			weight:        75.0,
			height:        1.75,
			want:          "",
			wantLogOutput: true,
		},
		{
			name:          "нулевая продолжительность",
			input:         "1000,0s",
			weight:        75.0,
			height:        1.75,
			want:          "",
			wantLogOutput: true,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			buf.Reset()

			got := DayActionInfo(tt.input, tt.weight, tt.height)

			assert.Equal(suite.T(), tt.want, got, "\nDayActionInfo() получено:\n%v\nожидается:\n%v\n(ввод: %q, вес: %.1f, рост: %.2f)",
				got, tt.want, tt.input, tt.weight, tt.height)

			if tt.wantLogOutput {
				assert.NotEmpty(suite.T(), buf.String(), "Ожидался вывод в лог, но его нет")
			} else {
				assert.Empty(suite.T(), buf.String(), "Неожиданный вывод в лог: %v", buf.String())
			}
		})
	}
}

// Дополнительные тесты для проверки точности вычислений
func (suite *DayStepsTestSuite) TestDayActionInfoPrecision() {
	tests := []struct {
		name   string
		input  string
		weight float64
		height float64
	}{
		{
			name:   "точность округления дистанции",
			input:  "1234,1h",
			weight: 70.0,
			height: 1.70,
		},
		{
			name:   "точность округления калорий",
			input:  "5678,45m",
			weight: 65.5,
			height: 1.80,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			result := DayActionInfo(tt.input, tt.weight, tt.height)

			// Проверяем, что результат не пустой и содержит ожидаемые части
			assert.NotEmpty(suite.T(), result)
			assert.Contains(suite.T(), result, "Количество шагов:")
			assert.Contains(suite.T(), result, "Дистанция составила")
			assert.Contains(suite.T(), result, "Вы сожгли")

			// Проверяем формат вывода
			assert.Regexp(suite.T(), `Количество шагов: \d+\.\nДистанция составила \d+\.\d{2} км\.\nВы сожгли \d+\.\d{2} ккал\.`, result)
		})
	}
}

// Бенчмарк-тест для измерения производительности
func BenchmarkDayActionInfo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DayActionInfo("6000,1h", 75.0, 1.75)
	}
}
