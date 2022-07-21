package main

import (
	"errors"
	"io/ioutil"
	"path"
)

var ErrNoAvatarURL = errors.New("Unable to get an avatar URL")

type Avatar interface {
	GetAvatarURL(c ChatUser) (string, error)
}

type AuthAvatar struct {
}

var UseAuthAvatar AuthAvatar

func (_ AuthAvatar) GetAvatarURL(u ChatUser) (string, error) {
	url := u.AvatarURL()
	if len(url) > 0 {
		return url, nil
	}
	return "", ErrNoAvatarURL
}

type GravatarAvatar struct {
}

var UseGravatarAvatar GravatarAvatar

func (_ GravatarAvatar) GetAvatarURL(u ChatUser) (string, error) {
	return "//www.gravatar.com/avatar/" + u.UniqueID(), nil
}

var UseFileSystemAvatar FileSystemAvatar

type FileSystemAvatar struct {
}

func (_ FileSystemAvatar) GetAvatarURL(u ChatUser) (string, error) {
	if files, err := ioutil.ReadDir("avatars"); err == nil {
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if match, _ := path.Match(u.UniqueID()+"*", file.Name()); match {
				return "/avatars/" + file.Name(), nil
			}
		}
	}
	return "", ErrNoAvatarURL
}

type TryAvatars []Avatar

func (t TryAvatars) GetAvatarURL(c ChatUser) (string, error) {
	for _, avatar := range t {
		if url, err := avatar.GetAvatarURL(c); err == nil {
			return url, nil
		}
	}
	return "", ErrNoAvatarURL
}
