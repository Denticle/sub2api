package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type publicAnnouncementRepoStub struct {
	announcementRepoStub
	active []Announcement
}

func (s *publicAnnouncementRepoStub) ListActive(context.Context, time.Time) ([]Announcement, error) {
	return s.active, nil
}

func TestAnnouncementServiceListPublicHidesTargetedAndInactive(t *testing.T) {
	past := time.Now().Add(-time.Hour)
	future := time.Now().Add(time.Hour)
	targeted := AnnouncementTargeting{AnyOf: []AnnouncementConditionGroup{{
		AllOf: []AnnouncementCondition{{Type: AnnouncementConditionTypeBalance, Operator: AnnouncementOperatorGT, Value: 10}},
	}}}
	repo := &publicAnnouncementRepoStub{active: []Announcement{
		{ID: 1, Title: "旧公告", Status: AnnouncementStatusActive},
		{ID: 2, Title: "定向公告", Status: AnnouncementStatusActive, Targeting: targeted},
		{ID: 3, Title: "已过期", Status: AnnouncementStatusActive, EndsAt: &past},
		{ID: 4, Title: "未开始", Status: AnnouncementStatusActive, StartsAt: &future},
		{ID: 5, Title: "新公告", Status: AnnouncementStatusActive},
	}}
	svc := NewAnnouncementService(repo, nil, nil, nil)

	got, err := svc.ListPublic(context.Background(), 5)
	require.NoError(t, err)
	require.Len(t, got, 2)
	require.Equal(t, int64(5), got[0].ID, "newest first")
	require.Equal(t, int64(1), got[1].ID)

	limited, err := svc.ListPublic(context.Background(), 1)
	require.NoError(t, err)
	require.Len(t, limited, 1)
	require.Equal(t, int64(5), limited[0].ID)
}
