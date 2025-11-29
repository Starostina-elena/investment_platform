package core

import "errors"

var ErrNicknameExists = errors.New("такое имя пользователя уже существует")
var ErrEmailExists = errors.New("такой email уже используется")
