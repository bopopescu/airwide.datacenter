/*
 *  Copyright (c) 2018, https://github.com/airwide-code/airwide.datacenter
 *  All rights reserved.
 *
 *
 *
 */

package rpc

import (
	"github.com/golang/glog"
	"github.com/airwide-code/airwide.datacenter/baselib/logger"
	"github.com/airwide-code/airwide.datacenter/baselib/grpc_util"
	"github.com/airwide-code/airwide.datacenter/mtproto"
	"golang.org/x/net/context"
	"github.com/airwide-code/airwide.datacenter/biz/base"
	contact2 "github.com/airwide-code/airwide.datacenter/biz/core/contact"
	"github.com/airwide-code/airwide.datacenter/biz/core/user"
	"github.com/airwide-code/airwide.datacenter/biz_server/sync_client"
	updates2 "github.com/airwide-code/airwide.datacenter/biz/core/update"
)

// contacts.importContacts#2c800be5 contacts:Vector<InputContact> = contacts.ImportedContacts;
func (s *ContactsServiceImpl) ContactsImportContacts(ctx context.Context, request *mtproto.TLContactsImportContacts) (*mtproto.Contacts_ImportedContacts, error) {
	md := grpc_util.RpcMetadataFromIncoming(ctx)
	glog.Infof("contacts.importContacts#2c800be5 - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

	var (
		err error
		importedContacts *mtproto.TLContactsImportedContacts
	)

	if len(request.Contacts) == 0 {
		err := mtproto.NewRpcError2(mtproto.TLRpcErrorCodes_BAD_REQUEST)
		glog.Error(err, ": contacts empty")
		return nil, err
	}

	// 注意: 目前只支持导入1个并且已经注册的联系人!!!!
	inputContact := request.Contacts[0].To_InputPhoneContact()

	if inputContact.GetFirstName() == "" {
		err := mtproto.NewRpcError2(mtproto.TLRpcErrorCodes_FIRSTNAME_INVALID)
		glog.Error(err, ": first_name is empty")
		return nil, err
	}

	pnumber, err := base.MakePhoneNumberUtil(inputContact.GetPhone(), "")
	if err != nil {
		region := user.GetCountryCodeByUser(md.UserId)
		pnumber, err = base.MakePhoneNumberUtil(inputContact.GetPhone(), region)
	}

	// var contactUser
	// phone, err := base.CheckAndGetPhoneNumber(inputContact.GetPhone())
	if err != nil {
		glog.Error(err)
		err := mtproto.NewRpcError2(mtproto.TLRpcErrorCodes_PHONE_CODE_INVALID)
		glog.Error(err, ": phone code invalid - ", inputContact.GetPhone())
		return nil, err
	}

	phoneNumber := pnumber.GetNormalizeDigits()
	contactUser := user.GetUserByPhoneNumber(md.UserId, phoneNumber)
	if contactUser == nil {
		// 该手机号未注册，我们认为手机号出错
		//err := mtproto.NewRpcError2(mtproto.TLRpcErrorCodes_PHONE_CODE_INVALID)
		//glog.Error(err, ": phone code invalid - ", inputContact.GetPhone())
		//return nil, err

		importedContacts = &mtproto.TLContactsImportedContacts{Data2: &mtproto.Contacts_ImportedContacts_Data{
			Imported:       []*mtproto.ImportedContact{},
			PopularInvites: []*mtproto.PopularContact{},
			RetryContacts:  []int64{},
			Users:          []*mtproto.User{},
		}}

		glog.Infof("contacts.importContacts#2c800be5 - reply: %s", logger.JsonDebugData(importedContacts))
		return importedContacts.To_Contacts_ImportedContacts(), nil
	}
	// contactUser.SetContact(true)
	// contactUser.SetMutualContact(true)
	contactLogic := contact2.MakeContactLogic(md.UserId)
	needUpdate := contactLogic.ImportContact(contactUser.GetId(), phoneNumber, inputContact.GetFirstName(), inputContact.GetLastName())
	// _ = needUpdate

	selfUpdates := updates2.NewUpdatesLogic(md.UserId)
	contactLink := &mtproto.TLUpdateContactLink{Data2: &mtproto.Update_Data{
		UserId:      contactUser.GetId(),
		MyLink:      mtproto.NewTLContactLinkContact().To_ContactLink(),
		ForeignLink: mtproto.NewTLContactLinkContact().To_ContactLink(),
	}}
	selfUpdates.AddUpdate(contactLink.To_Update())

	// 添加成功后是我的联系人
	contactUser.SetContact(true)
	selfUpdates.AddUser(contactUser.To_User())
	// TODO(@benqi): handle seq
	sync_client.GetSyncClient().PushToUserNotMeUpdatesData(md.AuthId, md.SessionId, md.UserId, selfUpdates.ToUpdates())

	// TODO(@benqi): 推给联系人逻辑需要再考虑考虑
	if needUpdate {
		// TODO(@benqi): push to contact user update contact link
		contactUpdates := updates2.NewUpdatesLogic(contactUser.GetId())
		contactLink2 := &mtproto.TLUpdateContactLink{Data2: &mtproto.Update_Data{
			UserId:      md.UserId,
			MyLink:      mtproto.NewTLContactLinkContact().To_ContactLink(),
			ForeignLink: mtproto.NewTLContactLinkContact().To_ContactLink(),
		}}
		contactUpdates.AddUpdate(contactLink2.To_Update())

		myUser := user.GetUserById(contactUser.GetId(), md.UserId)
		contactUpdates.AddUser(myUser.To_User())
		// TODO(@benqi): handle seq
		sync_client.GetSyncClient().PushToUserUpdatesData(contactUser.GetId(), contactUpdates.ToUpdates())
	}

	////////////////////////////////////////////////////////////////////////////////////////
	imported := &mtproto.TLImportedContact{Data2: &mtproto.ImportedContact_Data{
		UserId: contactUser.GetId(),
		ClientId: inputContact.GetClientId(),
	}}
	// contacts.importedContacts#77d01c3b imported:Vector<ImportedContact> popular_invites:Vector<PopularContact> retry_contacts:Vector<long> users:Vector<User> = contacts.ImportedContacts;
	importedContacts = &mtproto.TLContactsImportedContacts{Data2: &mtproto.Contacts_ImportedContacts_Data{
		Imported: []*mtproto.ImportedContact{imported.To_ImportedContact()},
		Users: []*mtproto.User{contactUser.To_User()},
	}}

	glog.Infof("contacts.importContacts#2c800be5 - reply: %s", logger.JsonDebugData(importedContacts))
	return importedContacts.To_Contacts_ImportedContacts(), nil
}
