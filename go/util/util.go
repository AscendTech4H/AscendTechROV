//Package util contains assorted utility code
package util

// UhOh throws error if there is error
func UhOh(err error) {
	if err != nil {
		panic(err)
	}
}
