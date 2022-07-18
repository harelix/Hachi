package storage

import "github.com/rills-ai/Hachi/pkg/messaging"

//Adds a new KV Store Bucket
func Add() {
	messaging.Get().NC.JetStream()
}

func Get() {

}

//Puts a value into a key
func Put() {

}

func Update() {

}

func Watch() {

}

func History() {

}

func Delete() {

}

func Purge() {

}
