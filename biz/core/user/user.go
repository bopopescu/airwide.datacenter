/*
 *  Copyright (c) 2018, https://github.com/airwide-code
 *  All rights reserved.
 *
 *
 *
 */

package user

import (
	"github.com/airwide-code/airwide.datacenter/biz/dal/dao"
	"github.com/airwide-code/airwide.datacenter/mtproto"
	"github.com/airwide-code/airwide.datacenter/biz/dal/dataobject"
	"github.com/airwide-code/airwide.datacenter/baselib/base"
	"time"
)

//type accountData struct {
//	userId int32
//	banned bool
//}
//
//func (this *accountData) IsBanned() bool {
//	return this.banned
//}

type userData struct {
	*mtproto.TLUser
}

func (this *userData) ToUser() *mtproto.User {
	return this.TLUser.To_User()
}

//func CheckBannedByPhoneNumber(phoneNumber string) bool {
//	do := dao.GetUsersDAO(dao.DB_SLAVE).SelectByPhoneNumber(phoneNumber)
//	return do != nil && do.Banned != 0
//}
//
//func GetUser(userId int32) (user* mtproto.TLUser) {
//	usersDAO := dao.GetUsersDAO(dao.DB_SLAVE)
//	userDO := usersDAO.SelectById(userId)
//
//	if userDO != nil {
//		// TODO(@benqi): fill bot, photo, about...
//		user = &mtproto.TLUser{ Data2: &mtproto.User_Data{
//			// user.Self由业务层进行判断
//			// user.Self = true
//			Id: userDO.Id,
//			AccessHash: userDO.AccessHash,
//			FirstName: userDO.FirstName,
//			LastName: userDO.LastName,
//			Username: userDO.Username,
//			Phone: userDO.Phone,
//		}}
//	}
//	return
//}

func GetUsersBySelfAndIDList(selfUserId int32, userIdList []int32) (users []*mtproto.User) {
	if len(userIdList) == 0 {
		users = []*mtproto.User{}
	} else {
		// TODO(@benqi):  需要优化，makeUserDataByDO需要查询用户状态以及获取Mutual和Contact状态信息而导致多次查询
		userDOList := dao.GetUsersDAO(dao.DB_SLAVE).SelectUsersByIdList(userIdList)
		users = make([]*mtproto.User, 0, len(userDOList))
		for i := 0; i < len(userDOList); i++ {
			user := makeUserDataByDO(selfUserId, &userDOList[i])
			//
			//// TODO(@benqi): fill bot, photo, about...
			//user := &mtproto.TLUser{Data2: &mtproto.User_Data{
			//	Self:          selfUserId == userDO.Id,
			//	Id:            userDO.Id,
			//	AccessHash:    userDO.AccessHash,
			//	FirstName:     userDO.FirstName,
			//	LastName:      userDO.LastName,
			//	Username:      userDO.Username,
			//	Phone:         userDO.Phone,
			//	Contact:       true,
			//	MutualContact: true,
			//}}
			users = append(users, user.To_User())
		}
	}
	return
}


func GetUserFull(selfUserId int32, userId int32) (userFull *mtproto.TLUserFull) {
	//TODO(@benqi): 等Link和NotifySettings实现后再来完善
	//fullUser := &mtproto.TLUserFull{}
	//fullUser.PhoneCallsAvailable = true
	//fullUser.PhoneCallsPrivate = true
	//fullUser.About = "@Benqi"
	//fullUser.CommonChatsCount = 0
	//
	//switch request.Id.Payload.(type) {
	//case *mtproto.InputUser_InputUserSelf:
	//	// User
	//	userDO, _ := s.UsersDAO.SelectById(2)
	//	user := &mtproto.TLUser{}
	//	user.Self = true
	//	user.Contact = false
	//	user.Id = userDO.Id
	//	user.FirstName = userDO.FirstName
	//	user.LastName = userDO.LastName
	//	user.Username = userDO.Username
	//	user.AccessHash = userDO.AccessHash
	//	user.Phone = userDO.Phone
	//	fullUser.User = mtproto.MakeUser(user)
	//
	//	// Link
	//	link := &mtproto.TLContactsLink{}
	//	link.MyLink = mtproto.MakeContactLink(&mtproto.TLContactLinkContact{})
	//	link.ForeignLink = mtproto.MakeContactLink(&mtproto.TLContactLinkContact{})
	//	link.User = mtproto.MakeUser(user)
	//	fullUser.Link = mtproto.MakeContacts_Link(link)
	//case *mtproto.InputUser_InputUser:
	//	inputUser := request.Id.Payload.(*mtproto.InputUser_InputUser).InputUser
	//	// User
	//	userDO, _ := s.UsersDAO.SelectById(inputUser.UserId)
	//	user := &mtproto.TLUser{}
	//	user.Self = false
	//	user.Contact = true
	//	user.Id = userDO.Id
	//	user.FirstName = userDO.FirstName
	//	user.LastName = userDO.LastName
	//	user.Username = userDO.Username
	//	user.AccessHash = userDO.AccessHash
	//	user.Phone = userDO.Phone
	//	fullUser.User = mtproto.MakeUser(user)
	//
	//	// Link
	//	link := &mtproto.TLContactsLink{}
	//	link.MyLink = mtproto.MakeContactLink(&mtproto.TLContactLinkContact{})
	//	link.ForeignLink = mtproto.MakeContactLink(&mtproto.TLContactLinkContact{})
	//	link.User = mtproto.MakeUser(user)
	//	fullUser.Link = mtproto.MakeContacts_Link(link)
	//case *mtproto.InputUser_InputUserEmpty:
	//	// TODO(@benqi): BAD_REQUEST: 400
	//}
	//
	//// NotifySettings
	//peerNotifySettings := &mtproto.TLPeerNotifySettings{}
	//peerNotifySettings.ShowPreviews = true
	//peerNotifySettings.MuteUntil = 0
	//peerNotifySettings.Sound = "default"
	//fullUser.NotifySettings = mtproto.MakePeerNotifySettings(peerNotifySettings)
	return nil
}

func UpdateUserStatus(userId int32, lastSeenAt int64) {
	presencesDAO := dao.GetUserPresencesDAO(dao.DB_MASTER)
	// now := time.Now().Unix()
	rows := presencesDAO.UpdateLastSeen(lastSeenAt, 0, userId)
	if rows == 0 {
		do := &dataobject.UserPresencesDO{
			UserId: userId,
			LastSeenAt: lastSeenAt,
			LastSeenAuthKeyId: 0,
			LastSeenIp: "",
			CreatedAt: base.NowFormatYMDHMS(),
		}
		presencesDAO.Insert(do)
	}
}

func GetUserStatus(userId int32) *mtproto.UserStatus {
	now := time.Now().Unix()
	do := dao.GetUserPresencesDAO(dao.DB_SLAVE).SelectByUserID(userId)
	if do == nil {
		return mtproto.NewTLUserStatusEmpty().To_UserStatus()
	}

	if now <= do.LastSeenAt + 5*60 {
		status := &mtproto.TLUserStatusOnline{Data2: &mtproto.UserStatus_Data{
			Expires: int32(do.LastSeenAt + 5*30),
		}}
		return status.To_UserStatus()
	} else {
		status := &mtproto.TLUserStatusOffline{Data2: &mtproto.UserStatus_Data{
			WasOnline: int32(do.LastSeenAt),
		}}
		return status.To_UserStatus()
	}
}

