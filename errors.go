package easyServer

type Error struct{
	error
	err string
	code int
}
