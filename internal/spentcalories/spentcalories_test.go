package spentcalories

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SpentCaloriesTestSuite struct {
	suite.Suite
}

func TestSpentCaloriesSuite(t *testing.T) {
	suite.Run(t, new(SpentCaloriesTestSuite))
}

func (suite *SpentCaloriesTestSuite) TestDistance() {
	tests := []struct {
		name   string
		steps  int
		height float64
		want   float64
	}{
		{
			name:   "нормальное количество шагов",
			steps:  1000,
			height: 1.75,
			// Расчет: 1000 * (1.75 * 0.45) / 1000 = 1000 * 0.7875 / 1000 = 0.7875
			want: 0.7875,
		},
		{
			name:   "большое количество шагов",
			steps:  10000,
			height: 1.75,
			// Расчет: 10000 * (1.75 * 0.45) / 1000 = 10000 * 0.7875 / 1000 = 7.875
			want: 7.875,
		},
		{
			name:   "маленькое количество шагов",
			steps:  100,
			height: 1.75,
			// Расчет: 100 * (1.75 * 0.45) / 1000 = 100 * 0.7875 / 1000 = 0.07875
			want: 0.07875,
		},
		{
			name:   "ноль шагов",
			steps:  0,
			height: 1.75,
			want:   0,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			got := distance(tt.steps, tt.height)
			assert.Equal(suite.T(), tt.want, got, "distance() получено: %v, ожидается: %v (шаги: %d, рост: %.2f)",
				got, tt.want, tt.steps, tt.height)
		})
	}
}

func (suite *SpentCaloriesTestSuite) TestMeanSpeed() {
	tests := []struct {
		name     string
		steps    int
		height   float64
		duration time.Duration
		want     float64
	}{
		{
			name:     "нормальная скорость - один час",
			steps:    10000,
			height:   1.75,
			duration: time.Hour,
			// Дистанция = 10000 * (1.75 * 0.45) / 1000 = 7.875 км
			// Скорость = 7.875 / 1 = 7.875 км/ч
			want: 7.875,
		},
		{
			name:     "нормальная скорость - полчаса",
			steps:    10000,
			height:   1.75,
			duration: 30 * time.Minute,
			// Дистанция = 7.875 км
			// Скорость = 7.875 / 0.5 = 15.75 км/ч
			want: 15.75,
		},
		{
			name:     "нормальная скорость - два часа",
			steps:    10000,
			height:   1.75,
			duration: 2 * time.Hour,
			// Дистанция = 7.875 км
			// Скорость = 7.875 / 2 = 3.9375 км/ч
			want: 3.9375,
		},
		{
			name:     "маленькая скорость",
			steps:    1000,
			height:   1.75,
			duration: 2 * time.Hour,
			// Дистанция = 0.7875 км
			// Скорость = 0.7875 / 2 = 0.39375 км/ч
			want: 0.39375,
		},
		{
			name:     "большая скорость",
			steps:    20000,
			height:   1.75,
			duration: time.Hour,
			// Дистанция = 15.75 км
			// Скорость = 15.75 / 1 = 15.75 км/ч
			want: 15.75,
		},
		{
			name:     "нулевая продолжительность",
			steps:    1000,
			height:   1.75,
			duration: 0,
			want:     0,
		},
		{
			name:     "отрицательная продолжительность",
			steps:    1000,
			height:   1.75,
			duration: -time.Hour,
			want:     0,
		},
		{
			name:     "ноль шагов",
			steps:    0,
			height:   1.75,
			duration: time.Hour,
			want:     0,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			got := meanSpeed(tt.steps, tt.height, tt.duration)
			assert.InEpsilon(suite.T(), tt.want, got, 0.001, "meanSpeed() получено: %v, ожидается: %v (шаги: %d, рост: %.2f, продолжительность: %v)",
				got, tt.want, tt.steps, tt.height, tt.duration)
		})
	}
}

func (suite *SpentCaloriesTestSuite) TestRunningSpentCalories() {
	tests := []struct {
		name     string
		steps    int
		weight   float64
		height   float64
		duration time.Duration
		want     float64
		wantErr  bool
	}{
		{
			name:     "нормальная нагрузка - один час",
			steps:    6000,
			weight:   75.0,
			height:   1.75,
			duration: time.Hour,
			// Дистанция = 6000 * (1.75 * 0.45) / 1000 = 4.725 км
			// Скорость = 4.725 / 1 = 4.725 км/ч
			// Калории = (75 * 4.725 * 60) / 60 = 354.375
			want:    354.38,
			wantErr: false,
		},
		{
			name:     "нормальная нагрузка - полчаса",
			steps:    6000,
			weight:   75.0,
			height:   1.75,
			duration: 30 * time.Minute,
			// Калории = (75 * 4.725 * 30) / 60 = 177.1875
			want:    177.19,
			wantErr: false,
		},
		{
			name:     "высокая скорость",
			steps:    20000,
			weight:   75.0,
			height:   1.75,
			duration: time.Hour,
			// Дистанция = 15.75 км, Скорость = 15.75 км/ч
			// Калории = (75 * 15.75 * 60) / 60 = 1181.25
			want:    1181.25,
			wantErr: false,
		},
		{
			name:     "низкая скорость",
			steps:    1000,
			weight:   75.0,
			height:   1.75,
			duration: 2 * time.Hour,
			// Дистанция = 0.7875 км, Скорость = 0.39375 км/ч
			// Калории = (75 * 0.39375 * 120) / 60 = 59.0625
			want:    59.06,
			wantErr: false,
		},
		{
			name:     "другой вес",
			steps:    6000,
			weight:   60.0,
			height:   1.75,
			duration: time.Hour,
			// Калории = (60 * 4.725 * 60) / 60 = 283.5
			want:    283.50,
			wantErr: false,
		},
		{
			name:     "нулевая продолжительность",
			steps:    1000,
			weight:   75.0,
			height:   1.75,
			duration: 0,
			want:     0,
			wantErr:  true,
		},
		{
			name:     "отрицательная продолжительность",
			steps:    1000,
			weight:   75.0,
			height:   1.75,
			duration: -time.Hour,
			want:     0,
			wantErr:  true,
		},
		{
			name:     "ноль шагов",
			steps:    0,
			weight:   75.0,
			height:   1.75,
			duration: time.Hour,
			want:     0,
			wantErr:  true,
		},
		{
			name:     "отрицательные шаги",
			steps:    -1000,
			weight:   75.0,
			height:   1.75,
			duration: time.Hour,
			want:     0,
			wantErr:  true,
		},
		{
			name:     "нулевой вес",
			steps:    1000,
			weight:   0,
			height:   1.75,
			duration: time.Hour,
			want:     0,
			wantErr:  true,
		},
		{
			name:     "отрицательный вес",
			steps:    1000,
			weight:   -75.0,
			height:   1.75,
			duration: time.Hour,
			want:     0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			got, err := RunningSpentCalories(tt.steps, tt.weight, tt.height, tt.duration)

			if tt.wantErr {
				assert.Error(suite.T(), err, "RunningSpentCalories() ожидалась ошибка, но её нет")
			} else {
				assert.NoError(suite.T(), err, "RunningSpentCalories() неожиданная ошибка: %v", err)
			}

			assert.InEpsilon(suite.T(), tt.want, got, 0.01, "RunningSpentCalories() получено: %v, ожидается: %v", got, tt.want)
		})
	}
}

func (suite *SpentCaloriesTestSuite) TestWalkingSpentCalories() {
	tests := []struct {
		name     string
		steps    int
		weight   float64
		height   float64
		duration time.Duration
		want     float64
		wantErr  bool
	}{
		{
			name:     "нормальная нагрузка",
			steps:    6000,
			weight:   75.0,
			height:   1.75,
			duration: time.Hour,
			// Расчет как для бега: 354.375
			// Умножаем на walkingCaloriesCoefficient = 0.5: 354.375 * 0.5 = 177.1875
			want:    177.19,
			wantErr: false,
		},
		{
			name:     "меньше шагов",
			steps:    3000,
			weight:   75.0,
			height:   1.75,
			duration: time.Hour,
			// Калории для бега: 177.1875
			// Умножаем на 0.5: 177.1875 * 0.5 = 88.59375
			want:    88.59,
			wantErr: false,
		},
		{
			name:     "больше шагов",
			steps:    20000,
			weight:   75.0,
			height:   1.75,
			duration: time.Hour,
			// Калории для бега: 1181.25
			// Умножаем на 0.5: 1181.25 * 0.5 = 590.625
			want:    590.62,
			wantErr: false,
		},
		{
			name:     "другой вес",
			steps:    6000,
			weight:   60.0,
			height:   1.75,
			duration: time.Hour,
			// Калории для бега: 283.5
			// Умножаем на 0.5: 283.5 * 0.5 = 141.75
			want:    141.75,
			wantErr: false,
		},
		{
			name:     "другой рост",
			steps:    6000,
			weight:   75.0,
			height:   1.90,
			duration: time.Hour,
			// Дистанция = 6000 * (1.90 * 0.45) / 1000 = 5.13 км
			// Скорость = 5.13 км/ч
			// Калории для бега = (75 * 5.13 * 60) / 60 = 384.75
			// Умножаем на 0.5: 384.75 * 0.5 = 192.375
			want:    192.38,
			wantErr: false,
		},
		{
			name:     "нулевые шаги",
			steps:    0,
			weight:   75.0,
			height:   1.75,
			duration: time.Hour,
			want:     0,
			wantErr:  true,
		},
		{
			name:     "отрицательные шаги",
			steps:    -1000,
			weight:   75.0,
			height:   1.75,
			duration: time.Hour,
			want:     0,
			wantErr:  true,
		},
		{
			name:     "нулевой вес",
			steps:    1000,
			weight:   0,
			height:   1.75,
			duration: time.Hour,
			want:     0,
			wantErr:  true,
		},
		{
			name:     "отрицательный вес",
			steps:    1000,
			weight:   -75.0,
			height:   1.75,
			duration: time.Hour,
			want:     0,
			wantErr:  true,
		},
		{
			name:     "нулевой рост",
			steps:    1000,
			weight:   75.0,
			height:   0,
			duration: time.Hour,
			want:     0,
			wantErr:  true,
		},
		{
			name:     "отрицательный рост",
			steps:    1000,
			weight:   75.0,
			height:   -1.75,
			duration: time.Hour,
			want:     0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			got, err := WalkingSpentCalories(tt.steps, tt.weight, tt.height, tt.duration)

			if tt.wantErr {
				assert.Error(suite.T(), err, "WalkingSpentCalories() ожидалась ошибка, но её нет")
			} else {
				assert.NoError(suite.T(), err, "WalkingSpentCalories() неожиданная ошибка: %v", err)
			}

			assert.InEpsilon(suite.T(), tt.want, got, 0.01, "WalkingSpentCalories() получено: %v, ожидается: %v", got, tt.want)
		})
	}
}

func (suite *SpentCaloriesTestSuite) TestTrainingInfo() {
	tests := []struct {
		name    string
		data    string
		weight  float64
		height  float64
		want    string
		wantErr bool
	}{
		{
			name:   "ходьба - нормальная нагрузка",
			data:   "6000,Ходьба,1h00m",
			weight: 75.0,
			height: 1.75,
			// Дистанция = 4.72 км, Скорость = 4.72 км/ч
			// Калории = 177.19 (как в WalkingSpentCalories)
			want:    "Тип тренировки: Ходьба\nДлительность: 1.00 ч.\nДистанция: 4.72 км.\nСкорость: 4.72 км/ч\nСожгли калорий: 177.19",
			wantErr: false,
		},
		{
			name:   "бег - нормальная нагрузка",
			data:   "6000,Бег,1h00m",
			weight: 75.0,
			height: 1.75,
			// Калории = 354.38 (как в RunningSpentCalories)
			want:    "Тип тренировки: Бег\nДлительность: 1.00 ч.\nДистанция: 4.72 км.\nСкорость: 4.72 км/ч\nСожгли калорий: 354.38",
			wantErr: false,
		},
		{
			name:   "ходьба - высокая скорость",
			data:   "20000,Ходьба,1h00m",
			weight: 75.0,
			height: 1.75,
			// Дистанция = 15.75 км, Скорость = 15.75 км/ч
			// Калории = 590.62
			want:    "Тип тренировки: Ходьба\nДлительность: 1.00 ч.\nДистанция: 15.75 км.\nСкорость: 15.75 км/ч\nСожгли калорий: 590.62",
			wantErr: false,
		},
		{
			name:   "бег - высокая скорость",
			data:   "20000,Бег,1h00m",
			weight: 75.0,
			height: 1.75,
			// Калории = 1181.25
			want:    "Тип тренировки: Бег\nДлительность: 1.00 ч.\nДистанция: 15.75 км.\nСкорость: 15.75 км/ч\nСожгли калорий: 1181.25",
			wantErr: false,
		},
		{
			name:   "ходьба - другой вес и рост",
			data:   "6000,Ходьба,1h00m",
			weight: 60.0,
			height: 1.85,
			// Дистанция = 5.00 км, Скорость = 5.00 км/ч
			// Калории = 141.75
			want:    "Тип тренировки: Ходьба\nДлительность: 1.00 ч.\nДистанция: 5.00 км.\nСкорость: 5.00 км/ч\nСожгли калорий: 141.75",
			wantErr: false,
		},
		{
			name:   "бег - другой вес",
			data:   "6000,Бег,1h00m",
			weight: 60.0,
			height: 1.75,
			// Калории = 283.50
			want:    "Тип тренировки: Бег\nДлительность: 1.00 ч.\nДистанция: 4.72 км.\nСкорость: 4.72 км/ч\nСожгли калорий: 283.50",
			wantErr: false,
		},
		{
			name:   "ходьба - полчаса",
			data:   "6000,Ходьба,30m",
			weight: 75.0,
			height: 1.75,
			// Дистанция = 4.72 км, Скорость = 9.44 км/ч (4.72 / 0.5)
			// Калории = 88.59
			want:    "Тип тренировки: Ходьба\nДлительность: 0.50 ч.\nДистанция: 4.72 км.\nСкорость: 9.44 км/ч\nСожгли калорий: 88.59",
			wantErr: false,
		},
		{
			name:   "бег - полчаса",
			data:   "6000,Бег,30m",
			weight: 75.0,
			height: 1.75,
			// Калории = 177.19
			want:    "Тип тренировки: Бег\nДлительность: 0.50 ч.\nДистанция: 4.72 км.\nСкорость: 9.44 км/ч\nСожгли калорий: 177.19",
			wantErr: false,
		},
		{
			name:    "неизвестный тип тренировки",
			data:    "6000,Плавание,1h00m",
			weight:  75.0,
			height:  1.75,
			want:    "",
			wantErr: true,
		},
		{
			name:    "некорректный формат данных",
			data:    "invalid,data,format",
			weight:  75.0,
			height:  1.75,
			want:    "",
			wantErr: true,
		},
		{
			name:    "некорректное количество шагов",
			data:    "abc,Ходьба,1h00m",
			weight:  75.0,
			height:  1.75,
			want:    "",
			wantErr: true,
		},
		{
			name:    "некорректная продолжительность",
			data:    "6000,Ходьба,invalid",
			weight:  75.0,
			height:  1.75,
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			got, err := TrainingInfo(tt.data, tt.weight, tt.height)

			if tt.wantErr {
				assert.Error(suite.T(), err, "TrainingInfo() ожидалась ошибка, но её нет")
			} else {
				assert.NoError(suite.T(), err, "TrainingInfo() неожиданная ошибка: %v", err)
			}

			assert.Equal(suite.T(), tt.want, got, "TrainingInfo() получено:\n%v\nожидается:\n%v", got, tt.want)
		})
	}
}
