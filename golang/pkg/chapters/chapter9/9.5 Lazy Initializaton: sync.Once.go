package chapter9

import (
	"image"
	"sync"
)

var (
	newIcons map[string]image.Image
	newRWMu  sync.RWMutex
)

func loadIconsUnsafe() {
	newIcons = map[string]image.Image{
		"spades.png": loadIcon("spades.png"),
		"hearts.png": loadIcon("hearts.png"),
		// ...
	}
}

func RawLazyIcon(name string) image.Image {
	if newIcons == nil {
		loadIconsUnsafe()
	}
	return newIcons[name]
}

func LazyIcon(name string) image.Image {
	mu.Lock()
	defer mu.Unlock()
	if newIcons == nil {
		loadIconsUnsafe()
	}
	return newIcons[name]
}

func ModifiedIcon(name string) image.Image {
	newRWMu.RLock()
	if icons != nil {
		icon := newIcons[name]
		newRWMu.RUnlock()
		return icon
	}
	newRWMu.RUnlock()

	mu.Lock()
	if newIcons == nil {
		loadIconsUnsafe()
	}
	icon := newIcons[name]
	mu.Unlock()
	return icon

}

var loadIconsOnce sync.Once

func IconWithOnce(name string) image.Image {
	loadIconsOnce.Do(loadIconsUnsafe)
	return newIcons[name]
}
