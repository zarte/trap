package config

import (
    "github.com/raincious/trap/trap/core/types"

    "net"
    "regexp"
    "os/exec"
    "io/ioutil"
    "encoding/json"
)

type rawCommandConfig [][]types.String

type rawServerKVMap map[types.String]types.String

type rawConfig struct {
    Log                 types.String        `json:"log"`

    Listens             []types.String      `json:"listen_ports"`

    AttemptThershold    types.UInt32        `json:"attempt_thershold"`
    AttemptExpire       types.UInt32        `json:"attempt_expire"`

    Commands            map[types.String]rawCommandConfig   `json:"commands"`

    StatusHost          types.String        `json:"status_host"`
    StatusInterface     types.IP            `json:"status_interface"`
    StatusPort          types.UInt16        `json:"status_port"`
    StatusAccounts      map[types.String][]types.String     `json:"status_accounts"`
    StatusAllowedIPs    []types.String      `json:"status_accessable_ip_prefixes"`

    SyncPort            types.UInt16        `json:"synchronize_port"`
    SyncPassphrase      types.String        `json:"synchronize_passphrase"`
    SyncWith            rawServerKVMap      `json:"synchronize_with"`
}

func Load(filePath string) (*Config, *types.Throw) {
    content, err := ioutil.ReadFile(filePath)

    if err != nil {
        return nil, types.ConvertError(err)
    }

    return Parse(content)
}

func Parse(configStr []byte) (*Config, *types.Throw) {
    config          := &Config{}
    rawConfig       := &rawConfig{}

    reg, err        := regexp.Compile("(\\/\\*\\*\\!(?msiU:.*)\\*\\/)")

    if err != nil {
        return nil, types.ConvertError(err)
    }

    rawParsedCfg    := reg.ReplaceAll(configStr, []byte(""))

    err             =  json.Unmarshal(rawParsedCfg, rawConfig)

    if err != nil {
        return nil, types.ConvertError(err)
    }

    // Parse `Log` Field
    config.Log      =  rawConfig.Log

    // Parse `Listens` Field
    config.Listens  =  Listens{}

    for _, lsten := range rawConfig.Listens {
        lSettings, lAdditional  :=  lsten.SpiltWith("|")

        lType, lHost            :=  lSettings.SpiltWith(":")

        listenItem              :=  Listen{}

        lPort, lIP              :=  types.String(""), types.String("")

        // If `lHost` is empty, use `lType` as the host
        if lHost == "" {
            lPort, lIP          =   lType.SpiltWith("@")
            listenItem.Type     =   "tcp"
        } else {
            lPort, lIP          =   lHost.SpiltWith("@")
            listenItem.Type     =   lType.Lower().Trim()
        }

        if lIP != "" {
            listenItem.IP       =   net.ParseIP(lIP.Lower().Trim().String())
            listenItem.Port     =   lPort.Trim().UInt16()
        } else {
            listenItem.IP       =   net.ParseIP("0.0.0.0")
            listenItem.Port     =   lPort.Trim().UInt16()
        }

        if listenItem.Type == "" || listenItem.Port == 0 || listenItem.IP == nil {
            return nil, ErrParseInvalidItem.Throw(lsten, "listen_ports")
        }

        listenItem.Additional   =   lAdditional

        config.Listens          =   append(config.Listens, listenItem)
    }

    // Parse `AttemptThershold` Field
    config.AttemptThershold = rawConfig.AttemptThershold

    if config.AttemptThershold < 0 {
        config.AttemptThershold = 0
    }

    // Parse `AttemptExpire` Field
    config.AttemptExpire = rawConfig.AttemptExpire

    if config.AttemptExpire < 0 {
        config.AttemptExpire = 0
    }

    // Parse `Commands` Fields
    config.Commands = Commands{}

    paramReg, err := regexp.Compile("\\$\\(\\([[:word:]]+\\)\\)")

    for cmdType, cmds := range rawConfig.Commands {
        cType := cmdType.Trim().Lower()

        for _, cmdParams := range cmds {
            cmdLen := len(cmdParams)

            if cmdLen < 1 {
                return nil, ErrInvalidCmdItem.Throw(cmdType.String(),
                    "commands")
            }

            cmdItem := Command{}

            file, fErr := exec.LookPath(cmdParams[0].String())

            if fErr != nil {
                return nil, types.ConvertError(fErr)
            }

            cmdItem.Command = types.String(file)

            if cmdLen > 1 {
                for cmdIdx := 1; cmdIdx < cmdLen; cmdIdx++ {
                    cmdString := cmdParams[cmdIdx].Trim()
                    cmdParamLables := []types.String{}

                    for _, maths := range paramReg.FindAllStringSubmatch(
                        cmdString.String(), -1) {
                        cmdParamLables = append(cmdParamLables,
                            types.String(maths[0]))
                    }

                    cmdParam := Parameter{
                        Format: cmdParams[cmdIdx].Trim(),
                        Labels: cmdParamLables,
                    }

                    if cmdParam.Format == "" {
                        continue
                    }

                    cmdItem.Parameters = append(cmdItem.Parameters, cmdParam)
                }
            }

            config.Commands[cType] = append(config.Commands[cType], cmdItem)
        }
    }

    // Parse `StatusHost` Field
    config.StatusHost = rawConfig.StatusHost

    // Parse `StatusInterface` Field
    config.StatusInterface = rawConfig.StatusInterface

    // Parse `StatusPort` Field
    config.StatusPort = rawConfig.StatusPort

    if config.StatusPort < 0 {
        config.StatusPort = 0
    }

    // Parse `StatusAccounts` Field
    config.StatusAccounts               =   map[types.String][]types.String{}
    for pass, permissions := range rawConfig.StatusAccounts {
        permissionList                  :=  []types.String{}

        for _, permission := range permissions {
            permissionList              =   append(permissionList, permission)
        }

        config.StatusAccounts[pass]     =   permissionList
    }

    // Parse `StatusAllowedIPs` Field
    for _, sIP := range rawConfig.StatusAllowedIPs {
        config.StatusAllowedIPs = append(config.StatusAllowedIPs,
                                            sIP.Trim().Lower())
    }

    // Parse `SyncPort` Field
    config.SyncPort = rawConfig.SyncPort

    // Parse `SyncPassphrase` Field
    config.SyncPassphrase = rawConfig.SyncPassphrase

    // Parse `SyncWith` Field
    config.SyncWith = Servers{}

    for syServer, syPass := range rawConfig.SyncWith {
        syncServer              := Server{}

        syHost, syPort          := syServer.SpiltWith(":")

        syncServer.Host         =  syHost.Trim().Lower()
        syncServer.Port         =  syPort.Trim().UInt16()
        syncServer.Passphrase   =  syPass

        if syncServer.Host == "" || syncServer.Port == 0 {
            return nil, ErrParseInvalidItem.Throw(syServer, "synchronize_with")
        }

        if syncServer.Passphrase == "" {
            return nil, ErrParseInvalidItem.Throw(syncServer.Passphrase,
                "synchronize_with")
        }

        config.SyncWith = append(config.SyncWith, syncServer)
    }

    return config, nil;
}