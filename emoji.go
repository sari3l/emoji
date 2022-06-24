package emoji

import (
	"github.com/brahma-adshonor/gohook"
	"syscall"
)

func emojiHook(fd int, p []byte) (n int, err error) {
	pPtr := &p
	pLen := len(p)
	var waitForWriteBytes []byte
	var emojiBytes []byte
	lock := false
	for i := 0; i < pLen; i++ {
		if p[i] != 58 && !lock {
			waitForWriteBytes = append(waitForWriteBytes, p[i])
			continue
		}
		if p[i] == 58 && !lock {
			lock = true
			emojiBytes = append(emojiBytes, p[i])
			continue
		}
		if p[i] == 58 && lock {
			if len(emojiBytes) == 1 {
				waitForWriteBytes = append(waitForWriteBytes, emojiBytes...)
				continue
			}
			lock = false
			emojiBytes = append(emojiBytes, p[i])
			emojiStr := emojiMap[string(emojiBytes)]
			if emojiStr != "" {
				waitForWriteBytes = append(waitForWriteBytes, []byte(emojiStr)...)
			} else {
				waitForWriteBytes = append(waitForWriteBytes, emojiBytes...)
			}
			emojiBytes = []byte{}
			continue
		}
		emojiBytes = append(emojiBytes, p[i])
		if p[i] == 32 {
			lock = false
			waitForWriteBytes = append(waitForWriteBytes, emojiBytes...)
			emojiBytes = []byte{}
		}
	}
	*pPtr = waitForWriteBytes
	wLen := len(waitForWriteBytes)
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
	gohook.Hook(syscall.Write, emojiHook, emojiHookTramp)
}

func InActivity() {
	gohook.UnHook(syscall.Write)
}
