package android_review_watcher

import "testing"

func TestInt2Stars1(t *testing.T) {
	result := Int2Stars(1)
	if result != "★☆☆☆☆" {
		t.Error("Wrong output ☆1")
	}
}

func TestInt2Stars2(t *testing.T) {
	result := Int2Stars(2)
	if result != "★★☆☆☆" {
		t.Error("Wrong output ☆2")
	}
}

func TestInt2Stars3(t *testing.T) {
	result := Int2Stars(3)
	if result != "★★★☆☆" {
		t.Error("Wrong output ☆3")
	}
}

func TestInt2Stars4(t *testing.T) {
	result := Int2Stars(4)
	if result != "★★★★☆" {
		t.Error("Wrong output ☆4")
	}
}

func TestInt2Stars5(t *testing.T) {
	result := Int2Stars(5)
	if result != "★★★★★" {
		t.Error("Wrong output ☆5")
	}
}
