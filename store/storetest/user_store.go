// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package storetest

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/store"
)

func TestUserStore(t *testing.T, ss store.Store) {
	result := <-ss.User().GetAll()
	require.Nil(t, result.Err, "failed cleaning up test users")
	users := result.Data.([]*model.User)

	for _, u := range users {
		result := <-ss.User().PermanentDelete(u.Id)
		require.Nil(t, result.Err, "failed cleaning up test user %s", u.Username)
	}

	// t.Run("Save", func(t *testing.T) { testUserStoreSave(t, ss) })
	// t.Run("Update", func(t *testing.T) { testUserStoreUpdate(t, ss) })
	// t.Run("UpdateUpdateAt", func(t *testing.T) { testUserStoreUpdateUpdateAt(t, ss) })
	// t.Run("UpdateFailedPasswordAttempts", func(t *testing.T) { testUserStoreUpdateFailedPasswordAttempts(t, ss) })
	// t.Run("Get", func(t *testing.T) { testUserStoreGet(t, ss) })
	// t.Run("UserCount", func(t *testing.T) { testUserCount(t, ss) })
	// t.Run("GetAllUsingAuthService", func(t *testing.T) { testGetAllUsingAuthService(t, ss) })
	// t.Run("GetAllProfiles", func(t *testing.T) { testUserStoreGetAllProfiles(t, ss) })
	// t.Run("GetProfiles", func(t *testing.T) { testUserStoreGetProfiles(t, ss) })
	// t.Run("GetProfilesInChannel", func(t *testing.T) { testUserStoreGetProfilesInChannel(t, ss) })
	// t.Run("GetProfilesInChannelByStatus", func(t *testing.T) { testUserStoreGetProfilesInChannelByStatus(t, ss) })
	// t.Run("GetProfilesWithoutTeam", func(t *testing.T) { testUserStoreGetProfilesWithoutTeam(t, ss) })
	// t.Run("GetAllProfilesInChannel", func(t *testing.T) { testUserStoreGetAllProfilesInChannel(t, ss) })
	// t.Run("GetProfilesNotInChannel", func(t *testing.T) { testUserStoreGetProfilesNotInChannel(t, ss) })
	// t.Run("GetProfilesByIds", func(t *testing.T) { testUserStoreGetProfilesByIds(t, ss) })
	// t.Run("GetProfilesByUsernames", func(t *testing.T) { testUserStoreGetProfilesByUsernames(t, ss) })
	// t.Run("GetSystemAdminProfiles", func(t *testing.T) { testUserStoreGetSystemAdminProfiles(t, ss) })
	t.Run("GetByEmail", func(t *testing.T) { testUserStoreGetByEmail(t, ss) })
	// t.Run("GetByAuthData", func(t *testing.T) { testUserStoreGetByAuthData(t, ss) })
	// t.Run("GetByUsername", func(t *testing.T) { testUserStoreGetByUsername(t, ss) })
	// t.Run("GetForLogin", func(t *testing.T) { testUserStoreGetForLogin(t, ss) })
	// t.Run("UpdatePassword", func(t *testing.T) { testUserStoreUpdatePassword(t, ss) })
	// t.Run("Delete", func(t *testing.T) { testUserStoreDelete(t, ss) })
	// t.Run("UpdateAuthData", func(t *testing.T) { testUserStoreUpdateAuthData(t, ss) })
	// t.Run("UserUnreadCount", func(t *testing.T) { testUserUnreadCount(t, ss) })
	// t.Run("UpdateMfaSecret", func(t *testing.T) { testUserStoreUpdateMfaSecret(t, ss) })
	// t.Run("UpdateMfaActive", func(t *testing.T) { testUserStoreUpdateMfaActive(t, ss) })
	// t.Run("GetRecentlyActiveUsersForTeam", func(t *testing.T) { testUserStoreGetRecentlyActiveUsersForTeam(t, ss) })
	// t.Run("GetNewUsersForTeam", func(t *testing.T) { testUserStoreGetNewUsersForTeam(t, ss) })
	// t.Run("Search", func(t *testing.T) { testUserStoreSearch(t, ss) })
	// t.Run("SearchNotInChannel", func(t *testing.T) { testUserStoreSearchNotInChannel(t, ss) })
	// t.Run("SearchInChannel", func(t *testing.T) { testUserStoreSearchInChannel(t, ss) })
	// t.Run("SearchNotInTeam", func(t *testing.T) { testUserStoreSearchNotInTeam(t, ss) })
	// t.Run("SearchWithoutTeam", func(t *testing.T) { testUserStoreSearchWithoutTeam(t, ss) })
	// t.Run("AnalyticsGetInactiveUsersCount", func(t *testing.T) { testUserStoreAnalyticsGetInactiveUsersCount(t, ss) })
	// t.Run("AnalyticsGetSystemAdminCount", func(t *testing.T) { testUserStoreAnalyticsGetSystemAdminCount(t, ss) })
	// t.Run("GetProfilesNotInTeam", func(t *testing.T) { testUserStoreGetProfilesNotInTeam(t, ss) })
	// t.Run("ClearAllCustomRoleAssignments", func(t *testing.T) { testUserStoreClearAllCustomRoleAssignments(t, ss) })
	// t.Run("GetAllAfter", func(t *testing.T) { testUserStoreGetAllAfter(t, ss) })
}

func testUserStoreSave(t *testing.T, ss store.Store) {
	teamId := model.NewId()
	maxUsersPerTeam := 50

	u1 := model.User{}
	u1.Email = MakeEmail()
	u1.Username = model.NewId()

	if err := (<-ss.User().Save(&u1)).Err; err != nil {
		t.Fatal("couldn't save user", err)
	}
	defer func() { store.Must(ss.User().PermanentDelete(u1.Id)) }()

	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, maxUsersPerTeam))

	if err := (<-ss.User().Save(&u1)).Err; err == nil {
		t.Fatal("shouldn't be able to update user from save")
	}

	u1.Id = ""
	if err := (<-ss.User().Save(&u1)).Err; err == nil {
		t.Fatal("should be unique email")
	}

	u1.Email = ""
	if err := (<-ss.User().Save(&u1)).Err; err == nil {
		t.Fatal("should be unique username")
	}

	u1.Email = strings.Repeat("0123456789", 20)
	u1.Username = ""
	if err := (<-ss.User().Save(&u1)).Err; err == nil {
		t.Fatal("should be unique username")
	}

	for i := 0; i < 49; i++ {
		u1.Id = ""
		u1.Email = MakeEmail()
		u1.Username = model.NewId()
		if err := (<-ss.User().Save(&u1)).Err; err != nil {
			t.Fatal("couldn't save item", err)
		}
		defer func() { store.Must(ss.User().PermanentDelete(u1.Id)) }()

		store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, maxUsersPerTeam))
	}

	u1.Id = ""
	u1.Email = MakeEmail()
	u1.Username = model.NewId()
	if err := (<-ss.User().Save(&u1)).Err; err != nil {
		t.Fatal("couldn't save item", err)
	}
	defer func() { store.Must(ss.User().PermanentDelete(u1.Id)) }()

	if err := (<-ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, maxUsersPerTeam)).Err; err == nil {
		t.Fatal("should be the limit")
	}
}

func testUserStoreUpdate(t *testing.T, ss store.Store) {
	u1 := &model.User{}
	u1.Email = MakeEmail()
	store.Must(ss.User().Save(u1))
	defer func() { store.Must(ss.User().PermanentDelete(u1.Id)) }()
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: model.NewId(), UserId: u1.Id}, -1))

	u2 := &model.User{}
	u2.Email = MakeEmail()
	u2.AuthService = "ldap"
	store.Must(ss.User().Save(u2))
	defer func() { store.Must(ss.User().PermanentDelete(u2.Id)) }()
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: model.NewId(), UserId: u2.Id}, -1))

	time.Sleep(100 * time.Millisecond)

	if err := (<-ss.User().Update(u1, false)).Err; err != nil {
		t.Fatal(err)
	}

	u1.Id = "missing"
	if err := (<-ss.User().Update(u1, false)).Err; err == nil {
		t.Fatal("Update should have failed because of missing key")
	}

	u1.Id = model.NewId()
	if err := (<-ss.User().Update(u1, false)).Err; err == nil {
		t.Fatal("Update should have faile because id change")
	}

	u2.Email = MakeEmail()
	if err := (<-ss.User().Update(u2, false)).Err; err == nil {
		t.Fatal("Update should have failed because you can't modify AD/LDAP fields")
	}

	u3 := &model.User{}
	u3.Email = MakeEmail()
	oldEmail := u3.Email
	u3.AuthService = "gitlab"
	store.Must(ss.User().Save(u3))
	defer func() { store.Must(ss.User().PermanentDelete(u3.Id)) }()
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: model.NewId(), UserId: u3.Id}, -1))

	u3.Email = MakeEmail()
	if result := <-ss.User().Update(u3, false); result.Err != nil {
		t.Fatal("Update should not have failed")
	} else {
		newUser := result.Data.([2]*model.User)[0]
		if newUser.Email != oldEmail {
			t.Fatal("Email should not have been updated as the update is not trusted")
		}
	}

	u3.Email = MakeEmail()
	if result := <-ss.User().Update(u3, true); result.Err != nil {
		t.Fatal("Update should not have failed")
	} else {
		newUser := result.Data.([2]*model.User)[0]
		if newUser.Email == oldEmail {
			t.Fatal("Email should have been updated as the update is trusted")
		}
	}

	if result := <-ss.User().UpdateLastPictureUpdate(u1.Id); result.Err != nil {
		t.Fatal("Update should not have failed")
	}
}

func testUserStoreUpdateUpdateAt(t *testing.T, ss store.Store) {
	u1 := &model.User{}
	u1.Email = MakeEmail()
	store.Must(ss.User().Save(u1))
	defer func() { store.Must(ss.User().PermanentDelete(u1.Id)) }()
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: model.NewId(), UserId: u1.Id}, -1))

	time.Sleep(10 * time.Millisecond)

	if err := (<-ss.User().UpdateUpdateAt(u1.Id)).Err; err != nil {
		t.Fatal(err)
	}

	if r1 := <-ss.User().Get(u1.Id); r1.Err != nil {
		t.Fatal(r1.Err)
	} else {
		if r1.Data.(*model.User).UpdateAt <= u1.UpdateAt {
			t.Fatal("UpdateAt not updated correctly")
		}
	}

}

func testUserStoreUpdateFailedPasswordAttempts(t *testing.T, ss store.Store) {
	u1 := &model.User{}
	u1.Email = MakeEmail()
	store.Must(ss.User().Save(u1))
	defer func() { store.Must(ss.User().PermanentDelete(u1.Id)) }()
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: model.NewId(), UserId: u1.Id}, -1))

	if err := (<-ss.User().UpdateFailedPasswordAttempts(u1.Id, 3)).Err; err != nil {
		t.Fatal(err)
	}

	if r1 := <-ss.User().Get(u1.Id); r1.Err != nil {
		t.Fatal(r1.Err)
	} else {
		if r1.Data.(*model.User).FailedAttempts != 3 {
			t.Fatal("FailedAttempts not updated correctly")
		}
	}

}

func testUserStoreGet(t *testing.T, ss store.Store) {
	u1 := &model.User{
		Email: MakeEmail(),
	}
	store.Must(ss.User().Save(u1))
	defer func() { store.Must(ss.User().PermanentDelete(u1.Id)) }()

	u2 := store.Must(ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: model.NewId(),
	})).(*model.User)
	store.Must(ss.Bot().Save(&model.Bot{
		UserId:    u2.Id,
		Username:  u2.Username,
		CreatorId: u1.Id,
	}))
	u2.IsBot = true
	defer func() { store.Must(ss.Bot().PermanentDelete(u2.Id)) }()
	defer func() { store.Must(ss.User().PermanentDelete(u2.Id)) }()

	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: model.NewId(), UserId: u1.Id}, -1))

	t.Run("fetch empty id", func(t *testing.T) {
		require.NotNil(t, (<-ss.User().Get("")).Err)
	})

	t.Run("fetch user 1", func(t *testing.T) {
		result := <-ss.User().Get(u1.Id)
		require.Nil(t, result.Err)

		actual := result.Data.(*model.User)
		require.Equal(t, u1, actual)
		require.False(t, actual.IsBot)
	})

	t.Run("fetch user 2, also a bot", func(t *testing.T) {
		result := <-ss.User().Get(u2.Id)
		require.Nil(t, result.Err)

		actual := result.Data.(*model.User)
		require.Equal(t, u2, actual)
		require.True(t, actual.IsBot)
	})
}

func testUserCount(t *testing.T, ss store.Store) {
	u1 := &model.User{}
	u1.Email = MakeEmail()
	store.Must(ss.User().Save(u1))
	defer func() { store.Must(ss.User().PermanentDelete(u1.Id)) }()
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: model.NewId(), UserId: u1.Id}, -1))

	if result := <-ss.User().GetTotalUsersCount(); result.Err != nil {
		t.Fatal(result.Err)
	} else {
		count := result.Data.(int64)
		require.False(t, count <= 0, "expected count > 0, got %d", count)
	}
}

func testGetAllUsingAuthService(t *testing.T, ss store.Store) {
	u1 := &model.User{}
	u1.Email = MakeEmail()
	u1.AuthService = "someservice"
	store.Must(ss.User().Save(u1))
	defer func() { store.Must(ss.User().PermanentDelete(u1.Id)) }()

	u2 := &model.User{}
	u2.Email = MakeEmail()
	u2.AuthService = "someservice"
	store.Must(ss.User().Save(u2))
	defer func() { store.Must(ss.User().PermanentDelete(u2.Id)) }()

	if r1 := <-ss.User().GetAllUsingAuthService(u1.AuthService); r1.Err != nil {
		t.Fatal(r1.Err)
	} else {
		users := r1.Data.([]*model.User)
		if len(users) < 2 {
			t.Fatal("invalid returned users")
		}
	}
}

func sanitized(user *model.User) *model.User {
	clonedUser := model.UserFromJson(strings.NewReader(user.ToJson()))
	clonedUser.AuthData = new(string)
	*clonedUser.AuthData = ""
	clonedUser.Props = model.StringMap{}

	return clonedUser
}

func testUserStoreGetAllProfiles(t *testing.T, ss store.Store) {
	u1 := store.Must(ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u1" + model.NewId(),
	})).(*model.User)
	defer func() { store.Must(ss.User().PermanentDelete(u1.Id)) }()

	u2 := store.Must(ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u2" + model.NewId(),
	})).(*model.User)
	defer func() { store.Must(ss.User().PermanentDelete(u2.Id)) }()

	u3 := store.Must(ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u3" + model.NewId(),
	})).(*model.User)
	store.Must(ss.Bot().Save(&model.Bot{
		UserId:    u3.Id,
		Username:  u3.Username,
		CreatorId: u1.Id,
	}))
	u3.IsBot = true
	defer func() { store.Must(ss.Bot().PermanentDelete(u3.Id)) }()
	defer func() { store.Must(ss.User().PermanentDelete(u3.Id)) }()

	t.Run("get offset 0, limit 100", func(t *testing.T) {
		result := <-ss.User().GetAllProfiles(0, 100)
		require.Nil(t, result.Err)

		actual := result.Data.([]*model.User)
		require.Equal(t, []*model.User{sanitized(u1), sanitized(u2), sanitized(u3)}, actual)
	})

	t.Run("get offset 0, limit 1", func(t *testing.T) {
		result := <-ss.User().GetAllProfiles(0, 1)
		require.Nil(t, result.Err)

		actual := result.Data.([]*model.User)
		require.Equal(t, []*model.User{sanitized(u1)}, actual)
	})

	t.Run("get all", func(t *testing.T) {
		result := <-ss.User().GetAll()
		require.Nil(t, result.Err)

		actual := result.Data.([]*model.User)
		require.Equal(t, []*model.User{u1, u2, u3}, actual)
	})

	t.Run("etag changes for all after user creation", func(t *testing.T) {
		result := <-ss.User().GetEtagForAllProfiles()
		require.Nil(t, result.Err)
		etag := result.Data.(string)

		u4 := &model.User{}
		u4.Email = MakeEmail()
		store.Must(ss.User().Save(u4))
		defer func() { store.Must(ss.User().PermanentDelete(u4.Id)) }()

		result = <-ss.User().GetEtagForAllProfiles()
		require.Nil(t, result.Err)
		updatedEtag := result.Data.(string)

		require.NotEqual(t, etag, updatedEtag)
	})
}

func testUserStoreGetProfiles(t *testing.T, ss store.Store) {
	teamId := model.NewId()

	u1 := store.Must(ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u1" + model.NewId(),
	})).(*model.User)
	defer func() { store.Must(ss.User().PermanentDelete(u1.Id)) }()
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, -1))

	u2 := store.Must(ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u2" + model.NewId(),
	})).(*model.User)
	defer func() { store.Must(ss.User().PermanentDelete(u2.Id)) }()
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u2.Id}, -1))

	u3 := store.Must(ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u3" + model.NewId(),
	})).(*model.User)
	store.Must(ss.Bot().Save(&model.Bot{
		UserId:    u3.Id,
		Username:  u3.Username,
		CreatorId: u1.Id,
	}))
	u3.IsBot = true
	defer func() { store.Must(ss.Bot().PermanentDelete(u3.Id)) }()
	defer func() { store.Must(ss.User().PermanentDelete(u3.Id)) }()
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u3.Id}, -1))

	t.Run("get offset 0, limit 100", func(t *testing.T) {
		result := <-ss.User().GetProfiles(teamId, 0, 100)
		require.Nil(t, result.Err)

		actual := result.Data.([]*model.User)
		require.Equal(t, []*model.User{sanitized(u1), sanitized(u2), sanitized(u3)}, actual)
	})

	t.Run("get offset 0, limit 1", func(t *testing.T) {
		result := <-ss.User().GetProfiles(teamId, 0, 1)
		require.Nil(t, result.Err)

		actual := result.Data.([]*model.User)
		require.Equal(t, []*model.User{sanitized(u1)}, actual)
	})

	t.Run("get unknown team id", func(t *testing.T) {
		result := <-ss.User().GetProfiles("123", 0, 100)
		require.Nil(t, result.Err)

		actual := result.Data.([]*model.User)
		require.Equal(t, []*model.User{}, actual)
	})

	t.Run("etag changes for all after user creation", func(t *testing.T) {
		result := <-ss.User().GetEtagForProfiles(teamId)
		require.Nil(t, result.Err)
		etag := result.Data.(string)

		u4 := &model.User{}
		u4.Email = MakeEmail()
		store.Must(ss.User().Save(u4))
		defer func() { store.Must(ss.User().PermanentDelete(u4.Id)) }()
		store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u4.Id}, -1))

		result = <-ss.User().GetEtagForProfiles(teamId)
		require.Nil(t, result.Err)
		updatedEtag := result.Data.(string)

		require.NotEqual(t, etag, updatedEtag)
	})
}

func testUserStoreGetProfilesInChannel(t *testing.T, ss store.Store) {
	teamId := model.NewId()

	u1 := store.Must(ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u1" + model.NewId(),
	})).(*model.User)
	defer func() { store.Must(ss.User().PermanentDelete(u1.Id)) }()
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, -1))

	u2 := store.Must(ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u2" + model.NewId(),
	})).(*model.User)
	defer func() { store.Must(ss.User().PermanentDelete(u2.Id)) }()
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u2.Id}, -1))

	u3 := store.Must(ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u3" + model.NewId(),
	})).(*model.User)
	defer func() { store.Must(ss.User().PermanentDelete(u3.Id)) }()
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u3.Id}, -1))
	store.Must(ss.Bot().Save(&model.Bot{
		UserId:    u3.Id,
		Username:  u3.Username,
		CreatorId: u1.Id,
	}))
	u3.IsBot = true

	c1 := store.Must(ss.Channel().Save(&model.Channel{
		TeamId:      teamId,
		DisplayName: "Profiles in channel",
		Name:        "profiles-" + model.NewId(),
		Type:        model.CHANNEL_OPEN,
	}, -1)).(*model.Channel)

	c2 := store.Must(ss.Channel().Save(&model.Channel{
		TeamId:      teamId,
		DisplayName: "Profiles in private",
		Name:        "profiles-" + model.NewId(),
		Type:        model.CHANNEL_PRIVATE,
	}, -1)).(*model.Channel)

	store.Must(ss.Channel().SaveMember(&model.ChannelMember{
		ChannelId:   c1.Id,
		UserId:      u1.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	}))

	store.Must(ss.Channel().SaveMember(&model.ChannelMember{
		ChannelId:   c1.Id,
		UserId:      u2.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	}))

	store.Must(ss.Channel().SaveMember(&model.ChannelMember{
		ChannelId:   c1.Id,
		UserId:      u3.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	}))

	store.Must(ss.Channel().SaveMember(&model.ChannelMember{
		ChannelId:   c2.Id,
		UserId:      u1.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	}))

	t.Run("get in channel 1, offset 0, limit 100", func(t *testing.T) {
		result := <-ss.User().GetProfilesInChannel(c1.Id, 0, 100)
		require.Nil(t, result.Err)
		assert.Equal(t, []*model.User{sanitized(u1), sanitized(u2), sanitized(u3)}, result.Data.([]*model.User))
	})

	t.Run("get in channel 1, offset 1, limit 2", func(t *testing.T) {
		result := <-ss.User().GetProfilesInChannel(c1.Id, 1, 2)
		require.Nil(t, result.Err)
		assert.Equal(t, []*model.User{sanitized(u2), sanitized(u3)}, result.Data.([]*model.User))
	})

	t.Run("get in channel 2, offset 0, limit 1", func(t *testing.T) {
		result := <-ss.User().GetProfilesInChannel(c2.Id, 0, 1)
		require.Nil(t, result.Err)
		assert.Equal(t, []*model.User{sanitized(u1)}, result.Data.([]*model.User))
	})
}

func testUserStoreGetProfilesInChannelByStatus(t *testing.T, ss store.Store) {
	teamId := model.NewId()

	u1 := store.Must(ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u1" + model.NewId(),
	})).(*model.User)
	defer func() { store.Must(ss.User().PermanentDelete(u1.Id)) }()
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, -1))

	u2 := store.Must(ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u2" + model.NewId(),
	})).(*model.User)
	defer func() { store.Must(ss.User().PermanentDelete(u2.Id)) }()
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u2.Id}, -1))

	u3 := store.Must(ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u3" + model.NewId(),
	})).(*model.User)
	defer func() { store.Must(ss.User().PermanentDelete(u3.Id)) }()
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u3.Id}, -1))
	store.Must(ss.Bot().Save(&model.Bot{
		UserId:    u3.Id,
		Username:  u3.Username,
		CreatorId: u1.Id,
	}))
	u3.IsBot = true

	c1 := store.Must(ss.Channel().Save(&model.Channel{
		TeamId:      teamId,
		DisplayName: "Profiles in channel",
		Name:        "profiles-" + model.NewId(),
		Type:        model.CHANNEL_OPEN,
	}, -1)).(*model.Channel)

	c2 := store.Must(ss.Channel().Save(&model.Channel{
		TeamId:      teamId,
		DisplayName: "Profiles in private",
		Name:        "profiles-" + model.NewId(),
		Type:        model.CHANNEL_PRIVATE,
	}, -1)).(*model.Channel)

	store.Must(ss.Channel().SaveMember(&model.ChannelMember{
		ChannelId:   c1.Id,
		UserId:      u1.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	}))

	store.Must(ss.Channel().SaveMember(&model.ChannelMember{
		ChannelId:   c1.Id,
		UserId:      u2.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	}))

	store.Must(ss.Channel().SaveMember(&model.ChannelMember{
		ChannelId:   c1.Id,
		UserId:      u3.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	}))

	store.Must(ss.Channel().SaveMember(&model.ChannelMember{
		ChannelId:   c2.Id,
		UserId:      u1.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	}))

	store.Must(ss.Status().SaveOrUpdate(&model.Status{
		UserId: u1.Id,
		Status: model.STATUS_DND,
	}))
	store.Must(ss.Status().SaveOrUpdate(&model.Status{
		UserId: u2.Id,
		Status: model.STATUS_AWAY,
	}))
	store.Must(ss.Status().SaveOrUpdate(&model.Status{
		UserId: u3.Id,
		Status: model.STATUS_ONLINE,
	}))

	t.Run("get in channel 1 by status, offset 0, limit 100", func(t *testing.T) {
		result := <-ss.User().GetProfilesInChannelByStatus(c1.Id, 0, 100)
		require.Nil(t, result.Err)
		assert.Equal(t, []*model.User{sanitized(u3), sanitized(u2), sanitized(u1)}, result.Data.([]*model.User))
	})

	t.Run("get in channel 2 by status, offset 0, limit 1", func(t *testing.T) {
		result := <-ss.User().GetProfilesInChannelByStatus(c2.Id, 0, 1)
		require.Nil(t, result.Err)
		assert.Equal(t, []*model.User{sanitized(u1)}, result.Data.([]*model.User))
	})
}

func testUserStoreGetProfilesWithoutTeam(t *testing.T, ss store.Store) {
	teamId := model.NewId()

	// These usernames need to appear in the first 100 users for this to work

	u1 := &model.User{}
	u1.Username = "a000000000" + model.NewId()
	u1.Email = MakeEmail()
	store.Must(ss.User().Save(u1))
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, -1))
	defer func() { store.Must(ss.User().PermanentDelete(u1.Id)) }()

	u2 := &model.User{}
	u2.Username = "a000000001" + model.NewId()
	u2.Email = MakeEmail()
	store.Must(ss.User().Save(u2))
	defer func() { store.Must(ss.User().PermanentDelete(u2.Id)) }()

	if r1 := <-ss.User().GetProfilesWithoutTeam(0, 100); r1.Err != nil {
		t.Fatal(r1.Err)
	} else {
		users := r1.Data.([]*model.User)

		found1 := false
		found2 := false
		for _, u := range users {
			if u.Id == u1.Id {
				found1 = true
			} else if u.Id == u2.Id {
				found2 = true
			}
		}

		if found1 {
			t.Fatal("shouldn't have returned user on team")
		} else if !found2 {
			t.Fatal("should've returned user without any teams")
		}
	}
}

func testUserStoreGetAllProfilesInChannel(t *testing.T, ss store.Store) {
	teamId := model.NewId()

	u1 := store.Must(ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u1" + model.NewId(),
	})).(*model.User)
	defer func() { store.Must(ss.User().PermanentDelete(u1.Id)) }()
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, -1))

	u2 := store.Must(ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u2" + model.NewId(),
	})).(*model.User)
	defer func() { store.Must(ss.User().PermanentDelete(u2.Id)) }()
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u2.Id}, -1))

	u3 := store.Must(ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u3" + model.NewId(),
	})).(*model.User)
	defer func() { store.Must(ss.User().PermanentDelete(u3.Id)) }()
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u3.Id}, -1))
	store.Must(ss.Bot().Save(&model.Bot{
		UserId:    u3.Id,
		Username:  u3.Username,
		CreatorId: u1.Id,
	}))
	u3.IsBot = true

	c1 := store.Must(ss.Channel().Save(&model.Channel{
		TeamId:      teamId,
		DisplayName: "Profiles in channel",
		Name:        "profiles-" + model.NewId(),
		Type:        model.CHANNEL_OPEN,
	}, -1)).(*model.Channel)

	c2 := store.Must(ss.Channel().Save(&model.Channel{
		TeamId:      teamId,
		DisplayName: "Profiles in private",
		Name:        "profiles-" + model.NewId(),
		Type:        model.CHANNEL_PRIVATE,
	}, -1)).(*model.Channel)

	store.Must(ss.Channel().SaveMember(&model.ChannelMember{
		ChannelId:   c1.Id,
		UserId:      u1.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	}))

	store.Must(ss.Channel().SaveMember(&model.ChannelMember{
		ChannelId:   c1.Id,
		UserId:      u2.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	}))

	store.Must(ss.Channel().SaveMember(&model.ChannelMember{
		ChannelId:   c1.Id,
		UserId:      u3.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	}))

	store.Must(ss.Channel().SaveMember(&model.ChannelMember{
		ChannelId:   c2.Id,
		UserId:      u1.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	}))

	t.Run("all profiles in channel 1, no caching", func(t *testing.T) {
		result := <-ss.User().GetAllProfilesInChannel(c1.Id, false)
		require.Nil(t, result.Err)
		assert.Equal(t, map[string]*model.User{
			u1.Id: sanitized(u1),
			u2.Id: sanitized(u2),
			u3.Id: sanitized(u3),
		}, result.Data.(map[string]*model.User))
	})

	t.Run("all profiles in channel 2, no caching", func(t *testing.T) {
		result := <-ss.User().GetAllProfilesInChannel(c2.Id, false)
		require.Nil(t, result.Err)
		assert.Equal(t, map[string]*model.User{
			u1.Id: sanitized(u1),
		}, result.Data.(map[string]*model.User))
	})

	t.Run("all profiles in channel 2, caching", func(t *testing.T) {
		result := <-ss.User().GetAllProfilesInChannel(c2.Id, true)
		require.Nil(t, result.Err)
		assert.Equal(t, map[string]*model.User{
			u1.Id: sanitized(u1),
		}, result.Data.(map[string]*model.User))
	})

	t.Run("all profiles in channel 2, caching [repeated]", func(t *testing.T) {
		result := <-ss.User().GetAllProfilesInChannel(c2.Id, true)
		require.Nil(t, result.Err)
		assert.Equal(t, map[string]*model.User{
			u1.Id: sanitized(u1),
		}, result.Data.(map[string]*model.User))
	})

	ss.User().InvalidateProfilesInChannelCacheByUser(u1.Id)
	ss.User().InvalidateProfilesInChannelCache(c2.Id)
}

func testUserStoreGetProfilesNotInChannel(t *testing.T, ss store.Store) {
	teamId := model.NewId()

	u1 := store.Must(ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u1" + model.NewId(),
	})).(*model.User)
	defer func() { store.Must(ss.User().PermanentDelete(u1.Id)) }()
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, -1))

	u2 := store.Must(ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u2" + model.NewId(),
	})).(*model.User)
	defer func() { store.Must(ss.User().PermanentDelete(u2.Id)) }()
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u2.Id}, -1))

	u3 := store.Must(ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u3" + model.NewId(),
	})).(*model.User)
	defer func() { store.Must(ss.User().PermanentDelete(u3.Id)) }()
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u3.Id}, -1))
	store.Must(ss.Bot().Save(&model.Bot{
		UserId:    u3.Id,
		Username:  u3.Username,
		CreatorId: u1.Id,
	}))
	u3.IsBot = true

	c1 := store.Must(ss.Channel().Save(&model.Channel{
		TeamId:      teamId,
		DisplayName: "Profiles in channel",
		Name:        "profiles-" + model.NewId(),
		Type:        model.CHANNEL_OPEN,
	}, -1)).(*model.Channel)

	c2 := store.Must(ss.Channel().Save(&model.Channel{
		TeamId:      teamId,
		DisplayName: "Profiles in private",
		Name:        "profiles-" + model.NewId(),
		Type:        model.CHANNEL_PRIVATE,
	}, -1)).(*model.Channel)

	t.Run("get team 1, channel 1, offset 0, limit 100", func(t *testing.T) {
		result := <-ss.User().GetProfilesNotInChannel(teamId, c1.Id, 0, 100)
		require.Nil(t, result.Err)
		assert.Equal(t, []*model.User{
			sanitized(u1),
			sanitized(u2),
			sanitized(u3),
		}, result.Data.([]*model.User))
	})

	t.Run("get team 1, channel 2, offset 0, limit 100", func(t *testing.T) {
		result := <-ss.User().GetProfilesNotInChannel(teamId, c2.Id, 0, 100)
		require.Nil(t, result.Err)
		assert.Equal(t, []*model.User{
			sanitized(u1),
			sanitized(u2),
			sanitized(u3),
		}, result.Data.([]*model.User))
	})

	store.Must(ss.Channel().SaveMember(&model.ChannelMember{
		ChannelId:   c1.Id,
		UserId:      u1.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	}))

	store.Must(ss.Channel().SaveMember(&model.ChannelMember{
		ChannelId:   c1.Id,
		UserId:      u2.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	}))

	store.Must(ss.Channel().SaveMember(&model.ChannelMember{
		ChannelId:   c1.Id,
		UserId:      u3.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	}))

	store.Must(ss.Channel().SaveMember(&model.ChannelMember{
		ChannelId:   c2.Id,
		UserId:      u1.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	}))

	t.Run("get team 1, channel 1, offset 0, limit 100, after update", func(t *testing.T) {
		result := <-ss.User().GetProfilesNotInChannel(teamId, c1.Id, 0, 100)
		require.Nil(t, result.Err)
		assert.Equal(t, []*model.User{}, result.Data.([]*model.User))
	})

	t.Run("get team 1, channel 2, offset 0, limit 100, after update", func(t *testing.T) {
		result := <-ss.User().GetProfilesNotInChannel(teamId, c2.Id, 0, 100)
		require.Nil(t, result.Err)
		assert.Equal(t, []*model.User{
			sanitized(u2),
			sanitized(u3),
		}, result.Data.([]*model.User))
	})
}

func testUserStoreGetProfilesByIds(t *testing.T, ss store.Store) {
	teamId := model.NewId()

	u1 := store.Must(ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u1" + model.NewId(),
	})).(*model.User)
	defer func() { store.Must(ss.User().PermanentDelete(u1.Id)) }()
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, -1))

	u2 := store.Must(ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u2" + model.NewId(),
	})).(*model.User)
	defer func() { store.Must(ss.User().PermanentDelete(u2.Id)) }()
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u2.Id}, -1))

	u3 := store.Must(ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u3" + model.NewId(),
	})).(*model.User)
	defer func() { store.Must(ss.User().PermanentDelete(u3.Id)) }()
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u3.Id}, -1))
	store.Must(ss.Bot().Save(&model.Bot{
		UserId:    u3.Id,
		Username:  u3.Username,
		CreatorId: u1.Id,
	}))
	u3.IsBot = true

	t.Run("get u1 by id, no caching", func(t *testing.T) {
		result := <-ss.User().GetProfileByIds([]string{u1.Id}, false)
		require.Nil(t, result.Err)
		assert.Equal(t, []*model.User{sanitized(u1)}, result.Data.([]*model.User))
	})

	t.Run("get u1 by id, caching", func(t *testing.T) {
		result := <-ss.User().GetProfileByIds([]string{u1.Id}, true)
		require.Nil(t, result.Err)
		assert.Equal(t, []*model.User{sanitized(u1)}, result.Data.([]*model.User))
	})

	t.Run("get u1, u2, u3 by id, no caching", func(t *testing.T) {
		result := <-ss.User().GetProfileByIds([]string{u1.Id, u2.Id, u3.Id}, false)
		require.Nil(t, result.Err)
		assert.Equal(t, []*model.User{sanitized(u1), sanitized(u2), sanitized(u3)}, result.Data.([]*model.User))
	})

	t.Run("get u1, u2, u3 by id, caching", func(t *testing.T) {
		result := <-ss.User().GetProfileByIds([]string{u1.Id, u2.Id, u3.Id}, true)
		require.Nil(t, result.Err)
		assert.Equal(t, []*model.User{sanitized(u1), sanitized(u2), sanitized(u3)}, result.Data.([]*model.User))
	})

	t.Run("get unknown id, caching", func(t *testing.T) {
		result := <-ss.User().GetProfileByIds([]string{"123"}, true)
		require.Nil(t, result.Err)
		assert.Equal(t, []*model.User{}, result.Data.([]*model.User))
	})
}

func testUserStoreGetProfilesByUsernames(t *testing.T, ss store.Store) {
	teamId := model.NewId()

	u1 := &model.User{}
	u1.Email = MakeEmail()
	u1.Username = "username1" + model.NewId()
	store.Must(ss.User().Save(u1))
	defer func() { store.Must(ss.User().PermanentDelete(u1.Id)) }()
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, -1))

	u2 := &model.User{}
	u2.Email = MakeEmail()
	u2.Username = "username2" + model.NewId()
	store.Must(ss.User().Save(u2))
	defer func() { store.Must(ss.User().PermanentDelete(u2.Id)) }()
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u2.Id}, -1))

	if r1 := <-ss.User().GetProfilesByUsernames([]string{u1.Username, u2.Username}, teamId); r1.Err != nil {
		t.Fatal(r1.Err)
	} else {
		users := r1.Data.([]*model.User)
		if len(users) != 2 {
			t.Fatal("invalid returned users")
		}

		if users[0].Id != u1.Id && users[1].Id != u1.Id {
			t.Fatal("invalid returned user 1")
		}

		if users[0].Id != u2.Id && users[1].Id != u2.Id {
			t.Fatal("invalid returned user 2")
		}
	}

	if r1 := <-ss.User().GetProfilesByUsernames([]string{u1.Username}, teamId); r1.Err != nil {
		t.Fatal(r1.Err)
	} else {
		users := r1.Data.([]*model.User)
		if len(users) != 1 {
			t.Fatal("invalid returned users")
		}

		if users[0].Id != u1.Id {
			t.Fatal("invalid returned user")
		}
	}

	team2Id := model.NewId()

	u3 := &model.User{}
	u3.Email = MakeEmail()
	u3.Username = "username3" + model.NewId()
	store.Must(ss.User().Save(u3))
	defer func() { store.Must(ss.User().PermanentDelete(u3.Id)) }()
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: team2Id, UserId: u3.Id}, -1))

	if r1 := <-ss.User().GetProfilesByUsernames([]string{u1.Username, u3.Username}, ""); r1.Err != nil {
		t.Fatal(r1.Err)
	} else {
		users := r1.Data.([]*model.User)
		if len(users) != 2 {
			t.Fatal("invalid returned users")
		}

		if users[0].Id != u1.Id && users[1].Id != u1.Id {
			t.Fatal("invalid returned user 1")
		}

		if users[0].Id != u3.Id && users[1].Id != u3.Id {
			t.Fatal("invalid returned user 3")
		}
	}

	if r1 := <-ss.User().GetProfilesByUsernames([]string{u1.Username, u3.Username}, teamId); r1.Err != nil {
		t.Fatal(r1.Err)
	} else {
		users := r1.Data.([]*model.User)
		if len(users) != 1 {
			t.Fatal("invalid returned users")
		}

		if users[0].Id != u1.Id {
			t.Fatal("invalid returned user")
		}
	}
}

func testUserStoreGetSystemAdminProfiles(t *testing.T, ss store.Store) {
	teamId := model.NewId()

	u1 := store.Must(ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Roles:    model.SYSTEM_USER_ROLE_ID + " " + model.SYSTEM_ADMIN_ROLE_ID,
		Username: "u1" + model.NewId(),
	})).(*model.User)
	defer func() { store.Must(ss.User().PermanentDelete(u1.Id)) }()
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, -1))

	u2 := store.Must(ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u2" + model.NewId(),
	})).(*model.User)
	defer func() { store.Must(ss.User().PermanentDelete(u2.Id)) }()
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u2.Id}, -1))

	u3 := store.Must(ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Roles:    model.SYSTEM_USER_ROLE_ID + " " + model.SYSTEM_ADMIN_ROLE_ID,
		Username: "u3" + model.NewId(),
	})).(*model.User)
	defer func() { store.Must(ss.User().PermanentDelete(u3.Id)) }()
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u3.Id}, -1))
	store.Must(ss.Bot().Save(&model.Bot{
		UserId:    u3.Id,
		Username:  u3.Username,
		CreatorId: u1.Id,
	}))
	u3.IsBot = true

	t.Run("all system admin profiles", func(t *testing.T) {
		result := <-ss.User().GetSystemAdminProfiles()
		require.Nil(t, result.Err)
		assert.Equal(t, map[string]*model.User{
			u1.Id: sanitized(u1),
			u3.Id: sanitized(u3),
		}, result.Data.(map[string]*model.User))
	})
}

func testUserStoreGetByEmail(t *testing.T, ss store.Store) {
	teamId := model.NewId()

	u1 := store.Must(ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Roles:    model.SYSTEM_USER_ROLE_ID + " " + model.SYSTEM_ADMIN_ROLE_ID,
		Username: "u1" + model.NewId(),
	})).(*model.User)
	defer func() { store.Must(ss.User().PermanentDelete(u1.Id)) }()
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, -1))

	u2 := store.Must(ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u2" + model.NewId(),
	})).(*model.User)
	defer func() { store.Must(ss.User().PermanentDelete(u2.Id)) }()
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u2.Id}, -1))

	u3 := store.Must(ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Roles:    model.SYSTEM_USER_ROLE_ID + " " + model.SYSTEM_ADMIN_ROLE_ID,
		Username: "u3" + model.NewId(),
	})).(*model.User)
	defer func() { store.Must(ss.User().PermanentDelete(u3.Id)) }()
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u3.Id}, -1))
	store.Must(ss.Bot().Save(&model.Bot{
		UserId:    u3.Id,
		Username:  u3.Username,
		CreatorId: u1.Id,
	}))
	u3.IsBot = true

	t.Run("get u1 by email", func(t *testing.T) {
		result := <-ss.User().GetByEmail(u1.Email)
		require.Nil(t, result.Err)
		assert.Equal(t, u1, result.Data.(*model.User))
	})

	t.Run("get u2 by email", func(t *testing.T) {
		result := <-ss.User().GetByEmail(u2.Email)
		require.Nil(t, result.Err)
		assert.Equal(t, u2, result.Data.(*model.User))
	})

	t.Run("get u3 by email", func(t *testing.T) {
		result := <-ss.User().GetByEmail(u3.Email)
		require.Nil(t, result.Err)
		assert.Equal(t, u3, result.Data.(*model.User))
	})

	t.Run("get by empty email", func(t *testing.T) {
		result := <-ss.User().GetByEmail("")
		require.NotNil(t, result.Err)
		require.Equal(t, result.Err.Id, store.MISSING_ACCOUNT_ERROR)
	})

	t.Run("get by unknown", func(t *testing.T) {
		result := <-ss.User().GetByEmail("unknown")
		require.NotNil(t, result.Err)
		require.Equal(t, result.Err.Id, store.MISSING_ACCOUNT_ERROR)
	})
}

func testUserStoreGetByAuthData(t *testing.T, ss store.Store) {
	teamId := model.NewId()

	auth := "123" + model.NewId()

	u1 := &model.User{}
	u1.Email = MakeEmail()
	u1.AuthData = &auth
	u1.AuthService = "service"
	store.Must(ss.User().Save(u1))
	defer func() { store.Must(ss.User().PermanentDelete(u1.Id)) }()
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, -1))

	if err := (<-ss.User().GetByAuth(u1.AuthData, u1.AuthService)).Err; err != nil {
		t.Fatal(err)
	}

	rauth := ""
	if err := (<-ss.User().GetByAuth(&rauth, "")).Err; err == nil {
		t.Fatal("Should have failed because of missing auth data")
	}
}

func testUserStoreGetByUsername(t *testing.T, ss store.Store) {
	teamId := model.NewId()

	u1 := &model.User{}
	u1.Email = MakeEmail()
	u1.Username = model.NewId()
	store.Must(ss.User().Save(u1))
	defer func() { store.Must(ss.User().PermanentDelete(u1.Id)) }()
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, -1))

	if err := (<-ss.User().GetByUsername(u1.Username)).Err; err != nil {
		t.Fatal(err)
	}

	if err := (<-ss.User().GetByUsername("")).Err; err == nil {
		t.Fatal("Should have failed because of missing username")
	}
}

func testUserStoreGetForLogin(t *testing.T, ss store.Store) {
	auth := model.NewId()

	u1 := &model.User{
		Email:       MakeEmail(),
		Username:    model.NewId(),
		AuthService: model.USER_AUTH_SERVICE_GITLAB,
		AuthData:    &auth,
	}
	store.Must(ss.User().Save(u1))
	defer func() { store.Must(ss.User().PermanentDelete(u1.Id)) }()

	auth2 := model.NewId()

	u2 := &model.User{
		Email:       MakeEmail(),
		Username:    model.NewId(),
		AuthService: model.USER_AUTH_SERVICE_LDAP,
		AuthData:    &auth2,
	}
	store.Must(ss.User().Save(u2))
	defer func() { store.Must(ss.User().PermanentDelete(u2.Id)) }()

	if result := <-ss.User().GetForLogin(u1.Username, true, true); result.Err != nil {
		t.Fatal("Should have gotten user by username", result.Err)
	} else if result.Data.(*model.User).Id != u1.Id {
		t.Fatal("Should have gotten user1 by username")
	}

	if result := <-ss.User().GetForLogin(u1.Email, true, true); result.Err != nil {
		t.Fatal("Should have gotten user by email", result.Err)
	} else if result.Data.(*model.User).Id != u1.Id {
		t.Fatal("Should have gotten user1 by email")
	}

	// prevent getting user when different login methods are disabled
	if result := <-ss.User().GetForLogin(u1.Username, false, true); result.Err == nil {
		t.Fatal("Should have failed to get user1 by username")
	}

	if result := <-ss.User().GetForLogin(u1.Email, true, false); result.Err == nil {
		t.Fatal("Should have failed to get user1 by email")
	}
}

func testUserStoreUpdatePassword(t *testing.T, ss store.Store) {
	teamId := model.NewId()

	u1 := &model.User{}
	u1.Email = MakeEmail()
	store.Must(ss.User().Save(u1))
	defer func() { store.Must(ss.User().PermanentDelete(u1.Id)) }()
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, -1))

	hashedPassword := model.HashPassword("newpwd")

	if err := (<-ss.User().UpdatePassword(u1.Id, hashedPassword)).Err; err != nil {
		t.Fatal(err)
	}

	if r1 := <-ss.User().GetByEmail(u1.Email); r1.Err != nil {
		t.Fatal(r1.Err)
	} else {
		user := r1.Data.(*model.User)
		if user.Password != hashedPassword {
			t.Fatal("Password was not updated correctly")
		}
	}
}

func testUserStoreDelete(t *testing.T, ss store.Store) {
	u1 := &model.User{}
	u1.Email = MakeEmail()
	store.Must(ss.User().Save(u1))
	defer func() { store.Must(ss.User().PermanentDelete(u1.Id)) }()
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: model.NewId(), UserId: u1.Id}, -1))

	if err := (<-ss.User().PermanentDelete(u1.Id)).Err; err != nil {
		t.Fatal(err)
	}
}

func testUserStoreUpdateAuthData(t *testing.T, ss store.Store) {
	teamId := model.NewId()

	u1 := &model.User{}
	u1.Email = MakeEmail()
	store.Must(ss.User().Save(u1))
	defer func() { store.Must(ss.User().PermanentDelete(u1.Id)) }()
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, -1))

	service := "someservice"
	authData := model.NewId()

	if err := (<-ss.User().UpdateAuthData(u1.Id, service, &authData, "", true)).Err; err != nil {
		t.Fatal(err)
	}

	if r1 := <-ss.User().GetByEmail(u1.Email); r1.Err != nil {
		t.Fatal(r1.Err)
	} else {
		user := r1.Data.(*model.User)
		if user.AuthService != service {
			t.Fatal("AuthService was not updated correctly")
		}
		if *user.AuthData != authData {
			t.Fatal("AuthData was not updated correctly")
		}
		if user.Password != "" {
			t.Fatal("Password was not cleared properly")
		}
	}
}

func testUserUnreadCount(t *testing.T, ss store.Store) {
	teamId := model.NewId()

	c1 := model.Channel{}
	c1.TeamId = teamId
	c1.DisplayName = "Unread Messages"
	c1.Name = "unread-messages-" + model.NewId()
	c1.Type = model.CHANNEL_OPEN

	c2 := model.Channel{}
	c2.TeamId = teamId
	c2.DisplayName = "Unread Direct"
	c2.Name = "unread-direct-" + model.NewId()
	c2.Type = model.CHANNEL_DIRECT

	u1 := &model.User{}
	u1.Username = "user1" + model.NewId()
	u1.Email = MakeEmail()
	store.Must(ss.User().Save(u1))
	defer func() { store.Must(ss.User().PermanentDelete(u1.Id)) }()
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, -1))

	u2 := &model.User{}
	u2.Email = MakeEmail()
	u2.Username = "user2" + model.NewId()
	store.Must(ss.User().Save(u2))
	defer func() { store.Must(ss.User().PermanentDelete(u2.Id)) }()
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u2.Id}, -1))

	if err := (<-ss.Channel().Save(&c1, -1)).Err; err != nil {
		t.Fatal("couldn't save item", err)
	}

	m1 := model.ChannelMember{}
	m1.ChannelId = c1.Id
	m1.UserId = u1.Id
	m1.NotifyProps = model.GetDefaultChannelNotifyProps()

	m2 := model.ChannelMember{}
	m2.ChannelId = c1.Id
	m2.UserId = u2.Id
	m2.NotifyProps = model.GetDefaultChannelNotifyProps()

	store.Must(ss.Channel().SaveMember(&m1))
	store.Must(ss.Channel().SaveMember(&m2))

	m1.ChannelId = c2.Id
	m2.ChannelId = c2.Id

	if err := (<-ss.Channel().SaveDirectChannel(&c2, &m1, &m2)).Err; err != nil {
		t.Fatal("couldn't save direct channel", err)
	}

	p1 := model.Post{}
	p1.ChannelId = c1.Id
	p1.UserId = u1.Id
	p1.Message = "this is a message for @" + u2.Username

	// Post one message with mention to open channel
	store.Must(ss.Post().Save(&p1))
	store.Must(ss.Channel().IncrementMentionCount(c1.Id, u2.Id))

	// Post 2 messages without mention to direct channel
	p2 := model.Post{}
	p2.ChannelId = c2.Id
	p2.UserId = u1.Id
	p2.Message = "first message"
	store.Must(ss.Post().Save(&p2))
	store.Must(ss.Channel().IncrementMentionCount(c2.Id, u2.Id))

	p3 := model.Post{}
	p3.ChannelId = c2.Id
	p3.UserId = u1.Id
	p3.Message = "second message"
	store.Must(ss.Post().Save(&p3))
	store.Must(ss.Channel().IncrementMentionCount(c2.Id, u2.Id))

	badge := (<-ss.User().GetUnreadCount(u2.Id)).Data.(int64)
	if badge != 3 {
		t.Fatal("should have 3 unread messages")
	}

	badge = (<-ss.User().GetUnreadCountForChannel(u2.Id, c1.Id)).Data.(int64)
	if badge != 1 {
		t.Fatal("should have 1 unread messages for that channel")
	}

	badge = (<-ss.User().GetUnreadCountForChannel(u2.Id, c2.Id)).Data.(int64)
	if badge != 2 {
		t.Fatal("should have 2 unread messages for that channel")
	}
}

func testUserStoreUpdateMfaSecret(t *testing.T, ss store.Store) {
	u1 := model.User{}
	u1.Email = MakeEmail()
	store.Must(ss.User().Save(&u1))
	defer func() { store.Must(ss.User().PermanentDelete(u1.Id)) }()

	time.Sleep(100 * time.Millisecond)

	if err := (<-ss.User().UpdateMfaSecret(u1.Id, "12345")).Err; err != nil {
		t.Fatal(err)
	}

	// should pass, no update will occur though
	if err := (<-ss.User().UpdateMfaSecret("junk", "12345")).Err; err != nil {
		t.Fatal(err)
	}
}

func testUserStoreUpdateMfaActive(t *testing.T, ss store.Store) {
	u1 := model.User{}
	u1.Email = MakeEmail()
	store.Must(ss.User().Save(&u1))
	defer func() { store.Must(ss.User().PermanentDelete(u1.Id)) }()

	time.Sleep(100 * time.Millisecond)

	if err := (<-ss.User().UpdateMfaActive(u1.Id, true)).Err; err != nil {
		t.Fatal(err)
	}

	if err := (<-ss.User().UpdateMfaActive(u1.Id, false)).Err; err != nil {
		t.Fatal(err)
	}

	// should pass, no update will occur though
	if err := (<-ss.User().UpdateMfaActive("junk", true)).Err; err != nil {
		t.Fatal(err)
	}
}

func testUserStoreGetRecentlyActiveUsersForTeam(t *testing.T, ss store.Store) {
	u1 := &model.User{}
	u1.Email = MakeEmail()
	store.Must(ss.User().Save(u1))
	defer func() { store.Must(ss.User().PermanentDelete(u1.Id)) }()
	store.Must(ss.Status().SaveOrUpdate(&model.Status{UserId: u1.Id, Status: model.STATUS_ONLINE, Manual: false, LastActivityAt: model.GetMillis(), ActiveChannel: ""}))
	tid := model.NewId()
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: tid, UserId: u1.Id}, -1))

	if r1 := <-ss.User().GetRecentlyActiveUsersForTeam(tid, 0, 100); r1.Err != nil {
		t.Fatal(r1.Err)
	}
}

func testUserStoreGetNewUsersForTeam(t *testing.T, ss store.Store) {
	u1 := &model.User{}
	u1.Email = MakeEmail()
	store.Must(ss.User().Save(u1))
	defer func() { store.Must(ss.User().PermanentDelete(u1.Id)) }()
	store.Must(ss.Status().SaveOrUpdate(&model.Status{UserId: u1.Id, Status: model.STATUS_ONLINE, Manual: false, LastActivityAt: model.GetMillis(), ActiveChannel: ""}))
	tid := model.NewId()
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: tid, UserId: u1.Id}, -1))

	if r1 := <-ss.User().GetNewUsersForTeam(tid, 0, 100); r1.Err != nil {
		t.Fatal(r1.Err)
	}
}

func assertUsers(t *testing.T, expected, actual []*model.User) {
	expectedUsernames := make([]string, 0, len(expected))
	for _, user := range expected {
		expectedUsernames = append(expectedUsernames, user.Username)
	}

	actualUsernames := make([]string, 0, len(actual))
	for _, user := range actual {
		actualUsernames = append(actualUsernames, user.Username)
	}

	if assert.Equal(t, expectedUsernames, actualUsernames) {
		assert.Equal(t, expected, actual)
	}
}

func testUserStoreSearch(t *testing.T, ss store.Store) {
	u1 := &model.User{
		Username:  "jimbo1" + model.NewId(),
		FirstName: "Tim",
		LastName:  "Bill",
		Nickname:  "Rob",
		Email:     "harold" + model.NewId() + "@simulator.amazonses.com",
	}
	store.Must(ss.User().Save(u1))
	defer func() { store.Must(ss.User().PermanentDelete(u1.Id)) }()

	u2 := &model.User{
		Username: "jim-bobby" + model.NewId(),
		Email:    MakeEmail(),
	}
	store.Must(ss.User().Save(u2))
	defer func() { store.Must(ss.User().PermanentDelete(u2.Id)) }()

	u3 := &model.User{
		Username: "jimbo3" + model.NewId(),
		Email:    MakeEmail(),
		DeleteAt: 1,
	}
	store.Must(ss.User().Save(u3))
	defer func() { store.Must(ss.User().PermanentDelete(u3.Id)) }()

	u5 := &model.User{
		Username:  "yu" + model.NewId(),
		FirstName: "En",
		LastName:  "Yu",
		Nickname:  "enyu",
		Email:     MakeEmail(),
	}
	store.Must(ss.User().Save(u5))
	defer func() { store.Must(ss.User().PermanentDelete(u5.Id)) }()

	u6 := &model.User{
		Username:  "underscore" + model.NewId(),
		FirstName: "Du_",
		LastName:  "_DE",
		Nickname:  "lodash",
		Email:     MakeEmail(),
	}
	store.Must(ss.User().Save(u6))
	defer func() { store.Must(ss.User().PermanentDelete(u6.Id)) }()

	tid := model.NewId()
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: tid, UserId: u1.Id}, -1))
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: tid, UserId: u2.Id}, -1))
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: tid, UserId: u3.Id}, -1))
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: tid, UserId: u5.Id}, -1))
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: tid, UserId: u6.Id}, -1))

	// The users returned from the database will have AuthData as an empty string.
	nilAuthData := new(string)
	*nilAuthData = ""

	u1.AuthData = nilAuthData
	u2.AuthData = nilAuthData
	u3.AuthData = nilAuthData
	u5.AuthData = nilAuthData
	u6.AuthData = nilAuthData

	testCases := []struct {
		Description string
		TeamId      string
		Term        string
		Options     *model.UserSearchOptions
		Expected    []*model.User
	}{
		{
			"search jimb",
			tid,
			"jimb",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{u1},
		},
		{
			"search en",
			tid,
			"en",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{u5},
		},
		{
			"search email",
			tid,
			u1.Email,
			&model.UserSearchOptions{
				AllowEmails:    true,
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{u1},
		},
		{
			"search maps * to space",
			tid,
			"jimb*",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{u1},
		},
		{
			"should not return spurious matches",
			tid,
			"harol",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{},
		},
		{
			"% should be escaped",
			tid,
			"h%",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{},
		},
		{
			"_ should be escaped",
			tid,
			"h_",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{},
		},
		{
			"_ should be escaped (2)",
			tid,
			"Du_",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{u6},
		},
		{
			"_ should be escaped (2)",
			tid,
			"_dE",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{u6},
		},
		{
			"search jimb, allowing inactive",
			tid,
			"jimb",
			&model.UserSearchOptions{
				AllowFullNames: true,
				AllowInactive:  true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{u1, u3},
		},
		{
			"search jimb, no team id",
			"",
			"jimb",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{u1},
		},
		{
			"search jim-bobb, no team id",
			"",
			"jim-bobb",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{u2},
		},

		{
			"search harol, search all fields",
			tid,
			"harol",
			&model.UserSearchOptions{
				AllowFullNames: true,
				AllowEmails:    true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{u1},
		},
		{
			"search Tim, search all fields",
			tid,
			"Tim",
			&model.UserSearchOptions{
				AllowFullNames: true,
				AllowEmails:    true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{u1},
		},
		{
			"search Tim, don't search full names",
			tid,
			"Tim",
			&model.UserSearchOptions{
				AllowFullNames: false,
				AllowEmails:    true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{},
		},
		{
			"search Bill, search all fields",
			tid,
			"Bill",
			&model.UserSearchOptions{
				AllowFullNames: true,
				AllowEmails:    true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{u1},
		},
		{
			"search Rob, search all fields",
			tid,
			"Rob",
			&model.UserSearchOptions{
				AllowFullNames: true,
				AllowEmails:    true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{u1},
		},
		{
			"leading @ should be ignored",
			tid,
			"@jimb",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{u1},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Description, func(t *testing.T) {
			result := <-ss.User().Search(testCase.TeamId, testCase.Term, testCase.Options)
			require.Nil(t, result.Err)
			assertUsers(t, testCase.Expected, result.Data.([]*model.User))
		})
	}

	t.Run("search empty string", func(t *testing.T) {
		searchOptions := &model.UserSearchOptions{
			AllowFullNames: true,
			Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
		}

		r1 := <-ss.User().Search(tid, "", searchOptions)
		require.Nil(t, r1.Err)
		assert.Len(t, r1.Data.([]*model.User), 4)
		// Don't assert contents, since Postgres' default collation order is left up to
		// the operating system, and jimbo1 might sort before or after jim-bo.
		// assertUsers(t, []*model.User{u2, u1, u6, u5}, r1.Data.([]*model.User))
	})

	t.Run("search empty string, limit 2", func(t *testing.T) {
		searchOptions := &model.UserSearchOptions{
			AllowFullNames: true,
			Limit:          2,
		}

		r1 := <-ss.User().Search(tid, "", searchOptions)
		require.Nil(t, r1.Err)
		assert.Len(t, r1.Data.([]*model.User), 2)
		// Don't assert contents, since Postgres' default collation order is left up to
		// the operating system, and jimbo1 might sort before or after jim-bo.
		// assertUsers(t, []*model.User{u2, u1, u6, u5}, r1.Data.([]*model.User))
	})
}

func testUserStoreSearchNotInChannel(t *testing.T, ss store.Store) {
	u1 := &model.User{
		Username:  "jimbo1" + model.NewId(),
		FirstName: "Tim",
		LastName:  "Bill",
		Nickname:  "Rob",
		Email:     "harold" + model.NewId() + "@simulator.amazonses.com",
	}
	store.Must(ss.User().Save(u1))
	defer func() { store.Must(ss.User().PermanentDelete(u1.Id)) }()

	u2 := &model.User{
		Username: "jim2-bobby" + model.NewId(),
		Email:    MakeEmail(),
	}
	store.Must(ss.User().Save(u2))
	defer func() { store.Must(ss.User().PermanentDelete(u2.Id)) }()

	u3 := &model.User{
		Username: "jimbo3" + model.NewId(),
		Email:    MakeEmail(),
		DeleteAt: 1,
	}
	store.Must(ss.User().Save(u3))
	defer func() { store.Must(ss.User().PermanentDelete(u3.Id)) }()

	tid := model.NewId()
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: tid, UserId: u1.Id}, -1))
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: tid, UserId: u2.Id}, -1))
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: tid, UserId: u3.Id}, -1))

	// The users returned from the database will have AuthData as an empty string.
	nilAuthData := new(string)
	*nilAuthData = ""

	u1.AuthData = nilAuthData
	u2.AuthData = nilAuthData
	u3.AuthData = nilAuthData

	c1 := model.Channel{
		TeamId:      tid,
		DisplayName: "NameName",
		Name:        "zz" + model.NewId() + "b",
		Type:        model.CHANNEL_OPEN,
	}
	c1 = *store.Must(ss.Channel().Save(&c1, -1)).(*model.Channel)

	c2 := model.Channel{
		TeamId:      tid,
		DisplayName: "NameName",
		Name:        "zz" + model.NewId() + "b",
		Type:        model.CHANNEL_OPEN,
	}
	c2 = *store.Must(ss.Channel().Save(&c2, -1)).(*model.Channel)

	store.Must(ss.Channel().SaveMember(&model.ChannelMember{
		ChannelId:   c2.Id,
		UserId:      u1.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	}))
	store.Must(ss.Channel().SaveMember(&model.ChannelMember{
		ChannelId:   c1.Id,
		UserId:      u3.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	}))
	store.Must(ss.Channel().SaveMember(&model.ChannelMember{
		ChannelId:   c2.Id,
		UserId:      u2.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	}))

	testCases := []struct {
		Description string
		TeamId      string
		ChannelId   string
		Term        string
		Options     *model.UserSearchOptions
		Expected    []*model.User
	}{
		{
			"search jimb, channel 1",
			tid,
			c1.Id,
			"jimb",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{u1},
		},
		{
			"search jimb, allow inactive, channel 1",
			tid,
			c1.Id,
			"jimb",
			&model.UserSearchOptions{
				AllowFullNames: true,
				AllowInactive:  true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{u1},
		},
		{
			"search jimb, channel 1, no team id",
			"",
			c1.Id,
			"jimb",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{u1},
		},
		{
			"search jimb, channel 1, junk team id",
			"junk",
			c1.Id,
			"jimb",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{},
		},
		{
			"search jimb, channel 2",
			tid,
			c2.Id,
			"jimb",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{},
		},
		{
			"search jimb, allow inactive, channel 2",
			tid,
			c2.Id,
			"jimb",
			&model.UserSearchOptions{
				AllowFullNames: true,
				AllowInactive:  true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{u3},
		},
		{
			"search jimb, channel 2, no team id",
			"",
			c2.Id,
			"jimb",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{},
		},
		{
			"search jimb, channel 2, junk team id",
			"junk",
			c2.Id,
			"jimb",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{},
		},
		{
			"search jim, channel 1",
			tid,
			c1.Id,
			"jim",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{u2, u1},
		},
		{
			"search jim, channel 1, limit 1",
			tid,
			c1.Id,
			"jim",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          1,
			},
			[]*model.User{u2},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Description, func(t *testing.T) {
			result := <-ss.User().SearchNotInChannel(
				testCase.TeamId,
				testCase.ChannelId,
				testCase.Term,
				testCase.Options,
			)
			require.Nil(t, result.Err)
			assertUsers(t, testCase.Expected, result.Data.([]*model.User))
		})
	}
}

func testUserStoreSearchInChannel(t *testing.T, ss store.Store) {
	u1 := &model.User{
		Username:  "jimbo1" + model.NewId(),
		FirstName: "Tim",
		LastName:  "Bill",
		Nickname:  "Rob",
		Email:     "harold" + model.NewId() + "@simulator.amazonses.com",
	}
	store.Must(ss.User().Save(u1))
	defer func() { store.Must(ss.User().PermanentDelete(u1.Id)) }()

	u2 := &model.User{
		Username: "jim-bobby" + model.NewId(),
		Email:    MakeEmail(),
	}
	store.Must(ss.User().Save(u2))
	defer func() { store.Must(ss.User().PermanentDelete(u2.Id)) }()

	u3 := &model.User{
		Username: "jimbo3" + model.NewId(),
		Email:    MakeEmail(),
		DeleteAt: 1,
	}
	store.Must(ss.User().Save(u3))
	defer func() { store.Must(ss.User().PermanentDelete(u3.Id)) }()

	tid := model.NewId()
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: tid, UserId: u1.Id}, -1))
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: tid, UserId: u2.Id}, -1))
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: tid, UserId: u3.Id}, -1))

	// The users returned from the database will have AuthData as an empty string.
	nilAuthData := new(string)
	*nilAuthData = ""

	u1.AuthData = nilAuthData
	u2.AuthData = nilAuthData
	u3.AuthData = nilAuthData

	c1 := model.Channel{
		TeamId:      tid,
		DisplayName: "NameName",
		Name:        "zz" + model.NewId() + "b",
		Type:        model.CHANNEL_OPEN,
	}
	c1 = *store.Must(ss.Channel().Save(&c1, -1)).(*model.Channel)

	c2 := model.Channel{
		TeamId:      tid,
		DisplayName: "NameName",
		Name:        "zz" + model.NewId() + "b",
		Type:        model.CHANNEL_OPEN,
	}
	c2 = *store.Must(ss.Channel().Save(&c2, -1)).(*model.Channel)

	store.Must(ss.Channel().SaveMember(&model.ChannelMember{
		ChannelId:   c1.Id,
		UserId:      u1.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	}))
	store.Must(ss.Channel().SaveMember(&model.ChannelMember{
		ChannelId:   c2.Id,
		UserId:      u2.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	}))
	store.Must(ss.Channel().SaveMember(&model.ChannelMember{
		ChannelId:   c1.Id,
		UserId:      u3.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	}))

	testCases := []struct {
		Description string
		ChannelId   string
		Term        string
		Options     *model.UserSearchOptions
		Expected    []*model.User
	}{
		{
			"search jimb, channel 1",
			c1.Id,
			"jimb",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{u1},
		},
		{
			"search jimb, allow inactive, channel 1",
			c1.Id,
			"jimb",
			&model.UserSearchOptions{
				AllowFullNames: true,
				AllowInactive:  true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{u1, u3},
		},
		{
			"search jimb, allow inactive, channel 1, limit 1",
			c1.Id,
			"jimb",
			&model.UserSearchOptions{
				AllowFullNames: true,
				AllowInactive:  true,
				Limit:          1,
			},
			[]*model.User{u1},
		},
		{
			"search jimb, channel 2",
			c2.Id,
			"jimb",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{},
		},
		{
			"search jimb, allow inactive, channel 2",
			c2.Id,
			"jimb",
			&model.UserSearchOptions{
				AllowFullNames: true,
				AllowInactive:  true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Description, func(t *testing.T) {
			result := <-ss.User().SearchInChannel(
				testCase.ChannelId,
				testCase.Term,
				testCase.Options,
			)
			require.Nil(t, result.Err)
			assertUsers(t, testCase.Expected, result.Data.([]*model.User))
		})
	}
}

func testUserStoreSearchNotInTeam(t *testing.T, ss store.Store) {
	u1 := &model.User{
		Username:  "jimbo1" + model.NewId(),
		FirstName: "Tim",
		LastName:  "Bill",
		Nickname:  "Rob",
		Email:     "harold" + model.NewId() + "@simulator.amazonses.com",
	}
	store.Must(ss.User().Save(u1))
	defer func() { store.Must(ss.User().PermanentDelete(u1.Id)) }()

	u2 := &model.User{
		Username: "jim-bobby" + model.NewId(),
		Email:    MakeEmail(),
	}
	store.Must(ss.User().Save(u2))
	defer func() { store.Must(ss.User().PermanentDelete(u2.Id)) }()

	u3 := &model.User{
		Username: "jimbo3" + model.NewId(),
		Email:    MakeEmail(),
		DeleteAt: 1,
	}
	store.Must(ss.User().Save(u3))
	defer func() { store.Must(ss.User().PermanentDelete(u3.Id)) }()

	u4 := &model.User{
		Username: "simon" + model.NewId(),
		Email:    MakeEmail(),
		DeleteAt: 0,
	}
	store.Must(ss.User().Save(u4))
	defer func() { store.Must(ss.User().PermanentDelete(u4.Id)) }()

	u5 := &model.User{
		Username:  "yu" + model.NewId(),
		FirstName: "En",
		LastName:  "Yu",
		Nickname:  "enyu",
		Email:     MakeEmail(),
	}
	store.Must(ss.User().Save(u5))
	defer func() { store.Must(ss.User().PermanentDelete(u5.Id)) }()

	u6 := &model.User{
		Username:  "underscore" + model.NewId(),
		FirstName: "Du_",
		LastName:  "_DE",
		Nickname:  "lodash",
		Email:     MakeEmail(),
	}
	store.Must(ss.User().Save(u6))
	defer func() { store.Must(ss.User().PermanentDelete(u6.Id)) }()

	teamId1 := model.NewId()
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: teamId1, UserId: u1.Id}, -1))
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: teamId1, UserId: u2.Id}, -1))
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: teamId1, UserId: u3.Id}, -1))
	// u4 is not in team 1
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: teamId1, UserId: u5.Id}, -1))
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: teamId1, UserId: u6.Id}, -1))

	teamId2 := model.NewId()
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: teamId2, UserId: u4.Id}, -1))

	// The users returned from the database will have AuthData as an empty string.
	nilAuthData := new(string)
	*nilAuthData = ""

	u1.AuthData = nilAuthData
	u2.AuthData = nilAuthData
	u3.AuthData = nilAuthData
	u4.AuthData = nilAuthData
	u5.AuthData = nilAuthData
	u6.AuthData = nilAuthData

	testCases := []struct {
		Description string
		TeamId      string
		Term        string
		Options     *model.UserSearchOptions
		Expected    []*model.User
	}{
		{
			"search simo, team 1",
			teamId1,
			"simo",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{u4},
		},

		{
			"search jimb, team 1",
			teamId1,
			"jimb",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{},
		},
		{
			"search jimb, allow inactive, team 1",
			teamId1,
			"jimb",
			&model.UserSearchOptions{
				AllowFullNames: true,
				AllowInactive:  true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{},
		},
		{
			"search simo, team 2",
			teamId2,
			"simo",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{},
		},
		{
			"search jimb, team2",
			teamId2,
			"jimb",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{u1},
		},
		{
			"search jimb, allow inactive, team 2",
			teamId2,
			"jimb",
			&model.UserSearchOptions{
				AllowFullNames: true,
				AllowInactive:  true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{u1, u3},
		},
		{
			"search jimb, allow inactive, team 2, limit 1",
			teamId2,
			"jimb",
			&model.UserSearchOptions{
				AllowFullNames: true,
				AllowInactive:  true,
				Limit:          1,
			},
			[]*model.User{u1},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Description, func(t *testing.T) {
			result := <-ss.User().SearchNotInTeam(
				testCase.TeamId,
				testCase.Term,
				testCase.Options,
			)
			require.Nil(t, result.Err)
			assertUsers(t, testCase.Expected, result.Data.([]*model.User))
		})
	}
}

func testUserStoreSearchWithoutTeam(t *testing.T, ss store.Store) {
	u1 := &model.User{
		Username:  "jimbo1" + model.NewId(),
		FirstName: "Tim",
		LastName:  "Bill",
		Nickname:  "Rob",
		Email:     "harold" + model.NewId() + "@simulator.amazonses.com",
	}
	store.Must(ss.User().Save(u1))
	defer func() { store.Must(ss.User().PermanentDelete(u1.Id)) }()

	u2 := &model.User{
		Username: "jim2-bobby" + model.NewId(),
		Email:    MakeEmail(),
	}
	store.Must(ss.User().Save(u2))
	defer func() { store.Must(ss.User().PermanentDelete(u2.Id)) }()

	u3 := &model.User{
		Username: "jimbo3" + model.NewId(),
		Email:    MakeEmail(),
		DeleteAt: 1,
	}
	store.Must(ss.User().Save(u3))
	defer func() { store.Must(ss.User().PermanentDelete(u3.Id)) }()

	tid := model.NewId()
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: tid, UserId: u3.Id}, -1))

	// The users returned from the database will have AuthData as an empty string.
	nilAuthData := new(string)
	*nilAuthData = ""

	u1.AuthData = nilAuthData
	u2.AuthData = nilAuthData
	u3.AuthData = nilAuthData

	testCases := []struct {
		Description string
		Term        string
		Options     *model.UserSearchOptions
		Expected    []*model.User
	}{
		{
			"empty string",
			"",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{u2, u1},
		},
		{
			"jim",
			"jim",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{u2, u1},
		},
		{
			"PLT-8354",
			"* ",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{u2, u1},
		},
		{
			"jim, limit 1",
			"jim",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          1,
			},
			[]*model.User{u2},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Description, func(t *testing.T) {
			result := <-ss.User().SearchWithoutTeam(
				testCase.Term,
				testCase.Options,
			)
			require.Nil(t, result.Err)
			assertUsers(t, testCase.Expected, result.Data.([]*model.User))
		})
	}
}

func testUserStoreAnalyticsGetInactiveUsersCount(t *testing.T, ss store.Store) {
	u1 := &model.User{}
	u1.Email = MakeEmail()
	store.Must(ss.User().Save(u1))
	defer func() { store.Must(ss.User().PermanentDelete(u1.Id)) }()

	var count int64

	if result := <-ss.User().AnalyticsGetInactiveUsersCount(); result.Err != nil {
		t.Fatal(result.Err)
	} else {
		count = result.Data.(int64)
	}

	u2 := &model.User{}
	u2.Email = MakeEmail()
	u2.DeleteAt = model.GetMillis()
	store.Must(ss.User().Save(u2))
	defer func() { store.Must(ss.User().PermanentDelete(u2.Id)) }()

	if result := <-ss.User().AnalyticsGetInactiveUsersCount(); result.Err != nil {
		t.Fatal(result.Err)
	} else {
		newCount := result.Data.(int64)
		if count != newCount-1 {
			t.Fatal("Expected 1 more inactive users but found otherwise.", count, newCount)
		}
	}
}

func testUserStoreAnalyticsGetSystemAdminCount(t *testing.T, ss store.Store) {
	var countBefore int64
	if result := <-ss.User().AnalyticsGetSystemAdminCount(); result.Err != nil {
		t.Fatal(result.Err)
	} else {
		countBefore = result.Data.(int64)
	}

	u1 := model.User{}
	u1.Email = MakeEmail()
	u1.Username = model.NewId()
	u1.Roles = "system_user system_admin"

	u2 := model.User{}
	u2.Email = MakeEmail()
	u2.Username = model.NewId()

	if err := (<-ss.User().Save(&u1)).Err; err != nil {
		t.Fatal("couldn't save user", err)
	}
	defer func() { store.Must(ss.User().PermanentDelete(u1.Id)) }()

	if err := (<-ss.User().Save(&u2)).Err; err != nil {
		t.Fatal("couldn't save user", err)
	}
	defer func() { store.Must(ss.User().PermanentDelete(u2.Id)) }()

	if result := <-ss.User().AnalyticsGetSystemAdminCount(); result.Err != nil {
		t.Fatal(result.Err)
	} else {
		// We expect to find 1 more system admin than there was at the start of this test function.
		if count := result.Data.(int64); count != countBefore+1 {
			t.Fatal("Did not get the expected number of system admins. Expected, got: ", countBefore+1, count)
		}
	}
}

func testUserStoreGetProfilesNotInTeam(t *testing.T, ss store.Store) {
	teamId := model.NewId()

	u1 := &model.User{}
	u1.Email = MakeEmail()
	store.Must(ss.User().Save(u1))
	defer func() { store.Must(ss.User().PermanentDelete(u1.Id)) }()
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, -1))
	store.Must(ss.User().UpdateUpdateAt(u1.Id))

	u2 := &model.User{}
	u2.Email = MakeEmail()
	store.Must(ss.User().Save(u2))
	defer func() { store.Must(ss.User().PermanentDelete(u2.Id)) }()
	store.Must(ss.User().UpdateUpdateAt(u2.Id))

	var initialUsersNotInTeam int
	var etag1, etag2, etag3 string

	if er1 := <-ss.User().GetEtagForProfilesNotInTeam(teamId); er1.Err != nil {
		t.Fatal(er1.Err)
	} else {
		etag1 = er1.Data.(string)
	}

	if r1 := <-ss.User().GetProfilesNotInTeam(teamId, 0, 100000); r1.Err != nil {
		t.Fatal(r1.Err)
	} else {
		users := r1.Data.([]*model.User)
		initialUsersNotInTeam = len(users)
		if initialUsersNotInTeam < 1 {
			t.Fatalf("Should be at least 1 user not in the team")
		}

		found := false
		for _, u := range users {
			if u.Id == u2.Id {
				found = true
			}
			if u.Id == u1.Id {
				t.Fatalf("Should not have found user1")
			}
		}

		if !found {
			t.Fatal("missing user2")
		}
	}

	time.Sleep(time.Millisecond * 10)
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u2.Id}, -1))
	store.Must(ss.User().UpdateUpdateAt(u2.Id))

	if er2 := <-ss.User().GetEtagForProfilesNotInTeam(teamId); er2.Err != nil {
		t.Fatal(er2.Err)
	} else {
		etag2 = er2.Data.(string)
		if etag1 == etag2 {
			t.Fatalf("etag should have changed")
		}
	}

	if r2 := <-ss.User().GetProfilesNotInTeam(teamId, 0, 100000); r2.Err != nil {
		t.Fatal(r2.Err)
	} else {
		users := r2.Data.([]*model.User)

		if len(users) != initialUsersNotInTeam-1 {
			t.Fatalf("Should be one less user not in team")
		}

		for _, u := range users {
			if u.Id == u2.Id {
				t.Fatalf("Should not have found user2")
			}
			if u.Id == u1.Id {
				t.Fatalf("Should not have found user1")
			}
		}
	}

	time.Sleep(time.Millisecond * 10)
	store.Must(ss.Team().RemoveMember(teamId, u1.Id))
	store.Must(ss.Team().RemoveMember(teamId, u2.Id))
	store.Must(ss.User().UpdateUpdateAt(u1.Id))
	store.Must(ss.User().UpdateUpdateAt(u2.Id))

	if er3 := <-ss.User().GetEtagForProfilesNotInTeam(teamId); er3.Err != nil {
		t.Fatal(er3.Err)
	} else {
		etag3 = er3.Data.(string)
		t.Log(etag3)
		if etag1 == etag3 || etag3 == etag2 {
			t.Fatalf("etag should have changed")
		}
	}

	if r3 := <-ss.User().GetProfilesNotInTeam(teamId, 0, 100000); r3.Err != nil {
		t.Fatal(r3.Err)
	} else {
		users := r3.Data.([]*model.User)
		found1, found2 := false, false
		for _, u := range users {
			if u.Id == u2.Id {
				found2 = true
			}
			if u.Id == u1.Id {
				found1 = true
			}
		}

		if !found1 || !found2 {
			t.Fatal("missing user1 or user2")
		}
	}

	time.Sleep(time.Millisecond * 10)
	u3 := &model.User{}
	u3.Email = MakeEmail()
	store.Must(ss.User().Save(u3))
	defer func() { store.Must(ss.User().PermanentDelete(u3.Id)) }()
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u3.Id}, -1))
	store.Must(ss.User().UpdateUpdateAt(u3.Id))

	if er4 := <-ss.User().GetEtagForProfilesNotInTeam(teamId); er4.Err != nil {
		t.Fatal(er4.Err)
	} else {
		etag4 := er4.Data.(string)
		t.Log(etag4)
		if etag4 != etag3 {
			t.Fatalf("etag should be the same")
		}
	}
}

func testUserStoreClearAllCustomRoleAssignments(t *testing.T, ss store.Store) {
	u1 := model.User{
		Email:    MakeEmail(),
		Username: model.NewId(),
		Roles:    "system_user system_admin system_post_all",
	}
	u2 := model.User{
		Email:    MakeEmail(),
		Username: model.NewId(),
		Roles:    "system_user custom_role system_admin another_custom_role",
	}
	u3 := model.User{
		Email:    MakeEmail(),
		Username: model.NewId(),
		Roles:    "system_user",
	}
	u4 := model.User{
		Email:    MakeEmail(),
		Username: model.NewId(),
		Roles:    "custom_only",
	}

	store.Must(ss.User().Save(&u1))
	defer func() { store.Must(ss.User().PermanentDelete(u1.Id)) }()
	store.Must(ss.User().Save(&u2))
	defer func() { store.Must(ss.User().PermanentDelete(u2.Id)) }()
	store.Must(ss.User().Save(&u3))
	defer func() { store.Must(ss.User().PermanentDelete(u3.Id)) }()
	store.Must(ss.User().Save(&u4))
	defer func() { store.Must(ss.User().PermanentDelete(u4.Id)) }()

	require.Nil(t, (<-ss.User().ClearAllCustomRoleAssignments()).Err)

	r1 := <-ss.User().GetByUsername(u1.Username)
	require.Nil(t, r1.Err)
	assert.Equal(t, u1.Roles, r1.Data.(*model.User).Roles)

	r2 := <-ss.User().GetByUsername(u2.Username)
	require.Nil(t, r2.Err)
	assert.Equal(t, "system_user system_admin", r2.Data.(*model.User).Roles)

	r3 := <-ss.User().GetByUsername(u3.Username)
	require.Nil(t, r3.Err)
	assert.Equal(t, u3.Roles, r3.Data.(*model.User).Roles)

	r4 := <-ss.User().GetByUsername(u4.Username)
	require.Nil(t, r4.Err)
	assert.Equal(t, "", r4.Data.(*model.User).Roles)
}

func testUserStoreGetAllAfter(t *testing.T, ss store.Store) {
	u1 := store.Must(ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: model.NewId(),
		Roles:    "system_user system_admin system_post_all",
	})).(*model.User)
	defer func() { store.Must(ss.User().PermanentDelete(u1.Id)) }()

	u2 := store.Must(ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u2" + model.NewId(),
	})).(*model.User)
	store.Must(ss.Bot().Save(&model.Bot{
		UserId:    u2.Id,
		Username:  u2.Username,
		CreatorId: u1.Id,
	}))
	u2.IsBot = true
	defer func() { store.Must(ss.Bot().PermanentDelete(u2.Id)) }()
	defer func() { store.Must(ss.User().PermanentDelete(u2.Id)) }()

	expected := []*model.User{u1, u2}
	if strings.Compare(u2.Id, u1.Id) < 0 {
		expected = []*model.User{u2, u1}
	}

	t.Run("get after lowest possible id", func(t *testing.T) {
		result := <-ss.User().GetAllAfter(10000, strings.Repeat("0", 26))
		require.Nil(t, result.Err)

		actual := result.Data.([]*model.User)
		assert.Equal(t, expected, actual)
	})

	t.Run("get after first user", func(t *testing.T) {
		result := <-ss.User().GetAllAfter(10000, expected[0].Id)
		require.Nil(t, result.Err)

		actual := result.Data.([]*model.User)
		assert.Equal(t, []*model.User{expected[1]}, actual)
	})

	t.Run("get after second user", func(t *testing.T) {
		result := <-ss.User().GetAllAfter(10000, expected[1].Id)
		require.Nil(t, result.Err)

		actual := result.Data.([]*model.User)
		assert.Equal(t, []*model.User{}, actual)
	})
}
