package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path"
)

var ErrNoAvatarURL = errors.New("Unable to get an avatar URL")

type Avatar interface {
	GetAvatarURL(c *client) (string, error)
}

type AuthAvatar struct {
}

var UseAuthAvatar AuthAvatar

func (AuthAvatar) GetAvatarURL(c *client) (string, error) {
	if url, ok := c.userData["avatar_url"]; ok {
		return url.(string), nil
	}
	return "", ErrNoAvatarURL
}

type GravatarAvatar struct {
}

var UseGravatarAvatar GravatarAvatar

func (a GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	if userId, ok := c.userData["userId"]; ok {
		if userIdStr, ok := userId.(string); ok {
			return fmt.Sprintf("//www.gravatar.com/avatar/%s", userIdStr), nil
		}
	}
	return "", ErrNoAvatarURL
}

var UseFileSystemAvatar FileSystemAvatar

type FileSystemAvatar struct {
}

func (FileSystemAvatar) GetAvatarURL(c *client) (string, error) {
	if userID, ok := c.userData["userId"]; ok {
		if userIdStr, ok := userID.(string); ok {
			if files, err := ioutil.ReadDir("avatars"); err == nil {
				for _, file := range files {
					if file.IsDir() {
						continue
					}
					if match, _ := path.Match(userIdStr+"*", file.Name()); match {
						return fmt.Sprintf("/avatars/" + file.Name()), nil
					}
				}
			}
		}
	}
	return "", ErrNoAvatarURL
}
