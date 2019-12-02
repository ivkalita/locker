package locker

func Semaphore(capacity int) semaphore {
	return semaphore{}
}

type semaphore struct{}
