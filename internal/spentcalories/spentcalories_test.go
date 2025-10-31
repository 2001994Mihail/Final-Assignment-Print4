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
			want:   0.7875,
		},
		{
			name:   "большое количество шагов",
			steps:  10000,
			height: 1.75,
			want:   7.875,
		},
		{
			name:   "маленькое количество шагов",
			steps:  100,
			height: 1.75,
			want:   0.07875,
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
			if tt.want == 0 {
				assert.Equal(suite.T(), tt.want, got)
			} else {
				assert.InEpsilon(suite.T(), tt.want, got, 1e-9)
			}
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
			want:     7.875,
		},
		{
			name:     "нормальная скорость - полчаса",
			steps:    10000,
			height:   1.75,
			duration: 30 * time.Minute,
			want:     15.75,
		},
		{
			name:     "нулевая продолжительность",
			steps:    1000,
			height:   1.75,
			duration: 0,
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
			if tt.want == 0 {
				assert.Equal(suite.T(), tt.want, got)
			} else {
				assert.InEpsilon(suite.T(), tt.want, got, 1e-9)
			}
		})
	}
}

func (suite *SpentCaloriesTestSuite) TestRunningSpentCalories_ValidInput() {
	tests := []struct {
		name     string
		steps    int
		weight   float64
		height   float64
		duration time.Duration
		want     float64
	}{
		{
			name:     "нормальная нагрузка - один час",
			steps:    6000,
			weight:   75.0,
			height:   1.75,
			duration: time.Hour,
			want:     354.375,
		},
		{
			name:     "нормальная нагрузка - полчаса",
			steps:    6000,
			weight:   75.0,
			height:   1.75,
			duration: 30 * time.Minute,
			want:     177.1875,
		},
		{
			name:     "высокая скорость",
			steps:    20000,
			weight:   75.0,
			height:   1.75,
			duration: time.Hour,
			want:     1181.25,
		},
		{
			name:     "другой вес",
			steps:    6000,
			weight:   60.0,
			height:   1.75,
			duration: time.Hour,
			want:     283.5,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			got, err := RunningSpentCalories(tt.steps, tt.weight, tt.height, tt.duration)
			assert.NoError(suite.T(), err)
			// Используем более широкий допуск для плавающей точки
			assert.InEpsilon(suite.T(), tt.want, got, 0.001)
		})
	}
}

func (suite *SpentCaloriesTestSuite) TestRunningSpentCalories_InvalidInput() {
	tests := []struct {
		name     string
		steps    int
		weight   float64
		height   float64
		duration time.Duration
	}{
		{
			name:     "нулевая продолжительность",
			steps:    1000,
			weight:   75.0,
			height:   1.75,
			duration: 0,
		},
		{
			name:     "отрицательная продолжительность",
			steps:    1000,
			weight:   75.0,
			height:   1.75,
			duration: -time.Hour,
		},
		{
			name:     "ноль шагов",
			steps:    0,
			weight:   75.0,
			height:   1.75,
			duration: time.Hour,
		},
		{
			name:     "отрицательные шаги",
			steps:    -1000,
			weight:   75.0,
			height:   1.75,
			duration: time.Hour,
		},
		{
			name:     "нулевой вес",
			steps:    1000,
			weight:   0,
			height:   1.75,
			duration: time.Hour,
		},
		{
			name:     "отрицательный вес",
			steps:    1000,
			weight:   -75.0,
			height:   1.75,
			duration: time.Hour,
		},
		{
			name:     "нулевой рост",
			steps:    1000,
			weight:   75.0,
			height:   0,
			duration: time.Hour,
		},
		{
			name:     "отрицательный рост",
			steps:    1000,
			weight:   75.0,
			height:   -1.75,
			duration: time.Hour,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			got, err := RunningSpentCalories(tt.steps, tt.weight, tt.height, tt.duration)
			assert.Error(suite.T(), err)
			assert.Equal(suite.T(), 0.0, got)
		})
	}
}

func (suite *SpentCaloriesTestSuite) TestWalkingSpentCalories_ValidInput() {
	tests := []struct {
		name     string
		steps    int
		weight   float64
		height   float64
		duration time.Duration
		want     float64
	}{
		{
			name:     "нормальная нагрузка",
			steps:    6000,
			weight:   75.0,
			height:   1.75,
			duration: time.Hour,
			// С коэффициентом 0.5 для ходьбы: 354.375 * 0.5 = 177.1875
			want: 177.1875,
		},
		{
			name:     "меньше шагов",
			steps:    3000,
			weight:   75.0,
			height:   1.75,
			duration: time.Hour,
			// 177.1875 * 0.5 = 88.59375
			want: 88.59375,
		},
		{
			name:     "больше шагов",
			steps:    20000,
			weight:   75.0,
			height:   1.75,
			duration: time.Hour,
			// 1181.25 * 0.5 = 590.625
			want: 590.625,
		},
		{
			name:     "другой вес",
			steps:    6000,
			weight:   60.0,
			height:   1.75,
			duration: time.Hour,
			// 283.5 * 0.5 = 141.75
			want: 141.75,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			got, err := WalkingSpentCalories(tt.steps, tt.weight, tt.height, tt.duration)
			assert.NoError(suite.T(), err)
			assert.InEpsilon(suite.T(), tt.want, got, 0.001)
		})
	}
}

func (suite *SpentCaloriesTestSuite) TestWalkingSpentCalories_InvalidInput() {
	invalidCases := []struct {
		name     string
		steps    int
		weight   float64
		height   float64
		duration time.Duration
	}{
		{
			name:     "нулевая продолжительность",
			steps:    1000,
			weight:   75.0,
			height:   1.75,
			duration: 0,
		},
		{
			name:     "ноль шагов",
			steps:    0,
			weight:   75.0,
			height:   1.75,
			duration: time.Hour,
		},
		{
			name:     "нулевой вес",
			steps:    1000,
			weight:   0,
			height:   1.75,
			duration: time.Hour,
		},
	}

	for _, tt := range invalidCases {
		suite.Run(tt.name, func() {
			got, err := WalkingSpentCalories(tt.steps, tt.weight, tt.height, tt.duration)
			assert.Error(suite.T(), err)
			assert.Equal(suite.T(), 0.0, got)
		})
	}
}

func (suite *SpentCaloriesTestSuite) TestTrainingInfo_ValidInput() {
	tests := []struct {
		name   string
		data   string
		weight float64
		height float64
		want   string
	}{
		{
			name:   "ходьба - нормальная нагрузка",
			data:   "6000,Ходьба,1h",
			weight: 75.0,
			height: 1.75,
			// Ходьба с коэффициентом 0.5: 354.375 * 0.5 = 177.1875 ≈ 177.19
			want: "Тип тренировки: Ходьба\nДлительность: 1.00 ч.\nДистанция: 4.72 км.\nСкорость: 4.72 км/ч\nСожгли калорий: 177.19",
		},
		{
			name:   "бег - нормальная нагрузка",
			data:   "6000,Бег,1h",
			weight: 75.0,
			height: 1.75,
			// Бег без коэффициента: 354.375 ≈ 354.38
			want: "Тип тренировки: Бег\nДлительность: 1.00 ч.\nДистанция: 4.72 км.\nСкорость: 4.72 км/ч\nСожгли калорий: 354.38",
		},
		{
			name:   "ходьба - полчаса",
			data:   "6000,Ходьба,30m",
			weight: 75.0,
			height: 1.75,
			// Ходьба 30 минут: 177.1875 * 0.5 = 88.59375 ≈ 88.59
			want: "Тип тренировки: Ходьба\nДлительность: 0.50 ч.\nДистанция: 4.72 км.\nСкорость: 9.45 км/ч\nСожгли калорий: 88.59",
		},
		{
			name:   "бег - полчаса",
			data:   "6000,Бег,30m",
			weight: 75.0,
			height: 1.75,
			// Бег 30 минут: 354.375 * 0.5 = 177.1875 ≈ 177.19
			want: "Тип тренировки: Бег\nДлительность: 0.50 ч.\nДистанция: 4.72 км.\nСкорость: 9.45 км/ч\nСожгли калорий: 177.19",
		},
		{
			name:   "разные регистры активности",
			data:   "6000,ходьба,1h",
			weight: 75.0,
			height: 1.75,
			want:   "Тип тренировки: ходьба\nДлительность: 1.00 ч.\nДистанция: 4.72 км.\nСкорость: 4.72 км/ч\nСожгли калорий: 177.19",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			got, err := TrainingInfo(tt.data, tt.weight, tt.height)
			assert.NoError(suite.T(), err)
			assert.Equal(suite.T(), tt.want, got)
		})
	}
}

func (suite *SpentCaloriesTestSuite) TestTrainingInfo_InvalidInput() {
	tests := []struct {
		name   string
		data   string
		weight float64
		height float64
	}{
		{
			name:   "неизвестный тип тренировки",
			data:   "6000,Плавание,1h",
			weight: 75.0,
			height: 1.75,
		},
		{
			name:   "некорректный формат данных",
			data:   "invalid,data,format",
			weight: 75.0,
			height: 1.75,
		},
		{
			name:   "некорректное количество шагов",
			data:   "abc,Ходьба,1h",
			weight: 75.0,
			height: 1.75,
		},
		{
			name:   "некорректная продолжительность",
			data:   "6000,Ходьба,invalid",
			weight: 75.0,
			height: 1.75,
		},
		{
			name:   "пустые данные",
			data:   "",
			weight: 75.0,
			height: 1.75,
		},
		{
			name:   "недостаточно данных",
			data:   "6000,Ходьба",
			weight: 75.0,
			height: 1.75,
		},
		{
			name:   "лишние данные",
			data:   "6000,Ходьба,1h,extra",
			weight: 75.0,
			height: 1.75,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			got, err := TrainingInfo(tt.data, tt.weight, tt.height)
			assert.Error(suite.T(), err)
			assert.Equal(suite.T(), "", got)
		})
	}
}

func (suite *SpentCaloriesTestSuite) TestParseTraining() {
	tests := []struct {
		name    string
		data    string
		want    int
		wantAct string
		wantDur time.Duration
		wantErr bool
	}{
		{
			name:    "валидные данные",
			data:    "6000,Ходьба,1h30m",
			want:    6000,
			wantAct: "Ходьба",
			wantDur: 90 * time.Minute,
			wantErr: false,
		},
		{
			name:    "отрицательные шаги",
			data:    "-100,Ходьба,1h",
			want:    -100,
			wantAct: "Ходьба",
			wantDur: time.Hour,
			wantErr: true,
		},
		{
			name:    "нулевая продолжительность",
			data:    "1000,Ходьба,0s",
			want:    1000,
			wantAct: "Ходьба",
			wantDur: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			steps, activity, duration, err := parseTraining(tt.data)
			if tt.wantErr {
				assert.Error(suite.T(), err)
			} else {
				assert.NoError(suite.T(), err)
				assert.Equal(suite.T(), tt.want, steps)
				assert.Equal(suite.T(), tt.wantAct, activity)
				assert.Equal(suite.T(), tt.wantDur, duration)
			}
		})
	}
}
