package config

import (
	"sync"

	"k8s.io/klog"

	"github.com/kubeedge/beehive/pkg/common/config"
	"github.com/kubeedge/kubeedge/edge/pkg/common/modules"
)

const (
	defaultSyncInterval = 60
)

var c Configure
var once sync.Once

var Connected = false

type Configure struct {
	// SendModuleGroupName is the name of the group to which we send the message
	SendModuleGroupName string

	// SendModuleName is the name of send module for remote query
	SendModuleName string

	SyncInterval int
}

func InitConfigure() {
	once.Do(func() {
		groupName, err := config.CONFIG.GetValue("metamanager.context-send-group").ToString()
		if err != nil || groupName == "" {
			// Guaranteed forward compatibility @kadisi
			groupName = modules.HubGroup
			klog.Infof("can not get metamanager.context-send-group key , use default %v", groupName)
		}

		edgeSite, err := config.CONFIG.GetValue("metamanager.edgesite").ToBool()
		if err != nil {
			// Guaranteed forward compatibility @kadisi
			edgeSite = false
			klog.Infof("can not get metamanager.edgesite key , use default %v", edgeSite)
		}
		moduleName, err := config.CONFIG.GetValue("metamanager.context-send-module").ToString()
		if err != nil || moduleName == "" {
			moduleName = "websocket"
		}

		Connected = edgeSite

		syncInterval, err := config.CONFIG.GetValue("meta.sync.podstatus.interval").ToInt()
		if err != nil || syncInterval < defaultSyncInterval {
			// Guaranteed forward compatibility @kadisi
			syncInterval = defaultSyncInterval
			klog.Infof("can not get meta.sync.podstatus.interval key, use default %v", syncInterval)
		}

		c = Configure{
			SendModuleGroupName: groupName,
			SendModuleName:      moduleName,
			SyncInterval:        syncInterval,
		}
		klog.Infof("init common config successfully，config info %++v", c)
	})
}

func Get() *Configure {
	return &c
}
