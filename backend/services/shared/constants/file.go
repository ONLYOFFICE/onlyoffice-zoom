package constants

import (
	"errors"
	"strings"
)

var ErrOnlyofficeExtensionNotSupported = errors.New("file extension is not supported")

const (
	_OnlyofficeWordType  string = "word"
	_OnlyofficeCellType  string = "cell"
	_OnlyofficeSlideType string = "slide"
)

var OnlyofficeEditableExtensions map[string]string = map[string]string{
	"xlsx": _OnlyofficeCellType,
	"pptx": _OnlyofficeSlideType,
	"docx": _OnlyofficeWordType,
}

var OnlyofficeFileExtensions map[string]string = map[string]string{
	"xls":  _OnlyofficeCellType,
	"xlsx": _OnlyofficeCellType,
	"xlsm": _OnlyofficeCellType,
	"xlt":  _OnlyofficeCellType,
	"xltx": _OnlyofficeCellType,
	"xltm": _OnlyofficeCellType,
	"ods":  _OnlyofficeCellType,
	"fods": _OnlyofficeCellType,
	"ots":  _OnlyofficeCellType,
	"csv":  _OnlyofficeCellType,
	"pps":  _OnlyofficeSlideType,
	"ppsx": _OnlyofficeSlideType,
	"ppsm": _OnlyofficeSlideType,
	"ppt":  _OnlyofficeSlideType,
	"pptx": _OnlyofficeSlideType,
	"pptm": _OnlyofficeSlideType,
	"pot":  _OnlyofficeSlideType,
	"potx": _OnlyofficeSlideType,
	"potm": _OnlyofficeSlideType,
	"odp":  _OnlyofficeSlideType,
	"fodp": _OnlyofficeSlideType,
	"otp":  _OnlyofficeSlideType,
	"doc":  _OnlyofficeWordType,
	"docx": _OnlyofficeWordType,
	"docm": _OnlyofficeWordType,
	"dot":  _OnlyofficeWordType,
	"dotx": _OnlyofficeWordType,
	"dotm": _OnlyofficeWordType,
	"odt":  _OnlyofficeWordType,
	"fodt": _OnlyofficeWordType,
	"ott":  _OnlyofficeWordType,
	"rtf":  _OnlyofficeWordType,
	"txt":  _OnlyofficeWordType,
	"html": _OnlyofficeWordType,
	"htm":  _OnlyofficeWordType,
	"mht":  _OnlyofficeWordType,
	"pdf":  _OnlyofficeWordType,
	"djvu": _OnlyofficeWordType,
	"fb2":  _OnlyofficeWordType,
	"epub": _OnlyofficeWordType,
	"xps":  _OnlyofficeWordType,
}

func IsExtensionSupported(fileExt string) bool {
	_, exists := OnlyofficeFileExtensions[strings.ToLower(fileExt)]
	if exists {
		return true
	}
	return false
}

func IsExtensionEditable(fileExt string) bool {
	_, exists := OnlyofficeEditableExtensions[strings.ToLower(fileExt)]
	if exists {
		return true
	}
	return false
}

func GetFileType(fileExt string) (string, error) {
	fileType, exists := OnlyofficeFileExtensions[strings.ToLower(fileExt)]
	if !exists {
		return "", ErrOnlyofficeExtensionNotSupported
	}
	return fileType, nil
}
