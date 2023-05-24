package service

type Subject interface {
	register(observer AuthStateObserver)
	deregister(observer AuthStateObserver)
	notifyAll()
}

type AuthState struct {
	observer     AuthStateObserver
	isAuthorized bool
}

type AuthStateObserver interface {
	update()
	getID() string
}
