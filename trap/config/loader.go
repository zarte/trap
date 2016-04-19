package config

import (
	"github.com/raincious/trap/trap/core/types"

	"encoding/json"
	"io/ioutil"
	"os/exec"
	"regexp"
	"time"
)

type rawCommandConfig [][]types.String

type rawServerKVMap map[types.String]types.String

type rawConfig struct {
	Log types.String `json:"log"`

	Listens            []types.String                    `json:"listens"`
	AttemptTimeout     types.UInt32                      `json:"attempt_timeout"`
	AttemptMaxBytes    types.UInt32                      `json:"attempt_max_bytes"`
	AttemptThershold   types.UInt32                      `json:"attempt_thershold"`
	AttemptExpire      types.UInt32                      `json:"attempt_expire"`
	AttemptRestrict    types.UInt32                      `json:"attempt_restrict"`
	Commands           map[types.String]rawCommandConfig `json:"commands"`
	StatusInterface    types.IP                          `json:"status_interface"`
	StatusPort         types.UInt16                      `json:"status_port"`
	StatusAccounts     map[types.String][]types.String   `json:"status_accounts"`
	StatusTLSCert      types.String                      `json:"status_tls_certificate"`
	StatusTLSCertKey   types.String                      `json:"status_tls_certificate_key"`
	SyncInterface      types.String                      `json:"synchronize_interface"`
	SyncPort           types.UInt16                      `json:"synchronize_port"`
	SyncCertificate    types.String                      `json:"synchronize_certificate"`
	SyncCertificateKey types.String                      `json:"synchronize_certificate_key"`
	SyncPassphrase     types.String                      `json:"synchronize_passphrase"`
	SyncConnTimeout    types.UInt16                      `json:"synchronize_connection_timeout"`
	SyncLooseTimeout   types.UInt16                      `json:"synchronize_loose_timeout"`
	SyncReqTimeout     types.UInt16                      `json:"synchronize_request_timeout"`
	SyncWith           rawServerKVMap                    `json:"synchronize_with"`
}

func Load(filePath string) (*Config, *types.Throw) {
	content, err := ioutil.ReadFile(filePath)

	if err != nil {
		return nil, types.ConvertError(err)
	}

	return Parse(content)
}

func Parse(configStr []byte) (*Config, *types.Throw) {
	config := &Config{}
	rawConfig := &rawConfig{}

	reg, err := regexp.Compile("(\\/\\*\\*\\!(?msiU:.*)\\*\\/)")

	if err != nil {
		return nil, types.ConvertError(err)
	}

	rawParsedCfg := reg.ReplaceAll(configStr, []byte(""))

	err = json.Unmarshal(rawParsedCfg, rawConfig)

	if err != nil {
		return nil, types.ConvertError(err)
	}

	// Parse `Log` Field
	config.Log = rawConfig.Log

	// Parse `Listens` Field
	config.Listens = Listens{}

	for _, lsten := range rawConfig.Listens {
		lProtocol, lSetting := lsten.SpiltWith(":")

		config.Listens = append(config.Listens, Listen{
			Protocol: lProtocol,
			Setting:  lSetting,
		})
	}

	// Parse `AttemptTimeout` Field
	config.AttemptTimeout = rawConfig.AttemptTimeout

	// Parse `AttemptMaxBytes` Field
	config.AttemptMaxBytes = rawConfig.AttemptMaxBytes

	// Parse `AttemptThershold` Field
	config.AttemptThershold = rawConfig.AttemptThershold

	if config.AttemptThershold < 0 {
		config.AttemptThershold = 0
	}

	// Parse `AttemptExpire` Field
	config.AttemptExpire = rawConfig.AttemptExpire

	// Parse `AttemptRestrict` Field
	config.AttemptRestrict = rawConfig.AttemptRestrict

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

	// Parse `StatusTLSCert` and `StatusTLSCertKey` Field
	config.StatusTLSCert = rawConfig.StatusTLSCert
	config.StatusTLSCertKey = rawConfig.StatusTLSCertKey

	// Parse `StatusInterface` Field
	config.StatusInterface = rawConfig.StatusInterface

	// Parse `StatusPort` Field
	config.StatusPort = rawConfig.StatusPort

	if config.StatusPort < 0 {
		config.StatusPort = 0
	}

	// Parse `StatusAccounts` Field
	config.StatusAccounts = map[types.String][]types.String{}
	for pass, permissions := range rawConfig.StatusAccounts {
		permissionList := []types.String{}

		for _, permission := range permissions {
			permissionList = append(permissionList, permission)
		}

		config.StatusAccounts[pass] = permissionList
	}

	// Parse `SyncInterface` Field
	syncIfaceIP, syncIfaceIpErr := types.ConvertIPFromString(
		rawConfig.SyncInterface)

	if syncIfaceIpErr != nil {
		return nil, ErrParseInvalidItem.Throw(syncIfaceIP,
			"synchronize_interface")
	}

	config.SyncInterface = syncIfaceIP

	// Parse `SyncPort` Field
	config.SyncPort = rawConfig.SyncPort

	// Parse `SyncCert` Field
	config.SyncCert = rawConfig.SyncCertificate

	// Parse `SyncCertKey` Field
	config.SyncCertKey = rawConfig.SyncCertificateKey

	// Parse `SyncPassphrase` Field
	config.SyncPass = rawConfig.SyncPassphrase

	// Parse `SyncConnTimeout` Field
	config.SyncConnTimeout = time.Duration(rawConfig.SyncConnTimeout) * time.Second

	// Parse `SyncLooseTimeout` Field
	config.SyncLooseTimeout = time.Duration(rawConfig.SyncLooseTimeout) * time.Second

	// Parse `SyncReqTimeout` Field
	config.SyncReqTimeout = time.Duration(rawConfig.SyncReqTimeout) * time.Second

	// Parse `SyncWith` Field
	config.SyncWith = Servers{}

	for syServer, syPass := range rawConfig.SyncWith {
		syncServer := Server{}

		syIPStr, syPortStr := syServer.SpiltWith(":")

		syIP, syIPErr := types.ConvertIPFromString(syIPStr)

		if syIPErr != nil {
			return nil, ErrParseInvalidItem.Throw(syIPStr, "synchronize_with")
		}

		syncServer.Address.IP = syIP
		syncServer.Address.Port = syPortStr.Trim().UInt16()
		syncServer.Passphrase = syPass

		if syncServer.Address.IsEmpty() || syncServer.Address.Port == 0 {
			return nil, ErrParseInvalidItem.Throw(syServer, "synchronize_with")
		}

		if syncServer.Passphrase == "" {
			return nil, ErrParseInvalidItem.Throw(syncServer.Passphrase,
				"synchronize_with")
		}

		config.SyncWith = append(config.SyncWith, syncServer)
	}

	return config, nil
}
