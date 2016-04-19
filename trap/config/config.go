package config

import (
	"github.com/raincious/trap/trap/core/types"

	"time"
)

type Listen struct {
	Protocol types.String
	Setting  types.String
}

type Parameter struct {
	Format types.String
	Labels []types.String
}

type Command struct {
	Command    types.String
	Parameters []Parameter
}

type Commands map[types.String][]Command

type Server struct {
	Address    types.IPAddress
	Passphrase types.String
}

type Servers []Server
type Listens []Listen

type Config struct {
	Log types.String

	Listens Listens

	AttemptTimeout   types.UInt32
	AttemptMaxBytes  types.UInt32
	AttemptThershold types.UInt32
	AttemptExpire    types.UInt32
	AttemptRestrict  types.UInt32

	Commands Commands

	StatusInterface  types.IP
	StatusPort       types.UInt16
	StatusAccounts   map[types.String][]types.String
	StatusTLSCert    types.String
	StatusTLSCertKey types.String

	SyncInterface    types.IP
	SyncPort         types.UInt16
	SyncCert         types.String
	SyncCertKey      types.String
	SyncPass         types.String
	SyncLooseTimeout time.Duration
	SyncConnTimeout  time.Duration
	SyncReqTimeout   time.Duration
	SyncWith         Servers
}
