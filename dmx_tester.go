package eltee

import (
	"github.com/eyethereal/go-config"
)

type DmxScene struct {
	data          []byte
	start_channel int
}

func NewDmxScene(node *config.AclNode) *DmxScene {
	scene := &DmxScene{}

	list := node.ChildAsIntList("data")
	scene.data = make([]byte, len(list))
	for i := 0; i < len(list); i++ {
		scene.data[i] = byte(list[i])
	}

	scene.start_channel = node.ChildAsInt("start_channel")
	return scene
}

type DmxTester struct {
	scenes map[string]*DmxScene

	toRun string
}

func NewDmxTester(node *config.AclNode) *DmxTester {
	tester := &DmxTester{
		scenes: make(map[string]*DmxScene),
	}

	scNode := node.Child("scenes")
	scNode.ForEachOrderedChild(func(name string, n *config.AclNode) {
		scene := NewDmxScene(n)
		tester.scenes[name] = scene
	})

	tester.toRun = node.ChildAsString("run")
	return tester
}

func (tester *DmxTester) HasTest() bool {
	if tester.toRun == "" {
		return false
	}

	return tester.scenes[tester.toRun] != nil
}

func (tester *DmxTester) UpdateFrame(frame []byte) {

	scene := tester.scenes[tester.toRun]
	if scene == nil {
		return
	}

	chIx := scene.start_channel
	if chIx > 0 {
		// Adjust to 0 offset
		chIx--
	}

	for i := 0; i < len(scene.data); i++ {
		frame[chIx] = scene.data[i]
		chIx++
	}
}
