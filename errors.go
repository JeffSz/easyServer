package esayServer

type Error struct{
	error
	err string
	code int
}
