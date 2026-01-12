package core

import "errors"

var ErrNotAuthorized = errors.New("пользователь не авторизован для выполнения этого действия")
var ErrOrgNotFound = errors.New("организация не найдена")
var ErrFileNotFound = errors.New("файл не найден")
