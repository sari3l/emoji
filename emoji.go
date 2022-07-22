package emoji

import (
	"github.com/brahma-adshonor/gohook"
	"log"
	"reflect"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

func doReplaceEmoji[E rune | byte](p []E) []E {
	pLen := len(p)
	var waitForWriteEs []E
	var emojiEs []E
	lock := false
	for i := 0; i < pLen; i++ {
		if p[i] != 58 && !lock {
			waitForWriteEs = append(waitForWriteEs, p[i])
			continue
		}
		if p[i] == 58 && !lock {
			lock = true
			emojiEs = append(emojiEs, p[i])
			continue
		}
		if p[i] == 58 && lock {
			if len(emojiEs) == 1 {
				waitForWriteEs = append(waitForWriteEs, emojiEs...)
				continue
			}
			lock = false
			emojiEs = append(emojiEs, p[i])
			emojiStr := emojiMap[esToString(emojiEs)]
			if emojiStr != "" {
				waitForWriteEs = append(waitForWriteEs, stringToEs(emojiStr, emojiEs[0])...)
			} else {
				waitForWriteEs = append(waitForWriteEs, emojiEs...)
			}
			emojiEs = []E{}
			continue
		}
		emojiEs = append(emojiEs, p[i])
		if p[i] == 32 {
			lock = false
			waitForWriteEs = append(waitForWriteEs, emojiEs...)
			emojiEs = []E{}
		}
	}
	return waitForWriteEs
}

func esToString[E byte | rune](data []E) string {
	s := ""
	for _, d := range data {
		s += string(d)
	}
	return s
}

func stringToEs[E byte | rune](data string, typeFlag E) []E {
	e := make([]E, 0)
	switch reflect.TypeOf(typeFlag).Kind() {
	case reflect.Uint8:
		dE := ([]byte)(data)
		for _, d := range dE {
			e = append(e, (E)(d))
		}
		return e
	case reflect.Uint32:
		dE := ([]rune)(data)
		for _, d := range dE {
			e = append(e, (E)(d))
		}
		return e
	}
	return nil
}

func emojiHookWindows(console syscall.Handle, buf *uint16, towrite uint32, written *uint32, reserved *byte) (err error) {
	bufSlice := make([]uint16, 0)
	//ptr := *buf
	for i := 0; i < int(towrite); i++ {
		c := unsafe.Pointer(uintptr(unsafe.Pointer(buf)) + uintptr(i)*unsafe.Sizeof(*buf))
		bufSlice = append(bufSlice, *(*uint16)(c))
	}
	chunk := utf16.Decode(bufSlice)
	p := utf16.Encode(doReplaceEmoji(chunk))
	bufPtr := &buf
	*bufPtr = &p[0]
	err = emojiHookWindowsTramp(console, buf, uint32(len(p)), written, reserved)
	if err != nil {
		return err
	} else {
		*written = uint32(len(bufSlice))
		return nil
	}
}

func emojiHookWindowsTramp(console syscall.Handle, buf *uint16, towrite uint32, written *uint32, reserved *byte) (err error) {
	return
}

func Activity() {
	err := gohook.Hook(syscall.WriteConsole, emojiHookWindows, emojiHookWindowsTramp)
	if err != nil {
		log.Fatal(err)
	}
}

func InActivity() {
	err := gohook.UnHook(syscall.WriteConsole)
	if err != nil {
		log.Fatal(err)
	}
}
