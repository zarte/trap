package config

import (
    "github.com/raincious/trap/trap/core/types"

    "net"
)

type Listen struct {
    Type                types.String
    IP                  net.IP
    Port                types.UInt16
    Additional          types.String
}

type Parameter struct {
    Format              types.String
    Labels              []types.String
}

type Command struct {
    Command             types.String
    Parameters          []Parameter
}

type Commands map[types.String][]Command

type Server struct {
    Host                types.String
    Port                types.UInt16
    Passphrase          types.String
}

type Servers []Server
type Listens []Listen

type Config struct {
    Log                 types.String

    Listens             Listens

    AttemptThershold    types.UInt32
    AttemptExpire       types.UInt32

    Commands            Commands

    StatusHost          types.String
    StatusInterface     types.IP
    StatusPort          types.UInt16
    StatusAccounts      map[types.String][]types.String
    StatusAllowedIPs    []types.String

    SyncPort            types.UInt16
    SyncPassphrase      types.String
    SyncWith            Servers
}