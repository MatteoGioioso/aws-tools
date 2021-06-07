package libs

import (
	"github.com/onsi/gomega"
	"testing"
	"time"
)

func mockNow() time.Time {
	return time.Date(2021, 6, 5, 5, 11, 0, 0, time.UTC)
}

var (
	insideOfficeHourAndDay     = time.Date(2021, 6, 3, 6, 11, 0, 0, time.UTC)
	outsideOfficeHourAndDay    = time.Date(2021, 6, 5, 1, 11, 0, 0, time.UTC)
	insideOfficeHourButNotDay  = time.Date(2021, 6, 5, 7, 11, 0, 0, time.UTC)
	outsideOfficeHourButNotDay = time.Date(2021, 6, 2, 22, 11, 0, 0, time.UTC)
)

func TestSchedulerConfig_getCurrentTimeFromTZ(t *testing.T) {
	g := gomega.NewWithT(t)

	t.Run("should convert to the given time zone", func(t *testing.T) {
		config := &SchedulerConfig{
			Period:   Period{},
			TimeZone: "Europe/Helsinki",
		}
		configClient := SchedulerConfigClient{
			now:    mockNow,
			Config: config,
		}

		got, err := configClient.getCurrentTimeFromTZ()
		if err != nil {
			t.Error(err)
		}

		g.Expect(got.Hour()).To(gomega.Equal(mockNow().Hour() + 3))
	})
}

func TestSchedulerConfig_shouldWakeup(t *testing.T) {
	g := gomega.NewWithT(t)

	tests := []struct {
		wants    bool
		mockTime func() time.Time
	}{
		{
			wants: true,
			mockTime: func() time.Time {
				return insideOfficeHourAndDay
			},
		},
		{
			wants: false,
			mockTime: func() time.Time {
				return outsideOfficeHourAndDay
			},
		},
		{
			wants: false,
			mockTime: func() time.Time {
				return insideOfficeHourButNotDay
			},
		},
		{
			wants: false,
			mockTime: func() time.Time {
				return outsideOfficeHourButNotDay
			},
		},
	}

	for _, tt := range tests {
		t.Run("should establish whether the resources should be awaken", func(t *testing.T) {
			config := &SchedulerConfig{
				Period: Period{
					Pattern: "office_hours",
				},
				TimeZone: "Europe/Helsinki",
			}
			configClient := SchedulerConfigClient{
				now:    tt.mockTime,
				Config: config,
			}

			got, err := configClient.ShouldWakeup()
			if err != nil {
				t.Error(err)
			}

			g.Expect(got).To(gomega.Equal(tt.wants))
		})
	}

	for _, tt := range tests {
		t.Run("should establish whether the resources should be awaken", func(t *testing.T) {
			config := &SchedulerConfig{
				Period: Period{
					Pattern: "permanent_shutdown",
				},
				TimeZone: "Europe/Helsinki",
			}
			configClient := SchedulerConfigClient{
				now:    tt.mockTime,
				Config: config,
			}

			got, err := configClient.ShouldWakeup()
			if err != nil {
				t.Error(err)
			}

			// This is always false
			g.Expect(got).To(gomega.Equal(false))
		})
	}

}
