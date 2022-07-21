package main

import (
	"github.com/stretchr/gomniauth/test"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar
	testUser := &test.TestUser{}
	testUser.On("AvatarURL").Return("", ErrNoAvatarURL)
	testChatUser := &chatUser{User: testUser}
	url, err := authAvatar.GetAvatarURL(testChatUser)
	if err != ErrNoAvatarURL {
		t.Error("AuthAvatar.GetAvatarURL should return ErrNoAvatarURL when no value present")
	}
	testUrl := "http://url-to-avatar/"
	testUser = &test.TestUser{}
	testChatUser.User = testUser
	testUser.On("AvatarURL").Return(testUrl, nil)
	url, err = authAvatar.GetAvatarURL(testChatUser)
	if err != nil {
		t.Error("AuthAvatar.GetAvatarURL should not return an error")
	} else {
		if url != testUrl {
			t.Errorf("AuthAvatar.GetAvatarURL wrongly returned %s", url)
		}
	}
}
func TestGravatarAvatar(t *testing.T) {
	var gravatarAvitar GravatarAvatar
	user := &chatUser{uniqueID: "abc"}
	url, err := gravatarAvitar.GetAvatarURL(user)
	if err != nil {
		t.Error("GravatarAvitar.GetAvatarURL should not return an 	error")
	}
	if url != "//www.gravatar.com/avatar/abc" {
		t.Errorf("GravatarAvitar.GetAvatarURL wrongly returned %s", url)
	}
}
func TestFileSystemAvatar(t *testing.T) {
	// make a test avatar file
	filename := path.Join("avatars", "abc.jpg")
	ioutil.WriteFile(filename, []byte{}, 0777)
	defer func() { os.Remove(filename) }()
	var fileSystemAvatar FileSystemAvatar
	user := &chatUser{uniqueID: "abc"}
	url, err := fileSystemAvatar.GetAvatarURL(user)
	if err != nil {
		t.Error("FileSystemAvatar.GetAvatarURL should not return an error")
	}
	if url != "/avatars/abc.jpg" {
		t.Errorf("FileSystemAvatar.GetAvatarURL wrongly returned %s",
			url)
	}
}
