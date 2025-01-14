package nettop

import (
	"fmt"

	"github.com/vishvananda/netns"
)

var (
	InitNetns = 0
)

// GetEntityByNetns get entity by netns, if netns was deleted asynchrously, return nil; otherwise return error
func GetEntityByNetns(nsinum int) (*Entity, error) {
	// if use nsinum 0, represent node level metrics
	if nsinum == 0 {
		return GetHostnetworlEntity()
	}
	v, found := nsCache.Get(fmt.Sprintf("%d", nsinum))
	if found {
		return v.(*Entity), nil
	}
	return nil, fmt.Errorf("entify for netns %d not found", nsinum)
}

func GetHostnetworlEntity() (*Entity, error) {
	return GetEntityByPid(1)
}

func GetHostnetworkNetnsFd() (int, error) {
	nsh, err := netns.GetFromPid(1)
	if err != nil {
		return 0, err
	}

	return int(nsh), nil
}

func GetEntityByPod(sandbox string) (*Entity, error) {
	v, found := podCache.Get(sandbox)
	if found {
		return v.(*Entity), nil
	}
	return nil, fmt.Errorf("entify for pod %s not found", sandbox)
}

func GetEntityByPid(pid int) (*Entity, error) {
	v, found := pidCache.Get(fmt.Sprintf("%d", pid))
	if found {
		return v.(*Entity), nil
	}
	return nil, fmt.Errorf("entify for process %d not found", pid)
}

func GetAllEntity() []*Entity {
	v := nsCache.Items()

	res := []*Entity{}
	for _, item := range v {
		et := item.Object.(*Entity)
		// filter unknow netns, such as extra test netns created by cni-plugin, cilium/calico etc..
		if et == nil || et.GetPodName() == "unknow" {
			continue
		}
		res = append(res, et)
	}

	return res
}

func GetEntityNetnsByPid(pid int) (int, error) {
	return getNsInumByPid(pid)
}
