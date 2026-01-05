package core

import "errors"

var ErrNicknameExists = errors.New("такое имя пользователя уже существует")
var ErrEmailExists = errors.New("такой email уже используется")
var ErrUserNotFound = errors.New("пользователь не найден")
var ErrInvalidToken = errors.New("невалидный токен")
