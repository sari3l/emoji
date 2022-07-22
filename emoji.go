package emoji

import (
	"github.com/brahma-adshonor/gohook"
	"log"
	"reflect"
	"syscall"
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

func emojiHook(fd int, p []byte) (n int, err error) {
	pPtr := &p
	pLen := len(p)
	chunk := doReplaceEmoji(p)
	*pPtr = chunk
	wLen := len(chunk)
	w, err := emojiHookTramp(fd, p)
	if w == wLen {
		return pLen, err
	} else {
		return w, err
	}
}

func emojiHookTramp(fd int, p []byte) (n int, err error) {
	return 0, nil
}

func Activity() {
	err := gohook.Hook(syscall.Write, emojiHook, emojiHookTramp)
	if err != nil {
		log.Fatal(err)
	}
}

func InActivity() {
	err := gohook.UnHook(syscall.Write)
	if err != nil {
		log.Fatal(err)
	}
}
