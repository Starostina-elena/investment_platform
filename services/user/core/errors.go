package core

import "errors"

var ErrNicknameExists = errors.New("Такое имя пользователя уже существует")
var ErrEmailExists = errors.New("Такой email уже используется")
